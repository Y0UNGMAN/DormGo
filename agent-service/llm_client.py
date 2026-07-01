import json
import re
import urllib.error
import urllib.request


class OpenAICompatibleClient:
    def __init__(self, api_url: str, api_key: str, model: str, timeout: int = 30):
        self.api_url = api_url
        self.api_key = api_key
        self.model = model
        self.timeout = timeout

    def complete_json(self, messages):
        payload = {
            "model": self.model,
            "messages": messages,
            "temperature": 0.2,
            "response_format": {"type": "json_object"},
        }
        body = json.dumps(payload, ensure_ascii=False).encode("utf-8")
        req = urllib.request.Request(
            self.api_url,
            data=body,
            headers={
                "Content-Type": "application/json",
                "Authorization": f"Bearer {self.api_key}",
            },
            method="POST",
        )
        try:
            with urllib.request.urlopen(req, timeout=self.timeout) as resp:
                result = json.loads(resp.read().decode("utf-8"))
        except urllib.error.HTTPError as exc:
            text = exc.read().decode("utf-8", errors="ignore")
            raise RuntimeError(f"LLM HTTP error {exc.code}: {text[:300]}") from exc
        except urllib.error.URLError as exc:
            raise RuntimeError(f"LLM service unavailable: {exc.reason}") from exc

        if result.get("error"):
            raise RuntimeError(result["error"].get("message") or "LLM returned an error")
        choices = result.get("choices") or []
        if not choices:
            raise RuntimeError("LLM returned no choices")

        content = ((choices[0].get("message") or {}).get("content") or "").strip()
        return _parse_json_content(content)


def _parse_json_content(content: str) -> dict:
    if not content:
        raise RuntimeError("LLM returned empty content")
    try:
        return json.loads(content)
    except json.JSONDecodeError:
        match = re.search(r"\{.*\}", content, flags=re.S)
        if not match:
            raise
        return json.loads(match.group(0))
