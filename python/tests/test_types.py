import pytest

from cadreen.types import (
    DirectResult,
    ClarifyResult,
    ExecutionResult,
    BlockedResult,
    ConnectRequiredResult,
    ResponseMessage,
    ClarificationQuestion,
    IntelligenceMeta,
    CapabilityTrace,
    ReasoningTrace,
    MemoryTrace,
    GovernanceTrace,
    HumilityTrace,
    ProcessTrace,
    FieldStability,
    IntentMessage,
    IntentRequest,
    IntentContext,
    intent_status,
    IntelligenceStage,
    NextAction,
)


def _default_intelligence(**overrides):
    kwargs = {
        "capability": CapabilityTrace(total_available=0, healthy_count=0),
        "reasoning": ReasoningTrace(),
        "memory": MemoryTrace(healthy=True),
        "governance": GovernanceTrace(active=False),
        "humility": HumilityTrace(),
        "process": ProcessTrace(started_at="", duration_ms=0),
        "field_stability": FieldStability(stable=[], evolving=[], internal=[]),
    }
    kwargs.update(overrides)
    return IntelligenceMeta(**kwargs)


class TestIntentStatus:
    def test_direct_result_status(self):
        msg = ResponseMessage(role="assistant", content="Hello")
        intel = _default_intelligence()
        result = DirectResult(type="direct", message=msg, intelligence=intel, trace_id="t1")
        status = intent_status(result)
        assert status == {"ready": True, "needs": [], "next": "done"}

    def test_clarify_result_status(self):
        q1 = ClarificationQuestion(id="q1", question="Budget?", type="text", required=True)
        intel = _default_intelligence()
        result = ClarifyResult(
            type="clarify", questions=[q1], conversation_id="conv1",
            intelligence=intel, trace_id="t1",
        )
        status = intent_status(result)
        assert status == {"ready": False, "needs": ["Budget?"], "next": "answer 1 question"}

    def test_execution_result_status(self):
        intel = _default_intelligence()
        result = ExecutionResult(
            type="execution",
            execution={"id": "mis1", "status": "running", "poll_url": "/poll/mis1"},
            intelligence=intel, trace_id="t1",
        )
        status = intent_status(result)
        assert status == {"ready": True, "needs": [], "next": "poll /poll/mis1"}

    def test_execution_result_with_stream_url(self):
        intel = _default_intelligence()
        result = ExecutionResult(
            type="execution",
            execution={"id": "mis1", "status": "running", "stream_url": "/stream/mis1"},
            intelligence=intel, trace_id="t1",
        )
        status = intent_status(result)
        assert status["next"] == "stream execution"

    def test_execution_result_without_urls(self):
        intel = _default_intelligence()
        result = ExecutionResult(
            type="execution",
            execution={"id": "mis1"},
            intelligence=intel, trace_id="t1",
        )
        status = intent_status(result)
        assert status["next"] == "poll mis1"

    def test_blocked_result_status(self):
        intel = _default_intelligence()
        result = BlockedResult(
            type="blocked", intelligence=intel, trace_id="t1",
            reason_code="policy_violation", policy_id="pol_123",
        )
        status = intent_status(result)
        assert status == {"ready": False, "needs": ["blocked: policy_violation"], "next": "pol_123"}

    def test_blocked_result_without_policy_id(self):
        intel = _default_intelligence()
        result = BlockedResult(
            type="blocked", intelligence=intel, trace_id="t1",
            reason_code="governance_gate",
        )
        status = intent_status(result)
        assert status["next"] == "resolve policy block"

    def test_connect_required_result_status(self):
        intel = _default_intelligence()
        result = ConnectRequiredResult(
            type="connect_required", intelligence=intel, trace_id="t1",
            endpoint="/connect/stripe",
        )
        status = intent_status(result)
        assert status == {"ready": False, "needs": ["connect /connect/stripe"], "next": "/connect/stripe"}

    def test_connect_required_empty_endpoint(self):
        intel = _default_intelligence()
        result = ConnectRequiredResult(
            type="connect_required", intelligence=intel, trace_id="t1",
        )
        status = intent_status(result)
        assert status["needs"] == []
        assert status["next"] == ""


