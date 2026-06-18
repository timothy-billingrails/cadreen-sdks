from __future__ import annotations

import json
from dataclasses import dataclass, field
from typing import Any, AsyncIterator, Optional

from ..client import HttpClient


# ── Chat Completions Types (OpenAI-compatible) ──


@dataclass
class ChatFunctionCall:
    name: str
    arguments: str  # JSON string


@dataclass
class ChatToolCall:
    id: str
    type: str = "function"
    function: ChatFunctionCall = field(default_factory=lambda: ChatFunctionCall(name="", arguments=""))


@dataclass
class ChatMessage:
    role: str
    content: Optional[str] = None
    name: Optional[str] = None
    tool_call_id: Optional[str] = None  # for "tool" role
    tool_calls: Optional[list[ChatToolCall]] = None  # for "assistant" role


@dataclass
class ChatFunctionDefinition:
    name: str
    description: Optional[str] = None
    parameters: Optional[Any] = None  # JSON Schema


@dataclass
class ChatToolDefinition:
    type: str = "function"
    function: ChatFunctionDefinition = field(default_factory=lambda: ChatFunctionDefinition(name=""))


@dataclass
class ChatCompletionRequest:
    messages: list[ChatMessage]
    model: Optional[str] = None
    stream: bool = False
    tools: Optional[list[ChatToolDefinition]] = None
    context: Optional[dict[str, Any]] = None
    conversation_id: Optional[str] = None


@dataclass
class ChatUsage:
    prompt_tokens: int = 0
    completion_tokens: int = 0
    total_tokens: int = 0


@dataclass
class ChatChoice:
    index: int = 0
    message: ChatMessage = field(default_factory=lambda: ChatMessage(role="assistant"))
    finish_reason: str = "stop"


@dataclass
class ChatCompletionResponse:
    id: str = ""
    object: str = "chat.completion"
    created: int = 0
    model: str = ""
    choices: list[ChatChoice] = field(default_factory=list)
    usage: Optional[ChatUsage] = None


@dataclass
class ChatDelta:
    role: Optional[str] = None
    content: Optional[str] = None
    tool_calls: Optional[list[ChatToolCall]] = None


@dataclass
class ChatChunkChoice:
    index: int = 0
    delta: ChatDelta = field(default_factory=lambda: ChatDelta())
    finish_reason: Optional[str] = None


@dataclass
class ChatCompletionChunk:
    id: str = ""
    object: str = "chat.completion.chunk"
    created: int = 0
    model: str = ""
    choices: list[ChatChunkChoice] = field(default_factory=list)
    usage: Optional[ChatUsage] = None


@dataclass
class ToolEntry:
    type: str = "function"
    function: ChatFunctionDefinition = field(default_factory=lambda: ChatFunctionDefinition(name=""))


@dataclass
class ListToolsResponse:
    object: str = "list"
    data: list[ToolEntry] = field(default_factory=list)


# ── Parsing helpers ──


def _parse_tool_calls(raw: list[dict[str, Any]] | None) -> list[ChatToolCall] | None:
    if not raw:
        return None
    return [
        ChatToolCall(
            id=tc.get("id", ""),
            type=tc.get("type", "function"),
            function=ChatFunctionCall(
                name=tc.get("function", {}).get("name", ""),
                arguments=tc.get("function", {}).get("arguments", ""),
            ),
        )
        for tc in raw
    ]


def _parse_chat_message(raw: dict[str, Any]) -> ChatMessage:
    return ChatMessage(
        role=raw.get("role", "assistant"),
        content=raw.get("content"),
        name=raw.get("name"),
        tool_call_id=raw.get("tool_call_id"),
        tool_calls=_parse_tool_calls(raw.get("tool_calls")),
    )


def _parse_chat_choice(raw: dict[str, Any]) -> ChatChoice:
    return ChatChoice(
        index=raw.get("index", 0),
        message=_parse_chat_message(raw.get("message", {})),
        finish_reason=raw.get("finish_reason", "stop"),
    )


def _parse_usage(raw: dict[str, Any] | None) -> ChatUsage | None:
    if not raw:
        return None
    return ChatUsage(
        prompt_tokens=raw.get("prompt_tokens", 0),
        completion_tokens=raw.get("completion_tokens", 0),
        total_tokens=raw.get("total_tokens", 0),
    )


