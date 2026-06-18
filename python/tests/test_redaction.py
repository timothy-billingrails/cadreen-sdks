import pytest
from dataclasses import dataclass

from cadreen.redaction import (
    redact_string,
    redact_value,
    redact_trace,
    redact_messages,
    RedactOptions,
)


class TestRedactString:
    def test_redact_email(self):
        result = redact_string("Contact user@example.com for help")
        assert "user@example.com" not in result
        assert "[email]" in result

    def test_redact_multiple_emails(self):
        result = redact_string("Send to alice@foo.com and bob@bar.org")
        assert "alice@foo.com" not in result
        assert "bob@bar.org" not in result
        assert result.count("[email]") == 2

    def test_redact_phone_us_format(self):
        result = redact_string("Call 555-123-4567 today")
        assert "555-123-4567" not in result
        assert "[phone]" in result

    def test_redact_phone_parens_format(self):
        result = redact_string("Call (555) 123-4567 or (555)123-4567")
        assert "(555) 123-4567" not in result
        assert "[phone]" in result

    def test_redact_phone_with_country_code(self):
        result = redact_string("Call +1-555-123-4567")
        assert "+1-555-123-4567" not in result
        assert "[phone]" in result

    def test_redact_api_key(self):
        result = redact_string("API key: sk_live_abc12345def")
        assert "sk_live_abc12345def" not in result
        assert "[api_key]" in result

    def test_redact_ip_address(self):
        result = redact_string("Request from 192.168.1.100")
        assert "192.168.1.100" not in result
        assert "[ip]" in result

    def test_redact_uuid_by_default(self):
        result = redact_string("Trace: aabbccdd-1111-2222-3333-abcdefabcdef")
        assert "aabbccdd-1111-2222-3333-abcdefabcdef" not in result
        assert "[id]" in result

    def test_preserve_uuids(self):
        opts = RedactOptions(preserve_uuids=True)
        result = redact_string("Trace: aabbccdd-1111-2222-3333-abcdefabcdef", opts)
        assert "aabbccdd-1111-2222-3333-abcdefabcdef" in result
        assert "[id]" not in result

    def test_default_options(self):
        """redact_string with no options uses defaults"""
        result = redact_string("user@test.com has phone 555-000-1111 from 10.0.0.1")
        assert "[email]" in result
        assert "[phone]" in result
        assert "[ip]" in result

    def test_no_sensitive_data_unchanged(self):
        text = "This is a normal sentence with no PII."
        assert redact_string(text) == text


class TestRedactValue:
    def test_redact_string_value(self):
        result = redact_value("Email: user@example.com")
        assert "[email]" in result

    def test_redact_list(self):
        data = ["user@foo.com", "Call 555-123-4567", "Normal text"]
        result = redact_value(data)
        assert "[email]" in result[0]
        assert "[phone]" in result[1]
        assert result[2] == "Normal text"

    def test_redact_nested_dict(self):
        data = {
            "user": "alice@example.com",
            "details": {
                "name": "Alice",
                "phone": "555-000-1111",
            },
            "tags": ["public", "alice@example.com"],
        }
        result = redact_value(data)
        assert "[email]" in result["user"]
        assert "[phone]" in result["details"]["phone"]
        assert "[email]" in result["tags"][1]

    def test_redact_dict_keys_not_redacted(self):
        """Values under keys NOT in keys_to_redact are only redacted if they match patterns."""
        data = {"id": "12345", "notes": "contact user@test.com"}
        result = redact_value(data)
        assert "12345" in result["id"]
        assert "[email]" in result["notes"]

    def test_redact_dataclass_values(self):
        @dataclass
        class User:
            name: str
            email: str
            phone: str

        user = User(name="Alice", email="alice@test.com", phone="555-000-1111")
        result = redact_value(user)
        assert "[email]" in result.email
        assert "[phone]" in result.phone

    def test_redact_dataclass_with_dict_fields(self):
        @dataclass
        class Payload:
            id: str
            data: dict

        payload = Payload(id="123", data={"message": "Call 555-000-1111 for help"})
        result = redact_value(payload)
        assert result.id == "123"
        assert "[phone]" in result.data["message"]

    def test_redact_dataclass_all_strings_redacted(self):
        @dataclass
        class Record:
            trace_id: str
            content: str
            priority: int

        record = Record(
            trace_id="aabbccdd-1111-2222-3333-abcdefabcdef",
            content="The API key is sk_test_abcdef123456",
            priority=10,
        )
        result = redact_value(record)
        assert "[id]" in result.trace_id
        assert "[api_key]" in result.content
        assert result.priority == 10

    def test_redact_non_special_types_passthrough(self):
        assert redact_value(42) == 42
        assert redact_value(None) is None
        assert redact_value(True) is True
        assert redact_value(3.14) == 3.14

    def test_redact_empty_list(self):
        assert redact_value([]) == []

    def test_redact_empty_dict(self):
        assert redact_value({}) == {}


class TestRedactTrace:
    def test_redact_trace_delegates_to_redact_value(self):
        trace = {"email": "user@test.com", "message": "Hello from 10.0.0.1"}
        result = redact_trace(trace)
        assert "[email]" in result["email"]
        assert "[ip]" in result["message"]


class TestRedactMessages:
    def test_redact_messages_list(self):
        messages = [
            {"role": "user", "content": "My email is alice@test.com"},
            {"role": "assistant", "content": "I see your IP is 192.168.1.1"},
        ]
        result = redact_messages(messages)
        assert "[email]" in result[0]["content"]
        assert "[ip]" in result[1]["content"]

    def test_redact_messages_empty(self):
        assert redact_messages([]) == []

    def test_redact_messages_with_options(self):
        opts = RedactOptions(preserve_uuids=True)
        messages = [
            {
                "id": "550e8400-e29b-41d4-a716-446655440000",
                "content": "Hello user@test.com",
            }
        ]
        result = redact_messages(messages, opts)
        assert "550e8400" in result[0]["id"]
        assert "[email]" in result[0]["content"]


class TestRedactOptions:
    def test_default_keys_to_redact(self):
        opts = RedactOptions()
        assert "content" in opts.keys_to_redact
        assert "message" in opts.keys_to_redact
        assert "email" in opts.keys_to_redact
        assert "phone" in opts.keys_to_redact

    def test_custom_keys_to_redact(self):
        opts = RedactOptions(keys_to_redact=["api_key", "secret"])
        assert opts.keys_to_redact == ["api_key", "secret"]
        assert "content" not in opts.keys_to_redact

    def test_preserve_uuids_default(self):
        opts = RedactOptions()
        assert opts.preserve_uuids is False
