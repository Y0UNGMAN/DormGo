import re
from dataclasses import dataclass, field
from typing import Callable, Any

MIN_VECTOR_SCORE = 0.55
RELAXED_VECTOR_SCORE = 0.52
EXACT_MATCH_STOP_TERMS = {
    "今天",
    "今晚",
    "明天",
    "最近",
    "一起",
    "有没有",
    "有",
    "想",
    "找",
    "帖子",
    "相关",
    "活动",
    "报名",
    "第一个",
    "第二个",
    "第三个",
    "这个",
    "那个",
    "宿舍",
}


@dataclass
class ToolSpec:
    name: str
    description: str
    parameters: dict
    executor: Callable[..., Any]

    def to_openai_schema(self) -> dict:
        return {
            "type": "function",
            "function": {
                "name": self.name,
                "description": self.description,
                "parameters": self.parameters,
            },
        }

    def describe(self) -> str:
        params = list(self.parameters.get("properties", {}).keys())
        required = self.parameters.get("required", [])
        param_strs = []
        for p in params:
            mark = "*" if p in required else ""
            param_strs.append(f"{p}{mark}")
        return f"- {self.name}({', '.join(param_strs)}): {self.description}"


@dataclass
class ToolResult:
    observation: str
    data: Any = None
    state_patch: dict = field(default_factory=dict)
    ok: bool = True
    error: str = ""

    def __str__(self) -> str:
        return self.observation

    def __contains__(self, item: str) -> bool:
        return item in self.observation


class ToolRegistry:
    def __init__(self):
        self._tools: dict[str, ToolSpec] = {}

    def register(self, spec: ToolSpec):
        self._tools[spec.name] = spec

    def get(self, name: str) -> ToolSpec:
        if name not in self._tools:
            raise KeyError(f"unknown tool: {name}. Available: {list(self._tools.keys())}")
        return self._tools[name]

    def describe(self) -> str:
        return "\n".join(spec.describe() for spec in self._tools.values())

    def openai_schemas(self) -> list:
        return [spec.to_openai_schema() for spec in self._tools.values()]

    def execute(self, name: str, user_id: int = 0, **kwargs) -> ToolResult:
        try:
            spec = self.get(name)
            result = spec.executor(user_id=user_id, **kwargs)
            if isinstance(result, dict):
                return _format_tool_result(spec.name, result)
            return ToolResult(observation=str(result), data=result)
        except Exception as exc:
            return ToolResult(
                observation=f"tool {name} failed: {exc}",
                ok=False,
                error=str(exc),
            )


def _format_tool_result(tool_name: str, result: dict) -> ToolResult:
    data = result.get("data") or {}
    if tool_name in ("search_posts", "search_signup_posts"):
        posts = data if isinstance(data, list) else []
        if not posts:
            return ToolResult(
                observation=f"{tool_name}: no matching posts found. Try different keywords or FINISH to tell the user.",
                data=[],
                state_patch={"last_candidates": []},
            )
        lines = [f"{tool_name}: found {len(posts)} candidates (semantic search). Use these candidates for reply, detail lookup, or signup preparation as requested."]
        for idx, post in enumerate(posts[:5], start=1):
            title = post.get("title") or "untitled"
            dorm = (post.get("Dorm") or post.get("dorm") or {}).get("dormname", "")
            post_id = post.get("id") or post.get("ID")
            is_limited = post.get("is_limited")
            score = post.get("match_score") or post.get("_vector_score")
            score_info = f" relevance={score:.2f}" if score else ""
            signup_info = " (signup)" if is_limited else ""
            lines.append(
                f"  {idx}. id={post_id} title={title} dorm={dorm}{signup_info}{score_info}"
            )
        return ToolResult(
            observation="\n".join(lines),
            data=posts,
            state_patch={"last_candidates": _candidate_refs(posts)},
        )
    if tool_name in ("prepare_signup", "prepare_signup_from_candidate"):
        return ToolResult(
            observation=f"{tool_name}: pending signup action created. Ask the user to confirm before executing signup. data={data}",
            data=data,
            state_patch={"pending_action": data if isinstance(data, dict) else {}},
        )
    if tool_name == "confirm_signup":
        return ToolResult(
            observation=f"{tool_name}: signup confirmed. data={data}",
            data=data,
            state_patch={"pending_action": None},
        )
    if tool_name in ("get_post_detail", "get_post_detail_from_candidate"):
        return ToolResult(
            observation=_format_post_detail(tool_name, data if isinstance(data, dict) else {}),
            data=data,
            state_patch={"selected_post": data if isinstance(data, dict) else {}},
        )
    return ToolResult(observation=f"{tool_name}: {result}", data=data)


