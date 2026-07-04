import tempfile
import unittest
from pathlib import Path

from config import load_go_api_config
from llm_client import LLMResult
from main import choose_candidate_for_signup, format_related_fallback, render_history_for_planner
from planner import LLMPlanner
from react_loop import ReActLoop
from reranker import LLMReranker
from tool_registry import ToolResult, ToolSpec, ToolRegistry, build_default_registry
from tracing import Tracer


class FakeLLM:
    def __init__(self, payload, usage=None):
        self.payload = payload
        self._usage = usage or {"prompt_tokens": 0, "completion_tokens": 0, "total_tokens": 0}
        self.calls = []

    def complete_json(self, messages):
        self.calls.append(messages)
        return self.payload

    def complete_json_with_usage(self, messages):
        self.calls.append(messages)
        return LLMResult(data=self.payload, usage=self._usage)


class ToolRegistryTest(unittest.TestCase):
    def test_registers_and_describes_tools(self):
        registry = ToolRegistry()
        registry.register(ToolSpec(
            name="echo",
            description="echo back",
            parameters={
                "type": "object",
                "properties": {"text": {"type": "string"}},
                "required": ["text"],
            },
            executor=lambda user_id=0, text="": f"echo: {text}",
        ))
        description = registry.describe()
        self.assertIn("echo", description)
        self.assertIn("text", description)

    def test_executes_registered_tool(self):
        registry = ToolRegistry()
        registry.register(ToolSpec(
            name="add",
            description="add two numbers",
            parameters={
                "type": "object",
                "properties": {"a": {"type": "integer"}, "b": {"type": "integer"}},
                "required": ["a", "b"],
            },
            executor=lambda user_id=0, a=0, b=0: {"data": a + b},
        ))
        result = registry.execute("add", user_id=1, a=3, b=5)
        self.assertIn("8", result)

    def test_unknown_tool_returns_error_result(self):
        registry = ToolRegistry()
        result = registry.execute("nonexistent")
        self.assertFalse(result.ok)
        self.assertIn("unknown tool", result.observation)

    def test_tool_error_returns_message(self):
        registry = ToolRegistry()
        registry.register(ToolSpec(
            name="fail",
            description="always fails",
            parameters={"type": "object", "properties": {}},
            executor=lambda user_id=0: (_ for _ in ()).throw(RuntimeError("boom")),
        ))
        result = registry.execute("fail", user_id=1)
        self.assertIn("tool fail failed", result)

    def test_search_result_exposes_candidate_state_patch(self):
        registry = ToolRegistry()
        registry.register(ToolSpec(
            name="search_posts",
            description="search posts",
            parameters={"type": "object", "properties": {}},
            executor=lambda user_id=0: {"data": [
                {"id": 51, "title": "周末想去爬山，有人一起吗", "is_limited": True},
            ]},
        ))

        result = registry.execute("search_posts", user_id=1)

        self.assertIsInstance(result, ToolResult)
        self.assertIn("found 1 candidates", result.observation)
        self.assertEqual(result.state_patch["last_candidates"][0]["post_id"], 51)
        self.assertTrue(result.state_patch["last_candidates"][0]["is_limited"])

    def test_vector_signup_search_filters_non_signup_posts(self):
        class FakeVectorStore:
            def is_ready(self):
                return True

            def search(self, query, k=20):
                return [{"id": 1, "score": 0.91}, {"id": 2, "score": 0.89}]

        class FakeGoClient:
            def get_all_posts(self, limit=100):
                return [
                    {"id": 1, "title": "普通欢迎帖", "is_limited": False},
                    {"id": 2, "title": "周末爬山报名", "is_limited": True, "max_enrollment": 10, "current_enrollment": 2},
                ]

        registry = build_default_registry(FakeGoClient(), vector_store=FakeVectorStore())

        result = registry.execute("search_signup_posts", user_id=1, query="出去玩")

        self.assertIn("周末爬山报名", result.observation)
        self.assertNotIn("普通欢迎帖", result.observation)
        self.assertEqual(result.state_patch["last_candidates"][0]["post_id"], 2)

    def test_vector_search_filters_low_similarity_noise(self):
        class FakeVectorStore:
            def is_ready(self):
                return True

            def search(self, query, k=20):
                return [
                    {"id": 1, "score": 0.50},
                    {"id": 2, "score": 0.49},
                    {"id": 3, "score": 0.46},
                ]

        class FakeGoClient:
            def get_all_posts(self, limit=100):
                return [
                    {"id": 1, "title": "周末想去爬山，有人一起吗", "content": "", "is_limited": True},
                    {"id": 2, "title": "英语展示找队友补位", "content": "", "is_limited": True},
                    {"id": 3, "title": "操作系统实验进程同步卡住了", "content": "", "is_limited": True},
                ]

        registry = build_default_registry(FakeGoClient(), vector_store=FakeVectorStore())

        result = registry.execute("search_signup_posts", user_id=1, query="一起去网吧")

        self.assertIn("no matching posts", result.observation)
        self.assertEqual(result.state_patch["last_candidates"], [])

    def test_vector_search_keeps_semantic_outing_matches_below_strict_threshold(self):
        class FakeVectorStore:
            def is_ready(self):
                return True

            def search(self, query, k=20):
                return [
                    {"id": 51, "score": 0.54},
                    {"id": 45, "score": 0.53},
                    {"id": 33, "score": 0.50},
                ]

        class FakeGoClient:
            def get_all_posts(self, limit=100):
                return [
                    {"id": 51, "title": "周末想去爬山，有人一起吗", "content": "天气好出去走走", "is_limited": True},
                    {"id": 45, "title": "找同学一起参加志愿活动", "content": "周末校外活动", "is_limited": True},
                    {"id": 33, "title": "操作系统实验进程同步卡住了", "content": "", "is_limited": True},
                ]

        registry = build_default_registry(FakeGoClient(), vector_store=FakeVectorStore())

        result = registry.execute("search_posts", user_id=1, query="一起出去玩")

        self.assertIn("周末想去爬山", result.observation)
        self.assertIn("找同学一起参加志愿活动", result.observation)
        self.assertNotIn("操作系统实验", result.observation)

    def test_vector_search_keeps_exact_term_match_below_threshold(self):
        class FakeVectorStore:
            def is_ready(self):
                return True

            def search(self, query, k=20):
                return [{"id": 1, "score": 0.50}]

        class FakeGoClient:
            def get_all_posts(self, limit=100):
                return [
                    {"id": 1, "title": "今晚一起去网吧开黑", "content": "三缺一", "is_limited": True},
                ]

        registry = build_default_registry(FakeGoClient(), vector_store=FakeVectorStore())

        result = registry.execute("search_signup_posts", user_id=1, query="一起去网吧")

        self.assertIn("今晚一起去网吧开黑", result.observation)

    def test_vector_search_does_not_keep_time_only_exact_match(self):
        class FakeVectorStore:
            def is_ready(self):
                return True

            def search(self, query, k=20):
                return [
                    {"id": 11, "score": 0.50},
                    {"id": 12, "score": 0.49},
                ]

        class FakeGoClient:
            def get_all_posts(self, limit=100):
                return [
                    {"id": 11, "title": "今晚八点操场慢跑有人一起吗", "content": "", "is_limited": True},
                    {"id": 12, "title": "有没有人今晚去快递站", "content": "", "is_limited": True},
                ]

        registry = build_default_registry(FakeGoClient(), vector_store=FakeVectorStore())

        result = registry.execute("search_signup_posts", user_id=1, query="榕9 一起吃饭 今晚")

        self.assertIn("no matching posts", result.observation)
        self.assertEqual(result.state_patch["last_candidates"], [])


