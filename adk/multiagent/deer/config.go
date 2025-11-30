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
	"os"
	"path/filepath"

	"github.com/RanFeng/ilog"
	"gopkg.in/yaml.v3"
)

// DeerConfig 定义一个结构体来解析 YAML 文件中的 mcp.servers 部分
type DeerConfig struct {
	MCP struct {
		Servers map[string]struct {
			Command string            `yaml:"command"`
			Args    []string          `yaml:"args"`
			Env     map[string]string `yaml:"env,omitempty"`
		} `yaml:"servers"`
	} `yaml:"mcp"`
	Model struct {
		DefaultModel string `yaml:"default_model"`
		APIKey       string `yaml:"api_key"`
		BaseURL      string `yaml:"base_url"`
	} `yaml:"model"`
	Setting struct {
		MaxPlanIterations int `yaml:"max_plan_iterations"`
		MaxStepNum        int `yaml:"max_step_num"`
	} `yaml:"setting"`
}

var (
	Config *DeerConfig = &DeerConfig{}
)

func LoadDeerConfig(ctx context.Context) {
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("获取当前工作目录失败: %w", err))
	}

	// 构建模板文件路径
	configPath := filepath.Join(dir, "conf", "deer-go.yaml")

	// 读取 YAML 文件内容
	configData, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("读取配置文件 %s 失败: %w", configPath, err))
	}

	var deerConfig DeerConfig
	if err := yaml.Unmarshal(configData, &deerConfig); err != nil {
		panic(fmt.Sprintf("解析配置文件 %s 失败: %w", configPath, err))
	}

	ilog.EventInfo(ctx, "load_config", "conf", deerConfig)

	Config = &deerConfig
}

// --- MCP Client Logic ---

const (
	transportStdio = "stdio"
	transportSSE   = "sse"
)

type MCPConfig struct {
	MCPServers map[string]ServerConfigWrapper `json:"mcpServers"`
}

type ServerConfig interface {
	GetType() string
}

type STDIOServerConfig struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env,omitempty"`
}

func (s STDIOServerConfig) GetType() string {
	return transportStdio
}

type SSEServerConfig struct {
	Url     string   `json:"url"`
	Headers []string `json:"headers,omitempty"`
}

func (s SSEServerConfig) GetType() string {
	return transportSSE
}

type ServerConfigWrapper struct {
	Config ServerConfig
}

func (w *ServerConfigWrapper) UnmarshalJSON(data []byte) error {
	var typeField struct {
		Url string `json:"url"`
	}

	if err := json.Unmarshal(data, &typeField); err != nil {
		return err
	}
	if typeField.Url != "" {
		var sse SSEServerConfig
		if err := json.Unmarshal(data, &sse); err != nil {
			return err
		}
		w.Config = sse
	} else {
		var stdio STDIOServerConfig
		if err := json.Unmarshal(data, &stdio); err != nil {
			return err
		}
		w.Config = stdio
	}

	return nil
}

func (w ServerConfigWrapper) MarshalJSON() ([]byte, error) {
	return json.Marshal(w.Config)
}
