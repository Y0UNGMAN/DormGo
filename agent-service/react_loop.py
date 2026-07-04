import json
from dataclasses import dataclass, field

from llm_client import OpenAICompatibleClient
from tool_registry import ToolRegistry
from tracing import Tracer

REACT_SYSTEM_PROMPT = """\
You are a ReAct agent for DormGo campus platform.
Your job: use available tools to help users find posts and sign up for events.
Work iteratively: THINK about the situation, then take one ACTION at a time.
When the task is done, set action to "FINISH".

CRITICAL: Output ONLY valid JSON. No markdown, no explanation outside JSON.

Format per step:
  thought: your reasoning about what you should do next
  action: either a tool_name or the word FINISH
  action_input: parameters for the tool, or for FINISH use {{"reply": "message to user"}}

Available tools:
__TOOLS__

Search strategy (IMPORTANT):
- The search is SEMANTIC vector search (FAISS + BGE). It finds conceptually related posts, NOT keyword match.
- "出去玩" will ALREADY match "郊游", "爬山", "逛街", "看电影" etc. — you do NOT need to re-search with synonyms.
- Do at MOST 2 searches total. After that, pick the best results and FINISH.
- If the first search returned posts and the user only asked to search, PRESENT them to the user.
- If the user asked for details or signup, continue with get_post_detail_from_candidate or prepare_signup_from_candidate.
- Only search a second time if the first search truly returned zero results.

General Rules:
- Only call one tool per step.
- When search results are returned, check them. If the user only asked to search, present the best matches via FINISH.
- If the user also asked for details or signup, continue with the appropriate detail/signup tool before FINISH.
- If you searched 2+ times with different keywords and found nothing, FINISH: tell the user "没有找到相关帖子，建议换个关键词试试".
- Before signup, ALWAYS call prepare_signup first. confirm_signup only after user explicit consent.
- If user refers to "第一个/第二个/这个帖子", use tools ending with _from_candidate instead of searching again.
- If user asks to introduce/detail a previous candidate, call get_post_detail_from_candidate.
- If user asks to sign up for a previous candidate, call prepare_signup_from_candidate, then FINISH asking the user to confirm.
- When answering the user: use FINISH with a clear reply in Chinese.
- If the user message is casual ("hello", "help"), FINISH directly with a friendly reply.
- NEVER repeat the same failed search. Each search must use different keywords from the last.
"""


@dataclass
class Step:
    thought: str
    action: str
    action_input: dict
    observation: str = ""