class ReActLoopTest(unittest.TestCase):
    def _build_loop(self, step_payloads, usage=None):
        fake = FakeLLM(step_payloads[0], usage=usage)
        fake.calls = []
        regimen = _CallSequenceRegimen(fake, step_payloads, usage)
        registry = ToolRegistry()
        registry.register(ToolSpec(
            name="search",
            description="search posts",
            parameters={
                "type": "object",
                "properties": {"query": {"type": "string"}},
                "required": ["query"],
            },
            executor=lambda user_id=0, query="": {"data": [{"id": 1, "title": "test"}]},
        ))
        return ReActLoop(regimen, registry, max_steps=3)

    def test_finish_on_help_message(self):
        loop = self._build_loop([
            {"thought": "user just said hello", "action": "FINISH", "action_input": {"reply": "Hi!"}},
        ])
        tracer = Tracer("trace-1")
        result = loop.run_sync(1, "hello", {}, tracer)
        self.assertEqual(result["reply"], "Hi!")

    def test_tool_call_then_finish(self):
        loop = self._build_loop([
            {
                "thought": "need to search",
                "action": "search",
                "action_input": {"query": "test"},
            },
            {
                "thought": "found results, done",
                "action": "FINISH",
                "action_input": {"reply": "Found 1 post"},
            },
        ])
        tracer = Tracer("trace-2")
        result = loop.run_sync(1, "find posts", {}, tracer)
        self.assertEqual(result["reply"], "Found 1 post")
        self.assertEqual(result["total_steps"], 2)

    def test_max_steps_limit(self):
        loop = self._build_loop([
            {"thought": "keep going", "action": "search", "action_input": {"query": "x"}},
            {"thought": "still going", "action": "search", "action_input": {"query": "x"}},
            {"thought": "never stops", "action": "search", "action_input": {"query": "x"}},
        ])
        tracer = Tracer("trace-3")
        result = loop.run_sync(1, "endless", {}, tracer)
        self.assertIn("抱歉", result["reply"])

    def test_yields_step_events(self):
        loop = self._build_loop([
            {"thought": "search", "action": "search", "action_input": {"query": "x"}},
            {"thought": "done", "action": "FINISH", "action_input": {"reply": "ok"}},
        ])
        tracer = Tracer("trace-4")
        events = list(loop.run(1, "find", {}, tracer))
        steps = [e for e in events if e["type"] == "step"]
        finish = [e for e in events if e["type"] == "finish"]
        self.assertEqual(len(steps), 2)
        self.assertEqual(len(finish), 1)

    def test_signup_request_continues_from_search_to_prepare_action(self):
        calls = []
        llm = _CallSequenceRegimen(None, [
            {
                "thought": "先搜索可报名活动",
                "action": "search_signup_posts",
                "action_input": {"query": "爬山", "limit": 5},
            },
            {
                "thought": "找到了候选，继续创建待确认报名动作",
                "action": "prepare_signup_from_candidate",
                "action_input": {"index": 1},
            },
            {
                "thought": "待确认动作已创建，等待用户确认",
                "action": "FINISH",
                "action_input": {"reply": "已为《周末想去爬山，有人一起吗》创建待确认报名动作，回复“确认报名”后才会真正报名。"},
            },
        ])
        registry = ToolRegistry()
        registry.register(ToolSpec(
            name="search_signup_posts",
            description="search signup posts",
            parameters={"type": "object", "properties": {"query": {"type": "string"}}},
            executor=lambda user_id=0, query="", limit=5: {"data": [
                {"id": 51, "title": "周末想去爬山，有人一起吗", "is_limited": True},
            ]},
        ))
        registry.register(ToolSpec(
            name="prepare_signup_from_candidate",
            description="prepare signup from candidate",
            parameters={"type": "object", "properties": {"index": {"type": "integer"}}},
            executor=lambda user_id=0, index=1: calls.append(("prepare", index)) or {
                "data": {"action_id": "act-1", "title": "周末想去爬山，有人一起吗"}
            },
        ))
        loop = ReActLoop(llm, registry, max_steps=4)

        result = loop.run_sync(1, "请详细介绍第一个帖子，并给我报名", {}, Tracer("trace-signup"))

        self.assertEqual(calls, [("prepare", 1)])
        self.assertIn("确认报名", result["reply"])
        self.assertEqual(result["state"]["last_candidates"][0]["post_id"], 51)
        self.assertEqual(result["state"]["pending_action"]["action_id"], "act-1")

    def test_prepare_signup_success_finishes_without_repeating_tool(self):
        calls = []
        llm = _CallSequenceRegimen(None, [
            {
                "thought": "创建待确认报名动作",
                "action": "prepare_signup_from_candidate",
                "action_input": {"index": 1},
            },
            {
                "thought": "错误地重复创建",
                "action": "prepare_signup_from_candidate",
                "action_input": {"index": 1},
            },
        ])
        registry = ToolRegistry()
        registry.register(ToolSpec(
            name="prepare_signup_from_candidate",
            description="prepare signup from candidate",
            parameters={"type": "object", "properties": {"index": {"type": "integer"}}},
            executor=lambda user_id=0, index=1: calls.append(("prepare", index)) or {
                "data": {"action_id": f"act-{len(calls)}", "title": "英语展示找队友补位"}
            },
        ))
        loop = ReActLoop(llm, registry, max_steps=4)

        result = loop.run_sync(1, "帮我报名第一个", {
            "last_candidates": [{"post_id": 28, "title": "英语展示找队友补位", "is_limited": True}]
        }, Tracer("trace-prepare-once"))

        self.assertEqual(calls, [("prepare", 1)])
        self.assertIn("英语展示找队友补位", result["reply"])
        self.assertIn("确认报名", result["reply"])
        self.assertEqual(result["state"]["pending_action"]["action_id"], "act-1")

    def test_search_only_request_finishes_after_first_successful_search(self):
        llm = _CallSequenceRegimen(None, [
            {
                "thought": "搜索出去玩的帖子",
                "action": "search_posts",
                "action_input": {"query": "一起出去玩", "limit": 5},
            },
            {
                "thought": "不应该再调用模型",
                "action": "FINISH",
                "action_input": {"reply": "slow path"},
            },
        ])
        registry = ToolRegistry()
        registry.register(ToolSpec(
            name="search_posts",
            description="search posts",
            parameters={"type": "object", "properties": {"query": {"type": "string"}}},
            executor=lambda user_id=0, query="", limit=5: {"data": [
                {"id": 51, "title": "周末想去爬山，有人一起吗", "is_limited": True},
            ]},
        ))
        loop = ReActLoop(llm, registry, max_steps=4)

        result = loop.run_sync(1, "帮我查找一些一起出去玩的帖子", {}, Tracer("trace-fast-search"))

        self.assertEqual(len(llm.calls), 1)
        self.assertIn("周末想去爬山", result["reply"])
        self.assertEqual(result["total_steps"], 1)

    def test_question_style_search_finishes_after_first_successful_search(self):
        llm = _CallSequenceRegimen(None, [
            {
                "thought": "搜索今晚打麻将的可报名活动",
                "action": "search_signup_posts",
                "action_input": {"query": "今晚打麻将", "limit": 5},
            },
            {
                "thought": "不应该继续搜索用户需求",
                "action": "search_signup_posts",
                "action_input": {"query": "用户需求", "limit": 5},
            },
        ])
        registry = ToolRegistry()
        registry.register(ToolSpec(
            name="search_signup_posts",
            description="search signup posts",
            parameters={"type": "object", "properties": {"query": {"type": "string"}}},
            executor=lambda user_id=0, query="", limit=5: {"data": [
                {"id": 29, "title": "三缺一麻将小游戏，纯娱乐", "is_limited": True},
            ]},
        ))
        loop = ReActLoop(llm, registry, max_steps=4)

        result = loop.run_sync(1, "今晚有想打麻将的吗？", {}, Tracer("trace-question-search"))

        self.assertEqual(len(llm.calls), 1)
        self.assertIn("三缺一麻将", result["reply"])
        self.assertEqual(result["total_steps"], 1)

    def test_empty_search_finishes_after_first_search(self):
        llm = _CallSequenceRegimen(None, [
            {
                "thought": "搜索今晚一起吃饭的可报名活动",
                "action": "search_signup_posts",
                "action_input": {"query": "今晚一起吃饭", "limit": 5},
            },
            {
                "thought": "不应该继续搜索用户需求",
                "action": "search_signup_posts",
                "action_input": {"query": "用户需求", "limit": 5},
            },
        ])
        registry = ToolRegistry()
        registry.register(ToolSpec(
            name="search_signup_posts",
            description="search signup posts",
            parameters={"type": "object", "properties": {"query": {"type": "string"}}},
            executor=lambda user_id=0, query="", limit=5: {"data": []},
        ))
        loop = ReActLoop(llm, registry, max_steps=4)

        result = loop.run_sync(1, "有没有今晚一起吃饭的朋友？", {}, Tracer("trace-empty-search"))

        self.assertEqual(len(llm.calls), 1)
        self.assertIn("没有找到匹配的帖子", result["reply"])
        self.assertEqual(result["total_steps"], 1)

    def test_detail_success_finishes_without_repeating_tool(self):
        calls = []
        llm = _CallSequenceRegimen(None, [
            {
                "thought": "读取第二个候选帖子的详情",
                "action": "get_post_detail_from_candidate",
                "action_input": {"index": 2},
            },
            {
                "thought": "错误地重复读取详情",
                "action": "get_post_detail_from_candidate",
                "action_input": {"index": 2},
            },
        ])
        registry = ToolRegistry()
        registry.register(ToolSpec(
            name="get_post_detail_from_candidate",
            description="get detail from candidate",
            parameters={"type": "object", "properties": {"index": {"type": "integer"}}},
            executor=lambda user_id=0, index=1: calls.append(("detail", index)) or {
                "data": {
                    "id": 11,
                    "title": "今晚八点操场慢跑有人一起吗",
                    "content": "不卷速度，跑两三公里就行。",
                    "is_limited": True,
                    "max_enrollment": 6,
                    "current_enrollment": 0,
                    "Dorm": {"dormname": "荔园8号"},
                }
            },
        ))
        loop = ReActLoop(llm, registry, max_steps=4)

        result = loop.run_sync(1, "详细介绍第二个", {
            "last_candidates": [
                {"post_id": 29, "title": "三缺一麻将小游戏，纯娱乐", "is_limited": True},
                {"post_id": 11, "title": "今晚八点操场慢跑有人一起吗", "is_limited": True},
            ]
        }, Tracer("trace-detail-once"))

        self.assertEqual(calls, [("detail", 2)])
        self.assertIn("今晚八点操场慢跑", result["reply"])
        self.assertIn("不卷速度", result["reply"])
        self.assertEqual(result["total_steps"], 1)
        self.assertEqual(result["state"]["selected_post"]["id"], 11)


