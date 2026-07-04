import json
import logging
import threading
import time
import traceback
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from typing import Optional

from config import load_go_api_config
from llm_client import OpenAICompatibleClient
from reranker import LLMReranker
from react_loop import ReActLoop
from tool_registry import build_default_registry
from tools import GoToolClient
from tracing import Tracer

logger = logging.getLogger("dormgo.agent")

# Vector store is optional — graceful fallback to keyword search if deps missing
try:
    import vector_store
    _VECTOR_STORE_AVAILABLE = True
except ImportError:
    vector_store = None
    _VECTOR_STORE_AVAILABLE = False
    logger = logging.getLogger("dormgo.agent")
    # logger will be configured later in main()

# Shared Go client for vector index refresh
_shared_go_client: Optional[GoToolClient] = None


def build_llm_client() -> OpenAICompatibleClient:
    config = load_go_api_config()
    return OpenAICompatibleClient(config.api_url, config.api_key, config.model)


def build_reranker() -> LLMReranker:
    return LLMReranker(build_llm_client())


def build_react_loop(reranker: LLMReranker = None, go_client: GoToolClient = None) -> ReActLoop:
    if reranker is None:
        reranker = build_reranker()
    registry = build_default_registry(go_client, reranker, vector_store)
    return ReActLoop(build_llm_client(), registry, max_steps=5)


def init_vector_index(go_client: GoToolClient):
    """Build the FAISS vector index from all active posts on startup."""
    try:
        logger.info("Building vector index from Go backend...")
        posts = go_client.get_all_posts(100)
        if posts:
            count = vector_store.build_index(posts)
            logger.info(f"Vector index ready: {count} posts indexed")
        else:
            logger.warning("No posts returned from backend, vector index empty")
    except Exception as exc:
        logger.error(f"Failed to build vector index: {exc}")


def refresh_vector_index():
    """Background thread: refresh the vector index every 5 minutes."""
    while True:
        time.sleep(300)  # 5 minutes
        try:
            if _shared_go_client is None:
                continue
            posts = _shared_go_client.get_all_posts(100)
            if posts:
                count = vector_store.build_index(posts)
                logger.info(f"Vector index refreshed: {count} posts")
        except Exception as exc:
            logger.error(f"Vector index refresh failed: {exc}")


def _extract_payload(body_bytes: bytes) -> dict:
    return json.loads(body_bytes.decode("utf-8") or "{}")


def _build_context(payload: dict) -> tuple:
    user_id = int(payload.get("user_id") or 0)
    message = payload.get("message") or ""
    go_base_url = payload.get("go_base_url") or "http://127.0.0.1:8080"
    internal_token = payload.get("internal_token") or "dormgo-agent-local-token"
    return user_id, message, go_base_url, internal_token


def handle_chat_render_only(payload: dict) -> dict:
    user_id, message, go_base_url, internal_token = _build_context(payload)
    client = GoToolClient(go_base_url, internal_token)
    session = client.get_session(user_id)
    reranker = build_reranker()
    react = build_react_loop(reranker, client)
    tracer = Tracer()
    result = react.run_sync(user_id, message, session, tracer)
    reply = result.get("reply") or "Agent returned no reply"
    state = result.get("state") or {}

    history = list(session.get("history") or [])
    history.append({"role": "user", "content": message})
    history.append({"role": "assistant", "content": reply})
    client.save_session(
        user_id,
        last_candidates=state.get("last_candidates") or session.get("last_candidates") or [],
        history=history[-20:],
        selected_post=state.get("selected_post"),
        pending_action=state.get("pending_action"),
    )

    return {
        "reply": reply,
        "type": "react_result",
        "trace": tracer.dump(),
    }


def handle_chat_stream(payload: dict):
    user_id, message, go_base_url, internal_token = _build_context(payload)
    client = GoToolClient(go_base_url, internal_token)
    session = client.get_session(user_id)
    reranker = build_reranker()
    react = build_react_loop(reranker, client)
    tracer = Tracer()
    final_reply = None
    final_state = {}

    for event in react.run(user_id, message, session, tracer):
        yield event
        if event["type"] == "finish":
            final_reply = event["data"].get("reply", "")
            final_state = event["data"].get("state") or {}

    history = list(session.get("history") or [])
    history.append({"role": "user", "content": message})
    history.append({"role": "assistant", "content": final_reply or "Agent finished"})
    client.save_session(
        user_id,
        last_candidates=final_state.get("last_candidates") or session.get("last_candidates") or [],
        history=history[-20:],
        selected_post=final_state.get("selected_post"),
        pending_action=final_state.get("pending_action"),
    )

    yield {"type": "trace", "data": tracer.dump()}


def choose_candidate_for_signup(session: dict, plan: dict) -> dict:
    candidates = session.get("last_candidates") or []
    try:
        index = int(((plan.get("target") or {}).get("index")) or 1)
    except (TypeError, ValueError):
        index = 1
    if index <= 0:
        index = 1
    if index > len(candidates):
        raise RuntimeError("没有可引用的候选帖子，请先搜索帖子")
    return candidates[index - 1]


