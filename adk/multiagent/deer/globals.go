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
	"fmt"
	"strings"

	"github.com/RanFeng/ilog"
	"github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/tool"
)

var (
	// Global Cache
	GlobalPrompts       map[string]string
	GlobalResearchTools []tool.BaseTool
	GlobalCoderTools    []tool.BaseTool
)

func InitSystem(ctx context.Context) error {
	// 1. Load Config
	LoadDeerConfig(ctx)

	// 2. Init MCP Clients
	InitMCP()

	// 3. Load Prompts
	if err := loadPrompts(ctx); err != nil {
		return err
	}

	// 4. Load Tools
	if err := loadTools(ctx); err != nil {
		return err
	}

	return nil
}

func loadPrompts(ctx context.Context) error {
	GlobalPrompts = make(map[string]string)
	promptNames := []string{
		"planner", "researcher", "coder", "reporter",
		"podcast_script_writer", "ppt_composer", "coordinator",
	}

	for _, name := range promptNames {
		content, err := GetPromptTemplate(ctx, name)
		if err != nil {
			return fmt.Errorf("failed to load prompt %s: %v", name, err)
		}
		GlobalPrompts[name] = content
	}
	return nil
}

func loadTools(ctx context.Context) error {
	GlobalResearchTools = []tool.BaseTool{}
	GlobalCoderTools = []tool.BaseTool{}

	for name, cli := range MCPServer {
		ts, err := mcp.GetTools(ctx, &mcp.Config{Cli: cli})
		if err != nil {
			ilog.EventError(ctx, err, "get_tools_error", "server", name)
			continue
		}

		// Add to Research Tools (All tools)
		GlobalResearchTools = append(GlobalResearchTools, ts...)

		// Add to Coder Tools (Python only)
		if strings.HasPrefix(name, "python") {
			GlobalCoderTools = append(GlobalCoderTools, ts...)
		}
	}

	if len(GlobalCoderTools) == 0 {
		ilog.EventWarn(ctx, "no_python_tools_found_for_coder")
	}

	return nil
}
