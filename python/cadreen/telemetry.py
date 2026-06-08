from __future__ import annotations

import time
from dataclasses import dataclass, field
from typing import Any, Callable, Optional, Protocol, runtime_checkable


@runtime_checkable
class TelemetrySpan(Protocol):
    def set_name(self, name: str) -> None: ...
    def set_attribute(self, key: str, value: str | int | float | bool) -> None: ...
    def set_status(self, status: str) -> None: ...
    def end(self) -> None: ...


@runtime_checkable
class TelemetryMeter(Protocol):
    def record_request(self, method: str, path: str, status: int, duration_ms: float) -> None: ...
    def record_retry(self, method: str, path: str, attempt: int) -> None: ...
    def record_stream_event(self, event_type: str) -> None: ...


@runtime_checkable
class TelemetryProvider(Protocol):
    def start_span(self, name: str, *, parent: Optional[Any] = None, attributes: Optional[dict[str, Any]] = None) -> TelemetrySpan: ...
    def get_meter(self) -> TelemetryMeter: ...


class NoOpSpan:
    def set_name(self, name: str) -> None: pass
    def set_attribute(self, key: str, value: str | int | float | bool) -> None: pass
    def set_status(self, status: str) -> None: pass
    def end(self) -> None: pass


class NoOpMeter:
    def record_request(self, method: str, path: str, status: int, duration_ms: float) -> None: pass
    def record_retry(self, method: str, path: str, attempt: int) -> None: pass
    def record_stream_event(self, event_type: str) -> None: pass


class NoOpProvider:
    def start_span(self, name: str, *, parent: Optional[Any] = None, attributes: Optional[dict[str, Any]] = None) -> TelemetrySpan:
        return NoOpSpan()
    def get_meter(self) -> TelemetryMeter:
        return NoOpMeter()


class OpenTelemetryAdapter:
    def __init__(self, tracer: Any, meter: Any) -> None:
        self._tracer = tracer
        self._request_counter = meter.create_counter(
            "cadreen.client.requests",
            description="Number of API requests made",
            unit="1",
        )
        self._retry_counter = meter.create_counter(
            "cadreen.client.retries",
            description="Number of request retries",
            unit="1",
        )
        self._stream_counter = meter.create_counter(
            "cadreen.client.stream_events",
            description="Number of SSE events received",
            unit="1",
        )
        self._duration_histogram = meter.create_histogram(
            "cadreen.client.request_duration",
            description="Request duration in milliseconds",
            unit="ms",
        )

    def start_span(self, name: str, *, parent: Optional[Any] = None, attributes: Optional[dict[str, Any]] = None) -> TelemetrySpan:
        kwargs: dict[str, Any] = {}
        if attributes:
            kwargs["attributes"] = attributes
        span = self._tracer.start_span(name, **kwargs)
        if parent:
            span.set_attribute("parent.trace_id", parent.trace_id if hasattr(parent, "trace_id") else str(parent))
            span.set_attribute("parent.span_id", parent.span_id if hasattr(parent, "span_id") else "")
        return _OtelSpanWrapper(span)

    def get_meter(self) -> TelemetryMeter:
        return _OtelMeterWrapper(
            self._request_counter,
            self._retry_counter,
            self._stream_counter,
            self._duration_histogram,
        )


class _OtelSpanWrapper:
    def __init__(self, span: Any) -> None:
        self._span = span

    def set_name(self, name: str) -> None:
        self._span.update_name(name)

    def set_attribute(self, key: str, value: str | int | float | bool) -> None:
        self._span.set_attribute(key, value)

    def set_status(self, status: str) -> None:
        from opentelemetry.trace import StatusCode
        code = StatusCode.ERROR if status == "error" else StatusCode.OK
        self._span.set_status(code)

    def end(self) -> None:
        self._span.end()


class _OtelMeterWrapper:
    def __init__(self, request_counter: Any, retry_counter: Any, stream_counter: Any, duration_histogram: Any) -> None:
        self._request_counter = request_counter
        self._retry_counter = retry_counter
        self._stream_counter = stream_counter
        self._duration_histogram = duration_histogram

    def record_request(self, method: str, path: str, status: int, duration_ms: float) -> None:
        attrs = {"http.method": method, "http.url": path, "http.status_code": status}
        self._request_counter.add(1, attrs)
        self._duration_histogram.record(duration_ms, attrs)

    def record_retry(self, method: str, path: str, attempt: int) -> None:
        self._retry_counter.add(1, {"http.method": method, "http.url": path, "attempt": attempt})

    def record_stream_event(self, event_type: str) -> None:
        self._stream_counter.add(1, {"event.type": event_type})


@dataclass
class TelemetryHooks:
    provider: TelemetryProvider = field(default_factory=NoOpProvider)

    def on_request_start(self, method: str, path: str) -> TelemetrySpan:
        span = self.provider.start_span(
            f"cadreen.{method.lower()}",
            attributes={"http.method": method, "http.url": path, "cadreen.version": "2026-06-03"},
        )
        return span

    def on_request_end(self, span: TelemetrySpan, method: str, path: str, status: int, duration_ms: float) -> None:
        span.set_attribute("http.status_code", status)
        span.set_status("ok" if status < 400 else "error")
        span.end()
        self.provider.get_meter().record_request(method, path, status, duration_ms)

    def on_retry(self, method: str, path: str, attempt: int) -> None:
        self.provider.get_meter().record_retry(method, path, attempt)

    def on_stream_event(self, event_type: str) -> None:
        self.provider.get_meter().record_stream_event(event_type)

    def on_error(self, span: TelemetrySpan, code: str, error_type: str) -> None:
        span.set_attribute("cadreen.error.code", code)
        span.set_attribute("cadreen.error.type", error_type)
        span.set_status("error")
        span.end()
