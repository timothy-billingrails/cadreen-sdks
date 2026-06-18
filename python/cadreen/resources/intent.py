from __future__ import annotations

from typing import Any

from ..client import HttpClient
from ..types import (
    DirectResult,
    ClarifyResult,
    ExecutionResult,
    BlockedResult,
    ConnectRequiredResult,
    IntentResult,
    IntentRequest,
    IntentMessage,
    IntentContext,
    IntelligenceMeta,
    CapabilityTrace,
    ReasoningTrace,
    MemoryTrace,
    GovernanceTrace,
    HumilityTrace,
    ProcessTrace,
    FieldStability,
    ResponseMessage,
    ClarificationQuestion,
)


def _default_intelligence() -> IntelligenceMeta:
    return IntelligenceMeta(
        capability=CapabilityTrace(total_available=0, healthy_count=0),
        reasoning=ReasoningTrace(),
        memory=MemoryTrace(healthy=True),
        governance=GovernanceTrace(active=False),
        humility=HumilityTrace(),
        process=ProcessTrace(started_at="", duration_ms=0),
        field_stability=FieldStability(stable=[], evolving=[], internal=[]),
    )


def _map_intent_response(raw: dict[str, Any]) -> IntentResult:
    intelligence = _default_intelligence()
    raw_meta = raw.get("intelligence")
    if raw_meta:
        intelligence = _parse_intelligence(raw_meta)

    trace_id = raw.get("trace_id", raw.get("id", ""))

    result_type = raw.get("type", "direct")

    if result_type == "direct":
        msg_raw = raw.get("message", {})
        message = ResponseMessage(role=msg_raw.get("role", "assistant"), content=msg_raw.get("content", ""))
        return DirectResult(type="direct", message=message, intelligence=intelligence, trace_id=trace_id)

    if result_type == "clarify":
        clar = raw.get("clarification", {})
        raw_questions = clar.get("questions", [])
        questions = []
        for q in raw_questions:
            if isinstance(q, str):
                questions.append(ClarificationQuestion(id="", question=q, type="open", required=False))
            else:
                questions.append(ClarificationQuestion(
                    id=q.get("id", ""),
                    question=q.get("question", ""),
                    type=q.get("type", "open"),
                    required=q.get("required", False),
                    reason=q.get("reason"),
                ))
        return ClarifyResult(
            type="clarify",
            questions=questions,
            conversation_id=clar.get("conversation_id", ""),
            intelligence=intelligence,
            trace_id=trace_id,
        )

    if result_type == "mission":
        mission = raw.get("mission", {})
        return ExecutionResult(
            type="execution",
            execution={
                "id": mission.get("id", ""),
                "status": mission.get("status", ""),
                "stream_url": mission.get("stream_url"),
                "poll_url": mission.get("poll_url"),
            },
            intelligence=intelligence,
            trace_id=trace_id,
        )

    if result_type == "execution":
        execution = raw.get("execution", {})
        mission = raw.get("mission", {})
        return ExecutionResult(
            type="execution",
            execution={
                "id": execution.get("id") or mission.get("id", ""),
                "status": execution.get("status") or mission.get("status", ""),
                "stream_url": execution.get("stream_url") or mission.get("stream_url"),
                "poll_url": execution.get("poll_url") or mission.get("poll_url"),
            },
            intelligence=intelligence,
            trace_id=trace_id,
        )

    if result_type == "blocked":
        gov = raw.get("meta", {}).get("governance", {})
        return BlockedResult(
            type="blocked",
            reason_code=gov.get("decision"),
            policy_id=gov.get("reason"),
            intelligence=intelligence,
            trace_id=trace_id,
        )

    if result_type == "connect_required":
        mission = raw.get("mission", {})
        gov = raw.get("meta", {}).get("governance", {})
        return ConnectRequiredResult(
            type="connect_required",
            endpoint=mission.get("stream_url", ""),
            reason=gov.get("reason", "connection required"),
            intelligence=intelligence,
            trace_id=trace_id,
        )

    msg_raw = raw.get("message", {})
    message = ResponseMessage(role=msg_raw.get("role", "assistant"), content=msg_raw.get("content", ""))
    return DirectResult(type="direct", message=message, intelligence=intelligence, trace_id=trace_id)


