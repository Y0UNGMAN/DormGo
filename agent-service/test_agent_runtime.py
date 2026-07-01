import tempfile
import unittest
from pathlib import Path

from config import load_go_api_config
from main import choose_candidate_for_signup, format_related_fallback, render_history_for_planner
from planner import LLMPlanner
from reranker import LLMReranker
from tools import GoToolClient


class FakeLLM:
    def __init__(self, payload):
        self.payload = payload
        self.calls = []

    def complete_json(self, messages):
        self.calls.append(messages)
        return self.payload


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
                "query": "一起学习",
                "filters": {"time_hint": "最近", "location": "", "require_signup": False},
                "next_action": "show_results",
                "reply_style": "concise",
            }
        )
        planner = LLMPlanner(fake_llm)

        plan = planner.plan("帮我找最近一起学习的帖子")

        self.assertEqual(plan["intent"], "search_posts")
        self.assertEqual(plan["query"], "一起学习")
        self.assertFalse(plan["filters"]["require_signup"])
        self.assertEqual(len(fake_llm.calls), 1)

    def test_planner_uses_llm_structured_plan_for_signup_search(self):
        fake_llm = FakeLLM(
            {
                "intent": "search_signup_posts",
                "query": "南校 拼车",
                "filters": {"time_hint": "今天下午", "location": "南校", "require_signup": True},
                "next_action": "prepare_signup",
                "reply_style": "concise",
            }
        )
        planner = LLMPlanner(fake_llm)

        plan = planner.plan("帮我找今天下午拼车去南校的帖子，如果有帮我报名")

        self.assertEqual(plan["intent"], "search_signup_posts")
        self.assertTrue(plan["filters"]["require_signup"])
        self.assertEqual(plan["next_action"], "prepare_signup")

    def test_choose_candidate_defaults_to_last_visible_item(self):
        session = {
            "last_candidates": [
                {"post_id": 91, "title": "周五想去看电影", "is_limited": True},
                {"post_id": 92, "title": "一起学习", "is_limited": False},
            ]
        }
        plan = {"target": {"index": 1}}

        candidate = choose_candidate_for_signup(session, plan)

        self.assertEqual(candidate["post_id"], 91)
        self.assertTrue(candidate["is_limited"])

    def test_history_rendering_keeps_recent_turns(self):
        session = {
            "history": [
                {"role": "user", "content": "你好"},
                {"role": "assistant", "content": "可以帮你找帖子"},
                {"role": "user", "content": "帮我找打车"},
            ]
        }

        text = render_history_for_planner(session)

        self.assertIn("user: 帮我找打车", text)
        self.assertIn("assistant: 可以帮你找帖子", text)

    def test_reranker_uses_llm_scores_instead_of_rule_expansion(self):
        posts = [
            {"id": 1, "title": "出二手台灯", "content": "宿舍自提"},
            {"id": 2, "title": "图书馆自习搭子", "content": "晚上一起复习高数"},
        ]
        fake_llm = FakeLLM(
            {
                "matches": [
                    {"id": 2, "score": 0.94, "reason": "自习搭子和一起学习高度相关"},
                    {"id": 1, "score": 0.05, "reason": "二手台灯不相关"},
                ]
            }
        )

        ranked = LLMReranker(fake_llm).rank("一起学习", posts)

        self.assertEqual(ranked[0]["id"], 2)
        self.assertEqual(ranked[0]["match_reason"], "自习搭子和一起学习高度相关")
        self.assertEqual(len(ranked), 1)

    def test_related_fallback_explains_relaxed_results(self):
        reply = format_related_fallback(
            [{"title": "周五想去看电影，有没有拼票的", "Dorm": {"dormname": "榕园8号"}}],
            {"query": "一起看电影"},
            signup_only=False,
        )

        self.assertIn("没有找到完全匹配", reply)
        self.assertIn("周五想去看电影", reply)


if __name__ == "__main__":
    unittest.main()
