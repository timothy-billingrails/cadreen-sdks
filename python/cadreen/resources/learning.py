from __future__ import annotations

from ..client import HttpClient
from ..types import (
    ListLearningPatternsResponse,
    ListLearningEpisodesResponse,
    ListLearningSuggestionsResponse,
    LearningPattern,
    LearningEpisode,
    LearningSuggestion,
    Pagination,
)


class LearningResource:
    def __init__(self, client: HttpClient) -> None:
        self._client = client

    async def patterns(self) -> ListLearningPatternsResponse:
        raw = await self._client.get("/api/v1/cadreen/learning/patterns")
        patterns = [
            LearningPattern(
                id=p["id"],
                pattern=p["pattern"],
                confidence=p.get("confidence", 0.0),
                occurrences=p.get("occurrences"),
                domain=p.get("domain"),
                tags=p.get("tags"),
                created_at=p.get("created_at"),
            )
            for p in raw.get("patterns", [])
        ]
        pagination = None
        if raw.get("pagination"):
            p = raw["pagination"]
            pagination = Pagination(limit=p["limit"], offset=p["offset"], has_more=p["has_more"])
        return ListLearningPatternsResponse(
            patterns=patterns,
            count=raw.get("count", 0),
            pagination=pagination,
        )

    async def episodes(self) -> ListLearningEpisodesResponse:
        raw = await self._client.get("/api/v1/cadreen/learning/episodes")
        episodes = [
            LearningEpisode(
                id=e["id"],
                description=e["description"],
                outcome=e.get("outcome"),
                trace_id=e.get("trace_id"),
                domain=e.get("domain"),
                created_at=e.get("created_at"),
            )
            for e in raw.get("episodes", [])
        ]
        pagination = None
        if raw.get("pagination"):
            p = raw["pagination"]
            pagination = Pagination(limit=p["limit"], offset=p["offset"], has_more=p["has_more"])
        return ListLearningEpisodesResponse(
            episodes=episodes,
            count=raw.get("count", 0),
            pagination=pagination,
        )

    async def suggestions(self) -> ListLearningSuggestionsResponse:
        raw = await self._client.get("/api/v1/cadreen/learning/suggestions")
        suggestions = [
            LearningSuggestion(
                id=s["id"],
                type=s["type"],
                description=s["description"],
                impact=s.get("impact"),
                domain=s.get("domain"),
            )
            for s in raw.get("suggestions", [])
        ]
        pagination = None
        if raw.get("pagination"):
            p = raw["pagination"]
            pagination = Pagination(limit=p["limit"], offset=p["offset"], has_more=p["has_more"])
        return ListLearningSuggestionsResponse(
            suggestions=suggestions,
            count=raw.get("count", 0),
            pagination=pagination,
        )
