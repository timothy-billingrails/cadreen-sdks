package commands

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	"github.com/spf13/cobra"
)

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Manage hardware devices (register, diagnose, query)",
	Long: `Devices are hardware endpoints connected via MQTT, OPC-UA, or other protocols.

Register a device, send telemetry, run diagnostics, and query the fused
occupancy grid — all from the command line.`,
}

var devicesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered devices",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/devices", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
			return nil
		}

		devices, _ := result["devices"].([]any)
		total, _ := result["total"].(float64)
		fmt.Printf("Devices: %d\n\n", int(total))
		for _, d := range devices {
			dev, _ := d.(map[string]any)
			id, _ := dev["id"].(string)
			pose, _ := dev["pose"].(map[string]any)
			battery, _ := dev["battery"].(map[string]any)
			fmt.Printf("  %s\n", id)
			if pose != nil {
				if x, ok := pose["x"].(float64); ok {
					fmt.Printf("    Position: (%.1f, %.1f, %.1f)\n", x, pose["y"], pose["z"])
				}
			}
			if battery != nil {
				if level, ok := battery["level"].(float64); ok {
					fmt.Printf("    Battery:  %.0f%%\n", level*100)
				}
			}
		}
		return nil
	},
}

var devicesGetCmd = &cobra.Command{
	Use:   "get [device-id]",
	Short: "Get device details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/devices/"+args[0], nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		output.Print(result, outputFormat())
		return nil
	},
}

var devicesStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get system status",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/devices/diagnostics/capabilities", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
			return nil
		}

		fmt.Printf("Diagnostic rules: %v\n", result["diagnostic_rules"])
		fmt.Printf("Devices:          %v\n", result["devices"])
		if msg, ok := result["message"].(string); ok {
			fmt.Printf("Status:           %s\n", msg)
		}
		return nil
	},
}

var devicesDiagnoseCmd = &cobra.Command{
	Use:   "diagnose",
	Short: "Diagnose sensor faults from readings",
	Long: `Diagnose sensor faults from one or more readings.

Examples:
  cadreen devices diagnose --reading motor_current=15 --reading temperature=90
  cadreen devices diagnose --reading vibration=5.2 --unit mm/s`,
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		readingsRaw, _ := cmd.Flags().GetStringSlice("reading")
		if len(readingsRaw) == 0 {
			return fmt.Errorf("--reading is required (e.g., --reading motor_current=15)")
		}

		unit, _ := cmd.Flags().GetString("unit")
		deviceID, _ := cmd.Flags().GetString("device-id")

		var readings []map[string]any
		for _, r := range readingsRaw {
			parts := strings.SplitN(r, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid reading format: %s (expected name=value)", r)
			}
			val, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return fmt.Errorf("invalid value for %s: %w", parts[0], err)
			}
			reading := map[string]any{
				"name":  parts[0],
				"value": val,
			}
			if unit != "" {
				reading["unit"] = unit
			}
			if deviceID != "" {
				reading["device_id"] = deviceID
			}
			readings = append(readings, reading)
		}

		payload := map[string]any{"readings": readings}
		resp, err := rawDo("POST", "/api/v1/cadreen/devices/diagnose", payload)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
			return nil
		}

		diagnoses, _ := result["diagnoses"].([]any)
		if len(diagnoses) == 0 {
			fmt.Println("No faults detected.")
			return nil
		}

		fmt.Printf("Found %d fault(s):\n\n", len(diagnoses))
		for _, d := range diagnoses {
			diag, _ := d.(map[string]any)
			rule, _ := diag["rule"].(string)
			severity, _ := diag["severity"].(string)
			message, _ := diag["message"].(string)
			fmt.Printf("  [%s] %s\n", severity, rule)
			fmt.Printf("    %s\n\n", message)
		}
		return nil
	},
}

var devicesAskCmd = &cobra.Command{
	Use:   "ask [question]",
	Short: "Ask a question (hybrid inference)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		payload := map[string]any{"question": args[0]}
		resp, err := rawDo("POST", "/api/v1/cadreen/devices/ask", payload)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
			return nil
		}

		answer, _ := result["answer"].(string)
		confidence, _ := result["confidence"].(float64)
		model, _ := result["model"].(string)

		fmt.Printf("%s\n\n", answer)
		fmt.Printf("Confidence: %.0f%%\n", confidence*100)
		fmt.Printf("Model:      %s\n", model)
		return nil
	},
}

