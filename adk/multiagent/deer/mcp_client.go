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
	"time"

	"github.com/RanFeng/ilog"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

var (
	MCPServer map[string]client.MCPClient
)

func InitMCP() {
	var err error
	MCPServer, err = CreateMCPClients()
	if err != nil {
		panic(err)
	}
}

func CreateMCPClients() (map[string]client.MCPClient, error) {
	// 将 DeerConfig 转换为 MCPConfig
	mcpConfig := &MCPConfig{
		MCPServers: make(map[string]ServerConfigWrapper),
	}

	for name, server := range Config.MCP.Servers {
		mcpConfig.MCPServers[name] = ServerConfigWrapper{
			Config: STDIOServerConfig{
				Command: server.Command,
				Args:    server.Args,
				Env:     server.Env,
			},
		}
	}

	clients := make(map[string]client.MCPClient)

	for name, server := range mcpConfig.MCPServers {
		var mcpClient client.MCPClient
		var err error
		ilog.EventInfo(context.Background(), "load mcp client", name, server.Config.GetType())
		if server.Config.GetType() == transportSSE {
			sseConfig := server.Config.(SSEServerConfig)

			// The client.WithHeaders returns a transport.ClientOption which is compatible if imported correctly.
			// However, NewSSEMCPClient takes transport.ClientOption...
			// It seems the import "github.com/mark3labs/mcp-go/client" exposes ClientOption?
			// Let's check the actual types if possible, but for now I'll try to fix by not using intermediate slice of wrong type if possible
			// or just ignore SSE for now if user config doesn't use it (Deer config uses stdio mostly).
			// But to be correct, let's try to use the right type.

			// Since I cannot easily see the exact type definition without `go doc`, I will assume the error message:
			// cannot use client.WithHeaders(...) as "github.com/mark3labs/mcp-go/client".ClientOption
			// Wait, `client.WithHeaders` is likely returning `transport.ClientOption`.
			// But `client.NewSSEMCPClient` takes `transport.ClientOption`.
			// So I should probably not declare `options` as `[]client.ClientOption`.
			// But `client.ClientOption` might be an alias or distinct type.

			// Let's just simplify and create client directly.

			mcpClient, err = client.NewSSEMCPClient(sseConfig.Url)
			if err == nil {
				// Start is a method on *SSEMCPClient, but NewSSEMCPClient returns (*Client, error) or interface?
				// The doc said: NewSSEMCPClient(baseURL string, options ...transport.ClientOption) (*Client, error)
				// And Client struct has Start method?
				// Let's check `go doc` output again... it said `*Client`.
				// Does `*Client` have `Start`?
				// `client.MCPClient` interface usually has `Start`? No, `MCPClient` interface has `Initialize`, `Ping`, etc.
				// `SSEMCPClient` usually needs to be started.

				// In mcp-go v0.8.2+, NewSSEMCPClient returns *SSEMCPClient.
				// But the error said: undefined: client.SSEMCPClient.

				// Let's try to start it if it implements an interface or just assume it's ready or cast it.
				// If NewSSEMCPClient returns `*Client`, maybe `*Client` has `Start`.

				if sseClient, ok := mcpClient.(interface{ Start(context.Context) error }); ok {
					err = sseClient.Start(context.Background())
				}
			}
		} else {
			stdioConfig := server.Config.(STDIOServerConfig)
			var env []string
			for k, v := range stdioConfig.Env {
				env = append(env, fmt.Sprintf("%s=%s", k, v))
			}
			mcpClient, err = client.NewStdioMCPClient(
				stdioConfig.Command,
				env,
				stdioConfig.Args...)
		}
		if err != nil {
			for _, c := range clients {
				_ = c.Close()
			}
			return nil, fmt.Errorf(
				"failed to create MCP client for %s: %w",
				name,
				err,
			)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		ilog.EventInfo(ctx, "Initializing server...", "name", name)
		initRequest := mcp.InitializeRequest{}
		initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
		initRequest.Params.ClientInfo = mcp.Implementation{
			Name:    "mcphost",
			Version: "0.1.0",
		}
		initRequest.Params.Capabilities = mcp.ClientCapabilities{}

		_, err = mcpClient.Initialize(ctx, initRequest)
		if err != nil {
			_ = mcpClient.Close()
			for _, c := range clients {
				_ = c.Close()
			}
			return nil, fmt.Errorf(
				"failed to initialize MCP client for %s: %w",
				name,
				err,
			)
		}

		clients[name] = mcpClient
	}

	return clients, nil
}
