import pytest

from cadreen.client import HttpClient
from cadreen.types import CadreenConfig
from cadreen.resources.executions import ExecutionsResource
from cadreen.types import ExecutionStatus, ExecutionEvent


EXECUTIONS_FIXTURES = {
    "GET /api/v1/cadreen/executions/exec_1": {
        "id": "exec_1",
        "status": "running",
        "progress": 0.65,
        "result": None,
        "error": None,
    },
    "GET /api/v1/cadreen/executions/exec_complete": {
        "id": "exec_complete",
        "status": "completed",
        "progress": 1.0,
        "result": {"output": "Success", "items_processed": 150},
        "error": None,
    },
    "GET /api/v1/cadreen/executions/exec_failed": {
        "id": "exec_failed",
        "status": "failed",
        "progress": 0.3,
        "result": None,
        "error": "Connection refused: upstream service unreachable",
    },
}


@pytest.fixture
def executions_client():
    config = CadreenConfig(api_key="key", sandbox=True, fixtures=EXECUTIONS_FIXTURES)
    return HttpClient(config)


class TestExecutionsResource:
    @pytest.mark.asyncio
    async def test_get_status_running(self, executions_client):
        resource = ExecutionsResource(executions_client)
        result = await resource.get_status("exec_1")
        assert isinstance(result, ExecutionStatus)
        assert result.id == "exec_1"
        assert result.status == "running"
        assert result.progress == 0.65
        assert result.result is None
        assert result.error is None

    @pytest.mark.asyncio
    async def test_get_status_completed(self, executions_client):
        resource = ExecutionsResource(executions_client)
        result = await resource.get_status("exec_complete")
        assert result.id == "exec_complete"
        assert result.status == "completed"
        assert result.progress == 1.0
        assert result.result == {"output": "Success", "items_processed": 150}

    @pytest.mark.asyncio
    async def test_get_status_failed(self, executions_client):
        resource = ExecutionsResource(executions_client)
        result = await resource.get_status("exec_failed")
        assert result.id == "exec_failed"
        assert result.status == "failed"
        assert result.progress == 0.3
        assert result.result is None
        assert result.error == "Connection refused: upstream service unreachable"

    @pytest.mark.asyncio
    async def test_get_status_minimal(self):
        """Response with only required fields"""
        fixtures = {
            "GET /api/v1/cadreen/executions/exec_min": {
                "id": "exec_min",
                "status": "pending",
            }
        }
        config = CadreenConfig(api_key="key", sandbox=True, fixtures=fixtures)
        client = HttpClient(config)
        resource = ExecutionsResource(client)
        result = await resource.get_status("exec_min")
        assert result.id == "exec_min"
        assert result.status == "pending"
        assert result.progress is None
        assert result.result is None
        assert result.error is None
