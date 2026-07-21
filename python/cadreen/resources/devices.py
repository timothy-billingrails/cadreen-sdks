from __future__ import annotations

from typing import Any

from ..client import HttpClient
from ..types import (
    Device,
    DeviceStatus,
    DiagnosisResponse,
    AskResponse,
    GridStats,
    Task,
    TaskStats,
    CollisionWarning,
    SyncStatus,
    BlackboardEntry,
    CreateDeviceRequest,
    DeviceDiagnoseRequest,
    CreateTaskRequest,
    ListDevicesResponse,
    ListTasksResponse,
    Pose,
    SensorReading,
)


class DevicesResource:
    def __init__(self, client: HttpClient) -> None:
        self._client = client

    async def list(
        self,
        *,
        limit: int | None = None,
        offset: int | None = None,
        type: str | None = None,
    ) -> ListDevicesResponse:
        params: dict[str, Any] = {}
        if limit is not None:
            params["limit"] = limit
        if offset is not None:
            params["offset"] = offset
        if type is not None:
            params["type"] = type
        raw = await self._client.get("/api/v1/cadreen/devices", params)
        return ListDevicesResponse(
            devices=[_parse_device(d) for d in raw.get("devices", [])],
            total=raw.get("total", 0),
        )

    async def create(self, request: CreateDeviceRequest) -> dict[str, Any]:
        body: dict[str, Any] = {}
        if request.id is not None:
            body["id"] = request.id
        if request.pose is not None:
            body["pose"] = _pose_to_dict(request.pose)
        return await self._client.post("/api/v1/cadreen/devices", body)

    async def get(self, device_id: str) -> Device:
        raw = await self._client.get(f"/api/v1/cadreen/devices/{device_id}")
        return _parse_device(raw)

    async def delete(self, device_id: str) -> dict[str, Any]:
        return await self._client.delete(f"/api/v1/cadreen/devices/{device_id}")

    async def get_status(self, device_id: str) -> DeviceStatus:
        raw = await self._client.get(f"/api/v1/cadreen/devices/{device_id}/status")
        return DeviceStatus(
            id=raw["id"],
            status=raw.get("status", "unknown"),
            pose=_parse_pose(raw.get("pose")) if raw.get("pose") else None,
            battery=_parse_battery(raw.get("battery")) if raw.get("battery") else None,
            last_update=raw.get("last_update"),
        )

    async def update_state(self, device_id: str, pose: Pose) -> dict[str, Any]:
        return await self._client.post(
            f"/api/v1/cadreen/devices/{device_id}/state",
            {"pose": _pose_to_dict(pose)},
        )

    async def get_map(self) -> dict[str, Any]:
        return await self._client.get("/api/v1/cadreen/devices/map")

    async def get_map_stats(self) -> GridStats:
        raw = await self._client.get("/api/v1/cadreen/devices/map/stats")
        return GridStats(
            total_cells=raw.get("total_cells", 0),
            observed_cells=raw.get("observed_cells", 0),
            free_cells=raw.get("free_cells", 0),
            occupied_cells=raw.get("occupied_cells", 0),
            total_sources=raw.get("total_sources", 0),
            avg_confidence=raw.get("avg_confidence", 0.0),
        )

    async def update_map(self, request: dict[str, Any]) -> dict[str, Any]:
        return await self._client.post("/api/v1/cadreen/devices/map", request)

    async def list_tasks(
        self,
        *,
        limit: int | None = None,
        offset: int | None = None,
        status: str | None = None,
    ) -> ListTasksResponse:
        params: dict[str, Any] = {}
        if limit is not None:
            params["limit"] = limit
        if offset is not None:
            params["offset"] = offset
        if status is not None:
            params["status"] = status
        raw = await self._client.get("/api/v1/cadreen/devices/tasks", params)
        return ListTasksResponse(
            tasks=[_parse_task(t) for t in raw.get("tasks", [])],
            total=raw.get("total", 0),
        )

    async def create_task(self, request: CreateTaskRequest) -> Task:
        body: dict[str, Any] = {
            "type": request.type,
            "target": {"x": request.target.x, "y": request.target.y},
        }
        if request.priority is not None:
            body["priority"] = request.priority
        raw = await self._client.post("/api/v1/cadreen/devices/tasks", body)
        return _parse_task(raw)

    async def complete_task(self, task_id: str) -> dict[str, Any]:
        return await self._client.post(f"/api/v1/cadreen/devices/tasks/{task_id}/complete")

    async def assign_tasks(self) -> dict[str, Any]:
        return await self._client.post("/api/v1/cadreen/devices/assign")

    async def detect_collisions(self) -> dict[str, Any]:
        return await self._client.get("/api/v1/cadreen/devices/collisions")

    async def get_avoidance(self) -> dict[str, Any]:
        return await self._client.get("/api/v1/cadreen/devices/avoidance")

    async def diagnose(self, request: DeviceDiagnoseRequest) -> DiagnosisResponse:
        body: dict[str, Any] = {
            "readings": [
                {
                    "name": r.name,
                    "value": r.value,
                    **({"unit": r.unit} if r.unit else {}),
                    **({"device_id": r.device_id} if r.device_id else {}),
                }
                for r in request.readings
            ]
        }
        raw = await self._client.post("/api/v1/cadreen/devices/diagnose", body)
        return DiagnosisResponse(
            diagnoses=[
                FaultDiagnosis(
                    fault_id=d.get("fault_id", ""),
                    fault_type=d.get("fault_type", ""),
                    severity=d.get("severity", "info"),
                    confidence=d.get("confidence", 0.0),
                    description=d.get("description", ""),
                    root_cause=d.get("root_cause"),
                    remediation=d.get("remediation"),
                )
                for d in raw.get("diagnoses", [])
            ],
            total=raw.get("total", 0),
        )

    async def ask(self, question: str) -> AskResponse:
        raw = await self._client.post("/api/v1/cadreen/devices/ask", {"question": question})
        return AskResponse(
            answer=raw.get("answer", ""),
            confidence=raw.get("confidence", 0.0),
            model=raw.get("model", ""),
            cost_cents=raw.get("cost_cents", 0.0),
        )

    async def get_model_stats(self) -> dict[str, Any]:
        return await self._client.get("/api/v1/cadreen/devices/diagnostics/stats")

    async def get_capabilities(self) -> dict[str, Any]:
        return await self._client.get("/api/v1/cadreen/devices/diagnostics/capabilities")

    async def get_sync_status(self) -> SyncStatus:
        raw = await self._client.get("/api/v1/cadreen/devices/sync/status")
        return SyncStatus(
            status=raw.get("status", "unknown"),
            message=raw.get("message"),
            connected=raw.get("connected"),
            latency=raw.get("latency"),
        )

    async def get_sync_pending(self) -> dict[str, Any]:
        return await self._client.get("/api/v1/cadreen/devices/sync/pending")

    async def get_sync_conflicts(self) -> dict[str, Any]:
        return await self._client.get("/api/v1/cadreen/devices/sync/conflicts")

    async def get_blackboard(
        self,
        *,
        category: str | None = None,
        hours: int | None = None,
        limit: int | None = None,
    ) -> dict[str, Any]:
        params: dict[str, Any] = {}
        if category is not None:
            params["category"] = category
        if hours is not None:
            params["hours"] = hours
        if limit is not None:
            params["limit"] = limit
        return await self._client.get("/api/v1/cadreen/devices/blackboard", params)


