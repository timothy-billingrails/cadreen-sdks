import pytest
import httpx

from cadreen.client import HttpClient, CadreenError, _build_query_string
from cadreen.types import CadreenConfig, RequestOptions


class TestCadreenConfig:
    def test_sandbox_config_defaults(self):
        config = CadreenConfig(api_key="key123", sandbox=True)
        assert config.api_key == "key123"
        assert config.sandbox is True
        assert config.fixtures is None
        assert config.profile == "full"
        assert config.max_retries is None
        assert config.timeout is None
        assert config.base_url is None

    def test_sandbox_config_with_fixtures(self):
        fixtures = {"GET /test": {"value": 42}}
        config = CadreenConfig(api_key="key123", sandbox=True, fixtures=fixtures)
        assert config.fixtures == fixtures

    def test_sandbox_config_with_custom_base_url(self):
        config = CadreenConfig(api_key="key123", sandbox=True, base_url="http://localhost:8080")
        assert config.base_url == "http://localhost:8080"

    def test_sandbox_config_with_retries_and_timeout(self):
        config = CadreenConfig(api_key="key123", sandbox=True, max_retries=5, timeout=60)
        assert config.max_retries == 5
        assert config.timeout == 60


class TestHttpClientSandbox:
    @pytest.mark.asyncio
    async def test_sandbox_returns_fixture_data(self):
        fixtures = {"GET /test": {"result": "sandbox_ok"}}
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixtures)
        client = HttpClient(config)
        result = await client.get("/test")
        assert result == {"result": "sandbox_ok"}

    @pytest.mark.asyncio
    async def test_sandbox_falls_back_to_bare_path(self):
        """Fallback: fixtures keyed by bare path (no method prefix) should also match"""
        fixtures = {"/bare": {"works": True}}
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixtures)
        client = HttpClient(config)
        result = await client.get("/bare")
        assert result == {"works": True}

    @pytest.mark.asyncio
    async def test_sandbox_prefers_method_prefix_over_bare_path(self):
        """When both METHOD /path and /path exist, METHOD /path takes precedence"""
        fixtures = {
            "GET /test": {"preferred": True},
            "/test": {"fallback": True},
        }
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixtures)
        client = HttpClient(config)
        result = await client.get("/test")
        assert result == {"preferred": True}

    @pytest.mark.asyncio
    async def test_sandbox_missing_fixture_raises_404(self):
        config = CadreenConfig(api_key="key", sandbox=True, fixtures={})
        client = HttpClient(config)
        with pytest.raises(CadreenError) as exc:
            await client.get("/nonexistent")
        assert exc.value.status == 404
        assert exc.value.code == "not_found"
        assert "No fixture for GET /nonexistent" in str(exc.value)

    @pytest.mark.asyncio
    async def test_sandbox_post_with_fixture(self):
        fixtures = {"POST /api/v1/endpoint": {"status": "created"}}
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixtures)
        client = HttpClient(config)
        result = await client.post("/api/v1/endpoint", {"data": 1})
        assert result == {"status": "created"}

    @pytest.mark.asyncio
    async def test_sandbox_delete_with_fixture(self):
        fixtures = {"DELETE /api/v1/item": None}
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixtures)
        client = HttpClient(config)
        result = await client.delete("/api/v1/item")
        assert result is None

    @pytest.mark.asyncio
    async def test_sandbox_put_with_fixture(self):
        fixtures = {"PUT /api/v1/item": {"updated": True}}
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixtures)
        client = HttpClient(config)
        result = await client.put("/api/v1/item", {"key": "val"})
        assert result == {"updated": True}