var devicesGridCmd = &cobra.Command{
	Use:   "grid",
	Short: "Get fused occupancy grid",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/devices/map", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
			return nil
		}

		width, _ := result["width"].(float64)
		height, _ := result["height"].(float64)
		resolution, _ := result["resolution"].(float64)
		fmt.Printf("Grid: %dx%d (resolution: %.2fm)\n", int(width), int(height), resolution)

		stats, _ := result["stats"].(map[string]any)
		if stats != nil {
			fmt.Printf("Total cells:     %v\n", stats["total_cells"])
			fmt.Printf("Observed cells:  %v\n", stats["observed_cells"])
			fmt.Printf("Free cells:      %v\n", stats["free_cells"])
			fmt.Printf("Occupied cells:  %v\n", stats["occupied_cells"])
		}
		return nil
	},
}

var devicesTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Manage swarm tasks",
}

var devicesTasksListCmd = &cobra.Command{
	Use:   "list",
	Short: "List swarm tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/devices/tasks", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
			return nil
		}

		tasks, _ := result["tasks"].([]any)
		total, _ := result["total"].(float64)
		fmt.Printf("Tasks: %d\n\n", int(total))
		for _, t := range tasks {
			task, _ := t.(map[string]any)
			id, _ := task["id"].(string)
			taskType, _ := task["type"].(string)
			status, _ := task["status"].(string)
			assigned, _ := task["assigned_to"].(string)
			fmt.Printf("  %s\n", id)
			fmt.Printf("    Type:     %s\n", taskType)
			fmt.Printf("    Status:   %s\n", status)
			if assigned != "" {
				fmt.Printf("    Assigned: %s\n", assigned)
			}
		}
		return nil
	},
}

var devicesTasksAddCmd = &cobra.Command{
	Use:   "add [type] [target-x] [target-y]",
	Short: "Add a swarm task",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		x, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			return fmt.Errorf("invalid target-x: %w", err)
		}
		y, err := strconv.ParseFloat(args[2], 64)
		if err != nil {
			return fmt.Errorf("invalid target-y: %w", err)
		}

		payload := map[string]any{
			"type": args[0],
			"target": map[string]any{
				"x": x,
				"y": y,
				"z": 0,
			},
		}

		resp, err := rawDo("POST", "/api/v1/cadreen/devices/tasks", payload)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 && resp.StatusCode != 201 {
			return handleHTTPError(resp)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
			return nil
		}

		fmt.Printf("Task created: %v\n", result["id"])
		fmt.Printf("  Type:   %s\n", args[0])
		fmt.Printf("  Target: (%s, %s)\n", args[1], args[2])
		return nil
	},
}

var devicesSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Edge-cloud sync status",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/devices/sync/status", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
			return nil
		}

		status, _ := result["status"].(string)
		message, _ := result["message"].(string)
		fmt.Printf("Sync status: %s\n", status)
		if message != "" {
			fmt.Printf("  %s\n", message)
		}
		return nil
	},
}

var devicesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Register a new device",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		id, _ := cmd.Flags().GetString("id")
		payload := map[string]any{}
		if id != "" {
			payload["id"] = id
		}

		resp, err := rawDo("POST", "/api/v1/cadreen/devices", payload)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 && resp.StatusCode != 201 {
			return handleHTTPError(resp)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
			return nil
		}

		fmt.Printf("Device created: %v\n", result["id"])
		return nil
	},
}

var devicesDeleteCmd = &cobra.Command{
	Use:   "delete [device-id]",
	Short: "Remove a device",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("DELETE", "/api/v1/cadreen/devices/"+args[0], nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
			return nil
		}

		fmt.Printf("Device removed: %v\n", result["id"])
		return nil
	},
}

var devicesCollisionsCmd = &cobra.Command{
	Use:   "collisions",
	Short: "Check for collision warnings",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/devices/collisions", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
			return nil
		}

		warnings, _ := result["warnings"].([]any)
		total, _ := result["total"].(float64)
		if int(total) == 0 {
			fmt.Println("No collision warnings.")
			return nil
		}

		fmt.Printf("Collision warnings: %d\n\n", int(total))
		for _, w := range warnings {
			warning, _ := w.(map[string]any)
			fmt.Printf("  %v ↔ %v\n", warning["device1"], warning["device2"])
			fmt.Printf("    Distance: %.2fm\n", warning["distance"])
			fmt.Printf("    Severity: %v\n\n", warning["severity"])
		}
		return nil
	},
}

var devicesAvoidanceCmd = &cobra.Command{
	Use:   "avoidance",
	Short: "Get avoidance maneuvers",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/devices/avoidance", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
			return nil
		}

		maneuvers, _ := result["maneuvers"].([]any)
		total, _ := result["total"].(float64)
		if int(total) == 0 {
			fmt.Println("No avoidance maneuvers needed.")
			return nil
		}

		fmt.Printf("Avoidance maneuvers: %d\n\n", int(total))
		for _, m := range maneuvers {
			maneuver, _ := m.(map[string]any)
			fmt.Printf("  Device: %v\n", maneuver["agent_id"])
			fmt.Printf("    Reason: %v\n", maneuver["reason"])
		}
		return nil
	},
}