class ReActLoop:
    def __init__(
        self,
        llm_client: OpenAICompatibleClient,
        tool_registry: ToolRegistry,
        max_steps: int = 6,
    ):
        self.llm = llm_client
        self.tools = tool_registry
        self.max_steps = max_steps

    def run(self, user_id: int, user_message: str, session: dict, tracer: Tracer):
        for event in self._iter(user_id, user_message, session, tracer):
            yield event

    def run_sync(self, user_id: int, user_message: str, session: dict, tracer: Tracer) -> dict:
        final = None
        for event in self._iter(user_id, user_message, session, tracer):
            if event["type"] == "finish":
                final = event["data"]
        return final or {"reply": "Agent did not reach a conclusion"}

    def _iter(self, user_id: int, user_message: str, session: dict, tracer: Tracer):
        import logging
        _log = logging.getLogger("dormgo.react")
        root_span = tracer.start_span("react_loop")
        scratchpad = ""
        state = {
            "last_candidates": list(session.get("last_candidates") or []),
            "pending_action": session.get("pending_action"),
            "selected_post": session.get("selected_post"),
        }
        empty_search_count = 0
        search_count = 0
        history_context = _render_history(session)
        candidates_context = _render_candidates(session)

        system_prompt = REACT_SYSTEM_PROMPT.replace("__TOOLS__", self.tools.describe())
        user_context = f"Context:\n{history_context}\n{candidates_context}\n\nUser: {user_message}"

        yield {"type": "start", "data": {"trace_id": tracer.trace_id}}

        for step_num in range(1, self.max_steps + 1):
            step_span = tracer.start_span("react_step", step=step_num)

            full_prompt = _build_react_prompt(system_prompt, scratchpad, user_context)

            llm_span = tracer.start_span("llm_call", step=step_num)
            try:
                result = self.llm.complete_json_with_usage([
                    {"role": "system", "content": full_prompt},
                    {"role": "user", "content": user_message},
                ])
                decision = result.data
                tracer.end_span(llm_span, tokens=result.usage.get("total_tokens", 0))
            except Exception as exc:
                tracer.end_span(llm_span, error=str(exc))
                tracer.end_span(step_span, error=str(exc))
                tracer.end_span(root_span)
                yield {"type": "error", "data": {"message": f"LLM error: {exc}"}}
                return

            thought = str(decision.get("thought") or "").strip()
            action = str(decision.get("action") or "").strip()
            action_input = decision.get("action_input") or {}

            if not isinstance(action_input, dict):
                action_input = {}

            if action == "FINISH":
                reply = action_input.get("reply", "Task completed")
                step = Step(thought=thought, action="FINISH", action_input=action_input)
                tracer.end_span(step_span, action="FINISH")
                tracer.end_span(root_span)
                yield {
                    "type": "step",
                    "data": _step_to_event(step_num, step, None),
                }
                yield {
                    "type": "finish",
                    "data": {"reply": reply, "total_steps": step_num, "state": state},
                }
                return

            tool_span = tracer.start_span("tool_call", tool=action)
            try:
                tool_result = self.tools.execute(action, user_id=user_id, **action_input)
                observation = tool_result.observation
                _merge_state(state, tool_result.state_patch)
                tracer.end_span(tool_span, status="ok")
            except Exception as exc:
                observation = f"tool call failed: {exc}"
                tracer.end_span(tool_span, error=str(exc))

            # Track search attempts to prevent unbounded retry loops.
            if action in ("search_posts", "search_signup_posts"):
                _log.info("search observation: %s", observation[:200])
                search_count += 1
                if "no matching posts" in observation:
                    empty_search_count += 1
                    _log.info("empty search %d/3", empty_search_count)
                else:
                    empty_search_count = 0

            step = Step(
                thought=thought,
                action=action,
                action_input=action_input,
                observation=observation,
            )
            tracer.end_span(step_span, action=action)

            yield {
                "type": "step",
                "data": _step_to_event(step_num, step, observation),
            }

            scratchpad += _format_step_for_scratchpad(step_num, step)

            if action in ("search_posts", "search_signup_posts"):
                if "no matching posts" in observation:
                    reply = _format_empty_search_reply()
                    tracer.end_span(root_span)
                    yield {
                        "type": "finish",
                        "data": {"reply": reply, "total_steps": step_num, "state": state},
                    }
                    return
                if _is_search_only_request(user_message):
                    reply = _format_search_reply(observation)
                    tracer.end_span(root_span)
                    yield {
                        "type": "finish",
                        "data": {"reply": reply, "total_steps": step_num, "state": state},
                    }
                    return

            if action in ("get_post_detail", "get_post_detail_from_candidate") and state.get("selected_post"):
                reply = _format_post_detail_reply(state["selected_post"])
                tracer.end_span(root_span)
                yield {
                    "type": "finish",
                    "data": {"reply": reply, "total_steps": step_num, "state": state},
                }
                return

            if action in ("prepare_signup", "prepare_signup_from_candidate") and state.get("pending_action"):
                reply = _format_pending_signup_reply(state["pending_action"])
                tracer.end_span(root_span)
                yield {
                    "type": "finish",
                    "data": {"reply": reply, "total_steps": step_num, "state": state},
                }
                return

            if action == "confirm_signup":
                reply = _format_confirm_signup_reply(observation)
                tracer.end_span(root_span)
                yield {
                    "type": "finish",
                    "data": {"reply": reply, "total_steps": step_num, "state": state},
                }
                return

            # Force FINISH after 3 consecutive empty searches
            if empty_search_count >= 3:
                reply = "抱歉，我用不同关键词搜索了多次，但没有找到匹配的帖子。建议你换个更具体的关键词试试，比如具体的地点、时间或活动类型。"
                tracer.end_span(root_span)
                yield {
                    "type": "step",
                    "data": _step_to_event(step_num + 1, Step(
                        thought="多次搜索无结果，停止搜索并告知用户",
                        action="FINISH",
                        action_input={"reply": reply},
                    ), None),
                }
                yield {
                    "type": "finish",
                    "data": {"reply": reply, "total_steps": step_num + 1, "state": state},
                }
                return

            candidates_context = _render_candidates(state)
            user_context = f"Context:\n{history_context}\n{candidates_context}\n\nUser: {user_message}"

        tracer.end_span(root_span)
        yield {
            "type": "finish",
            "data": {"reply": "抱歉，我尝试了多次但没能完成任务。请尝试用更简单的方式描述你的需求。",
                     "total_steps": self.max_steps, "state": state},
        }


def _render_history(session: dict) -> str:
    items = session.get("history") or []
    if not items:
        return "Previous conversation: none"
    lines = ["Previous conversation:"]
    for item in items[-8:]:
        role = item.get("role") or "unknown"
        content = (item.get("content") or "").replace("\n", " ")
        lines.append(f"  {role}: {content}")
    return "\n".join(lines)


def _render_candidates(session: dict) -> str:
    candidates = session.get("last_candidates") or []
    if not candidates:
        return "Last search candidates: none"
    lines = ["Last search candidates:"]
    for idx, c in enumerate(candidates[:5], start=1):
        flag = "signup" if c.get("is_limited") else "no-signup"
        lines.append(f"  {idx}. post_id={c.get('post_id')} title={c.get('title')} {flag}")
    return "\n".join(lines)


