from __future__ import annotations

import time
import uuid
from typing import Any, Optional
from urllib.parse import urlencode, quote

import httpx
from httpx_sse import aconnect_sse

from .types import CadreenConfig, RequestOptions
from .telemetry import TelemetryHooks, NoOpProvider


class CadreenError(Exception):
    status: int
    code: str
    error_type: str
    details: Optional[list[dict[str, str]]]
    intelligence: Optional[Any]

    def __init__(
        self,
        status: int,
        code: str,
        error_type: str,
        message: str,
        details: Optional[list[dict[str, str]]] = None,
        intelligence: Optional[Any] = None,
    ) -> None:
        super().__init__(message)
        self.status = status
        self.code = code
        self.error_type = error_type
        self.details = details
        self.intelligence = intelligence


DEFAULT_BASE_URL = "https://accomplishanything.today"
DEFAULT_MAX_RETRIES = 2
DEFAULT_TIMEOUT = 30
RETRYABLE_STATUS_CODES = {408, 429, 502, 503, 504}
IDEMPOTENT_METHODS = {"GET", "PUT"}


def _build_query_string(params: dict[str, Any]) -> str:
    parts: list[str] = []
    for key, value in params.items():
        if value is not None and value != "":
            parts.append(f"{quote(key)}={quote(str(value))}")
    return f"?{'&'.join(parts)}" if parts else ""


class HttpClient:
    def __init__(self, config: CadreenConfig) -> None:
        self._base_url = (config.base_url or DEFAULT_BASE_URL).rstrip("/")
        self._api_key = config.api_key
        self._max_retries = config.max_retries if config.max_retries is not None else DEFAULT_MAX_RETRIES
        self._timeout = config.timeout if config.timeout is not None else DEFAULT_TIMEOUT
        self._sandbox = getattr(config, "sandbox", False) or False
        self._fixtures = getattr(config, "fixtures", None) or {}
        self._profile = getattr(config, "profile", None) or "full"
        provider = config.telemetry if config.telemetry else NoOpProvider()
        self._telemetry = TelemetryHooks(provider=provider)

    async def request(
        self,
        method: str,
        path: str,
        body: Optional[Any] = None,
        options: Optional[RequestOptions] = None,
    ) -> Any:
        if self._sandbox:
            fixture_key = f"{method} {path}"
            if fixture_key in self._fixtures:
                return self._fixtures[fixture_key]
            if path in self._fixtures:
                return self._fixtures[path]
            raise CadreenError(404, "not_found", "not_found", f"No fixture for {fixture_key}. Provide fixtures via CadreenConfig.fixtures keyed by 'METHOD /path' or '/path'.")

        url = f"{self._base_url}{path}"
        span = self._telemetry.on_request_start(method, path)
        start_time = time.monotonic()
        headers: dict[str, str] = {
            "Content-Type": "application/json",
            "Authorization": f"Bearer {self._api_key}",
            "Accept": f'application/json; profile="{self._profile}"',
            **(options.headers if options and options.headers else {}),
        }

        if method in ("POST", "PUT", "PATCH"):
            headers["Idempotency-Key"] = (
                (options.idempotency_key if options else None) or str(uuid.uuid4())
            )

        is_idempotent = method in IDEMPOTENT_METHODS or "Idempotency-Key" in headers
        max_attempts = (self._max_retries + 1) if is_idempotent else 1

        last_error: Optional[CadreenError] = None

        for attempt in range(max_attempts):
            if attempt > 0:
                import asyncio
                delay = min(1.0 * (2 ** (attempt - 1)), 10.0)
                self._telemetry.on_retry(method, path, attempt)
                await asyncio.sleep(delay)

            try:
                async with httpx.AsyncClient(timeout=self._timeout) as client:
                    response = await client.request(
                        method=method,
                        url=url,
                        headers=headers,
                        json=body if body is not None else None,
                    )

                if not response.is_success:
                    error_body: Optional[dict[str, Any]] = None
                    try:
                        error_body = response.json()
                    except Exception:
                        pass

                    err = CadreenError(
                        status=response.status_code,
                        code=error_body.get("error", {}).get("code", "unknown") if error_body else "unknown",
                        error_type=error_body.get("error", {}).get("type", "error") if error_body else "error",
                        message=error_body.get("error", {}).get("message", response.text) if error_body else response.text,
                        details=error_body.get("error", {}).get("details") if error_body else None,
                        intelligence=error_body.get("intelligence") if error_body else None,
                    )

                    if response.status_code in RETRYABLE_STATUS_CODES and is_idempotent and attempt < max_attempts - 1:
                        last_error = err
                        continue

                    self._telemetry.on_error(span, err.code, err.error_type)
                    raise err

                if response.status_code == 204:
                    duration_ms = (time.monotonic() - start_time) * 1000
                    self._telemetry.on_request_end(span, method, path, 204, duration_ms)
                    return None

                duration_ms = (time.monotonic() - start_time) * 1000
                self._telemetry.on_request_end(span, method, path, response.status_code, duration_ms)
                return response.json()

            except CadreenError:
                raise
            except httpx.TimeoutException:
                if is_idempotent and attempt < max_attempts - 1:
                    last_error = CadreenError(408, "timeout", "timeout", "Request timed out")
                    continue
                raise CadreenError(408, "timeout", "timeout", "Request timed out")
            except Exception as exc:
                if is_idempotent and attempt < max_attempts - 1:
                    last_error = CadreenError(0, "network_error", "network", str(exc))
                    continue
                raise CadreenError(0, "network_error", "network", str(exc))

        raise last_error or CadreenError(0, "network_error", "network", "Request failed after retries")

    async def get(self, path: str, params: Optional[dict[str, Any]] = None) -> Any:
        qs = _build_query_string(params) if params else ""
        return await self.request("GET", f"{path}{qs}")

    async def post(self, path: str, body: Optional[Any] = None, options: Optional[RequestOptions] = None) -> Any:
        return await self.request("POST", path, body, options)

    async def put(self, path: str, body: Optional[Any] = None, options: Optional[RequestOptions] = None) -> Any:
        return await self.request("PUT", path, body, options)

    async def delete(self, path: str) -> Any:
        return await self.request("DELETE", path)

    async def stream(self, path: str) -> Any:
        url = f"{self._base_url}{path}"
        headers = {
            "Authorization": f"Bearer {self._api_key}",
            "Accept": "text/event-stream",
        }

        async with httpx.AsyncClient(timeout=None) as client:
            async with aconnect_sse(client, "GET", url, headers=headers) as event_source:
                async for event in event_source.aiter_sse():
                    try:
                        import json
                        data = json.loads(event.data)
                    except Exception:
                        data = {"raw": event.data}
                    yield {"type": event.event or "message", "data": data}
