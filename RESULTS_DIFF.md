--- ../../../adk/multiagent/deer/adk_result.txt	2025-11-30 12:22:26
+++ deer_result.txt	2025-11-30 12:22:35
@@ -1,12 +1,21 @@
-[90m2025-11-30 12:22:23.489262[0m [32mINF[0m [1mconfig.go:74[0m[36m >[0m [36mevent=[0mload_config [36mlog_id=[0m [36mpayload=[0m{"conf":"{\"MCP\":{\"Servers\":{\"python\":{\"Command\":\"uv\",\"Args\":[\"--directory\",\"biz/mcps/python\",\"run\",\"server.py\"],\"Env\":null},\"tavily\":{\"Command\":\"npx\",\"Args\":[\"-y\",\"tavily-mcp@0.1.3\"],\"Env\":{\"TAVILY_API_KEY\":\"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\"}}}},\"Model\":{\"DefaultModel\":\"\\u003cyour model\\u003e\",\"APIKey\":\"\\u003cyour api key\\u003e\",\"BaseURL\":\"\\u003cyour base url\\u003e\"},\"Setting\":{\"MaxPlanIterations\":1,\"MaxStepNum\":3}}"} [36msuffix=[0mnull
-[90m2025-11-30 12:22:23.489341[0m [32mINF[0m [1mmcp_client.go:62[0m[36m >[0m [36mevent=[0m"load mcp client" [36mlog_id=[0m [36mpayload=[0m{"python":"stdio"} [36msuffix=[0mnull
-[90m2025-11-30 12:22:23.490487[0m [32mINF[0m [1mmcp_client.go:127[0m[36m >[0m [36mevent=[0m"Initializing server..." [36mlog_id=[0m [36mpayload=[0m{"name":"python"} [36msuffix=[0mnull
-[90m2025-11-30 12:22:23.779369[0m [32mINF[0m [1mmcp_client.go:62[0m[36m >[0m [36mevent=[0m"load mcp client" [36mlog_id=[0m [36mpayload=[0m{"tavily":"stdio"} [36msuffix=[0mnull
-[90m2025-11-30 12:22:23.780466[0m [32mINF[0m [1mmcp_client.go:127[0m[36m >[0m [36mevent=[0m"Initializing server..." [36mlog_id=[0m [36mpayload=[0m{"name":"tavily"} [36msuffix=[0mnull
-请输入你的需求： --------------------------------------------------
-name: coordinator
-path: [{coordinator}]
-error: [NodeRunError] Error code: 429 - {"code":"SetLimitExceeded","message":"Your account [2101141013] has reached the set inference limit for the [doubao-seed-1-6] model, and the model service has been paused. To continue using this model, please visit the Model Activation page to adjust or close the \"Safe Experience Mode\". Request id: 021764476546557db96079787cbb1eedcd8d537f824414587f0fd","param":"","type":"TooManyRequests","request_id":"202511301222260000D5BF7EC18A7CC09C"}
+[90m2025-11-30 12:22:32.358643[0m [32mINF[0m [1mconf/config.go:73[0m[36m >[0m [36mevent=[0mload_config [36mlog_id=[0m [36mpayload=[0m{"conf":"{\"MCP\":{\"Servers\":{\"python\":{\"Command\":\"uv\",\"Args\":[\"--directory\",\"biz/mcps/python\",\"run\",\"server.py\"],\"Env\":null},\"tavily\":{\"Command\":\"npx\",\"Args\":[\"-y\",\"tavily-mcp@0.1.3\"],\"Env\":{\"TAVILY_API_KEY\":\"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\"}}}},\"Model\":{\"DefaultModel\":\"\",\"APIKey\":\"\",\"BaseURL\":\"\"},\"Setting\":{\"MaxPlanIterations\":2,\"MaxStepNum\":3}}"} [36msuffix=[0mnull
+[90m2025-11-30 12:22:32.360477[0m [32mINF[0m [1mbiz/infra/mcp.go:132[0m[36m >[0m [36mevent=[0m"load mcp client" [36mlog_id=[0m [36mpayload=[0m{"tavily":"stdio"} [36msuffix=[0mnull
+[90m2025-11-30 12:22:32.361842[0m [32mINF[0m [1mbiz/infra/mcp.go:169[0m[36m >[0m [36mevent=[0m"Initializing server..." [36mlog_id=[0m [36mpayload=[0m{"name":"tavily"} [36msuffix=[0mnull
+[90m2025-11-30 12:22:33.674792[0m [32mINF[0m [1mbiz/infra/mcp.go:132[0m[36m >[0m [36mevent=[0m"load mcp client" [36mlog_id=[0m [36mpayload=[0m{"python":"stdio"} [36msuffix=[0mnull
+[90m2025-11-30 12:22:33.677433[0m [32mINF[0m [1mbiz/infra/mcp.go:169[0m[36m >[0m [36mevent=[0m"Initializing server..." [36mlog_id=[0m [36mpayload=[0m{"name":"python"} [36msuffix=[0mnull
+请输入你的需求： 
+==================
+ [OnStart] coordinator 
+==================
+=========[OnError]=========
+Error code: 429 - {"code":"SetLimitExceeded","message":"Your account [2101141013] has reached the set inference limit for the [doubao-seed-1-6] model, and the model service has been paused. To continue using this model, please visit the Model Activation page to adjust or close the \"Safe Experience Mode\". Request id: 021764476555326c85c17782ebc45169ced08d4429f63daf2cd87","param":"","type":"TooManyRequests","request_id":"2025113012223500001B0CCA2F6AB15C50"}
+=========[OnError]=========
+[NodeRunError] Error code: 429 - {"code":"SetLimitExceeded","message":"Your account [2101141013] has reached the set inference limit for the [doubao-seed-1-6] model, and the model service has been paused. To continue using this model, please visit the Model Activation page to adjust or close the \"Safe Experience Mode\". Request id: 021764476555326c85c17782ebc45169ced08d4429f63daf2cd87","param":"","type":"TooManyRequests","request_id":"2025113012223500001B0CCA2F6AB15C50"}
 ------------------------
-node path: [node_1, ChatModel]
-
+node path: [agent]
+=========[OnError]=========
+[NodeRunError] Error code: 429 - {"code":"SetLimitExceeded","message":"Your account [2101141013] has reached the set inference limit for the [doubao-seed-1-6] model, and the model service has been paused. To continue using this model, please visit the Model Activation page to adjust or close the \"Safe Experience Mode\". Request id: 021764476555326c85c17782ebc45169ced08d4429f63daf2cd87","param":"","type":"TooManyRequests","request_id":"2025113012223500001B0CCA2F6AB15C50"}
+------------------------
+node path: [coordinator, agent]
+[90m2025-11-30 12:22:35.364345[0m [31mERR[0m [1mmain.go:133[0m[36m >[0m [36mevent=[0m"run failed" [36mlog_id=[0m [36mpayload=[0m{"err":"[NodeRunError] Error code: 429 - {\"code\":\"SetLimitExceeded\",\"message\":\"Your account [2101141013] has reached the set inference limit for the [doubao-seed-1-6] model, and the model service has been paused. To continue using this model, please visit the Model Activation page to adjust or close the \\\"Safe Experience Mode\\\". Request id: 021764476555326c85c17782ebc45169ced08d4429f63daf2cd87\",\"param\":\"\",\"type\":\"TooManyRequests\",\"request_id\":\"2025113012223500001B0CCA2F6AB15C50\"}\n------------------------\nnode path: [coordinator, agent]"} [36msuffix=[0mnull
+[90m2025-11-30 12:22:35.364584[0m [32mINF[0m [1mmain.go:135[0m[36m >[0m [36mevent=[0m"run console finish" [36mlog_id=[0m [36mpayload=[0m{"time.Date(2025, time.November, 30, 12, 22, 35, 364573000, time.Local)":""} [36msuffix=[0mnull