def _merge_state(state: dict, patch: dict):
    if not isinstance(patch, dict):
        return
    for key, value in patch.items():
        state[key] = value


def _is_search_only_request(user_message: str) -> bool:
    text = user_message or ""
    action_words = ("报名", "参加", "加入", "确认", "详细", "介绍", "详情", "第一个", "第二个", "第三个", "这个", "那个")
    if any(word in text for word in action_words):
        return False
    search_words = ("找", "查", "搜索", "有没有", "看看", "推荐", "有谁", "有人", "有想")
    if any(word in text for word in search_words):
        return True
    return "吗" in text and any(word in text for word in ("有", "想", "约", "一起"))


def _format_empty_search_reply() -> str:
    return "没有找到匹配的帖子。你可以换个更具体的关键词试试，比如地点、时间或活动类型。"


def _format_search_reply(observation: str) -> str:
    titles = []
    for line in observation.splitlines():
        if " title=" not in line:
            continue
        title_part = line.split(" title=", 1)[1]
        title = title_part.split(" dorm=", 1)[0].strip()
        if title:
            titles.append(title)
    if not titles:
        return "我找到了一些相关帖子。你可以告诉我想查看哪一个，或让我帮你继续筛选。"
    lines = ["我找到了一些相关帖子，你看看有没有感兴趣的："]
    for idx, title in enumerate(titles[:5], start=1):
        lines.append(f"{idx}. {title}")
    lines.append("需要我详细介绍某个帖子，或者帮你处理报名的话，直接告诉我编号。")
    return "\n".join(lines)


def _format_post_detail_reply(post: dict) -> str:
    title = post.get("title") or "未命名帖子"
    content = post.get("content") or "这个帖子暂时没有正文内容。"
    dorm = (post.get("Dorm") or post.get("dorm") or {}).get("dormname", "")
    current = post.get("current_enrollment", 0)
    max_count = post.get("max_enrollment") or 0
    lines = [f"《{title}》"]
    if dorm:
        lines.append(f"地点/宿舍：{dorm}")
    lines.append(f"内容：{content}")
    if post.get("is_limited"):
        cap = max_count if max_count else "不限"
        lines.append(f"报名状态：可报名，当前 {current}/{cap} 人。")
        lines.append("需要报名的话，告诉我“帮我报名这个”。")
    else:
        lines.append("报名状态：这个帖子没有开启报名。")
    return "\n".join(lines)


def _format_pending_signup_reply(pending_action: dict) -> str:
    title = pending_action.get("title") or "该帖子"
    action_id = pending_action.get("action_id") or ""
    suffix = f"\n确认编号：{action_id}" if action_id else ""
    return f"已为《{title}》创建待确认报名动作。回复“确认报名”后才会真正报名。{suffix}"


def _format_confirm_signup_reply(observation: str) -> str:
    if "signup confirmed" in observation:
        return "报名成功。"
    return observation or "报名确认已处理。"


def _build_react_prompt(system_prompt: str, scratchpad: str, user_context: str) -> str:
    if scratchpad:
        return (
            f"{system_prompt}\n\n"
            f"=== Steps so far ===\n{scratchpad}\n"
            f"=== Current task ===\n{user_context}\n\n"
            f"What is your next thought and action? Output JSON."
        )
    return f"{system_prompt}\n\n{user_context}\n\nWhat is your first thought and action? Output JSON."


def _format_step_for_scratchpad(step_num: int, step: Step) -> str:
    return (
        f"Step {step_num}:\n"
        f"Thought: {step.thought}\n"
        f"Action: {step.action}\n"
        f"Action Input: {json.dumps(step.action_input, ensure_ascii=False)}\n"
        f"Observation: {step.observation}\n"
    )


def _step_to_event(step_num: int, step: Step, observation: str):
    return {
        "step": step_num,
        "thought": step.thought,
        "action": step.action,
        "action_input": step.action_input,
        "observation": observation or step.observation,
    }


def _force_finish_after_search(llm, observation: str, user_message: str) -> str:
    """After a search returns results, use one LLM call to compose a friendly reply.
    If the LLM fails, fall back to a generic reply based on the observation."""
    prompt = (
        "You are a campus assistant. You searched for posts and found results.\n"
        f"User request: {user_message}\n"
        f"Search results:\n{observation}\n\n"
        "Write a friendly reply in Chinese listing the top posts and asking if the user wants details.\n"
        "Keep it brief. Output ONLY valid JSON: {\"reply\": \"你的回复\"}"
    )
    try:
        result = llm.complete_json([
            {"role": "system", "content": prompt},
        ])
        reply = (result.get("reply") or "").strip()
        if reply:
            return reply
    except Exception:
        pass
    # Fallback: generic reply based on observation
    return f"我找到了一些相关帖子，你看有没有感兴趣的？如果有需要可以告诉我，我帮你查看详情或报名。\n\n{observation.split(chr(10), 1)[0] if chr(10) in observation else observation}"