var devicesAssignCmd = &cobra.Command{
	Use:   "assign",
	Short: "Assign tasks to nearest devices",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("POST", "/api/v1/cadreen/devices/assign", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
			return nil
		}

		assigned, _ := result["assigned"].(float64)
		fmt.Printf("Assigned %d tasks\n", int(assigned))
		return nil
	},
}

var devicesMapStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Get map coverage statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/devices/map/stats", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
			return nil
		}

		fmt.Printf("Total cells:     %v\n", result["total_cells"])
		fmt.Printf("Observed cells:  %v\n", result["observed_cells"])
		fmt.Printf("Free cells:      %v\n", result["free_cells"])
		fmt.Printf("Occupied cells:  %v\n", result["occupied_cells"])
		return nil
	},
}

var devicesTasksCompleteCmd = &cobra.Command{
	Use:   "complete [task-id]",
	Short: "Mark a task as completed",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("POST", "/api/v1/cadreen/devices/tasks/"+args[0]+"/complete", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
			return nil
		}

		fmt.Printf("Task completed: %v\n", result["id"])
		return nil
	},
}

var devicesBlackboardCmd = &cobra.Command{
	Use:   "blackboard",
	Short: "Read blackboard entries (shared state)",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		category, _ := cmd.Flags().GetString("category")
		hours, _ := cmd.Flags().GetInt("hours")
		limit, _ := cmd.Flags().GetInt("limit")

		path := "/api/v1/cadreen/devices/blackboard"
		q := url.Values{}
		if category != "" {
			q.Set("category", category)
		}
		if hours > 0 {
			q.Set("hours", strconv.Itoa(hours))
		}
		if limit > 0 {
			q.Set("limit", strconv.Itoa(limit))
		}
		if len(q) > 0 {
			path += "?" + q.Encode()
		}

		resp, err := rawDo("GET", path, nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
			return nil
		}

		entries, _ := result["entries"].([]any)
		total, _ := result["total"].(float64)
		fmt.Printf("Blackboard entries: %d\n\n", int(total))
		for _, e := range entries {
			entry, _ := e.(map[string]any)
			fmt.Printf("  [%v] %v = %v\n", entry["category"], entry["key"], entry["value"])
		}
		return nil
	},
}

var devicesMapCmd = &cobra.Command{
	Use:   "map",
	Short: "View the fused occupancy grid",
}

func init() {
	devicesCmd.AddCommand(devicesListCmd)
	devicesCmd.AddCommand(devicesGetCmd)
	devicesCmd.AddCommand(devicesCreateCmd)
	devicesCmd.AddCommand(devicesDeleteCmd)
	devicesCmd.AddCommand(devicesStatusCmd)
	devicesCmd.AddCommand(devicesDiagnoseCmd)
	devicesCmd.AddCommand(devicesAskCmd)
	devicesCmd.AddCommand(devicesGridCmd)
	devicesCmd.AddCommand(devicesMapCmd)
	devicesCmd.AddCommand(devicesTasksCmd)
	devicesCmd.AddCommand(devicesCollisionsCmd)
	devicesCmd.AddCommand(devicesAvoidanceCmd)
	devicesCmd.AddCommand(devicesAssignCmd)
	devicesCmd.AddCommand(devicesSyncCmd)
	devicesCmd.AddCommand(devicesBlackboardCmd)

	devicesTasksCmd.AddCommand(devicesTasksListCmd)
	devicesTasksCmd.AddCommand(devicesTasksAddCmd)
	devicesTasksCmd.AddCommand(devicesTasksCompleteCmd)

	devicesMapCmd.AddCommand(devicesGridCmd)
	devicesMapCmd.AddCommand(devicesMapStatsCmd)

	devicesCreateCmd.Flags().StringP("id", "i", "", "Device ID (generated if empty)")

	devicesDiagnoseCmd.Flags().StringSliceP("reading", "r", nil, "Sensor reading in name=value format (repeatable)")
	devicesDiagnoseCmd.Flags().StringP("unit", "u", "", "Unit for all readings (e.g., A, C, mm/s)")
	devicesDiagnoseCmd.Flags().StringP("device-id", "d", "", "Device ID (optional)")

	devicesBlackboardCmd.Flags().StringP("category", "c", "", "Filter by category")
	devicesBlackboardCmd.Flags().IntP("hours", "H", 24, "Hours to look back")
	devicesBlackboardCmd.Flags().IntP("limit", "l", 50, "Max entries")

	rootCmd.AddCommand(devicesCmd)
}
