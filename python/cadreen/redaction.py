from __future__ import annotations

import re
from dataclasses import replace
from typing import Any, Sequence

_EMAIL_RE = re.compile(r"[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}")
_PHONE_RE = re.compile(r"(?:\+?\d{1,3}[-.\s]?)?\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}")
_UUID_RE = re.compile(r"[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}", re.IGNORECASE)
_API_KEY_RE = re.compile(r"sk_[a-zA-Z]+_[a-zA-Z0-9]{8,}")
_IP_RE = re.compile(r"\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b")


class RedactOptions:
    def __init__(
        self,
        *,
        preserve_uuids: bool = False,
        keys_to_redact: Sequence[str] | None = None,
    ) -> None:
        self.preserve_uuids = preserve_uuids
        self.keys_to_redact = list(keys_to_redact) if keys_to_redact else [
            "content", "message", "text", "body", "email", "phone", "address", "name",
        ]


def redact_string(text: str, options: RedactOptions | None = None) -> str:
    opts = options or RedactOptions()
    result = _EMAIL_RE.sub("[email]", text)
    result = _PHONE_RE.sub("[phone]", result)
    result = _API_KEY_RE.sub("[api_key]", result)
    result = _IP_RE.sub("[ip]", result)
    if not opts.preserve_uuids:
        result = _UUID_RE.sub("[id]", result)
    return result


def redact_value(value: Any, options: RedactOptions | None = None) -> Any:
    opts = options or RedactOptions()
    if isinstance(value, str):
        return redact_string(value, opts)
    if isinstance(value, list):
        return [redact_value(v, opts) for v in value]
    if isinstance(value, dict):
        result: dict[str, Any] = {}
        for key, val in value.items():
            if key.lower() in opts.keys_to_redact and isinstance(val, str):
                result[key] = redact_string(val, opts)
            else:
                result[key] = redact_value(val, opts)
        return result
    if hasattr(value, "__dict__") and hasattr(value, "__dataclass_fields__"):
        changes: dict[str, Any] = {}
        for field_name in value.__dataclass_fields__:
            field_val = getattr(value, field_name)
            if isinstance(field_val, str) and field_name.lower() in opts.keys_to_redact:
                changes[field_name] = redact_string(field_val, opts)
            elif isinstance(field_val, (list, dict)) or hasattr(field_val, "__dataclass_fields__"):
                changes[field_name] = redact_value(field_val, opts)
            elif isinstance(field_val, str):
                changes[field_name] = redact_string(field_val, opts)
        return replace(value, **changes)
    return value


def redact_trace(intel: Any, options: RedactOptions | None = None) -> Any:
    return redact_value(intel, options)


def redact_messages(messages: Sequence[Any], options: RedactOptions | None = None) -> list[Any]:
    return [redact_value(msg, options) for msg in messages]
