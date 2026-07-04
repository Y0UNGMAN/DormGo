import time
import uuid
from dataclasses import dataclass, field
from typing import Optional


@dataclass
class Span:
    span_id: str
    trace_id: str
    parent_id: Optional[str]
    name: str
    start_time: float
    end_time: float = 0
    attributes: dict = field(default_factory=dict)

    @property
    def duration_ms(self) -> float:
        return round((self.end_time - self.start_time) * 1000, 2)


class Tracer:
    def __init__(self, trace_id: str = None):
        self.trace_id = trace_id or uuid.uuid4().hex
        self.spans: list[Span] = []
        self._stack: list[Span] = []

    def start_span(self, name: str, **attrs) -> Span:
        parent_id = self._stack[-1].span_id if self._stack else None
        span = Span(
            span_id=uuid.uuid4().hex[:8],
            trace_id=self.trace_id,
            parent_id=parent_id,
            name=name,
            start_time=time.time(),
            attributes=attrs,
        )
        self._stack.append(span)
        return span

    def end_span(self, span: Span, **attrs):
        span.end_time = time.time()
        span.attributes.update(attrs)
        if self._stack and self._stack[-1].span_id == span.span_id:
            self._stack.pop()
        self.spans.append(span)

    def dump(self) -> dict:
        spans = []
        for s in self.spans:
            spans.append({
                "span_id": s.span_id,
                "parent_id": s.parent_id,
                "name": s.name,
                "duration_ms": s.duration_ms,
                "attributes": s.attributes,
            })
        return {
            "trace_id": self.trace_id,
            "total_duration_ms": round(
                sum(s.duration_ms for s in self.spans if s.parent_id is None), 2
            ),
            "step_count": sum(1 for s in self.spans if s.name == "react_step"),
            "spans": spans,
        }