def _candidate_refs(posts: list) -> list:
    refs = []
    for post in posts[:8]:
        post_id = _pick_id(post)
        if post_id is None:
            continue
        refs.append({
            "post_id": post_id,
            "title": post.get("title") or "",
            "is_limited": bool(post.get("is_limited")),
        })
    return refs


def _format_post_detail(tool_name: str, post: dict) -> str:
    if not post:
        return f"{tool_name}: no post detail found"
    title = post.get("title") or "untitled"
    content = post.get("content") or ""
    dorm = (post.get("Dorm") or post.get("dorm") or {}).get("dormname", "")
    post_id = post.get("id") or post.get("ID")
    signup = "signup enabled" if post.get("is_limited") else "signup disabled"
    current = post.get("current_enrollment", 0)
    max_count = post.get("max_enrollment", 0)
    return (
        f"{tool_name}: id={post_id} title={title} dorm={dorm} {signup} "
        f"enrollment={current}/{max_count or 'unlimited'} content={content}"
    )


def _pick_id(post: dict):
    """Extract post ID from a post dict, handling both 'id' and 'ID' keys."""
    val = post.get("id") or post.get("ID")
    try:
        return int(val)
    except (TypeError, ValueError):
        return None


def build_default_registry(go_client, reranker=None, vector_store=None) -> ToolRegistry:
    registry = ToolRegistry()

    def _vector_search_posts(query: str, limit: int, signup_only: bool) -> list:
        """
        Industrial-grade search pipeline:
          user query → vector embed → FAISS ANN(top-k) → fetch post details → LLM rerank → return
        Falls back to SQL LIKE if vector index is not ready.
        """
        use_vector = vector_store is not None and vector_store.is_ready()

        if use_vector:
            # Step 1: Semantic vector search
            top_k = 20
            scored = vector_store.search(query, k=top_k)

            if not scored:
                return []

            # Step 2: Fetch full post details
            all_posts = go_client.get_all_posts(100)
            post_by_id = {}
            for p in all_posts:
                pid = _pick_id(p)
                if pid is not None:
                    post_by_id[pid] = p

            # Step 3: Match vector results to full posts
            candidates = []
            relaxed_candidates = []
            for item in scored:
                pid = item["id"]
                if pid in post_by_id:
                    post = dict(post_by_id[pid])
                    if signup_only and not _is_signup_available(post):
                        continue
                    score = float(item.get("score") or 0)
                    post["_vector_score"] = score
                    if score >= MIN_VECTOR_SCORE or _has_exact_query_term(query, post):
                        candidates.append(post)
                    elif score >= RELAXED_VECTOR_SCORE:
                        relaxed_candidates.append(post)
                    if len(candidates) >= limit:
                        break

            if not candidates:
                candidates = relaxed_candidates[:limit]

            # Vector scores are already cosine-similarity based. Reranker is optional extra.
            return candidates

        # Fallback: SQL LIKE keyword search (old path)
        if signup_only:
            candidates = go_client.search_signup_posts(0, query, {}, limit)
        else:
            candidates = go_client.search_posts(0, query, {}, limit)
        if reranker and candidates and query:
            candidates = reranker.rank(query, candidates)
        return candidates

    def search_posts(user_id: int, query: str = "", limit: int = 10, search_terms: list = None):
        # Combine query with search_terms for richer semantic search
        search_text = query
        if search_terms:
            search_text = query + " " + " ".join(search_terms)
        candidates = _vector_search_posts(search_text, int(limit), signup_only=False)
        return {"data": candidates}

    def search_signup_posts(user_id: int, query: str = "", limit: int = 10, search_terms: list = None):
        search_text = query
        if search_terms:
            search_text = query + " " + " ".join(search_terms)
        candidates = _vector_search_posts(search_text, int(limit), signup_only=True)
        return {"data": candidates}

    def prepare_signup(user_id: int, post_id: int = 0):
        return go_client.prepare_signup(int(user_id), int(post_id))

    def confirm_signup(user_id: int, action_id: str = ""):
        return go_client.confirm_signup(int(user_id), action_id)

    def prepare_signup_from_candidate(user_id: int, index: int = 1):
        return go_client.prepare_signup_from_candidate(int(user_id), int(index))

    def get_post_detail(user_id: int, post_id: int = 0):
        return go_client.get_post_detail(int(user_id), int(post_id))

    def get_post_detail_from_candidate(user_id: int, index: int = 1):
        return go_client.get_post_detail_from_candidate(int(user_id), int(index))

    registry.register(ToolSpec(
        name="search_posts",
        description="semantic search for non-signup posts. Uses vector search to find conceptually related content — natural language queries work well.",
        parameters={
            "type": "object",
            "properties": {
                "query": {"type": "string", "description": "search keywords"},
                "limit": {"type": "integer", "description": "max results to return", "default": 10},
                "search_terms": {"type": "array", "items": {"type": "string"}, "description": "related keywords to broaden search, e.g. ['mall', 'shopping', 'weekend']"},
            },
            "required": ["query"],
        },
        executor=search_posts,
    ))

    registry.register(ToolSpec(
        name="search_signup_posts",
        description="semantic search for signup-capable posts. Uses vector search to match user intent with activities — natural language queries work well.",
        parameters={
            "type": "object",
            "properties": {
                "query": {"type": "string", "description": "search keywords"},
                "limit": {"type": "integer", "description": "max results to return", "default": 10},
                "search_terms": {"type": "array", "items": {"type": "string"}, "description": "related keywords to broaden search, e.g. ['hiking', 'outing', 'weekend trip']"},
            },
            "required": ["query"],
        },
        executor=search_signup_posts,
    ))

    registry.register(ToolSpec(
        name="prepare_signup",
        description="create a pending signup action for a post, requires user confirmation before actual signup",
        parameters={
            "type": "object",
            "properties": {
                "post_id": {"type": "integer", "description": "post id to prepare signup for"},
            },
            "required": ["post_id"],
        },
        executor=prepare_signup,
    ))

    registry.register(ToolSpec(
        name="prepare_signup_from_candidate",
        description="prepare signup by candidate index (1-based) from last search results",
        parameters={
            "type": "object",
            "properties": {
                "index": {"type": "integer", "description": "1-based index of candidate from previous search results"},
            },
            "required": ["index"],
        },
        executor=prepare_signup_from_candidate,
    ))

    registry.register(ToolSpec(
        name="get_post_detail",
        description="get full post details by post id, including content and signup status",
        parameters={
            "type": "object",
            "properties": {
                "post_id": {"type": "integer", "description": "post id to inspect"},
            },
            "required": ["post_id"],
        },
        executor=get_post_detail,
    ))

    registry.register(ToolSpec(
        name="get_post_detail_from_candidate",
        description="get full post details by candidate index (1-based) from last search results",
        parameters={
            "type": "object",
            "properties": {
                "index": {"type": "integer", "description": "1-based index of candidate from previous search results"},
            },
            "required": ["index"],
        },
        executor=get_post_detail_from_candidate,
    ))

    registry.register(ToolSpec(
        name="confirm_signup",
        description="confirm and execute a pending signup action",
        parameters={
            "type": "object",
            "properties": {
                "action_id": {"type": "string", "description": "action id from prepare_signup"},
            },
            "required": [],
        },
        executor=confirm_signup,
    ))

    return registry


def _is_signup_available(post: dict) -> bool:
    if not post.get("is_limited"):
        return False
    max_count = int(post.get("max_enrollment") or 0)
    current = int(post.get("current_enrollment") or 0)
    return max_count == 0 or current < max_count


def _has_exact_query_term(query: str, post: dict) -> bool:
    normalized_query = " ".join((query or "").lower().split())
    if not normalized_query:
        return False
    searchable = f"{post.get('title') or ''} {post.get('content') or ''}".lower()
    compact_searchable = searchable.replace(" ", "")
    compact_query = _compact_term(normalized_query)
    if len(compact_query) >= 4 and compact_query in compact_searchable:
        return True
    return any(term in searchable or _compact_term(term) in compact_searchable for term in _meaningful_exact_terms(normalized_query))


def _meaningful_exact_terms(query: str) -> list[str]:
    terms = []
    for raw in re.split(r"\s+", query):
        term = _compact_term(raw)
        if len(term) < 2 or term in EXACT_MATCH_STOP_TERMS:
            continue
        terms.append(term)
    return terms


def _compact_term(value: str) -> str:
    return re.sub(r"[\s,，。.!！?？、:：;；（）()【】\[\]\"'“”‘’]+", "", (value or "").lower())