class TestResultProperties:
    def test_direct_result_explain(self):
        msg = ResponseMessage(role="assistant", content="Here is the answer.")
        intel = _default_intelligence()
        result = DirectResult(type="direct", message=msg, intelligence=intel, trace_id="t1")
        assert result.explain() == "Here is the answer."
        assert result.ready is True
        assert result.needs == []
        assert result.next == "done"

    def test_direct_result_next_with_action(self):
        msg = ResponseMessage(role="assistant", content="Done.")
        next_action = NextAction(type="review", label="Review changes", reason="needs approval", endpoint="/review")
        intel = _default_intelligence(next_action=next_action)
        result = DirectResult(type="direct", message=msg, intelligence=intel, trace_id="t1")
        assert result.next == "Review changes"

    def test_clarify_result_explain(self):
        q1 = ClarificationQuestion(id="q1", question="What is the budget?", type="text", required=True)
        q2 = ClarificationQuestion(id="q2", question="Which region?", type="choice", required=False)
        intel = _default_intelligence()
        result = ClarifyResult(
            type="clarify", questions=[q1, q2], conversation_id="conv1",
            intelligence=intel, trace_id="t1",
        )
        assert "Clarification needed:" in result.explain()
        assert "What is the budget?" in result.explain()
        assert "Which region?" in result.explain()
        assert result.ready is False
        assert result.needs == ["What is the budget?", "Which region?"]
        assert result.next == "answer 2 questions"

    def test_clarify_result_single_question(self):
        q1 = ClarificationQuestion(id="q1", question="Your name?", type="text", required=True)
        intel = _default_intelligence()
        result = ClarifyResult(
            type="clarify", questions=[q1], conversation_id="conv1",
            intelligence=intel, trace_id="t1",
        )
        assert result.next == "answer 1 question"

    def test_execution_result_explain(self):
        intel = _default_intelligence()
        result = ExecutionResult(
            type="execution",
            execution={"id": "mis_abc", "status": "running"},
            intelligence=intel, trace_id="t1",
        )
        assert "Execution started: mis_abc" in result.explain()
        assert result.ready is True
        assert result.needs == []

    def test_blocked_result_explain_with_status(self):
        intel = _default_intelligence()
        result = BlockedResult(
            type="blocked", intelligence=intel, trace_id="t1",
            reason_code="rate_limit", status="Too many requests — try again in 30s",
        )
        assert result.explain() == "Too many requests — try again in 30s"

    def test_blocked_result_explain_without_status(self):
        intel = _default_intelligence()
        result = BlockedResult(
            type="blocked", intelligence=intel, trace_id="t1",
            reason_code="policy_block",
        )
        assert "Blocked by policy: policy_block" in result.explain()
        assert result.ready is False
        assert result.needs == ["blocked: policy_block"]

    def test_blocked_result_without_reason_code(self):
        intel = _default_intelligence()
        result = BlockedResult(type="blocked", intelligence=intel, trace_id="t1")
        assert "governance gate" in result.explain()

    def test_connect_required_result_explain(self):
        intel = _default_intelligence()
        result = ConnectRequiredResult(
            type="connect_required", intelligence=intel, trace_id="t1",
            endpoint="/connect/stripe", status="Stripe connection needed",
        )
        assert result.explain() == "Stripe connection needed"

    def test_connect_required_result_explain_no_status(self):
        intel = _default_intelligence()
        result = ConnectRequiredResult(
            type="connect_required", intelligence=intel, trace_id="t1",
            endpoint="/connect/stripe",
        )
        assert result.explain() == "Connection required: /connect/stripe"


