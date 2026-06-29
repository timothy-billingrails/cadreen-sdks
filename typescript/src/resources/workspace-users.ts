import type {
  WorkspaceUser,
  InviteUserRequest,
  UpdateRoleRequest,
  ListWorkspaceUsersResponse,
} from "../types";
import { HttpClient } from "../client";

export class WorkspaceUsersResource {
  constructor(private client: HttpClient) {}

  async list(): Promise<ListWorkspaceUsersResponse> {
    return this.client.get<ListWorkspaceUsersResponse>("/api/v1/cadreen/workspace/users");
  }

  async invite(request: InviteUserRequest): Promise<WorkspaceUser> {
    return this.client.post<WorkspaceUser>("/api/v1/cadreen/workspace/users", request);
  }

  async updateRole(id: string, role: string): Promise<WorkspaceUser> {
    return this.client.patch<WorkspaceUser>(
      `/api/v1/cadreen/workspace/users/${encodeURIComponent(id)}`,
      { role } as UpdateRoleRequest
    );
  }

  async remove(id: string): Promise<void> {
    await this.client.delete(`/api/v1/cadreen/workspace/users/${encodeURIComponent(id)}`);
  }
}
