import json
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

from config import load_go_api_config
from llm_client import OpenAICompatibleClient
from planner import LLMPlanner
from reranker import LLMReranker
from tools import GoToolClient


def build_llm_client() -> OpenAICompatibleClient:
    config = load_go_api_config()
    return OpenAICompatibleClient(config.api_url, config.api_key, config.model)


def build_planner() -> LLMPlanner:
    return LLMPlanner(build_llm_client())


def build_reranker() -> LLMReranker:
    return LLMReranker(build_llm_client())


def handle_chat(payload: dict, planner: LLMPlanner = None, reranker: LLMReranker = None) -> dict:
    user_id = int(payload.get("user_id") or 0)
    message = payload.get("message") or ""
    go_base_url = payload.get("go_base_url") or "http://127.0.0.1:8080"
    internal_token = payload.get("internal_token") or "dormgo-agent-local-token"

    planner = planner or build_planner()
    reranker = reranker or build_reranker()
    client = GoToolClient(go_base_url, internal_token)
    session = client.get_session(user_id)
    planner_input = _planner_input(message, session)
    plan = planner.plan(planner_input)

    if plan["intent"] == "confirm_signup":
        try:
            result = client.confirm_signup(user_id)
        except RuntimeError:
            index = int((plan.get("target") or {}).get("index") or 1)
            prepared = client.prepare_signup_from_candidate(user_id, index)
            action_id = ((prepared.get("data") or {}).get("action_id")) or ""
            result = client.confirm_signup(user_id, action_id)
        title = ((result.get("data") or {}).get("title")) or "该帖子"
        response = {
            "reply": f"报名成功：{title}",
            "type": "signup_success",
            "data": result.get("data"),
            "plan": _public_plan(plan),
        }
        _save_turn(client, user_id, session, message, response["reply"], session.get("last_candidates") or [])
        return response

    if plan["intent"] == "search_posts":
        candidates = _recall_candidates(client.search_posts, user_id, plan, limit=20)
        posts = reranker.rank(_semantic_query(plan), candidates)
        if not posts:
            fallback_posts = client.search_posts(user_id, "", {}, 5)
            if fallback_posts:
                response = {
                    "reply": format_related_fallback(fallback_posts, plan, signup_only=False),
                    "type": "fallback_candidates",
                    "candidates": fallback_posts,
                    "plan": _public_plan(plan),
                }
                _save_turn(client, user_id, session, message, response["reply"], _candidate_refs(fallback_posts))
                return response
            response = {
                "reply": _empty_reply(plan, signup_only=False),
                "type": "empty",
                "candidates": [],
                "plan": _public_plan(plan),
            }
            _save_turn(client, user_id, session, message, response["reply"], [])
            return response
        response = {
            "reply": _format_candidates(posts, plan, prepared_action=None),
            "type": "candidates",
            "candidates": posts,
            "plan": _public_plan(plan),
        }
        _save_turn(client, user_id, session, message, response["reply"], _candidate_refs(posts))
        return response

    if plan["intent"] == "search_signup_posts":
        candidates = _recall_candidates(client.search_signup_posts, user_id, plan, limit=20)
        posts = reranker.rank(_semantic_query(plan), candidates)
        if not posts:
            related_candidates = _recall_candidates(client.search_posts, user_id, plan, limit=20)
            related_posts = reranker.rank(_semantic_query(plan), related_candidates)
            if related_posts:
                response = {
                    "reply": format_related_fallback(related_posts, plan, signup_only=True),
                    "type": "fallback_candidates",
                    "candidates": related_posts,
                    "plan": _public_plan(plan),
                }
                _save_turn(client, user_id, session, message, response["reply"], _candidate_refs(related_posts))
                return response
            response = {
                "reply": _empty_reply(plan, signup_only=True),
                "type": "empty",
                "candidates": [],
                "plan": _public_plan(plan),
            }
            _save_turn(client, user_id, session, message, response["reply"], [])
            return response

        prepared_action = None
        first = posts[0]
        post_id = int(first.get("id") or first.get("ID") or 0)
        if plan["next_action"] == "prepare_signup" and post_id > 0:
            prepared = client.prepare_signup(user_id, post_id)
            prepared_action = prepared.get("data") or {}
        response = {
            "reply": _format_candidates(posts, plan, prepared_action=prepared_action),
            "type": "candidates",
            "candidates": posts,
            "data": {"pending_action": prepared_action} if prepared_action else {},
            "plan": _public_plan(plan),
        }
        _save_turn(client, user_id, session, message, response["reply"], _candidate_refs(posts))
        return response

    response = {
        "reply": "我可以帮你查找帖子、筛选可报名活动，并在你确认后完成报名。请直接描述你要找什么。",
        "type": "help",
        "plan": _public_plan(plan),
    }
    _save_turn(client, user_id, session, message, response["reply"], session.get("last_candidates") or [])
    return response


