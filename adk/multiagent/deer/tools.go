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

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

type SearchReq struct {
	Query string `json:"query"`
}

type SearchResp struct {
	Result string `json:"result"`
}

func NewWebSearchTool() tool.BaseTool {
	// Mock search tool
	return utils.NewTool[*SearchReq, *SearchResp](&schema.ToolInfo{
		Name: "web_search_tool",
		Desc: "Search the web for information",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"query": {
				Type: schema.String,
				Desc: "The search query",
			},
		}),
	}, func(ctx context.Context, req *SearchReq) (*SearchResp, error) {
		return &SearchResp{
			Result: "Mock search result for: " + req.Query + ". Eino is an ultimate LLM/AI application development framework in Golang.",
		}, nil
	})
}

type PythonReq struct {
	Code string `json:"code"`
}

type PythonResp struct {
	Output string `json:"output"`
}

func NewPythonTool() tool.BaseTool {
	// Mock python tool
	return utils.NewTool[*PythonReq, *PythonResp](&schema.ToolInfo{
		Name: "python_runner",
		Desc: "Run python code",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"code": {
				Type: schema.String,
				Desc: "The python code to run",
			},
		}),
	}, func(ctx context.Context, req *PythonReq) (*PythonResp, error) {
		return &PythonResp{
			Output: "Mock python output: Executed code successfully.",
		}, nil
	})
}
