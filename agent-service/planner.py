ALLOWED_INTENTS = {
    "search_posts",
    "search_signup_posts",
    "confirm_signup",
    "get_help",
}

SYSTEM_PROMPT = """
你是 DormGo 校园互助平台的任务型 Agent Planner。
你的职责是把用户自然语言转换成严格 JSON 计划，不要执行工具，不要输出解释。

可用工具:
1. search_posts: 搜索普通帖子。用于“找学习帖子、闲置、求助、讨论”等不一定需要报名的请求。
2. search_signup_posts: 搜索可报名帖子。用于用户明确想报名、参加、加入活动，或请求“有就帮我报名”。
3. confirm_signup: 确认上一步已经创建的待确认报名动作。只有用户明确说“确定、确认、报名第一个、可以”等确认语时使用。
4. get_help: 用户意图不属于以上任务时使用。

必须返回 JSON，字段固定:
{
  "intent": "search_posts | search_signup_posts | confirm_signup | get_help",
  "query": "用于数据库搜索的核心关键词，不要包含帮我/找/帖子/报名等动作词",
  "filters": {
    "time_hint": "最近/今天下午/明天/空字符串",
    "location": "南校/北校/宿舍楼等地点，未知为空字符串",
    "search_terms": ["用于放宽搜索的同义词或相关词，最多 6 个"],
    "require_signup": true 或 false
  },
  "next_action": "show_results | prepare_signup | confirm_signup | none",
  "target": {"index": 1},
  "reply_style": "concise"
}

规则:
- “帮我找最近一起学习的帖子” => search_posts，query 应聚焦“一起学习”，search_terms 可包含“学习、自习、学习搭子、图书馆、组队学习”。
- “有没有一起看电影的相关帖子” => search_posts，query 应聚焦“看电影”，search_terms 可包含“电影、看电影、拼票、影院”。
- “有没有一起打球的相关帖子” => search_posts，query 应聚焦“打球”，search_terms 可包含“打球、篮球、羽毛球、足球、乒乓球”。
- “帮我找今天下午拼车去南校的帖子，如果有帮我报名” => search_signup_posts，require_signup=true，next_action=prepare_signup。
- “确定”“报名第一个”“报名这个”“好的，就报名这个”“可以” => confirm_signup。若用户说“这个”，默认 target.index=1，表示上一轮展示的第 1 个候选。
- 不要编造数据库结果。
""".strip()


class LLMPlanner:
    def __init__(self, llm_client):
        self.llm_client = llm_client

    def plan(self, message: str) -> dict:
        raw = self.llm_client.complete_json(
            [
                {"role": "system", "content": SYSTEM_PROMPT},
                {"role": "user", "content": message},
            ]
        )
        return normalize_plan(raw)


def normalize_plan(raw: dict) -> dict:
    if not isinstance(raw, dict):
        raise RuntimeError("LLM plan must be a JSON object")

    intent = raw.get("intent")
    if intent not in ALLOWED_INTENTS:
        intent = "get_help"

    filters = raw.get("filters")
    if not isinstance(filters, dict):
        filters = {}

    require_signup = bool(filters.get("require_signup"))
    if intent == "search_signup_posts":
        require_signup = True
    if intent == "search_posts":
        require_signup = False

    query = str(raw.get("query") or "").strip()
    next_action = raw.get("next_action") or "none"
    if intent == "confirm_signup":
        next_action = "confirm_signup"
    if intent == "search_signup_posts" and next_action not in ("prepare_signup", "show_results"):
        next_action = "show_results"
    if intent == "search_posts":
        next_action = "show_results"

    return {
        "intent": intent,
        "query": query,
        "filters": {
            "time_hint": str(filters.get("time_hint") or "").strip(),
            "location": str(filters.get("location") or "").strip(),
            "search_terms": _normalize_terms(filters.get("search_terms")),
            "require_signup": require_signup,
        },
        "next_action": next_action,
        "target": _normalize_target(raw.get("target")),
        "reply_style": raw.get("reply_style") or "concise",
    }


def _normalize_target(raw_target) -> dict:
    if not isinstance(raw_target, dict):
        return {"index": 1}
    try:
        index = int(raw_target.get("index") or 1)
    except (TypeError, ValueError):
        index = 1
    if index <= 0:
        index = 1
    return {"index": index}


def _normalize_terms(raw_terms) -> list:
    if not isinstance(raw_terms, list):
        return []
    terms = []
    for item in raw_terms:
        term = str(item or "").strip()
        if term and term not in terms:
            terms.append(term)
        if len(terms) >= 8:
            break
    return terms
