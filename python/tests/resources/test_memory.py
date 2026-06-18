import pytest

from cadreen.client import HttpClient
from cadreen.types import CadreenConfig
from cadreen.resources.memory import MemoryResource, _parse_atom, _parse_create_response
from cadreen.types import Atom, AtomContent, CreateMemoryResponse, SearchMemoryResponse, MemoryTypesResponse


class TestParseAtom:
    def test_parse_atom_full(self):
        raw = {
            "id": "atom_1",
            "type": "fact",
            "domain": "science",
            "authority": 10,
            "version": 2,
            "scope": "global",
            "content": {"text": "Water boils at 100C", "source": "physics"},
            "tags": ["physics", "temperature"],
            "created_at": "2026-06-17T00:00:00Z",
        }
        atom = _parse_atom(raw)
        assert atom.id == "atom_1"
        assert atom.type == "fact"
        assert atom.domain == "science"
        assert atom.authority == 10
        assert atom.version == 2
        assert atom.scope == "global"
        assert atom.content.text == "Water boils at 100C"
        assert atom.content.source == "physics"
        assert atom.tags == ["physics", "temperature"]
        assert atom.created_at == "2026-06-17T00:00:00Z"

    def test_parse_atom_minimal(self):
        raw = {"id": "atom_2", "type": "note", "domain": "personal"}
        atom = _parse_atom(raw)
        assert atom.id == "atom_2"
        assert atom.type == "note"
        assert atom.domain == "personal"
        assert atom.authority == 0
        assert atom.version == 0
        assert atom.content is None
        assert atom.tags is None
        assert atom.created_at is None


class TestParseCreateResponse:
    def test_parse_create_full(self):
        raw = {
            "id": "mem_1",
            "type": "preference",
            "domain": "settings",
            "authority": 5,
            "version": 1,
            "scope": "tenant",
            "content": {"text": "Dark mode enabled"},
            "indexed": True,
            "tags": ["ui"],
            "created_at": "2026-06-17T00:00:00Z",
        }
        resp = _parse_create_response(raw)
        assert resp.id == "mem_1"
        assert resp.type == "preference"
        assert resp.domain == "settings"
        assert resp.indexed is True
        assert resp.content.text == "Dark mode enabled"

    def test_parse_create_minimal(self):
        raw = {"id": "mem_2", "type": "note", "domain": "notes"}
        resp = _parse_create_response(raw)
        assert resp.id == "mem_2"
        assert resp.authority == 0
        assert resp.content is None
        assert resp.indexed is None


MEMORY_FIXTURES = {
    "POST /api/v1/cadreen/memory": {
        "id": "mem_abc",
        "type": "fact",
        "domain": "general",
        "authority": 5,
        "version": 1,
        "scope": "tenant",
        "content": {"text": "Remembered information", "source": "user"},
        "indexed": True,
        "tags": ["test"],
        "created_at": "2026-06-17T00:00:00Z",
    },
    "GET /api/v1/cadreen/memory/search?query=test%20query": {
        "results": [
            {"id": "atom_1", "type": "fact", "domain": "general", "content": {"text": "Result 1"}},
            {"id": "atom_2", "type": "reference", "domain": "docs", "content": {"text": "Result 2"}},
        ],
        "count": 2,
    },
    "GET /api/v1/cadreen/memory/search?query=test%20query&limit=1": {
        "results": [{"id": "atom_1", "type": "fact", "domain": "general", "content": {"text": "Single result"}}],
        "count": 1,
    },
    "GET /api/v1/cadreen/memory/mem_xyz": {
        "id": "mem_xyz",
        "type": "preference",
        "domain": "settings",
        "authority": 3,
        "version": 2,
        "scope": "personal",
        "content": {"text": "User prefers dark mode"},
        "tags": ["ui", "preferences"],
    },
    "GET /api/v1/cadreen/memory/types": {
        "type_values": ["fact", "preference", "reference", "episode", "precedent", "note"],
        "kind_values": ["text", "code", "image"],
        "description": "Available memory types in the system",
    },
}


@pytest.fixture
def memory_client():
    config = CadreenConfig(api_key="key", sandbox=True, fixtures=MEMORY_FIXTURES)
    return HttpClient(config)


class TestMemoryResource:
    @pytest.mark.asyncio
    async def test_remember(self, memory_client):
        resource = MemoryResource(memory_client)
        result = await resource.remember("fact", {"text": "Test content"}, domain="general", tags=["test"])
        assert isinstance(result, CreateMemoryResponse)
        assert result.id == "mem_abc"
        assert result.type == "fact"
        assert result.indexed is True
        assert result.content.text == "Remembered information"

    @pytest.mark.asyncio
    async def test_remember_with_all_params(self, memory_client):
        resource = MemoryResource(memory_client)
        result = await resource.remember(
            "fact", {"text": "Content"},
            domain="general",
            scope="tenant",
            authority=5,
            tags=["tag1", "tag2"],
        )
        assert result.id == "mem_abc"
        assert result.scope == "tenant"
        assert result.authority == 5

    @pytest.mark.asyncio
    async def test_search(self, memory_client):
        resource = MemoryResource(memory_client)
        result = await resource.search("test query")
        assert isinstance(result, SearchMemoryResponse)
        assert result.count == 2
        assert len(result.results) == 2
        assert result.results[0].id == "atom_1"
        assert result.results[0].content.text == "Result 1"

    @pytest.mark.asyncio
    async def test_search_with_limit(self, memory_client):
        resource = MemoryResource(memory_client)
        result = await resource.search("test query", limit=1)
        assert isinstance(result, SearchMemoryResponse)
        assert result.count == 1
        assert len(result.results) == 1
        assert result.results[0].id == "atom_1"

    @pytest.mark.asyncio
    async def test_get(self, memory_client):
        resource = MemoryResource(memory_client)
        result = await resource.get("mem_xyz")
        assert isinstance(result, Atom)
        assert result.id == "mem_xyz"
        assert result.type == "preference"
        assert result.domain == "settings"
        assert result.content.text == "User prefers dark mode"
        assert result.tags == ["ui", "preferences"]

    @pytest.mark.asyncio
    async def test_types(self, memory_client):
        resource = MemoryResource(memory_client)
        result = await resource.types()
        assert isinstance(result, MemoryTypesResponse)
        assert result.type_values == ["fact", "preference", "reference", "episode", "precedent", "note"]
        assert result.kind_values == ["text", "code", "image"]
        assert "Available memory types" in result.description
