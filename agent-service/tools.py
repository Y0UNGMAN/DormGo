import json
import urllib.error
import urllib.request


class GoToolClient:
    def __init__(self, base_url: str, internal_token: str):
        self.base_url = base_url.rstrip("/")
        self.internal_token = internal_token

    def search_posts(self, user_id: int, query: str, filters=None, limit: int = 5):
        payload = {"user_id": user_id, "keyword": self._keyword(query, filters), "limit": limit}
        data = self._post("/internal/agent/tools/search_posts", payload)
        return data.get("data") or []

    def search_signup_posts(self, user_id: int, query: str, filters=None, limit: int = 5):
        payload = {"user_id": user_id, "keyword": self._keyword(query, filters), "limit": limit}
        data = self._post("/internal/agent/tools/search_signup_posts", payload)
        return data.get("data") or []

    def prepare_signup(self, user_id: int, post_id: int):
        payload = {"user_id": user_id, "post_id": post_id}
        return self._post("/internal/agent/tools/prepare_signup", payload)

    def prepare_signup_from_candidate(self, user_id: int, index: int = 1):
        payload = {"user_id": user_id, "index": index}
        return self._post("/internal/agent/tools/prepare_signup_from_candidate", payload)

    def get_post_detail(self, user_id: int, post_id: int):
        payload = {"user_id": user_id, "post_id": post_id}
        return self._post("/internal/agent/tools/post_detail", payload)

    def get_post_detail_from_candidate(self, user_id: int, index: int = 1):
        payload = {"user_id": user_id, "index": index}
        return self._post("/internal/agent/tools/post_detail_from_candidate", payload)

    def confirm_signup(self, user_id: int, action_id: str = ""):
        payload = {"user_id": user_id, "action_id": action_id}
        return self._post("/internal/agent/tools/confirm_signup", payload)

    def get_session(self, user_id: int):
        payload = {"user_id": user_id}
        data = self._post("/internal/agent/tools/get_session", payload)
        return data.get("data") or {"user_id": user_id, "last_candidates": [], "history": []}

    def get_all_posts(self, limit: int = 100):
        payload = {"limit": limit}
        data = self._post("/internal/agent/tools/all_posts", payload)
        return data.get("data") or []

    def save_session(self, user_id: int, last_candidates=None, history=None, selected_post=None, pending_action=None):
        payload = {
            "user_id": user_id,
            "last_candidates": last_candidates or [],
            "history": history or [],
            "selected_post": selected_post or {},
            "pending_action": pending_action or {},
        }
        return self._post("/internal/agent/tools/save_session", payload)

    def _keyword(self, query: str, filters) -> str:
        parts = []
        if query:
            parts.append(query)
        filters = filters or {}
        location = filters.get("location") or ""
        time_hint = filters.get("time_hint") or ""
        for term in filters.get("search_terms") or []:
            self._append_unique(parts, term)
        base = parts[0] if parts else ""
        if location and location not in base:
            self._append_unique(parts, location)
        if time_hint and time_hint not in base:
            self._append_unique(parts, time_hint)
        return " ".join(part.strip() for part in parts if part and part.strip())

    def _append_unique(self, parts, value):
        value = (value or "").strip()
        if value and value not in parts:
            parts.append(value)

    def _post(self, path: str, payload: dict):
        body = json.dumps(payload).encode("utf-8")
        req = urllib.request.Request(
            self.base_url + path,
            data=body,
            headers={
                "Content-Type": "application/json",
                "X-Agent-Token": self.internal_token,
            },
            method="POST",
        )
        try:
            with urllib.request.urlopen(req, timeout=10) as resp:
                result = json.loads(resp.read().decode("utf-8"))
        except urllib.error.HTTPError as exc:
            raise RuntimeError(f"Go tool HTTP error: {exc.code}") from exc
        except urllib.error.URLError as exc:
            raise RuntimeError(f"Go tool unavailable: {exc.reason}") from exc

        if result.get("code") != 200:
            raise RuntimeError(result.get("msg") or "Go tool returned failure")
        return result
