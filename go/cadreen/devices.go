package cadreen

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// Device represents a registered hardware device.
type Device struct {
	ID         string        `json:"id"`
	Pose       *Pose         `json:"pose,omitempty"`
	Twist      *Twist        `json:"twist,omitempty"`
	Battery    *BatteryState `json:"battery,omitempty"`
	LastUpdate string        `json:"last_update,omitempty"`
}

// Twist represents linear and angular velocity.
type Twist struct {
	Linear  Point3D `json:"linear"`
	Angular Point3D `json:"angular"`
}

// BatteryState represents a device's battery status.
type BatteryState struct {
	Percentage *float64 `json:"percentage,omitempty"`
	Voltage    *float64 `json:"voltage,omitempty"`
	Current    *float64 `json:"current,omitempty"`
}

// Point2D represents a 2D point.
type Point2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// DeviceStatus represents a device's current status.
type DeviceStatus struct {
	ID        string         `json:"id"`
	Pose      map[string]any `json:"pose,omitempty"`
	Twist     map[string]any `json:"twist,omitempty"`
	Battery   map[string]any `json:"battery,omitempty"`
	LastUpdate string        `json:"last_update,omitempty"`
}

// DiagnosisResult represents a sensor diagnosis.
type DiagnosisResult struct {
	Rule     string  `json:"rule"`
	Severity string  `json:"severity"`
	Message  string  `json:"message"`
	SensorID string  `json:"sensor_id,omitempty"`
	Value    float64 `json:"value,omitempty"`
	Threshold float64 `json:"threshold,omitempty"`
}

// DiagnosisResponse is the response from the diagnose endpoint.
type DiagnosisResponse struct {
	DeviceID  string            `json:"device_id"`
	Diagnoses []DiagnosisResult `json:"diagnoses"`
	Count     int               `json:"count"`
}

// AskResponse is the response from the ask endpoint.
type AskResponse struct {
	Answer     string  `json:"answer"`
	Confidence float64 `json:"confidence"`
	Model      string  `json:"model"`
	CostCents  float64 `json:"cost_cents"`
}

// GridStats represents occupancy grid statistics.
type GridStats struct {
	TotalCells     int     `json:"total_cells"`
	ObservedCells  int     `json:"observed_cells"`
	FreeCells      int     `json:"free_cells"`
	OccupiedCells  int     `json:"occupied_cells"`
	TotalSources   int     `json:"total_sources"`
	AvgConfidence  float64 `json:"avg_confidence"`
}

// ListDevicesParams contains parameters for listing devices.
type ListDevicesParams struct {
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Type   string `json:"type,omitempty"`
}

// ListDevicesResponse is the response from the list devices endpoint.
type ListDevicesResponse struct {
	Devices []Device `json:"devices"`
	Total   int      `json:"total"`
}

// CreateDeviceRequest contains parameters for creating a device.
type CreateDeviceRequest struct {
	ID   string `json:"id,omitempty"`
	Pose *Pose  `json:"pose,omitempty"`
}

// Pose represents a 3D pose with position and orientation.
type Pose struct {
	Position    Point3D     `json:"position"`
	Orientation *Quaternion `json:"orientation,omitempty"`
}

// Point3D represents a 3D point.
type Point3D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// Quaternion represents a rotation quaternion.
type Quaternion struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
	W float64 `json:"w"`
}

// SensorReading represents a sensor reading for diagnosis.
type SensorReading struct {
	Name      string  `json:"name"`
	Value     float64 `json:"value"`
	Unit      string  `json:"unit"`
	DeviceID  string  `json:"device_id,omitempty"`
	Timestamp string  `json:"timestamp,omitempty"`
}

// DiagnoseDeviceRequest contains parameters for diagnosing a device.
type DiagnoseDeviceRequest struct {
	Readings []SensorReading `json:"readings"`
}

// ListTasksParams contains parameters for listing tasks.
type ListTasksParams struct {
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Status string `json:"status,omitempty"`
}

// Task represents a swarm task.
type Task struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Status     string    `json:"status"`
	Target     *Point2D  `json:"target,omitempty"`
	AssignedTo string    `json:"assigned_to,omitempty"`
	CreatedAt  string    `json:"created_at,omitempty"`
}

// ListTasksResponse is the response from the list tasks endpoint.
type ListTasksResponse struct {
	Tasks []Task `json:"tasks"`
	Total int    `json:"total"`
}