def format_related_fallback(posts, plan, signup_only: bool) -> str:
    query = plan.get("query") or "你的条件"
    if signup_only:
        lines = [f'没有找到完全匹配「{query}」的可报名活动，但找到一些相关普通帖子：']
    else:
        lines = [f'没有找到完全匹配「{query}」的帖子，先给你看一些最近的相关/可参考帖子：']
    for idx, post in enumerate(posts[:5], start=1):
        title = post.get("title") or "未命名帖子"
        dorm = (post.get("Dorm") or post.get("dorm") or {}).get("dormname", "")
        marker = "，可报名" if post.get("is_limited") else ""
        lines.append(f"{idx}. {title} {dorm}{marker}".strip())
    if signup_only:
        lines.append("这些结果不一定能直接报名；如果你想继续，我可以帮你重新放宽条件搜索。")
    else:
        lines.append("如果这些不相关，可以告诉我更具体的关键词，比如地点、时间或活动类型。")
    return "\n".join(lines)


def render_history_for_planner(session: dict) -> str:
    items = session.get("history") or []
    lines = []
    for item in items[-8:]:
        role = item.get("role") or "unknown"
        content = (item.get("content") or "").replace("\n", " ")
        lines.append(f"{role}: {content}")
    return "\n".join(lines)


def _event_label(event: dict) -> str:
    event_type = event.get("type", "")
    if event_type == "start":
        return "start"
    if event_type == "step":
        return "step"
    if event_type == "finish":
        return "finish"
    if event_type == "trace":
        return "trace"
    if event_type == "error":
        return "error"
    return "message"


def _write_agent_step_event(wfile, event: dict):
    label = _event_label(event)
    blob = json.dumps({label: event.get("data")}, ensure_ascii=False)
    wfile.write(f"event: agent_step\ndata: {blob}\n\n".encode("utf-8"))
    wfile.flush()


class AgentHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        length = int(self.headers.get("Content-Length", "0"))
        body_bytes = self.rfile.read(length)

        if self.path == "/chat":
            self._handle_chat(body_bytes)
        elif self.path == "/chat/stream":
            self._handle_chat_stream(body_bytes)
        else:
            self.send_error(404)

    def _handle_chat(self, body_bytes: bytes):
        try:
            payload = _extract_payload(body_bytes)
            response = handle_chat_render_only(payload)
            self._json(200, response)
        except Exception as exc:
            traceback.print_exc()
            self._json(200, {"reply": f"Agent error: {exc}", "type": "error"})

    def _handle_chat_stream(self, body_bytes: bytes):
        try:
            payload = _extract_payload(body_bytes)
        except Exception:
            self.send_error(400)
            return

        self.send_response(200)
        self.send_header("Content-Type", "text/event-stream")
        self.send_header("Cache-Control", "no-cache")
        self.send_header("Connection", "keep-alive")
        self.end_headers()

        try:
            for event in handle_chat_stream(payload):
                _write_agent_step_event(self.wfile, event)
        except Exception as exc:
            traceback.print_exc()
            self.wfile.write(
                f"event: agent_error\ndata: {{\"error\": \"{exc}\"}}\n\n".encode("utf-8")
            )
            self.wfile.flush()

    def log_message(self, format, *args):
        return

    def _json(self, status: int, payload: dict):
        body = json.dumps(payload, ensure_ascii=False).encode("utf-8")
        self.send_response(status)
        self.send_header("Content-Type", "application/json; charset=utf-8")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)


def main():
    global _shared_go_client

    logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(name)s] %(levelname)s: %(message)s")

    # Initialize vector index on startup
    config = load_go_api_config()
    go_base_url = config.api_url.replace("/chat/completions", "").replace("/v1", "")
    # Read Go base URL from config or use default
    try:
        from pathlib import Path
        import sys
        config_path = Path(__file__).resolve().parents[1] / "backend" / "config.yaml"
        # Use default go_base_url for indexing — the handler will override per-request
        go_base_url = "http://127.0.0.1:8080"
        internal_token = "dormgo-agent-local-token"
    except Exception:
        go_base_url = "http://127.0.0.1:8080"
        internal_token = "dormgo-agent-local-token"

    _shared_go_client = GoToolClient(go_base_url, internal_token)
    init_vector_index(_shared_go_client)

    # Start background refresh thread
    refresh_thread = threading.Thread(target=refresh_vector_index, daemon=True)
    refresh_thread.start()

    server = ThreadingHTTPServer(("127.0.0.1", 8090), AgentHandler)
    print("DormGo Agent service (ReAct + FAISS vector search) listening on http://127.0.0.1:8090")
    server.serve_forever()


if __name__ == "__main__":
    main()