class _CallSequenceRegimen:
    def __init__(self, fake_llm, payloads, usage_base=None):
        self._inner = fake_llm
        self._payloads = list(payloads)
        self._usage_base = usage_base or {"prompt_tokens": 10, "completion_tokens": 5, "total_tokens": 15}
        self._call_index = 0
        self.calls = []

    def complete_json(self, messages):
        self.calls.append(messages)
        payload = self._payloads[self._call_index % len(self._payloads)]
        self._call_index += 1
        return payload

    def complete_json_with_usage(self, messages):
        self.calls.append(messages)
        payload = self._payloads[self._call_index % len(self._payloads)]
        self._call_index += 1
        return LLMResult(data=payload, usage=dict(self._usage_base))


class AgentRuntimeTest(unittest.TestCase):
    def test_loads_deepseek_config_from_go_yaml(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            config_path = Path(tmpdir) / "config.yaml"
            config_path.write_text(
                """
app:
  port: 8080
api:
  ApiUrl: "https://api.deepseek.com/chat/completions"
  ApiKey: "sk-test"
  ApiModel: "deepseek-chat"
""".strip(),
                encoding="utf-8",
            )

            config = load_go_api_config(str(config_path))

        self.assertEqual(config.api_url, "https://api.deepseek.com/chat/completions")
        self.assertEqual(config.api_key, "sk-test")
        self.assertEqual(config.model, "deepseek-chat")

    def test_planner_uses_llm_structured_plan_for_regular_search(self):
        fake_llm = FakeLLM(
            {
                "intent": "search_posts",
                "query": "\u4e00\u8d77\u5b66\u4e60",
                "filters": {"time_hint": "\u6700\u8fd1", "location": "", "require_signup": False},
                "next_action": "show_results",
                "reply_style": "concise",
            }
        )
        planner = LLMPlanner(fake_llm)

        plan = planner.plan("\u5e2e\u6211\u627e\u6700\u8fd1\u4e00\u8d77\u5b66\u4e60\u7684\u5e16\u5b50")

        self.assertEqual(plan["intent"], "search_posts")
        self.assertEqual(plan["query"], "\u4e00\u8d77\u5b66\u4e60")
        self.assertFalse(plan["filters"]["require_signup"])
        self.assertEqual(len(fake_llm.calls), 1)

    def test_planner_uses_llm_structured_plan_for_signup_search(self):
        fake_llm = FakeLLM(
            {
                "intent": "search_signup_posts",
                "query": "\u5357\u6821 \u62fc\u8f66",
                "filters": {"time_hint": "\u4eca\u5929\u4e0b\u5348", "location": "\u5357\u6821", "require_signup": True},
                "next_action": "prepare_signup",
                "reply_style": "concise",
            }
        )
        planner = LLMPlanner(fake_llm)

        plan = planner.plan("\u5e2e\u6211\u627e\u4eca\u5929\u4e0b\u5348\u62fc\u8f66\u53bb\u5357\u6821\u7684\u5e16\u5b50\uff0c\u5982\u679c\u6709\u5e2e\u6211\u62a5\u540d")

        self.assertEqual(plan["intent"], "search_signup_posts")
        self.assertTrue(plan["filters"]["require_signup"])
        self.assertEqual(plan["next_action"], "prepare_signup")

    def test_choose_candidate_defaults_to_last_visible_item(self):
        session = {
            "last_candidates": [
                {"post_id": 91, "title": "\u5468\u4e94\u60f3\u53bb\u770b\u7535\u5f71", "is_limited": True},
                {"post_id": 92, "title": "\u4e00\u8d77\u5b66\u4e60", "is_limited": False},
            ]
        }
        plan = {"target": {"index": 1}}

        candidate = choose_candidate_for_signup(session, plan)

        self.assertEqual(candidate["post_id"], 91)
        self.assertTrue(candidate["is_limited"])

    def test_history_rendering_keeps_recent_turns(self):
        session = {
            "history": [
                {"role": "user", "content": "\u4f60\u597d"},
                {"role": "assistant", "content": "\u53ef\u4ee5\u5e2e\u4f60\u627e\u5e16\u5b50"},
                {"role": "user", "content": "\u5e2e\u6211\u627e\u6253\u8f66"},
            ]
        }

        text = render_history_for_planner(session)

        self.assertIn("user: \u5e2e\u6211\u627e\u6253\u8f66", text)
        self.assertIn("assistant: \u53ef\u4ee5\u5e2e\u4f60\u627e\u5e16\u5b50", text)

    def test_reranker_uses_llm_scores_instead_of_rule_expansion(self):
        posts = [
            {"id": 1, "title": "\u51fa\u4e8c\u624b\u53f0\u706f", "content": "\u5bbf\u820d\u81ea\u63d0"},
            {"id": 2, "title": "\u56fe\u4e66\u9986\u81ea\u4e60\u642d\u5b50", "content": "\u665a\u4e0a\u4e00\u8d77\u590d\u4e60\u9ad8\u6570"},
        ]
        fake_llm = FakeLLM(
            {
                "matches": [
                    {"id": 2, "score": 0.94, "reason": "\u81ea\u4e60\u642d\u5b50\u548c\u4e00\u8d77\u5b66\u4e60\u9ad8\u5ea6\u76f8\u5173"},
                    {"id": 1, "score": 0.05, "reason": "\u4e8c\u624b\u53f0\u706f\u4e0d\u76f8\u5173"},
                ]
            }
        )

        ranked = LLMReranker(fake_llm).rank("\u4e00\u8d77\u5b66\u4e60", posts)

        self.assertEqual(ranked[0]["id"], 2)
        self.assertEqual(ranked[0]["match_reason"], "\u81ea\u4e60\u642d\u5b50\u548c\u4e00\u8d77\u5b66\u4e60\u9ad8\u5ea6\u76f8\u5173")
        self.assertEqual(len(ranked), 1)

    def test_related_fallback_explains_relaxed_results(self):
        reply = format_related_fallback(
            [{"title": "\u5468\u4e94\u60f3\u53bb\u770b\u7535\u5f71\uff0c\u6709\u6ca1\u6709\u62fc\u7968\u7684", "Dorm": {"dormname": "\u6995\u56ed8\u53f7"}}],
            {"query": "\u4e00\u8d77\u770b\u7535\u5f71"},
            signup_only=False,
        )

        self.assertIn("\u6ca1\u6709\u627e\u5230\u5b8c\u5168\u5339\u914d", reply)
        self.assertIn("\u5468\u4e94\u60f3\u53bb\u770b\u7535\u5f71", reply)


class TracerTest(unittest.TestCase):
    def test_creates_nested_spans(self):
        tracer = Tracer("trace-x")
        root = tracer.start_span("loop")
        child = tracer.start_span("tool")
        tracer.end_span(child, status="ok")
        tracer.end_span(root, steps=1)
        dump = tracer.dump()
        self.assertEqual(dump["trace_id"], "trace-x")
        self.assertEqual(len(dump["spans"]), 2)
        self.assertGreaterEqual(dump["total_duration_ms"], 0)

    def test_dump_has_step_count(self):
        tracer = Tracer()
        s1 = tracer.start_span("react_step", step=1)
        tracer.end_span(s1)
        s2 = tracer.start_span("react_step", step=2)
        tracer.end_span(s2)
        dump = tracer.dump()
        self.assertEqual(dump["step_count"], 2)


class LLMResultTest(unittest.TestCase):
    def test_llm_result_dataclass(self):
        result = LLMResult(data={"hello": "world"}, usage={"total_tokens": 42})
        self.assertEqual(result.data["hello"], "world")
        self.assertEqual(result.usage["total_tokens"], 42)


if __name__ == "__main__":
    unittest.main()
