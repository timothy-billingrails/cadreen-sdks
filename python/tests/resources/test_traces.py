import pytest

from cadreen.client import HttpClient
from cadreen.types import CadreenConfig
from cadreen.resources.traces import TracesResource
from cadreen.types import (
    IntelligenceTraceEntry,
    ListIntelligenceResponse,
    IntelligenceStats,
    ReplayResult,
    HandoffPacket,
    PromoteResult,
    Pagination,
)


TRACES_FIXTURES = {
    "GET /api/v1/cadreen/intelligence/trace_1": {
        "id": "trace_1",
        "domain": "finance",
        "request_path": "/api/v1/cadreen/intent",
        "request_method": "POST",
        "meta": {
            "version": "2026-06-03",
            "capability": {"total_available": 10, "healthy_count": 8},
            "memory": {"healthy": True},
            "governance": {"active": True, "decision": "auto", "confidence": 0.95},
            "humility": {"gaps_detected": 0},
            "process": {"started_at": "2026-06-17T00:00:00Z", "duration_ms": 120},
        },
        "created_at": "2026-06-17T00:00:00Z",
    },
    "GET /api/v1/cadreen/intelligence?limit=10": {
        "traces": [
            {
                "id": "trace_1",
                "domain": "finance",
                "request_path": "/api/v1/cadreen/intent",
                "request_method": "POST",
                "meta": {"capability": {"total_available": 5}, "memory": {"healthy": True}, "governance": {}, "humility": {}, "process": {"started_at": "", "duration_ms": 0}},
            },
            {
                "id": "trace_2",
                "domain": "sales",
                "request_path": "/api/v1/cadreen/intent",
                "request_method": "POST",
                "meta": {"capability": {"total_available": 3}, "memory": {"healthy": True}, "governance": {}, "humility": {}, "process": {"started_at": "", "duration_ms": 0}},
            },
        ],
        "count": 2,
        "pagination": {"limit": 10, "offset": 0, "has_more": True},
    },
    "GET /api/v1/cadreen/intelligence/stats": {
        "traces_24h": 150,
        "traces_7d": 980,
        "traces_30d": 4200,
        "avg_confidence_by_domain": {"finance": 0.92, "sales": 0.88},
        "gap_detection_rate": 0.15,
        "governance_decisions": {"auto": 3500, "handoff": 500, "blocked": 200},
    },
    "POST /api/v1/cadreen/intelligence/trace_1/replay": {
        "trace_id": "trace_1",
        "mode": "current",
        "domain": "finance",
        "original_gate": "auto",
        "original_confidence": 0.95,
        "current_gate": "auto",
        "current_confidence": 0.97,
        "gate_changed": False,
        "change_summary": "Confidence increased by 0.02 due to new capabilities",
        "current_capability": {"total_available": 12, "healthy_count": 11},
        "current_memory": {"healthy": True, "knowledge_queried": 5},
        "current_gaps": {"detected": 0},
        "replay_note": "Replay completed successfully",
    },
    "GET /api/v1/cadreen/intelligence/trace_1/handoff": {
        "trace_id": "trace_1",
        "domain": "finance",
        "created_at": "2026-06-17T00:00:00Z",
        "governance": {"decision": "handoff", "reason": "Amount exceeds threshold"},
        "what_the_system_knew": {"capabilities": 10, "memory_items": 5},
        "what_the_system_didnt_know": {"approval_required": True},
        "what_happened": {"blocked": True, "escalation_id": "esc_5"},
        "suggested_actions": [{"action": "approve", "endpoint": "/policies/confirm"}],
        "next_action": {"type": "approval", "label": "Approve transfer"},
        "trace_url": "https://app.example.com/traces/trace_1",
    },
    "POST /api/v1/cadreen/intelligence/trace_1/promote": {
        "id": "tool_1",
        "kind": "tool",
        "status": "created",
        "source_trace_id": "trace_1",
        "tool_name": "CalculateTax",
        "tool_sequence": ["validate", "compute", "format"],
    },
}


@pytest.fixture
def traces_client():
    config = CadreenConfig(api_key="key", sandbox=True, fixtures=TRACES_FIXTURES)
    return HttpClient(config)


