import type {
  ListLearningPatternsResponse,
  ListLearningEpisodesResponse,
  ListLearningSuggestionsResponse,
} from "../types";
import { HttpClient } from "../client";

export class LearningResource {
  constructor(private client: HttpClient) {}

  async patterns(): Promise<ListLearningPatternsResponse> {
    return this.client.get<ListLearningPatternsResponse>("/api/v1/cadreen/learning/patterns");
  }

  async episodes(): Promise<ListLearningEpisodesResponse> {
    return this.client.get<ListLearningEpisodesResponse>("/api/v1/cadreen/learning/episodes");
  }

  async suggestions(): Promise<ListLearningSuggestionsResponse> {
    return this.client.get<ListLearningSuggestionsResponse>("/api/v1/cadreen/learning/suggestions");
  }
}
