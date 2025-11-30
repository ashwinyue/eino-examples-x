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
	"time"

	"github.com/cloudwego/eino-examples/adk/common/model"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/supervisor"
	"github.com/cloudwego/eino/compose"
)

// NewPlanner creates a planner agent
func NewPlanner(ctx context.Context) (adk.Agent, error) {
	m := model.NewChatModel()
	prompt, ok := GlobalPrompts["planner"]
	if !ok {
		return nil, fmt.Errorf("planner prompt not found")
	}

	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "planner",
		Description: "Responsible for breaking down complex tasks into executable steps. Call this agent when you need to create a plan.",
		Instruction: prompt,
		Model:       m,
	})
}

// NewResearcher creates a researcher agent
func NewResearcher(ctx context.Context) (adk.Agent, error) {
	m := model.NewChatModel()
	prompt, ok := GlobalPrompts["researcher"]
	if !ok {
		return nil, fmt.Errorf("researcher prompt not found")
	}

	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "researcher",
		Description: "Responsible for executing research tasks. Call this agent when you need to gather information.",
		Instruction: prompt,
		Model:       m,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: GlobalResearchTools,
			},
		},
	})
}

// NewCoder creates a coder agent
func NewCoder(ctx context.Context) (adk.Agent, error) {
	m := model.NewChatModel()
	prompt, ok := GlobalPrompts["coder"]
	if !ok {
		return nil, fmt.Errorf("coder prompt not found")
	}

	// Replace current time dynamically for each request (since NewCoder is called per request)
	instruction := strings.ReplaceAll(prompt, "{{ CURRENT_TIME }}", time.Now().Format(time.RFC3339))

	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "coder",
		Description: "Responsible for writing code. Call this agent when you need to generate code.",
		Instruction: instruction,
		Model:       m,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: GlobalCoderTools,
			},
		},
	})
}

// NewReporter creates a reporter agent
func NewReporter(ctx context.Context) (adk.Agent, error) {
	m := model.NewChatModel()
	prompt, ok := GlobalPrompts["reporter"]
	if !ok {
		return nil, fmt.Errorf("reporter prompt not found")
	}
	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "reporter",
		Description: "Responsible for generating the final report based on the findings.",
		Instruction: prompt,
		Model:       m,
	})
}

// NewPodcastScriptWriter creates a podcast script writer agent
func NewPodcastScriptWriter(ctx context.Context) (adk.Agent, error) {
	m := model.NewChatModel()
	prompt, ok := GlobalPrompts["podcast_script_writer"]
	if !ok {
		return nil, fmt.Errorf("podcast_script_writer prompt not found")
	}
	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "podcast_writer",
		Description: "Responsible for converting content into a podcast script.",
		Instruction: prompt,
		Model:       m,
	})
}

// NewPPTComposer creates a PPT composer agent
func NewPPTComposer(ctx context.Context) (adk.Agent, error) {
	m := model.NewChatModel()
	prompt, ok := GlobalPrompts["ppt_composer"]
	if !ok {
		return nil, fmt.Errorf("ppt_composer prompt not found")
	}
	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "ppt_composer",
		Description: "Responsible for creating a markdown presentation from content.",
		Instruction: prompt,
		Model:       m,
	})
}

// NewCoordinator creates the supervisor/coordinator
func NewCoordinator(ctx context.Context, subAgents []adk.Agent) (adk.Agent, error) {
	m := model.NewChatModel()

	prompt, ok := GlobalPrompts["coordinator"]
	if !ok {
		return nil, fmt.Errorf("coordinator prompt not found")
	}

	// Replace CURRENT_TIME in CoordinatorPrompt
	baseInstruction := strings.ReplaceAll(prompt, "{{ CURRENT_TIME }}", time.Now().Format(time.RFC3339))

	// Append the workflow instruction
	workflowInstruction := fmt.Sprintf(`

# Workflow Management

You also manage the research and content generation workflow.
Your available agents are:
- investigator: conducts initial background investigation
- planner: breaks down tasks into a research plan
- researcher: gathers information for research steps
- coder: processes data or writes code for processing steps
- reporter: generates a final report
- podcast_writer: converts content into a podcast script
- ppt_composer: creates a presentation

If the user request requires research (Category 3), follow this workflow:
1. (Optional) If you need initial context, ask the **investigator** to search.
2. Ask the **planner** to create a plan (Max steps: %d).
3. **Human Feedback**: Output the plan to the user and ask for approval.
    - If the user requests changes, ask the **planner** to revise the plan.
    - Only proceed after the user approves.
4. Iterate through the steps in the plan:
    - If it is a research step, assign it to the **researcher**.
    - If it is a processing step, assign it to the **coder**.
    - Pass the result of previous steps to the next step if needed.
5. Once all steps are complete (or you have enough info), choose the appropriate output format:
    - Default: ask the **reporter** to generate a final report.
    - If user asked for a podcast: ask the **podcast_writer**.
    - If user asked for a PPT/presentation: ask the **ppt_composer**.
6. Return the final result to the user.

Always follow this flow for complex tasks. Do not skip planning.
`, Config.Setting.MaxStepNum)

	instruction := baseInstruction + workflowInstruction

	sv, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "coordinator",
		Description: "Coordinates the planning and execution process.",
		Instruction: instruction,
		Model:       m,
		Exit:        &adk.ExitTool{}, // Allow supervisor to finish the conversation
	})
	if err != nil {
		return nil, err
	}

	return supervisor.New(ctx, &supervisor.Config{
		Supervisor: sv,
		SubAgents:  subAgents,
	})
}