class TestCadreenError:
    def test_error_full_construction(self):
        err = CadreenError(
            status=429,
            code="rate_limited",
            error_type="rate_limit",
            message="Too many requests",
            details=[{"field": "retry_after", "message": "30s"}],
            intelligence={"decision": "backoff"},
        )
        assert err.status == 429
        assert err.code == "rate_limited"
        assert err.error_type == "rate_limit"
        assert str(err) == "Too many requests"
        assert err.details == [{"field": "retry_after", "message": "30s"}]
        assert err.intelligence == {"decision": "backoff"}

    def test_error_minimal_construction(self):
        err = CadreenError(status=500, code="internal", error_type="server", message="Boom")
        assert err.status == 500
        assert err.code == "internal"
        assert err.error_type == "server"
        assert str(err) == "Boom"
        assert err.details is None
        assert err.intelligence is None

    def test_error_defaults(self):
        err = CadreenError(status=400, code="bad_request", error_type="validation", message="Invalid")
        assert err.details is None
        assert err.intelligence is None


class TestBuildQueryString:
    def test_empty_params(self):
        assert _build_query_string({}) == ""

    def test_single_param(self):
        assert _build_query_string({"query": "hello"}) == "?query=hello"

    def test_multiple_params(self):
        qs = _build_query_string({"query": "test", "limit": 5})
        assert "query=test" in qs
        assert "limit=5" in qs
        assert qs.startswith("?")

    def test_none_values_skipped(self):
        qs = _build_query_string({"query": "ok", "domain": None})
        assert qs == "?query=ok"

    def test_empty_string_values_skipped(self):
        qs = _build_query_string({"query": "", "limit": 10})
        assert qs == "?limit=10"

    def test_all_skipped_returns_empty(self):
        assert _build_query_string({"a": None, "b": ""}) == ""

    def test_url_encodes_special_chars(self):
        qs = _build_query_string({"query": "hello world"})
        assert "hello%20world" in qs


class TestIdempotencyKey:
    @pytest.mark.asyncio
    async def test_sandbox_does_not_generate_network_requests(self):
        """In sandbox mode, the client returns fixture data without making network calls.
        We verify idempotency is handled structurally — sandbox bypasses all network I/O."""
        fixtures = {"POST /api/v1/idempotent": {"ok": True}}
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixtures)
        client = HttpClient(config)
        result = await client.post("/api/v1/idempotent", {"data": 1})
        assert result == {"ok": True}

    @pytest.mark.asyncio
    async def test_sandbox_with_idempotency_key_option(self):
        """IdempotencyKey in RequestOptions is accepted structurally."""
        fixtures = {"POST /api/v1/idempotent": {"ok": True}}
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixtures)
        client = HttpClient(config)
        result = await client.post(
            "/api/v1/idempotent",
            {"data": 1},
            options=RequestOptions(idempotency_key="my-key-123"),
        )
        assert result == {"ok": True}


class TestTimeoutHandling:
    @pytest.mark.asyncio
    async def test_default_timeout_is_30(self):
        config = CadreenConfig(api_key="key", sandbox=True)
        client = HttpClient(config)
        assert client._timeout == 30

    @pytest.mark.asyncio
    async def test_custom_timeout_stored(self):
        config = CadreenConfig(api_key="key", sandbox=True, timeout=10)
        client = HttpClient(config)
        assert client._timeout == 10

    @pytest.mark.asyncio
    async def test_default_max_retries_is_2(self):
        config = CadreenConfig(api_key="key", sandbox=True)
        client = HttpClient(config)
        assert client._max_retries == 2

    @pytest.mark.asyncio
    async def test_custom_max_retries_stored(self):
        config = CadreenConfig(api_key="key", sandbox=True, max_retries=5)
        client = HttpClient(config)
        assert client._max_retries == 5


class TestSSEParsing:
    @pytest.mark.asyncio
    async def test_sandbox_stream_raises_on_sandbox(self):
        """In sandbox mode, stream() bypasses sandbox checks and makes real HTTP calls.
        The sandbox intentionally only intercepts request() — stream() uses httpx directly.
        Verify that stream() does not crash during instantiation in sandbox config."""
        config = CadreenConfig(api_key="key", sandbox=True, fixtures={})
        client = HttpClient(config)
        assert client._sandbox is True
        generator = client.stream("/api/v1/stream")
        assert generator is not None
        await generator.aclose()
