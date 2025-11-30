/*
 * Copyright 2025 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/RanFeng/ilog"
	"github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// NewInvestigator creates a background investigator agent
// This agent is responsible for initial background search if enabled
func NewInvestigator(ctx context.Context) (adk.Agent, error) {
	return &InvestigatorAgent{}, nil
}

type InvestigatorAgent struct{}

func (a *InvestigatorAgent) Run(ctx context.Context, input *adk.AgentInput, opts ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent] {
	iter, gen := adk.NewAsyncIteratorPair[*adk.AgentEvent]()

	// Use goroutine to avoid blocking
	go func() {
		defer gen.Close()

		// 1. Get User Query
		if len(input.Messages) == 0 {
			gen.Send(&adk.AgentEvent{
				Output: &adk.AgentOutput{
					MessageOutput: &adk.MessageVariant{
						Message: &schema.Message{
							Role:    schema.Assistant,
							Content: "",
						},
					},
				},
			})
			return
		}
		lastMsg := input.Messages[len(input.Messages)-1]
		query := lastMsg.Content

		// 2. Find Search Tool
		var searchTool tool.InvokableTool
		for name, cli := range MCPServer {
			if searchTool != nil {
				break
			}
			ts, err := mcp.GetTools(ctx, &mcp.Config{Cli: cli})
			if err != nil {
				ilog.EventError(ctx, err, "investigator_get_tools_error", "server", name)
				continue
			}
			for _, t := range ts {
				info, _ := t.Info(ctx)
				if strings.HasSuffix(info.Name, "search") {
					searchTool, _ = t.(tool.InvokableTool)
					break
				}
			}
		}

		if searchTool == nil {
			ilog.EventWarn(ctx, "investigator_no_search_tool_found")
			gen.Send(&adk.AgentEvent{
				Output: &adk.AgentOutput{
					MessageOutput: &adk.MessageVariant{
						Message: &schema.Message{
							Role:    schema.Assistant,
							Content: "No search tool available for background investigation.",
						},
					},
				},
			})
			return
		}

		// 3. Run Search
		args := map[string]any{
			"query": query,
		}
		argsBytes, err := json.Marshal(args)
		if err != nil {
			gen.Send(&adk.AgentEvent{
				Err: err,
			})
			return
		}

		result, err := searchTool.InvokableRun(ctx, string(argsBytes))
		if err != nil {
			ilog.EventError(ctx, err, "investigator_search_error")
			gen.Send(&adk.AgentEvent{
				Output: &adk.AgentOutput{
					MessageOutput: &adk.MessageVariant{
						Message: &schema.Message{
							Role:    schema.Assistant,
							Content: fmt.Sprintf("Background investigation failed: %v", err),
						},
					},
				},
			})
			return
		}

		// 4. Return Result
		gen.Send(&adk.AgentEvent{
			Output: &adk.AgentOutput{
				MessageOutput: &adk.MessageVariant{
					Message: &schema.Message{
						Role:    schema.Assistant,
						Content: fmt.Sprintf("Background Investigation Result: %s", result),
					},
				},
			},
		})
	}()

	return iter
}

func (a *InvestigatorAgent) Name(ctx context.Context) string {
	return "investigator"
}

func (a *InvestigatorAgent) Description(ctx context.Context) string {
	return "Conducts initial background investigation using search tools. Useful for gathering context before planning."
}
