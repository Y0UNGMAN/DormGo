SYSTEM_PROMPT = """
你是 DormGo 校园帖子检索重排器。
给定用户查询和候选帖子列表，请判断哪些帖子与用户意图相关，并按相关性排序。
不要编造候选列表以外的帖子。

必须返回 JSON:
{
  "matches": [
    {"id": 123, "score": 0.0到1.0, "reason": "一句话说明为什么相关"}
  ]
}

规则:
- 只返回 score >= 0.35 的帖子。
- “一起学习”可以匹配“自习搭子、图书馆、复习、组队学习”等语义相关内容。
- “一起看电影”可以匹配“电影、拼票、影院、电影票”等语义相关内容。
- “一起打球”可以匹配“篮球、羽毛球、足球、乒乓球、约球”等语义相关内容。
- 如果没有相关帖子，返回空 matches。
""".strip()


class LLMReranker:
    def __init__(self, llm_client, min_score: float = 0.35):
        self.llm_client = llm_client
        self.min_score = min_score

    def rank(self, query: str, posts: list) -> list:
        if not posts:
            return []
        post_map = {_post_id(post): post for post in posts if _post_id(post) is not None}
        if not post_map:
            return []

        raw = self.llm_client.complete_json(
            [
                {"role": "system", "content": SYSTEM_PROMPT},
                {"role": "user", "content": _build_user_prompt(query, posts)},
            ]
        )
        matches = raw.get("matches") if isinstance(raw, dict) else []
        if not isinstance(matches, list):
            return []

        ranked = []
        seen = set()
        for item in matches:
            if not isinstance(item, dict):
                continue
            post_id = _safe_int(item.get("id"))
            if post_id is None or post_id in seen or post_id not in post_map:
                continue
            score = _safe_float(item.get("score"))
            if score < self.min_score:
                continue
            post = dict(post_map[post_id])
            post["match_score"] = score
            post["match_reason"] = str(item.get("reason") or "").strip()
            ranked.append(post)
            seen.add(post_id)
        ranked.sort(key=lambda post: post.get("match_score", 0), reverse=True)
        return ranked


def _build_user_prompt(query: str, posts: list) -> str:
    lines = [f"用户查询: {query}", "候选帖子:"]
    for post in posts[:30]:
        lines.append(
            f"- id={_post_id(post)} title={post.get('title') or ''} content={_short(post.get('content') or '')}"
        )
    return "\n".join(lines)


def _post_id(post: dict):
    return _safe_int(post.get("id") or post.get("ID"))


def _safe_int(value):
    try:
        return int(value)
    except (TypeError, ValueError):
        return None


def _safe_float(value):
    try:
        return float(value)
    except (TypeError, ValueError):
        return 0.0


def _short(text: str, limit: int = 160) -> str:
    text = " ".join(text.split())
    if len(text) <= limit:
        return text
    return text[:limit] + "..."
