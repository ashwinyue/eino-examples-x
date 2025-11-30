# Deer-Go ADK Implementation

This directory contains the Eino ADK implementation of the Deer-Go multi-agent system.
It demonstrates how to use the `Supervisor` pattern to coordinate multiple agents (Planner, Researcher, Coder).

## Structure

- `main.go`: The entry point, sets up the agents and runs the query.
- `agents.go`: Definitions of the agents (Planner, Researcher, Coder, Coordinator).
- `types.go`: Shared type definitions (State, Plan, etc.) - *Note: mostly for compatibility, ADK handles state internally via AgentInput/Output*.

## How to Run

Make sure you have your LLM environment variables set (e.g., `OPENAI_API_KEY` or `ARK_API_KEY`).

```bash
# Run the example
go run .
```

## Logic

1.  **Coordinator (Supervisor)**: Receives the user request and delegates tasks to sub-agents.
2.  **Planner**: Breaks down complex tasks.
3.  **Researcher**: Gathers information (mocked in this example).
4.  **Coder**: Generates code.

The Coordinator uses the LLM to decide which agent to call based on the conversation history and the agents' descriptions.