def _parse_device(raw: dict[str, Any]) -> Device:
    return Device(
        id=raw["id"],
        pose=_parse_pose(raw.get("pose")) if raw.get("pose") else None,
        battery=_parse_battery(raw.get("battery")) if raw.get("battery") else None,
        last_update=raw.get("last_update"),
    )


def _parse_pose(raw: dict[str, Any] | None) -> Pose | None:
    if raw is None:
        return None
    pos = raw.get("position", {})
    orient = raw.get("orientation")
    return Pose(
        position=Point3D(x=pos.get("x", 0), y=pos.get("y", 0), z=pos.get("z", 0)),
        orientation=Quaternion(x=orient["x"], y=orient["y"], z=orient["z"], w=orient["w"]) if orient else None,
    )


def _parse_battery(raw: dict[str, Any] | None) -> BatteryState | None:
    if raw is None:
        return None
    return BatteryState(
        percentage=raw.get("percentage"),
        voltage=raw.get("voltage"),
        current=raw.get("current"),
    )


def _parse_task(raw: dict[str, Any]) -> Task:
    target_raw = raw.get("target")
    return Task(
        id=raw["id"],
        type=raw.get("type", ""),
        status=raw.get("status", "pending"),
        target=Point2D(x=target_raw["x"], y=target_raw["y"]) if target_raw else None,
        assigned_to=raw.get("assigned_to"),
        created_at=raw.get("created_at"),
    )


def _pose_to_dict(pose: Pose) -> dict[str, Any]:
    result: dict[str, Any] = {
        "position": {"x": pose.position.x, "y": pose.position.y, "z": pose.position.z},
    }
    if pose.orientation:
        result["orientation"] = {
            "x": pose.orientation.x,
            "y": pose.orientation.y,
            "z": pose.orientation.z,
            "w": pose.orientation.w,
        }
    return result
