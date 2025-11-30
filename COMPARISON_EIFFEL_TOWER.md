# 测试对比：Eiffel Tower 与世界最高建筑的高度倍数

## 测试设置
- 问题：How many times taller is the Eiffel Tower than the tallest building in the world?
- 模型：统一使用 ADK 项目环境配置（`MODEL_TYPE=ark`，`ARK_MODEL=doubao-seed-1-6-251015`，`ARK_BASE_URL=https://ark.cn-beijing.volces.com/api/v3`）
- 运行方式：mac + zsh，均使用 `go run` 控制台模式
- 项目路径：
  - 原版：`flow/agent/deer-go`
  - ADK 版：`adk/multiagent/deer`
- 原始输出文件：
  - ADK 输出：`adk/multiagent/deer/adk_result.txt`
  - 原版输出：`flow/agent/deer-go/deer_result.txt`

## 结果摘要
- ADK 版
  - 通过 Tavily MCP 搜索与研究、报告代理完整工作流，给出数值结论：
    - 艾菲尔铁塔高度：330m（2022 天线后）
    - 世界最高建筑：Burj Khalifa，高度 828m
    - 计算：330 ÷ 828 ≈ 0.40
    - 结论：Eiffel Tower ≈ 0.40 倍于 Burj Khalifa；换言之 Burj Khalifa ≈ 2.5 倍 Eiffel Tower
- 原版（当前配置下未启用 MCP 搜索）
  - 规划与研究环节运行正常，给出 Eiffel Tower 高度的多来源事实与历史变更，以及世界最高建筑的高度与历史对比
  - 未执行最终“处理/报告”环节的数值计算与汇总，因此没有输出倍数结论

## 关键差异
- 工具使用：
  - ADK 版启用 MCP（Tavily + Python），检索并组织权威来源，流程完成至“Reporter”给出比值结论
  - 原版当前 YAML 中未配置 MCP 服务（`servers: {}`），仅依靠模型生成研究文本，无工具检索，流程停留在“Researcher/ResearchTeam”阶段
- 计划执行：
  - ADK 版由 Supervisor + 提示词流程自调度到报告代理完成计算与表格输出
  - 原版严格依赖 `model.Plan` + 路由分派，当前运行设置下未进入“处理/报告”最终计算

## 一致性结论
- 在相同模型配置下，ADK 版与原版在“事实收集”层面一致：
  - Eiffel Tower 现高 330m；Burj Khalifa 高 828m
- 差异主要在“是否完成最终数值计算与呈现”：
  - ADK 版完成且输出了倍数关系（0.40 倍/2.5 倍）
  - 原版当前配置未完成该步，如需一致输出，建议为原版补上 MCP（Tavily + Python）并允许进入“处理/报告”阶段

## 建议与复现
- 为原版 `flow/agent/deer-go/conf/deer-go.yaml` 增加 MCP 与 API Key（可参考 ADK 的 `conf/deer-go.yaml`），再运行即可对齐 ADK 的完整结果路径。
- 统一 `setting.max_plan_iterations` 与 `setting.max_step_num`，避免研究团队在一次迭代内无法完成“处理/报告”。

## 附：原始输出摘录
- ADK 版（节选）：
  - 结论：`The Eiffel Tower is approximately 0.40 times the height of the Burj Khalifa ...`（来源：`adk/multiagent/deer/adk_result.txt`）
- 原版（节选）：
  - Eiffel Tower 高度演变与来源列表；Burj Khalifa 高度与历史对比（来源：`flow/agent/deer-go/deer_result.txt`）

## 备注
- 本对比仅统一模型设置，工具配置以各项目当前 YAML 为准；若统一工具配置，预期两者输出结论可进一步对齐。
