import type { CadreenError } from "./client";

export interface SpanContext {
  traceId: string;
  spanId: string;
}

export interface TelemetrySpan {
  setName(name: string): void;
  setAttribute(key: string, value: string | number | boolean): void;
  setStatus(status: "ok" | "error"): void;
  end(): void;
}

export interface TelemetryMeter {
  recordRequest(method: string, path: string, status: number, durationMs: number): void;
  recordRetry(method: string, path: string, attempt: number): void;
  recordStreamEvent(eventType: string): void;
}

export interface TelemetryProvider {
  startSpan(name: string, options?: SpanOptions): TelemetrySpan;
  getMeter(): TelemetryMeter;
}

export interface SpanOptions {
  parent?: SpanContext;
  attributes?: Record<string, string | number | boolean>;
}

export class NoOpSpan implements TelemetrySpan {
  setName() {}
  setAttribute() {}
  setStatus() {}
  end() {}
}

export class NoOpMeter implements TelemetryMeter {
  recordRequest() {}
  recordRetry() {}
  recordStreamEvent() {}
}

export class NoOpProvider implements TelemetryProvider {
  startSpan(_name: string, _options?: SpanOptions): TelemetrySpan {
    return new NoOpSpan();
  }
  getMeter(): TelemetryMeter {
    return new NoOpMeter();
  }
}

export class OpenTelemetryAdapter implements TelemetryProvider {
  private tracer: any;
  private requestCounter: any;
  private retryCounter: any;
  private streamCounter: any;
  private durationHistogram: any;

  constructor(tracer: any, meter: any) {
    this.tracer = tracer;

    this.requestCounter = meter.createCounter("cadreen.client.requests", {
      description: "Number of API requests made",
      unit: "1",
    });
    this.retryCounter = meter.createCounter("cadreen.client.retries", {
      description: "Number of request retries",
      unit: "1",
    });
    this.streamCounter = meter.createCounter("cadreen.client.stream_events", {
      description: "Number of SSE events received",
      unit: "1",
    });
    this.durationHistogram = meter.createHistogram("cadreen.client.request_duration", {
      description: "Request duration in milliseconds",
      unit: "ms",
    });
  }

  startSpan(name: string, options?: SpanOptions): TelemetrySpan {
    const span = this.tracer.startSpan(name, {
      attributes: options?.attributes,
    });
    if (options?.parent) {
      span.setAttribute("parent.trace_id", options.parent.traceId);
      span.setAttribute("parent.span_id", options.parent.spanId);
    }
    return new OtelSpanAdapter(span);
  }

  getMeter(): TelemetryMeter {
    return new OtelMeterAdapter(
      this.requestCounter,
      this.retryCounter,
      this.streamCounter,
      this.durationHistogram
    );
  }
}

class OtelSpanAdapter implements TelemetrySpan {
  private span: any;

  constructor(span: any) {
    this.span = span;
  }

  setName(name: string): void {
    this.span.updateName(name);
  }

  setAttribute(key: string, value: string | number | boolean): void {
    this.span.setAttribute(key, value);
  }

  setStatus(status: "ok" | "error"): void {
    if (status === "error") {
      this.span.setStatus({ code: 2 });
    } else {
      this.span.setStatus({ code: 1 });
    }
  }

  end(): void {
    this.span.end();
  }
}

class OtelMeterAdapter implements TelemetryMeter {
  private requestCounter: any;
  private retryCounter: any;
  private streamCounter: any;
  private durationHistogram: any;

  constructor(requestCounter: any, retryCounter: any, streamCounter: any, durationHistogram: any) {
    this.requestCounter = requestCounter;
    this.retryCounter = retryCounter;
    this.streamCounter = streamCounter;
    this.durationHistogram = durationHistogram;
  }

  recordRequest(method: string, path: string, status: number, durationMs: number): void {
    this.requestCounter.add(1, { "http.method": method, "http.url": path, "http.status_code": status });
    this.durationHistogram.record(durationMs, { "http.method": method, "http.url": path, "http.status_code": status });
  }

  recordRetry(method: string, path: string, attempt: number): void {
    this.retryCounter.add(1, { "http.method": method, "http.url": path, attempt });
  }

  recordStreamEvent(eventType: string): void {
    this.streamCounter.add(1, { "event.type": eventType });
  }
}

export function wrapWithTelemetry(provider: TelemetryProvider) {
  return {
    onRequestStart(method: string, path: string): TelemetrySpan {
      const span = provider.startSpan(`cadreen.${method.toLowerCase()}`, {
        attributes: {
          "http.method": method,
          "http.url": path,
          "cadreen.version": "2026-06-03",
        },
      });
      return span;
    },
    onRequestEnd(span: TelemetrySpan, method: string, path: string, status: number, durationMs: number): void {
      span.setAttribute("http.status_code", status);
      span.setStatus(status < 400 ? "ok" : "error");
      span.end();
      provider.getMeter().recordRequest(method, path, status, durationMs);
    },
    onRetry(method: string, path: string, attempt: number): void {
      provider.getMeter().recordRetry(method, path, attempt);
    },
    onStreamEvent(eventType: string): void {
      provider.getMeter().recordStreamEvent(eventType);
    },
    onError(span: TelemetrySpan, error: CadreenError): void {
      span.setAttribute("cadreen.error.code", error.code);
      span.setAttribute("cadreen.error.type", error.errorType);
      span.setStatus("error");
      span.end();
    },
  };
}