def _parse_intelligence(raw: dict[str, Any]) -> IntelligenceMeta:
    cap = raw.get("capability", {})
    fs_raw = raw.get("field_stability", {})
    field_stability = FieldStability(
        stable=fs_raw.get("stable", []),
        evolving=fs_raw.get("evolving", []),
        internal=fs_raw.get("internal", []),
    )
    return IntelligenceMeta(
        version=raw.get("version"),
        summary=raw.get("summary"),
        capability=CapabilityTrace(
            total_available=cap.get("total_available", 0),
            healthy_count=cap.get("healthy_count", 0),
            active_integrations=cap.get("active_integrations"),
        ),
        reasoning=ReasoningTrace(capability_matches=raw.get("reasoning", {}).get("capability_matches")),
        memory=MemoryTrace(
            healthy=raw.get("memory", {}).get("healthy", True),
            knowledge_queried=raw.get("memory", {}).get("knowledge_queried"),
        ),
        governance=GovernanceTrace(
            active=raw.get("governance", {}).get("active", False),
            decision=raw.get("governance", {}).get("decision"),
            confidence=raw.get("governance", {}).get("confidence"),
            reason_code=raw.get("governance", {}).get("reason_code"),
            policy_id=raw.get("governance", {}).get("policy_id"),
            next_actions=raw.get("governance", {}).get("next_actions"),
        ),
        humility=HumilityTrace(
            gaps_detected=raw.get("humility", {}).get("gaps_detected"),
            blocking=raw.get("humility", {}).get("blocking"),
        ),
        process=ProcessTrace(
            started_at=raw.get("process", {}).get("started_at", ""),
            duration_ms=raw.get("process", {}).get("duration_ms", 0),
            components=raw.get("process", {}).get("components"),
        ),
        field_stability=field_stability,
    )


class IntentResource:
    def __init__(self, client: HttpClient) -> None:
        self._client = client

    async def __call__(
        self,
        prompt: str,
        *,
        domain: str | None = None,
        mode: str | None = None,
        conversation_id: str | None = None,
        context: dict[str, Any] | None = None,
        stream: bool | None = None,
    ) -> IntentResult:
        messages = [IntentMessage(role="user", content=prompt)]
        ctx = None
        if domain or context:
            ctx = IntentContext(domain=domain, constraints=context)
        request = IntentRequest(
            messages=messages,
            conversation_id=conversation_id,
            context=ctx,
            mode=mode,
            stream=stream,
        )
        return await self.invoke(request)

    async def invoke(self, request: IntentRequest) -> IntentResult:
        body: dict[str, Any] = {
            "messages": [{"role": m.role, "content": m.content} for m in request.messages],
        }
        if request.conversation_id:
            body["conversation_id"] = request.conversation_id
        if request.context:
            ctx: dict[str, Any] = {}
            if request.context.existing_connectors is not None:
                ctx["existing_connectors"] = request.context.existing_connectors
            if request.context.constraints is not None:
                ctx["constraints"] = request.context.constraints
            if request.context.domain is not None:
                ctx["domain"] = request.context.domain
            body["context"] = ctx
        if request.mode:
            body["mode"] = request.mode
        if request.stream is not None:
            body["stream"] = request.stream

        raw = await self._client.post("/api/v1/cadreen/intent", body)
        return _map_intent_response(raw)

    async def invoke_stream(
        self,
        request: IntentRequest,
    ) -> AsyncIterator[dict[str, Any]]:
        """Stream intent processing via SSE. Yields events as they arrive."""
        import json as _json

        body: dict[str, Any] = {
            "messages": [{"role": m.role, "content": m.content} for m in request.messages],
            "stream": True,
        }
        if request.conversation_id:
            body["conversation_id"] = request.conversation_id
        if request.context:
            ctx: dict[str, Any] = {}
            if request.context.existing_connectors is not None:
                ctx["existing_connectors"] = request.context.existing_connectors
            if request.context.constraints is not None:
                ctx["constraints"] = request.context.constraints
            if request.context.domain is not None:
                ctx["domain"] = request.context.domain
            body["context"] = ctx
        if request.mode:
            body["mode"] = request.mode

        url = f"{self._client._base_url}/api/v1/cadreen/intent"
        async with self._client._session.post(
            url,
            json=body,
            headers={
                "Authorization": f"Bearer {self._client._api_key}",
                "Accept": "text/event-stream",
            },
        ) as resp:
            resp.raise_for_status()
            current_event = "message"
            async for line in resp.content:
                line_str = line.decode("utf-8").rstrip("\r\n")
                if line_str.startswith("event:"):
                    current_event = line_str[6:].strip()
                elif line_str.startswith("data:"):
                    data = line_str[5:].strip()
                    if data == "[DONE]":
                        return
                    try:
                        yield {"type": current_event, "data": _json.loads(data)}
                    except _json.JSONDecodeError:
                        yield {"type": current_event, "data": {"raw": data}}
                    current_event = "message"
