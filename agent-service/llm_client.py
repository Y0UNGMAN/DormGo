import json
import re
import ssl
import time
from dataclasses import dataclass

import requests
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry


@dataclass
class LLMResult:
    data: dict
    usage: dict


def _create_session(max_retries: int = 3) -> requests.Session:
    """Create a requests.Session with retry logic and TLS adapter."""
    session = requests.Session()

    retry_strategy = Retry(
        total=max_retries,
        backoff_factor=0.5,
        status_forcelist=[429, 500, 502, 503, 504],
        allowed_methods=["POST"],
    )
    adapter = HTTPAdapter(max_retries=retry_strategy)
    session.mount("https://", adapter)
    session.mount("http://", adapter)

    return session


class OpenAICompatibleClient:
    def __init__(self, api_url: str, api_key: str, model: str, timeout: int = 30):
        self.api_url = api_url
        self.api_key = api_key
        self.model = model
        self.timeout = timeout
        self._session = _create_session()

    def complete_json(self, messages):
        result = self._raw_complete(messages)
        return result.data

    def complete_json_with_usage(self, messages) -> LLMResult:
        return self._raw_complete(messages)

    def _raw_complete(self, messages) -> LLMResult:
        payload = {
            "model": self.model,
            "messages": messages,
            "temperature": 0.2,
            "response_format": {"type": "json_object"},
        }

        last_error = None
        for attempt in range(3):
            try:
                resp = self._session.post(
                    self.api_url,
                    json=payload,
                    headers={
                        "Content-Type": "application/json",
                        "Authorization": f"Bearer {self.api_key}",
                    },
                    timeout=self.timeout,
                )
                if resp.status_code >= 400:
                    text = resp.text[:300]
                    if resp.status_code == 429:
                        time.sleep(2 ** attempt)
                        last_error = RuntimeError(f"LLM rate limited: {text}")
                        continue
                    raise RuntimeError(f"LLM HTTP error {resp.status_code}: {text}")

                result = resp.json()
                break
            except requests.ConnectionError as exc:
                last_error = RuntimeError(f"LLM connection error: {exc}")
                if attempt < 2:
                    time.sleep(1 * (attempt + 1))
                    continue
            except requests.Timeout as exc:
                last_error = RuntimeError(f"LLM request timeout: {exc}")
                if attempt < 2:
                    time.sleep(1 * (attempt + 1))
                    continue
        else:
            raise last_error or RuntimeError("LLM request failed after retries")

        if result.get("error"):
            raise RuntimeError(result["error"].get("message") or "LLM returned an error")
        choices = result.get("choices") or []
        if not choices:
            raise RuntimeError("LLM returned no choices")

        content = ((choices[0].get("message") or {}).get("content") or "").strip()
        usage_raw = result.get("usage") or {}
        usage = {
            "prompt_tokens": usage_raw.get("prompt_tokens", 0),
            "completion_tokens": usage_raw.get("completion_tokens", 0),
            "total_tokens": usage_raw.get("total_tokens", 0),
        }
        return LLMResult(data=_parse_json_content(content), usage=usage)


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
