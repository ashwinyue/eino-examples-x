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
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/cloudwego/eino-examples/adk/common/prints"
	"github.com/cloudwego/eino-examples/adk/common/trace"
)

// createAgents creates a new set of agents for a request, using cached resources
func createAgents(ctx context.Context) (adk.Agent, error) {
	investigator, err := NewInvestigator(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create investigator: %v", err)
	}

	planner, err := NewPlanner(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create planner: %v", err)
	}

	researcher, err := NewResearcher(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create researcher: %v", err)
	}

	coder, err := NewCoder(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create coder: %v", err)
	}

	reporter, err := NewReporter(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create reporter: %v", err)
	}

	podcastWriter, err := NewPodcastScriptWriter(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create podcast writer: %v", err)
	}

	pptComposer, err := NewPPTComposer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create ppt composer: %v", err)
	}

	sv, err := NewCoordinator(ctx, []adk.Agent{
		planner,
		researcher,
		coder,
		reporter,
		podcastWriter,
		pptComposer,
		investigator,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create coordinator: %v", err)
	}
	return sv, nil
}

func runServer() {
	r := gin.Default()

	// CORS configuration to match the original Hertz behavior
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Authorization", "Content-Length", "X-CSRF-Token", "Token", "session", "X_Requested_With", "Accept", "Origin", "Host", "Connection", "Accept-Encoding", "Accept-Language", "DNT", "X-CustomHeader", "Keep-Alive", "User-Agent", "X-Requested-With", "If-Modified-Since", "Cache-Control", "Content-Type", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Cache-Control", "Content-Language", "Content-Type", "Expires", "Last-Modified", "Pragma", "FooBar"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	r.Use(cors.New(config))

	r.GET("/api/config", func(c *gin.Context) {
		basic := []string{}
		if m := os.Getenv("OPENAI_MODEL"); m != "" {
			basic = append(basic, m)
		}
		if m := os.Getenv("ARK_MODEL"); m != "" {
			basic = append(basic, m)
		}
		c.JSON(200, gin.H{
			"rag":    gin.H{"provider": ""},
			"models": gin.H{"basic": basic, "reasoning": []string{}},
		})
	})

	r.POST("/v1/chat/completions", func(c *gin.Context) {
		var req struct {
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Get the last user message
		var query string
		for i := len(req.Messages) - 1; i >= 0; i-- {
			if req.Messages[i].Role == "user" {
				query = req.Messages[i].Content
				break
			}
		}
		if query == "" {
			c.JSON(400, gin.H{"error": "no user message found"})
			return
		}

		ctx := c.Request.Context()

		// Setup tracing
		traceCloseFn, startSpanFn := trace.AppendCozeLoopCallbackIfConfigured(ctx)
		defer traceCloseFn(ctx)

		sv, err := createAgents(ctx)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		runner := adk.NewRunner(ctx, adk.RunnerConfig{
			Agent:           sv,
			EnableStreaming: true,
		})

		ctx, endSpanFn := startSpanFn(ctx, "DeerGo-ADK-Server", query)
		iter := runner.Query(ctx, query)

		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		c.Stream(func(w io.Writer) bool {
			event, hasEvent := iter.Next()
			if !hasEvent {
				endSpanFn(ctx, nil) // End span when done
				return false
			}

			if event.Output != nil && event.Output.MessageOutput != nil {
				if m := event.Output.MessageOutput.Message; m != nil {
					if m.Content != "" {
						c.SSEvent("message", m.Content)
					}
					if len(m.ToolCalls) > 0 {
						for _, tc := range m.ToolCalls {
							payload := map[string]any{
								"id":   tc.ID,
								"type": tc.Type,
								"name": tc.Function.Name,
								"args": tc.Function.Arguments,
							}
							b, _ := json.Marshal(payload)
							c.SSEvent("tool_call_start", string(b))
							c.SSEvent("tool_call_end", string(b))
						}
					}
				} else if s := event.Output.MessageOutput.MessageStream; s != nil {
					toolMap := map[int][]*schema.Message{}
					for {
						chunk, err := s.Recv()
						if err != nil {
							if err == io.EOF {
								break
							}
							c.SSEvent("error", err.Error())
							break
						}
						if chunk.Content != "" {
							c.SSEvent("message", chunk.Content)
							if chunk.Role == schema.Tool {
								c.SSEvent("tool_result", chunk.Content)
							}
						}
						if len(chunk.ToolCalls) > 0 {
							for _, tc := range chunk.ToolCalls {
								if tc.Index == nil {
									continue
								}
								payload := map[string]any{
									"id":    tc.ID,
									"type":  tc.Type,
									"name":  tc.Function.Name,
									"args":  tc.Function.Arguments,
									"index": *tc.Index,
								}
								b, _ := json.Marshal(payload)
								c.SSEvent("tool_call_delta", string(b))
								toolMap[*tc.Index] = append(toolMap[*tc.Index], &schema.Message{
									Role: chunk.Role,
									ToolCalls: []schema.ToolCall{{
										ID:    tc.ID,
										Type:  tc.Type,
										Index: tc.Index,
										Function: schema.FunctionCall{
											Name:      tc.Function.Name,
											Arguments: tc.Function.Arguments,
										},
									}},
								})
							}
						}
					}
					for _, msgs := range toolMap {
						m, err := schema.ConcatMessages(msgs)
						if err != nil {
							c.SSEvent("error", err.Error())
							continue
						}
						payload := map[string]any{
							"name": m.ToolCalls[0].Function.Name,
							"args": m.ToolCalls[0].Function.Arguments,
						}
						b, _ := json.Marshal(payload)
						c.SSEvent("tool_call_start", string(b))
						c.SSEvent("tool_call_end", string(b))
					}
				}
			}
			if event.Action != nil {
				if event.Action.TransferToAgent != nil {
					c.SSEvent("transfer_to_agent", event.Action.TransferToAgent.DestAgentName)
				}
				if event.Action.Interrupted != nil {
					b, _ := json.Marshal(event.Action.Interrupted.InterruptContexts)
					c.SSEvent("interrupt_options", string(b))
				}
				if event.Action.Exit {
					c.SSEvent("exit", "")
				}
			}
			if event.Err != nil {
				c.SSEvent("error", event.Err.Error())
			}
			return true
		})
	})

	r.POST("/api/chat/stream", func(c *gin.Context) {
		var req struct {
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
			ThreadID string `json:"thread_id"`
			Locale   string `json:"locale"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		var query string
		for i := len(req.Messages) - 1; i >= 0; i-- {
			if req.Messages[i].Role == "user" {
				query = req.Messages[i].Content
				break
			}
		}
		if query == "" {
			c.JSON(400, gin.H{"error": "no user message found"})
			return
		}
		if req.ThreadID == "" {
			req.ThreadID = fmt.Sprintf("thread_%d", time.Now().UnixNano())
		}

		if os.Getenv("OPENAI_API_KEY") == "" && os.Getenv("ARK_API_KEY") == "" {
			c.Header("Content-Type", "text/event-stream")
			c.Header("Cache-Control", "no-cache")
			c.Header("Connection", "keep-alive")
			write := func(event string, data any) {
				b, _ := json.Marshal(data)
				c.SSEvent(event, string(b))
			}
			c.Stream(func(w io.Writer) bool {
				write("message_chunk", map[string]any{
					"id":        fmt.Sprintf("run-coordinator-%d", time.Now().UnixNano()),
					"thread_id": req.ThreadID,
					"agent":     "coordinator",
					"role":      "assistant",
					"content":   "模型未配置，无法连接。请在 zsh 中导出环境变量，例如：\nexport OPENAI_API_KEY=你的key OPENAI_MODEL=gpt-4o\n或使用 Ark：\nexport MODEL_TYPE=ark ARK_API_KEY=你的key ARK_MODEL=ep-xxx ARK_BASE_URL=https://ark.cn-beijing.volces.com\n并重启后端服务。",
				})
				return false
			})
			return
		}

		ctx := c.Request.Context()
		traceCloseFn, startSpanFn := trace.AppendCozeLoopCallbackIfConfigured(ctx)
		defer traceCloseFn(ctx)
		sv, err := createAgents(ctx)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		runner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: sv, EnableStreaming: true})
		ctx, endSpanFn := startSpanFn(ctx, "DeerFlow-ADK-Server", query)
		iter := runner.Query(ctx, query)

		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		write := func(event string, data any) {
			b, _ := json.Marshal(data)
			c.SSEvent(event, string(b))
		}

		msgIDs := map[string]string{}
		getMsgID := func(agent string) string {
			if id, ok := msgIDs[agent]; ok {
				return id
			}
			id := fmt.Sprintf("run-%s-%d", agent, time.Now().UnixNano())
			msgIDs[agent] = id
			return id
		}

		c.Stream(func(w io.Writer) bool {
			e, ok := iter.Next()
			if !ok {
				write("message_chunk", map[string]any{
					"id":            getMsgID("coordinator"),
					"thread_id":     req.ThreadID,
					"agent":         "coordinator",
					"role":          "assistant",
					"finish_reason": "stop",
				})
				endSpanFn(ctx, nil)
				return false
			}
			base := map[string]any{
				"id": getMsgID(func() string {
					if e.AgentName == "" {
						return "coordinator"
					} else {
						return e.AgentName
					}
				}()),
				"thread_id": req.ThreadID,
				"agent": func() string {
					if e.AgentName == "" {
						return "coordinator"
					} else {
						return e.AgentName
					}
				}(),
				"role":          "assistant",
				"finish_reason": "",
			}
			if e.Output != nil && e.Output.MessageOutput != nil {
				if m := e.Output.MessageOutput.Message; m != nil {
					role := "assistant"
					if m.Role == schema.Tool {
						if e.AgentName == "coordinator" {
							role = "assistant"
						} else {
							role = "tool"
						}
					}
					write("message_chunk", map[string]any{
						"id": base["id"], "thread_id": base["thread_id"], "agent": base["agent"], "role": role,
						"content": m.Content,
					})
					if len(m.ToolCalls) > 0 {
						toolCalls := make([]map[string]any, 0, len(m.ToolCalls))
						for _, tc := range m.ToolCalls {
							args := tc.Function.Arguments
							var argsObj map[string]any
							if err := json.Unmarshal([]byte(args), &argsObj); err != nil {
								argsObj = map[string]any{}
							}
							toolCalls = append(toolCalls, map[string]any{
								"type": "tool_call", "id": tc.ID, "name": tc.Function.Name, "args": argsObj,
							})
						}
						write("tool_calls", map[string]any{
							"id": base["id"], "thread_id": base["thread_id"], "agent": base["agent"], "role": "assistant",
							"finish_reason":    "tool_calls",
							"tool_calls":       toolCalls,
							"tool_call_chunks": []any{},
						})
					}
				} else if s := e.Output.MessageOutput.MessageStream; s != nil {
					for {
						chunk, err := s.Recv()
						if err != nil {
							if err == io.EOF {
								break
							}
							write("message_chunk", map[string]any{
								"id": base["id"], "thread_id": base["thread_id"], "agent": base["agent"], "role": "assistant",
								"content": err.Error(),
							})
							break
						}
						role := "assistant"
						if chunk.Role == schema.Tool {
							if e.AgentName == "coordinator" {
								role = "assistant"
							} else {
								role = "tool"
							}
						}
						if chunk.Content != "" {
							write("message_chunk", map[string]any{
								"id": base["id"], "thread_id": base["thread_id"], "agent": base["agent"], "role": role,
								"content": chunk.Content,
							})
						}
						if len(chunk.ToolCalls) > 0 {
							chunks := make([]map[string]any, 0, len(chunk.ToolCalls))
							for _, tc := range chunk.ToolCalls {
								idx := 0
								if tc.Index != nil {
									idx = *tc.Index
								}
								chunks = append(chunks, map[string]any{
									"type": "tool_call_chunk", "index": idx, "id": tc.ID, "name": tc.Function.Name, "args": tc.Function.Arguments,
								})
							}
							write("tool_call_chunks", map[string]any{
								"id": base["id"], "thread_id": base["thread_id"], "agent": base["agent"], "role": "assistant",
								"tool_call_chunks": chunks,
							})
						}
					}
				}
			}
			if e.Action != nil {
				if e.Action.Interrupted != nil {
					write("interrupt", map[string]any{
						"id": base["id"], "thread_id": base["thread_id"], "agent": base["agent"], "role": "assistant",
						"options": e.Action.Interrupted.InterruptContexts,
					})
				}
				if e.Action.Exit {
					write("message_chunk", map[string]any{
						"id": base["id"], "thread_id": base["thread_id"], "agent": base["agent"], "role": "assistant",
						"finish_reason": "stop",
					})
				}
			}
			if e.Err != nil {
				write("message_chunk", map[string]any{
					"id": base["id"], "thread_id": base["thread_id"], "agent": base["agent"], "role": "assistant",
					"content": e.Err.Error(),
				})
			}
			return true
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	r.Run(":" + port)
}

func runConsole() {
	ctx := context.Background()
	if os.Getenv("OPENAI_API_KEY") == "" && os.Getenv("ARK_API_KEY") == "" {
		fmt.Println("Please set OPENAI_API_KEY or ARK_API_KEY environment variable.")
		return
	}

	// Setup tracing
	traceCloseFn, startSpanFn := trace.AppendCozeLoopCallbackIfConfigured(ctx)
	defer traceCloseFn(ctx)

	sv, err := createAgents(ctx)
	if err != nil {
		log.Fatalf("failed to initialize agents: %v", err)
	}

	// Create Runner
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           sv,
		EnableStreaming: true,
	})

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("请输入你的需求： ")
	query, _ := reader.ReadString('\n')
	query = strings.TrimSpace(query)

	if query == "" {
		return
	}

	fmt.Println("--------------------------------------------------")

	ctx, endSpanFn := startSpanFn(ctx, "DeerGo-ADK-Console", query)
	iter := runner.Query(ctx, query)

	var lastMessage adk.Message
	for {
		event, hasEvent := iter.Next()
		if !hasEvent {
			break
		}

		// Print events using the common helper
		prints.Event(event)

		if event.Output != nil {
			lastMessage, _, err = adk.GetMessage(event)
		}
	}

	endSpanFn(ctx, lastMessage)

	// wait for all span to be ended
	time.Sleep(1 * time.Second)
}

func main() {
	// Initialize Global System (Config, MCP, Prompts, Tools)
	if err := InitSystem(context.Background()); err != nil {
		log.Fatalf("Failed to initialize system: %v", err)
	}

	if len(os.Args) == 2 && os.Args[1] == "-s" {
		runServer()
		return
	}
	runConsole()
}