class TestIntelligenceMeta:
    def test_requires_human_false_by_default(self):
        intel = _default_intelligence()
        assert intel.requires_human() is False

    def test_requires_human_with_blocked_decision_and_blocking_gaps(self):
        governance = GovernanceTrace(active=True, decision="blocked")
        humility = HumilityTrace(blocking=3)
        intel = _default_intelligence(governance=governance, humility=humility)
        assert intel.requires_human() is True

    def test_requires_human_with_escalation_stage(self):
        stage = IntelligenceStage(name="escalation", status="active", detail="Needs approval")
        intel = _default_intelligence(stages=[stage])
        assert intel.requires_human() is True

    def test_requires_human_with_human_handoff_stage(self):
        stage = IntelligenceStage(name="human_handoff", status="pending")
        intel = _default_intelligence(stages=[stage])
        assert intel.requires_human() is True

    def test_requires_human_blocked_no_gaps(self):
        governance = GovernanceTrace(active=True, decision="blocked")
        humility = HumilityTrace(blocking=0)
        intel = _default_intelligence(governance=governance, humility=humility)
        assert intel.requires_human() is False

    def test_requires_human_skipped_stages(self):
        stage = IntelligenceStage(name="escalation", status="skipped")
        intel = _default_intelligence(stages=[stage])
        assert intel.requires_human() is False

    def test_handoff_reason_from_stage(self):
        stage = IntelligenceStage(name="escalation", status="active", detail="User requested review")
        intel = _default_intelligence(stages=[stage])
        assert intel.handoff_reason() == "User requested review"

    def test_handoff_reason_from_next_action(self):
        next_action = NextAction(type="review", label="Review", reason="Budget exceeds limit")
        intel = _default_intelligence(next_action=next_action)
        assert intel.handoff_reason() == "Budget exceeds limit"

    def test_handoff_reason_none_when_no_sources(self):
        intel = _default_intelligence()
        assert intel.handoff_reason() is None

    def test_handoff_reason_skipped_stage_is_ignored(self):
        stage = IntelligenceStage(name="escalation", status="skipped", detail="Nothing")
        intel = _default_intelligence(stages=[stage])
        assert intel.handoff_reason() is None

    def test_explain_basic(self):
        intel = _default_intelligence()
        explanation = intel.explain()
        assert "Intelligence:" in explanation.summary
        assert isinstance(explanation.steps, list)
        assert explanation.recommendations is None

    def test_explain_with_capabilities(self):
        capability = CapabilityTrace(total_available=5, healthy_count=3)
        intel = _default_intelligence(capability=capability)
        explanation = intel.explain()
        assert "3/5 capabilities healthy" in explanation.summary

    def test_explain_with_governance(self):
        governance = GovernanceTrace(active=True, decision="blocked", confidence=0.85)
        intel = _default_intelligence(governance=governance)
        explanation = intel.explain()
        assert "Governance:" in explanation.summary
        assert "blocked" in explanation.summary
        assert "confidence: 0.85" in explanation.summary

    def test_explain_with_humility_gaps(self):
        humility = HumilityTrace(gaps_detected=4, blocking=2)
        intel = _default_intelligence(humility=humility)
        explanation = intel.explain()
        assert "4 gaps detected (2 blocking)" in explanation.summary

    def test_explain_with_memory(self):
        memory = MemoryTrace(healthy=True, knowledge_queried=12)
        intel = _default_intelligence(memory=memory)
        explanation = intel.explain()
        assert "12 knowledge items queried" in explanation.summary

    def test_explain_with_stages(self):
        stage = IntelligenceStage(
            name="classification", status="ok", duration_ms=45, detail="Intent matched"
        )
        intel = _default_intelligence(stages=[stage])
        explanation = intel.explain()
        stage_line = [s for s in explanation.steps if "classification" in s][0]
        assert "ok" in stage_line
        assert "Intent matched" in stage_line

    def test_explain_with_blocking_recommendations(self):
        capability = CapabilityTrace(total_available=5, healthy_count=3)
        humility = HumilityTrace(gaps_detected=4, blocking=3)
        governance = GovernanceTrace(active=True, decision="blocked")
        memory = MemoryTrace(healthy=False)
        intel = _default_intelligence(
            capability=capability, humility=humility, governance=governance, memory=memory
        )
        explanation = intel.explain()
        assert explanation.recommendations is not None
        recs = explanation.recommendations
        assert any("blocking" in r.lower() for r in recs)
        assert any("degraded" in r.lower() for r in recs)
        assert any("human intervention" in r.lower() for r in recs)

    def test_explain_no_capabilities(self):
        intel = _default_intelligence(capability=CapabilityTrace(total_available=0, healthy_count=0))
        explanation = intel.explain()
        assert "capabilities" not in explanation.summary

    def test_field_stability(self):
        intel = _default_intelligence()
        assert intel.field_stability.stable == []
        assert intel.field_stability.evolving == []
        assert intel.field_stability.internal == []


class TestIntentRequest:
    def test_intent_request_construction(self):
        msg = IntentMessage(role="user", content="Hello")
        req = IntentRequest(messages=[msg], conversation_id="conv1", mode="chat")
        assert len(req.messages) == 1
        assert req.messages[0].role == "user"
        assert req.messages[0].content == "Hello"
        assert req.conversation_id == "conv1"
        assert req.mode == "chat"
        assert req.stream is None
        assert req.context is None

    def test_intent_request_with_context(self):
        ctx = IntentContext(domain="finance", constraints={"budget": 1000})
        msg = IntentMessage(role="user", content="Analyze")
        req = IntentRequest(messages=[msg], context=ctx, stream=False)
        assert req.context.domain == "finance"
        assert req.context.constraints == {"budget": 1000}
        assert req.stream is False