class TestTracesResource:
    @pytest.mark.asyncio
    async def test_get(self, traces_client):
        resource = TracesResource(traces_client)
        result = await resource.get("trace_1")
        assert isinstance(result, IntelligenceTraceEntry)
        assert result.id == "trace_1"
        assert result.domain == "finance"
        assert result.request_path == "/api/v1/cadreen/intent"
        assert result.request_method == "POST"
        assert result.intelligence.version == "2026-06-03"
        assert result.intelligence.capability.total_available == 10
        assert result.intelligence.governance.confidence == 0.95
        assert result.intelligence.process.duration_ms == 120

    @pytest.mark.asyncio
    async def test_get_intelligence_property(self, traces_client):
        """intelligence property aliases meta"""
        resource = TracesResource(traces_client)
        result = await resource.get("trace_1")
        assert result.intelligence is result.meta

    @pytest.mark.asyncio
    async def test_list(self, traces_client):
        resource = TracesResource(traces_client)
        result = await resource.list(limit=10)
        assert isinstance(result, ListIntelligenceResponse)
        assert result.count == 2
        assert len(result.traces) == 2
        assert result.traces[0].id == "trace_1"
        assert result.traces[0].domain == "finance"
        assert result.traces[1].id == "trace_2"
        assert result.pagination.has_more is True

    @pytest.mark.asyncio
    async def test_list_with_filters(self):
        """List with domain and decision filters"""
        fixture_key = "GET /api/v1/cadreen/intelligence?domain=finance&decision=auto&limit=5"
        fixtures = {
            fixture_key: {
                "traces": [{"id": "t1", "domain": "finance", "request_path": "/api", "request_method": "GET", "meta": {"capability": {}, "memory": {}, "governance": {}, "humility": {}, "process": {"started_at": "", "duration_ms": 0}}}],
                "count": 1,
            }
        }
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixtures)
        client = HttpClient(config)
        resource = TracesResource(client)
        result = await resource.list(domain="finance", decision="auto", limit=5)
        assert result.count == 1
        assert len(result.traces) == 1

    @pytest.mark.asyncio
    async def test_stats(self, traces_client):
        resource = TracesResource(traces_client)
        result = await resource.stats()
        assert isinstance(result, IntelligenceStats)
        assert result.traces_24h == 150
        assert result.traces_7d == 980
        assert result.traces_30d == 4200
        assert result.avg_confidence_by_domain == {"finance": 0.92, "sales": 0.88}
        assert result.gap_detection_rate == 0.15
        assert result.governance_decisions == {"auto": 3500, "handoff": 500, "blocked": 200}

    @pytest.mark.asyncio
    async def test_replay(self, traces_client):
        resource = TracesResource(traces_client)
        result = await resource.replay("trace_1")
        assert isinstance(result, ReplayResult)
        assert result.trace_id == "trace_1"
        assert result.mode == "current"
        assert result.original_gate == "auto"
        assert result.original_confidence == 0.95
        assert result.current_confidence == 0.97
        assert result.gate_changed is False
        assert "Confidence increased" in result.change_summary
        assert result.current_capability["total_available"] == 12

    @pytest.mark.asyncio
    async def test_replay_with_mode(self):
        """Replay with explicit mode parameter"""
        fixtures = {
            "POST /api/v1/cadreen/intelligence/trace_x/replay": {
                "trace_id": "trace_x",
                "mode": "historical",
                "domain": "sales",
                "original_gate": "handoff",
                "original_confidence": 0.5,
                "current_gate": "auto",
                "current_confidence": 0.9,
                "gate_changed": True,
                "change_summary": "Gate changed from handoff to auto",
                "current_capability": {},
                "current_memory": {},
                "current_gaps": {},
                "replay_note": "Historical replay with current state",
            }
        }
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixtures)
        client = HttpClient(config)
        resource = TracesResource(client)
        result = await resource.replay("trace_x", mode="historical")
        assert result.trace_id == "trace_x"
        assert result.mode == "historical"
        assert result.gate_changed is True

    @pytest.mark.asyncio
    async def test_handoff(self, traces_client):
        resource = TracesResource(traces_client)
        result = await resource.handoff("trace_1")
        assert isinstance(result, HandoffPacket)
        assert result.trace_id == "trace_1"
        assert result.domain == "finance"
        assert result.governance["decision"] == "handoff"
        assert result.what_the_system_knew["capabilities"] == 10
        assert len(result.suggested_actions) == 1
        assert result.next_action["type"] == "approval"
        assert result.trace_url == "https://app.example.com/traces/trace_1"

    @pytest.mark.asyncio
    async def test_promote(self, traces_client):
        resource = TracesResource(traces_client)
        result = await resource.promote("trace_1", "tool", name="CalculateTax")
        assert isinstance(result, PromoteResult)
        assert result.id == "tool_1"
        assert result.kind == "tool"
        assert result.status == "created"
        assert result.source_trace_id == "trace_1"
        assert result.tool_name == "CalculateTax"
        assert result.tool_sequence == ["validate", "compute", "format"]