def _parse_chat_response(raw: dict[str, Any]) -> ChatCompletionResponse:
    return ChatCompletionResponse(
        id=raw.get("id", ""),
        object=raw.get("object", "chat.completion"),
        created=raw.get("created", 0),
        model=raw.get("model", ""),
        choices=[_parse_chat_choice(c) for c in raw.get("choices", [])],
        usage=_parse_usage(raw.get("usage")),
    )


def _parse_chunk_choice(raw: dict[str, Any]) -> ChatChunkChoice:
    delta_raw = raw.get("delta", {})
    delta = ChatDelta(
        role=delta_raw.get("role"),
        content=delta_raw.get("content"),
        tool_calls=_parse_tool_calls(delta_raw.get("tool_calls")),
    )
    return ChatChunkChoice(
        index=raw.get("index", 0),
        delta=delta,
        finish_reason=raw.get("finish_reason"),
    )


def _parse_chat_chunk(raw: dict[str, Any]) -> ChatCompletionChunk:
    return ChatCompletionChunk(
        id=raw.get("id", ""),
        object=raw.get("object", "chat.completion.chunk"),
        created=raw.get("created", 0),
        model=raw.get("model", ""),
        choices=[_parse_chunk_choice(c) for c in raw.get("choices", [])],
        usage=_parse_usage(raw.get("usage")),
    )


def _parse_tool_entry(raw: dict[str, Any]) -> ToolEntry:
    fn = raw.get("function", {})
    return ToolEntry(
        type=raw.get("type", "function"),
        function=ChatFunctionDefinition(
            name=fn.get("name", ""),
            description=fn.get("description"),
            parameters=fn.get("parameters"),
        ),
    )


# ── Serialization helpers ──


def _message_to_dict(msg: ChatMessage) -> dict[str, Any]:
    d: dict[str, Any] = {"role": msg.role}
    if msg.content is not None:
        d["content"] = msg.content
    if msg.name is not None:
        d["name"] = msg.name
    if msg.tool_call_id is not None:
        d["tool_call_id"] = msg.tool_call_id
    if msg.tool_calls:
        d["tool_calls"] = [
            {
                "id": tc.id,
                "type": tc.type,
                "function": {"name": tc.function.name, "arguments": tc.function.arguments},
            }
            for tc in msg.tool_calls
        ]
    return d


def _tool_def_to_dict(td: ChatToolDefinition) -> dict[str, Any]:
    return {
        "type": td.type,
        "function": {
            "name": td.function.name,
            "description": td.function.description,
            "parameters": td.function.parameters,
        },
    }


# ── Chat Resource ──


class ChatResource:
    def __init__(self, client: HttpClient) -> None:
        self._client = client

    async def completions(
        self,
        request: ChatCompletionRequest,
    ) -> ChatCompletionResponse:
        body: dict[str, Any] = {
            "messages": [_message_to_dict(m) for m in request.messages],
            "stream": False,
        }
        if request.model:
            body["model"] = request.model
        if request.tools:
            body["tools"] = [_tool_def_to_dict(t) for t in request.tools]
        if request.context:
            body["context"] = request.context
        if request.conversation_id:
            body["conversation_id"] = request.conversation_id

        raw = await self._client.post("/v1/chat/completions", body)
        return _parse_chat_response(raw)

    async def completions_stream(
        self,
        request: ChatCompletionRequest,
    ) -> AsyncIterator[ChatCompletionChunk]:
        body: dict[str, Any] = {
            "messages": [_message_to_dict(m) for m in request.messages],
            "stream": True,
        }
        if request.model:
            body["model"] = request.model
        if request.tools:
            body["tools"] = [_tool_def_to_dict(t) for t in request.tools]
        if request.context:
            body["context"] = request.context
        if request.conversation_id:
            body["conversation_id"] = request.conversation_id

        url = f"{self._client._base_url}/v1/chat/completions"
        headers = {
            "Content-Type": "application/json",
            "Authorization": f"Bearer {self._client._api_key}",
            "Accept": "text/event-stream",
        }

        import httpx
        from httpx_sse import aconnect_sse

        async with httpx.AsyncClient(timeout=None) as client:
            async with aconnect_sse(client, "POST", url, headers=headers, json=body) as event_source:
                async for event in event_source.aiter_sse():
                    if event.data == "[DONE]":
                        return
                    try:
                        data = json.loads(event.data)
                        yield _parse_chat_chunk(data)
                    except Exception:
                        continue

    async def list_tools(self) -> ListToolsResponse:
        raw = await self._client.get("/v1/tools")
        return ListToolsResponse(
            object=raw.get("object", "list"),
            data=[_parse_tool_entry(t) for t in raw.get("data", [])],
        )
