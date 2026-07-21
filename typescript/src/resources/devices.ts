import type {
  Device,
  DeviceStatus,
  DiagnosisResponse,
  AskResponse,
  GridStats,
  Task,
  CollisionWarning,
  AvoidanceManeuver,
  SyncStatus,
  BlackboardEntry,
  CreateDeviceRequest,
  DeviceDiagnoseRequest,
  CreateTaskRequest,
  ListDevicesResponse,
  ListTasksResponse,
  Pose,
} from "../types";
import { HttpClient } from "../client";

export class DevicesResource {
  constructor(private client: HttpClient) {}

  async list(params?: { limit?: number; offset?: number; type?: string }): Promise<ListDevicesResponse> {
    return this.client.get<ListDevicesResponse>("/api/v1/cadreen/devices", params as Record<string, string | number | boolean | undefined>);
  }

  async create(request: CreateDeviceRequest): Promise<{ id: string; status: string; message: string }> {
    return this.client.post("/api/v1/cadreen/devices", request);
  }

  async get(deviceId: string): Promise<Device> {
    return this.client.get<Device>(`/api/v1/cadreen/devices/${encodeURIComponent(deviceId)}`);
  }

  async delete(deviceId: string): Promise<{ id: string; status: string; message: string }> {
    return this.client.delete(`/api/v1/cadreen/devices/${encodeURIComponent(deviceId)}`);
  }

  async getStatus(deviceId: string): Promise<DeviceStatus> {
    return this.client.get<DeviceStatus>(`/api/v1/cadreen/devices/${encodeURIComponent(deviceId)}/status`);
  }

  async updateState(deviceId: string, request: { pose: Pose }): Promise<{ id: string; status: string; message: string }> {
    return this.client.post(`/api/v1/cadreen/devices/${encodeURIComponent(deviceId)}/state`, request);
  }

  async getMap(): Promise<{ grid: unknown; devices: number; resolution: number }> {
    return this.client.get("/api/v1/cadreen/devices/map");
  }

  async getMapStats(): Promise<GridStats> {
    return this.client.get<GridStats>("/api/v1/cadreen/devices/map/stats");
  }

  async updateMap(request: { device_id: string; device_pose: Pose; laser_scan?: unknown; point_cloud?: unknown; occupancy_grid?: unknown }): Promise<{ status: string; message: string; devices: number }> {
    return this.client.post("/api/v1/cadreen/devices/map", request);
  }

  async listTasks(params?: { limit?: number; offset?: number; status?: string }): Promise<ListTasksResponse> {
    return this.client.get<ListTasksResponse>("/api/v1/cadreen/devices/tasks", params as Record<string, string | number | boolean | undefined>);
  }

  async createTask(request: CreateTaskRequest): Promise<Task> {
    return this.client.post<Task>("/api/v1/cadreen/devices/tasks", request);
  }

  async completeTask(taskId: string): Promise<{ id: string; status: string; message: string }> {
    return this.client.post(`/api/v1/cadreen/devices/tasks/${encodeURIComponent(taskId)}/complete`);
  }

  async assignTasks(): Promise<{ assigned: number; message: string }> {
    return this.client.post("/api/v1/cadreen/devices/assign");
  }

  async detectCollisions(): Promise<{ warnings: CollisionWarning[]; total: number }> {
    return this.client.get("/api/v1/cadreen/devices/collisions");
  }

  async getAvoidance(): Promise<{ maneuvers: AvoidanceManeuver[]; total: number }> {
    return this.client.get("/api/v1/cadreen/devices/avoidance");
  }

  async diagnose(request: DeviceDiagnoseRequest): Promise<DiagnosisResponse> {
    return this.client.post<DiagnosisResponse>("/api/v1/cadreen/devices/diagnose", request);
  }

  async ask(question: string): Promise<AskResponse> {
    return this.client.post<AskResponse>("/api/v1/cadreen/devices/ask", { question });
  }

  async getModelStats(): Promise<{ status: string; message?: string; total_requests?: number; local_requests?: number; remote_requests?: number }> {
    return this.client.get("/api/v1/cadreen/devices/diagnostics/stats");
  }

  async getCapabilities(): Promise<{ diagnostic_rules: number; devices: number; message: string }> {
    return this.client.get("/api/v1/cadreen/devices/diagnostics/capabilities");
  }

  async getSyncStatus(): Promise<SyncStatus> {
    return this.client.get<SyncStatus>("/api/v1/cadreen/devices/sync/status");
  }

  async getSyncPending(): Promise<{ pending: number; conflicts: number; message: string }> {
    return this.client.get("/api/v1/cadreen/devices/sync/pending");
  }

  async getSyncConflicts(): Promise<{ conflicts: unknown[]; total: number }> {
    return this.client.get("/api/v1/cadreen/devices/sync/conflicts");
  }

  async getBlackboard(params?: { category?: string; hours?: number; limit?: number }): Promise<{ entries: BlackboardEntry[]; total: number; message: string }> {
    return this.client.get("/api/v1/cadreen/devices/blackboard", params as Record<string, string | number | boolean | undefined>);
  }
}
