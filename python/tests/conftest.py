import pytest


@pytest.fixture
def sandbox_config():
    """Base config for sandbox tests — no network calls"""
    return {
        "api_key": "test_key_sandbox",
        "sandbox": True,
        "fixtures": {},
    }


@pytest.fixture
def direct_result_fixture():
    """Fixture for a direct intent response"""
    return {
        "POST /api/v1/cadreen/intent": {
            "type": "direct",
            "message": {"role": "assistant", "content": "Here is your answer."},
            "trace_id": "trace_abc123",
            "status": "completed",
        }
    }


@pytest.fixture
def blocked_result_fixture():
    return {
        "POST /api/v1/cadreen/intent": {
            "type": "blocked",
            "trace_id": "trace_blocked_001",
            "meta": {
                "governance": {
                    "decision": "human_approval_required",
                    "reason": "pol_wallet_lockdown",
                }
            },
        }
    }


@pytest.fixture
def clarify_result_fixture():
    return {
        "POST /api/v1/cadreen/intent": {
            "type": "clarify",
            "trace_id": "trace_clarify_001",
            "clarification": {
                "conversation_id": "conv_abc",
                "questions": [
                    {"id": "q1", "question": "What is the budget?", "type": "text", "required": True, "reason": "Budget constraint missing"},
                    {"id": "q2", "question": "Which region?", "type": "choice", "required": False},
                ],
            },
        }
    }


@pytest.fixture
def mission_result_fixture():
    return {
        "POST /api/v1/cadreen/intent": {
            "type": "mission",
            "trace_id": "trace_mission_001",
            "mission": {
                "id": "mis_abc",
                "status": "running",
                "stream_url": "/stream/mis_abc",
                "poll_url": "/poll/mis_abc",
            },
        }
    }


@pytest.fixture
def connect_required_result_fixture():
    return {
        "POST /api/v1/cadreen/intent": {
            "type": "connect_required",
            "trace_id": "trace_connect_001",
            "mission": {"stream_url": "/connect/stripe"},
            "meta": {"governance": {"reason": "connection required"}},
        }
    }
