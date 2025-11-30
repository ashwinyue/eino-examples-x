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

package infra

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino-ext/components/model/openai"
	einomodel "github.com/cloudwego/eino/components/model"
	arkModel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"

	"github.com/cloudwego/eino-examples/flow/agent/deer-go/conf"
)

var (
	ChatModel einomodel.ToolCallingChatModel
	PlanModel einomodel.ToolCallingChatModel
)

func InitModel() {
	if os.Getenv("MODEL_TYPE") == "ark" {
		cm, _ := ark.NewChatModel(context.Background(), &ark.ChatModelConfig{
			APIKey:  os.Getenv("ARK_API_KEY"),
			Model:   os.Getenv("ARK_MODEL"),
			BaseURL: os.Getenv("ARK_BASE_URL"),
			Thinking: &arkModel.Thinking{
				Type: arkModel.ThinkingTypeDisabled,
			},
		})
		ChatModel = cm
		PlanModel = cm
		return
	}

	base := conf.Config.Model.BaseURL
	key := conf.Config.Model.APIKey
	mdl := conf.Config.Model.DefaultModel
	if key == "" {
		key = os.Getenv("OPENAI_API_KEY")
	}
	if mdl == "" {
		mdl = os.Getenv("OPENAI_MODEL")
	}
	if base == "" {
		base = os.Getenv("OPENAI_BASE_URL")
	}
	cm, _ := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
		BaseURL: base,
		APIKey:  key,
		Model:   mdl,
		ByAzure: func() bool { return os.Getenv("OPENAI_BY_AZURE") == "true" }(),
	})
	ChatModel = cm
	PlanModel = cm
}
