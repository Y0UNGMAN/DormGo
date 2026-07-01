# DormGo Agent Phase 1 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a hybrid Agent flow where the Vue frontend calls Go, Python uses the configured LLM to plan task execution, and Go remains responsible for all post search, signup confirmation, permissions, and database writes.

**Architecture:** Go exposes `/api/v1/agent/chat` for authenticated users and `/internal/agent/tools/*` for the Python Agent. Python reads `backend/config.yaml` for the DeepSeek-compatible API config, asks the model for a structured plan, calls Go internal tools for broad candidate recall, then uses an LLM reranker to score candidates semantically. Redis stores pending signup actions with an `action_id` and server-side session context so destructive actions require user confirmation.

**Tech Stack:** Go/Gin/GORM/Redis, Python standard-library HTTP service, Vue 3.

---

### Task 1: Backend Agent Client and Internal Tools

**Files:**
- Create: `backend/logic/agent.go`
- Create: `backend/controller/agent_con.go`
- Modify: `backend/router/app.go`
- Modify: `backend/model/dg_post.go`

- [ ] Add Go tests for agent request payloads and pending-action key generation.
- [ ] Implement Go client that forwards user messages to Python Agent.
- [ ] Implement internal tool endpoints for search, prepare signup, and confirm signup.
- [ ] Add a model search function for signup-capable posts.

### Task 2: Python Agent Service

**Files:**
- Create: `agent-service/main.py`
- Create: `agent-service/config.py`
- Create: `agent-service/llm_client.py`
- Create: `agent-service/planner.py`
- Create: `agent-service/reranker.py`
- Create: `agent-service/tools.py`
- Create: `agent-service/test_agent_runtime.py`
- Create: `agent-service/README.md`

- [ ] Add Python tests for Go config loading and model-driven structured planning.
- [ ] Implement OpenAI-compatible LLM calls using the `api` section from `backend/config.yaml`.
- [ ] Implement broad candidate recall from Go tools and semantic reranking with the LLM.
- [ ] Implement a small HTTP Agent service that asks the LLM for a plan, calls Go internal tools, reranks candidates, and stores session context.

### Task 3: Frontend Entry Point

**Files:**
- Create: `frontend/src/views/AgentAssistant.vue`
- Modify: `frontend/src/router/index.ts`
- Modify: `frontend/src/views/DormgoHome.vue`

- [ ] Add a simple Agent chat page and route.
- [ ] Add a compact entry button from the DormGo home screen.
- [ ] Keep API calls through the existing Axios wrapper.

### Task 4: Verification

**Commands:**
- `go test ./...` from `backend`
- `python3 -m unittest discover agent-service`
- `npm run build` from `frontend`

**Expected:** Go unit tests pass, Python intent tests pass, and frontend build completes.