def _planner_input(message: str, session: dict) -> str:
    history = render_history_for_planner(session)
    candidates = session.get("last_candidates") or []
    candidate_lines = []
    for idx, candidate in enumerate(candidates[:5], start=1):
        flag = "可报名" if candidate.get("is_limited") else "不可报名"
        candidate_lines.append(f"{idx}. post_id={candidate.get('post_id')} title={candidate.get('title')} {flag}")
    context = "\n".join(
        [
            "最近对话:",
            history or "无",
            "上一轮候选:",
            "\n".join(candidate_lines) or "无",
            "当前用户消息:",
            message,
        ]
    )
    return context


def _recall_candidates(search_func, user_id: int, plan: dict, limit: int = 20):
    candidates = []
    seen = set()
    for query, filters in [
        (plan.get("query") or "", plan.get("filters") or {}),
        ("", {}),
    ]:
        for post in search_func(user_id, query, filters, limit):
            post_id = post.get("id") or post.get("ID")
            if post_id in seen:
                continue
            seen.add(post_id)
            candidates.append(post)
            if len(candidates) >= limit:
                return candidates
    return candidates


def _semantic_query(plan: dict) -> str:
    filters = plan.get("filters") or {}
    parts = [
        plan.get("query") or "",
        filters.get("location") or "",
        filters.get("time_hint") or "",
    ]
    return " ".join(part for part in parts if part)


def render_history_for_planner(session: dict) -> str:
    items = session.get("history") or []
    lines = []
    for item in items[-8:]:
        role = item.get("role") or "unknown"
        content = (item.get("content") or "").replace("\n", " ")
        lines.append(f"{role}: {content}")
    return "\n".join(lines)


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


def _save_turn(client: GoToolClient, user_id: int, session: dict, user_message: str, assistant_message: str, candidates):
    history = list(session.get("history") or [])
    history.append({"role": "user", "content": user_message})
    history.append({"role": "assistant", "content": assistant_message})
    client.save_session(user_id, last_candidates=candidates, history=history[-20:])


def _candidate_refs(posts):
    refs = []
    for post in posts[:8]:
        refs.append(
            {
                "post_id": int(post.get("id") or post.get("ID") or 0),
                "title": post.get("title") or "",
                "is_limited": bool(post.get("is_limited")),
            }
        )
    return refs


def _format_candidates(posts, plan, prepared_action=None) -> str:
    lines = [f"我按“{plan['query'] or '你的描述'}”找到了这些帖子："]
    for idx, post in enumerate(posts[:5], start=1):
        title = post.get("title") or "未命名帖子"
        dorm = (post.get("Dorm") or post.get("dorm") or {}).get("dormname", "")
        current = post.get("current_enrollment", 0)
        max_count = post.get("max_enrollment", 0)
        progress = ""
        if post.get("is_limited"):
            progress = f"，报名 {current}/{max_count}" if max_count else f"，报名 {current}/不限"
        lines.append(f"{idx}. {title} {dorm}{progress}".strip())

    if prepared_action:
        lines.append(f"已为第 1 个帖子创建待确认报名动作，确认编号 {prepared_action.get('action_id')}。回复“确定”后才会真正报名。")
    elif plan["intent"] == "search_signup_posts":
        lines.append("如果你要报名其中某个帖子，请回复“报名第一个”或说明编号。")
    return "\n".join(lines)


def format_related_fallback(posts, plan, signup_only: bool) -> str:
    query = plan.get("query") or "你的条件"
    if signup_only:
        lines = [f"没有找到完全匹配“{query}”的可报名活动，但找到一些相关普通帖子："]
    else:
        lines = [f"没有找到完全匹配“{query}”的帖子，先给你看一些最近的相关/可参考帖子："]
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


def _empty_reply(plan, signup_only: bool) -> str:
    scope = "可报名帖子" if signup_only else "帖子"
    query = plan["query"] or "你的条件"
    return f"没有找到匹配“{query}”的{scope}。你可以换一个关键词，或者去掉时间/地点限制再试。"


def _public_plan(plan: dict) -> dict:
    return {
        "intent": plan.get("intent"),
        "query": plan.get("query"),
        "filters": plan.get("filters"),
        "next_action": plan.get("next_action"),
        "target": plan.get("target"),
    }


class AgentHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        if self.path != "/chat":
            self.send_error(404)
            return
        length = int(self.headers.get("Content-Length", "0"))
        payload = json.loads(self.rfile.read(length).decode("utf-8") or "{}")
        try:
            response = handle_chat(payload)
            self._json(200, response)
        except Exception as exc:
            self._json(200, {"reply": f"Agent 执行失败：{exc}", "type": "error"})

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
    server = ThreadingHTTPServer(("127.0.0.1", 8090), AgentHandler)
    print("DormGo Agent service listening on http://127.0.0.1:8090")
    server.serve_forever()


if __name__ == "__main__":
    main()
