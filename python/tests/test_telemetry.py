import pytest

from cadreen.telemetry import (
    NoOpSpan,
    NoOpMeter,
    NoOpProvider,
    TelemetryHooks,
)


class TestNoOpSpan:
    def test_noop_span_does_not_throw(self):
        span = NoOpSpan()
        span.set_name("test")
        span.set_attribute("key", "value")
        span.set_attribute("count", 42)
        span.set_attribute("flag", True)
        span.set_attribute("ratio", 0.5)
        span.set_status("ok")
        span.set_status("error")
        span.end()


class TestNoOpMeter:
    def test_noop_meter_does_not_throw(self):
        meter = NoOpMeter()
        meter.record_request("GET", "/test", 200, 42.5)
        meter.record_retry("POST", "/api", 3)
        meter.record_stream_event("message")


class TestNoOpProvider:
    def test_noop_provider_start_span(self):
        provider = NoOpProvider()
        span = provider.start_span("test_span")
        assert isinstance(span, NoOpSpan)
        span.set_name("renamed")
        span.set_attribute("custom", "attr")
        span.set_status("ok")
        span.end()

    def test_noop_provider_start_span_with_parent(self):
        provider = NoOpProvider()
        parent_span = provider.start_span("parent")
        child_span = provider.start_span("child", parent=parent_span)
        assert isinstance(child_span, NoOpSpan)

    def test_noop_provider_start_span_with_attributes(self):
        provider = NoOpProvider()
        span = provider.start_span("test", attributes={"http.method": "GET"})
        assert isinstance(span, NoOpSpan)

    def test_noop_provider_get_meter(self):
        provider = NoOpProvider()
        meter = provider.get_meter()
        assert isinstance(meter, NoOpMeter)


class TestTelemetryHooks:
    def test_hooks_use_noop_by_default(self):
        hooks = TelemetryHooks()
        assert isinstance(hooks.provider, NoOpProvider)

    def test_on_request_start_returns_span(self):
        hooks = TelemetryHooks()
        span = hooks.on_request_start("GET", "/test")
        assert span is not None

    def test_on_request_end_success(self):
        hooks = TelemetryHooks()
        span = hooks.provider.start_span("test")
        hooks.on_request_end(span, "GET", "/test", 200, 100.0)

    def test_on_request_end_error(self):
        hooks = TelemetryHooks()
        span = hooks.provider.start_span("test")
        hooks.on_request_end(span, "POST", "/bad", 500, 5000.0)

    def test_on_retry(self):
        hooks = TelemetryHooks()
        hooks.on_retry("POST", "/flaky", 2)

    def test_on_stream_event(self):
        hooks = TelemetryHooks()
        hooks.on_stream_event("message")

    def test_on_error(self):
        hooks = TelemetryHooks()
        span = hooks.provider.start_span("failing")
        hooks.on_error(span, "timeout", "network")

    def test_on_request_end_4xx_is_error(self):
        hooks = TelemetryHooks()
        span = hooks.provider.start_span("test")
        hooks.on_request_end(span, "GET", "/unauthorized", 401, 200.0)

    def test_hooks_with_custom_provider(self):
        """TelemetryHooks should accept a custom provider and delegate to it"""
        provider = NoOpProvider()
        hooks = TelemetryHooks(provider=provider)
        span = hooks.on_request_start("GET", "/custom")
        assert span is not None
        hooks.on_request_end(span, "GET", "/custom", 200, 50.0)