// CreateTaskRequest contains parameters for creating a task.
type CreateTaskRequest struct {
	Type   string   `json:"type"`
	Target *Point2D `json:"target"`
}

// ListDevices lists all registered devices.
func (c *Client) ListDevices(ctx context.Context, params ListDevicesParams, opts ...RequestOption) (*ListDevicesResponse, error) {
	path := "/api/v1/cadreen/devices"
	q := url.Values{}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset > 0 {
		q.Set("offset", strconv.Itoa(params.Offset))
	}
	if params.Type != "" {
		q.Set("type", params.Type)
	}
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	var result ListDevicesResponse
	if err := c.do(ctx, "GET", path, nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list devices: %w", err)
	}
	return &result, nil
}

// GetDevice gets a device by ID.
func (c *Client) GetDevice(ctx context.Context, deviceID string, opts ...RequestOption) (*Device, error) {
	var result Device
	if err := c.do(ctx, "GET", "/api/v1/cadreen/devices/"+deviceID, nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("get device: %w", err)
	}
	return &result, nil
}

// CreateDevice registers a new device.
func (c *Client) CreateDevice(ctx context.Context, req CreateDeviceRequest, opts ...RequestOption) (*Device, error) {
	var result Device
	if err := c.do(ctx, "POST", "/api/v1/cadreen/devices", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("create device: %w", err)
	}
	return &result, nil
}

// DeleteDevice removes a device.
func (c *Client) DeleteDevice(ctx context.Context, deviceID string, opts ...RequestOption) error {
	if err := c.do(ctx, "DELETE", "/api/v1/cadreen/devices/"+deviceID, nil, nil, opts...); err != nil {
		return fmt.Errorf("delete device: %w", err)
	}
	return nil
}

// GetDeviceStatus gets a device's current status.
func (c *Client) GetDeviceStatus(ctx context.Context, deviceID string, opts ...RequestOption) (*DeviceStatus, error) {
	var result DeviceStatus
	if err := c.do(ctx, "GET", "/api/v1/cadreen/devices/"+deviceID+"/status", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("get device status: %w", err)
	}
	return &result, nil
}

// DiagnoseDevice diagnoses sensor faults for a device.
func (c *Client) DiagnoseDevice(ctx context.Context, readings []SensorReading, opts ...RequestOption) (*DiagnosisResponse, error) {
	var result DiagnosisResponse
	if err := c.do(ctx, "POST", "/api/v1/cadreen/devices/diagnose", DiagnoseDeviceRequest{Readings: readings}, &result, opts...); err != nil {
		return nil, fmt.Errorf("diagnose device: %w", err)
	}
	return &result, nil
}

// AskDevice asks a question using hybrid inference.
func (c *Client) AskDevice(ctx context.Context, question string, opts ...RequestOption) (*AskResponse, error) {
	var result AskResponse
	if err := c.do(ctx, "POST", "/api/v1/cadreen/devices/ask", map[string]any{"question": question}, &result, opts...); err != nil {
		return nil, fmt.Errorf("ask device: %w", err)
	}
	return &result, nil
}

// GetGridStats gets occupancy grid statistics.
func (c *Client) GetGridStats(ctx context.Context, opts ...RequestOption) (*GridStats, error) {
	var result GridStats
	if err := c.do(ctx, "GET", "/api/v1/cadreen/devices/map/stats", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("get grid stats: %w", err)
	}
	return &result, nil
}

// ListTasks lists swarm tasks.
func (c *Client) ListTasks(ctx context.Context, params ListTasksParams, opts ...RequestOption) (*ListTasksResponse, error) {
	path := "/api/v1/cadreen/devices/tasks"
	q := url.Values{}
	if params.Limit > 0 {
		q.Set("limit", strconv.Itoa(params.Limit))
	}
	if params.Offset > 0 {
		q.Set("offset", strconv.Itoa(params.Offset))
	}
	if params.Status != "" {
		q.Set("status", params.Status)
	}
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	var result ListTasksResponse
	if err := c.do(ctx, "GET", path, nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	return &result, nil
}

// CreateTask creates a swarm task.
func (c *Client) CreateTask(ctx context.Context, req CreateTaskRequest, opts ...RequestOption) (*Task, error) {
	var result Task
	if err := c.do(ctx, "POST", "/api/v1/cadreen/devices/tasks", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	return &result, nil
}
