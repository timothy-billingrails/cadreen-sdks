import { describe, it, expect } from "vitest";
import {
  NoOpProvider,
  NoOpSpan,
  NoOpMeter,
  OpenTelemetryAdapter,
  wrapWithTelemetry,
} from "../telemetry";
import { CadreenError } from "../client";

describe("NoOpProvider", () => {
  it("startSpan returns NoOpSpan without throwing", () => {
    const provider = new NoOpProvider();
    const span = provider.startSpan("test");
    expect(span).toBeInstanceOf(NoOpSpan);
  });

  it("getMeter returns NoOpMeter without throwing", () => {
    const provider = new NoOpProvider();
    const meter = provider.getMeter();
    expect(meter).toBeInstanceOf(NoOpMeter);
  });
});

describe("NoOpSpan", () => {
  it("setName, setAttribute, setStatus, end do not throw", () => {
    const span = new NoOpSpan();
    expect(() => span.setName("test")).not.toThrow();
    expect(() => span.setAttribute("key", "value")).not.toThrow();
    expect(() => span.setStatus("ok")).not.toThrow();
    expect(() => span.setStatus("error")).not.toThrow();
    expect(() => span.end()).not.toThrow();
  });
});

describe("NoOpMeter", () => {
  it("recordRequest, recordRetry, recordStreamEvent do not throw", () => {
    const meter = new NoOpMeter();
    expect(() => meter.recordRequest("GET", "/test", 200, 100)).not.toThrow();
    expect(() => meter.recordRetry("GET", "/test", 1)).not.toThrow();
    expect(() => meter.recordStreamEvent("message")).not.toThrow();
  });
});

describe("OpenTelemetryAdapter", () => {
  it("constructs adapter pattern with mock tracer and meter", () => {
    const mockSpan = {
      setName: () => {},
      setAttribute: () => {},
      setStatus: () => {},
      end: () => {},
      updateName: () => {},
    };

    const mockTracer = {
      startSpan: (_name: string, _opts?: unknown) => mockSpan,
    };

    const mockMeter = {
      createCounter: (_name: string, _opts?: unknown) => ({
        add: (_value: number, _attrs?: unknown) => {},
      }),
      createHistogram: (_name: string, _opts?: unknown) => ({
        record: (_value: number, _attrs?: unknown) => {},
      }),
    };

    const adapter = new OpenTelemetryAdapter(mockTracer, mockMeter);
    expect(adapter).toBeDefined();

    const span = adapter.startSpan("test.op", { attributes: { "http.method": "GET" } });
    expect(span).toBeDefined();
    expect(() => span.setName("updated")).not.toThrow();
    expect(() => span.setAttribute("key", "val")).not.toThrow();
    expect(() => span.setStatus("ok")).not.toThrow();
    expect(() => span.end()).not.toThrow();

    const meter = adapter.getMeter();
    expect(meter).toBeDefined();
    expect(() => meter.recordRequest("GET", "/test", 200, 150)).not.toThrow();
    expect(() => meter.recordRetry("GET", "/test", 1)).not.toThrow();
    expect(() => meter.recordStreamEvent("data")).not.toThrow();
  });

  it("passes parent span context as attributes", () => {
    const attrLog: Array<Record<string, unknown>> = [];
    const mockSpan = {
      updateName: () => {},
      setAttribute: (key: string, val: string) => attrLog.push({ [key]: val }),
      setStatus: () => {},
      end: () => {},
    };

    const mockTracer = {
      startSpan: (_name: string, opts?: { attributes?: Record<string, unknown> }) => {
        attrLog.push(opts?.attributes || {});
        return mockSpan;
      },
    };

    const mockMeter = {
      createCounter: () => ({ add: () => {} }),
      createHistogram: () => ({ record: () => {} }),
    };

    const adapter = new OpenTelemetryAdapter(mockTracer, mockMeter);
    const parent = { traceId: "trace-parent", spanId: "span-parent" };
    adapter.startSpan("child.op", { parent, attributes: { custom: "val" } });

    // First push is startSpan attributes
    expect(attrLog[0]).toEqual({ custom: "val" });
    // Second and third are setAttribute calls for parent context
    expect(attrLog[1]).toEqual({ "parent.trace_id": "trace-parent" });
    expect(attrLog[2]).toEqual({ "parent.span_id": "span-parent" });
  });
});

describe("wrapWithTelemetry", () => {
  it("returns hooks object", () => {
    const hooks = wrapWithTelemetry(new NoOpProvider());
    expect(hooks).toHaveProperty("onRequestStart");
    expect(hooks).toHaveProperty("onRequestEnd");
    expect(hooks).toHaveProperty("onRetry");
    expect(hooks).toHaveProperty("onStreamEvent");
    expect(hooks).toHaveProperty("onError");
  });

  it("onRequestStart returns a TelemetrySpan", () => {
    const hooks = wrapWithTelemetry(new NoOpProvider());
    const span = hooks.onRequestStart("GET", "/api/test");
    expect(span).toBeInstanceOf(NoOpSpan);
  });

  it("onRequestEnd does not throw", () => {
    const hooks = wrapWithTelemetry(new NoOpProvider());
    const span = hooks.onRequestStart("GET", "/api/test");
    expect(() => hooks.onRequestEnd(span, "GET", "/api/test", 200, 100)).not.toThrow();
  });

  it("onRetry does not throw", () => {
    const hooks = wrapWithTelemetry(new NoOpProvider());
    expect(() => hooks.onRetry("GET", "/api/test", 1)).not.toThrow();
  });

  it("onStreamEvent does not throw", () => {
    const hooks = wrapWithTelemetry(new NoOpProvider());
    expect(() => hooks.onStreamEvent("message")).not.toThrow();
  });

  it("onError does not throw", () => {
    const hooks = wrapWithTelemetry(new NoOpProvider());
    const span = hooks.onRequestStart("GET", "/api/test");
    const err = new CadreenError(500, "internal", "error", "Server error");
    expect(() => hooks.onError(span, err)).not.toThrow();
  });
});
