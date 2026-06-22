package commands

import (
	"encoding/json"
	"fmt"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	"github.com/spf13/cobra"
)

type blueprintSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status"`
	Version     int    `json:"version"`
	Intent      string `json:"intent,omitempty"`
	CreatedAt   string `json:"created_at"`
}

type blueprintRun struct {
	ID            string `json:"id"`
	BlueprintID   string `json:"blueprint_id"`
	Status        string `json:"status"`
	ResultSummary string `json:"result_summary,omitempty"`
	CreatedAt     string `json:"created_at"`
}

var blueprintsCmd = &cobra.Command{
	Use:   "blueprints",
	Short: "Manage reusable workflows (list, create, run)",
	Long: `Blueprints are reusable templates for work Cadreen does.

Create one from a trace or from scratch, then run it with different inputs.`,
}

var blueprintsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List blueprints",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		status, _ := cmd.Flags().GetString("status")
		limit, _ := cmd.Flags().GetInt("limit")

		resp, err := rawDo("GET", fmt.Sprintf("/api/v1/cadreen/blueprints?status=%s&limit=%d", status, limit), nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result struct {
			Blueprints []blueprintSummary `json:"blueprints"`
			Count      int                `json:"count"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(result, format)
		default:
			printBlueprintsText(result.Blueprints)
		}

		return nil
	},
}

var blueprintsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a blueprint",
	Long: `Create a reusable blueprint.

Examples:
  cadreen blueprints create --name "Process refund" --description "Validate and process refunds"
  cadreen blueprints create --name "Daily report" --param-schema '{"type":"object","properties":{"date":{"type":"string"}}}'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		paramSchema, _ := cmd.Flags().GetString("param-schema")

		if name == "" {
			return fmt.Errorf("--name is required")
		}

		payload := map[string]any{
			"name": name,
		}
		if description != "" {
			payload["description"] = description
		}
		if paramSchema != "" {
			var schema map[string]any
			if err := json.Unmarshal([]byte(paramSchema), &schema); err != nil {
				return fmt.Errorf("invalid --param-schema JSON: %w", err)
			}
			payload["parameter_schema"] = schema
		}

		resp, err := rawDo("POST", "/api/v1/cadreen/blueprints", payload)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 && resp.StatusCode != 201 {
			return handleHTTPError(resp)
		}

		var result blueprintSummary
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
		} else {
			fmt.Printf("Blueprint created: %s\n", result.ID)
			fmt.Printf("  Name:    %s\n", result.Name)
			fmt.Printf("  Status:  %s\n", result.Status)
			fmt.Printf("  Version: %d\n", result.Version)
		}

		return nil
	},
}

var blueprintsRunCmd = &cobra.Command{
	Use:   "run [blueprint-id]",
	Short: "Run a blueprint",
	Long: `Execute a blueprint with optional parameter overrides.

Examples:
  cadreen blueprints run bp_01JX
  cadreen blueprints run bp_01JX --param date=2026-06-22
  cadreen blueprints run bp_01JX --param date=2026-06-22 --param recipients=team@example.com`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		blueprintID := args[0]
		params, _ := cmd.Flags().GetStringToString("param")

		payload := map[string]any{}
		if len(params) > 0 {
			payload["params"] = params
		}

		resp, err := rawDo("POST", fmt.Sprintf("/api/v1/cadreen/blueprints/%s/runs", blueprintID), payload)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 && resp.StatusCode != 201 {
			return handleHTTPError(resp)
		}

		var result struct {
			ID          string `json:"id"`
			BlueprintID string `json:"blueprint_id"`
			Status      string `json:"status"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
		} else {
			fmt.Printf("Run started: %s\n", result.ID)
			fmt.Printf("  Blueprint: %s\n", result.BlueprintID)
			fmt.Printf("  Status:    %s\n", result.Status)
		}

		return nil
	},
}

var blueprintsShowCmd = &cobra.Command{
	Use:   "show [blueprint-id]",
	Short: "View blueprint details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", fmt.Sprintf("/api/v1/cadreen/blueprints/%s", args[0]), nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result blueprintSummary
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
		} else {
			printBlueprintDetail(&result)
		}

		return nil
	},
}

var blueprintsRunsCmd = &cobra.Command{
	Use:   "runs [blueprint-id]",
	Short: "View run history for a blueprint",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		limit, _ := cmd.Flags().GetInt("limit")

		resp, err := rawDo("GET", fmt.Sprintf("/api/v1/cadreen/blueprints/%s/runs?limit=%d", args[0], limit), nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result struct {
			Runs  []blueprintRun `json:"runs"`
			Count int            `json:"count"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(result, format)
		default:
			printBlueprintRunsText(result.Runs)
		}

		return nil
	},
}

var blueprintsArchiveCmd = &cobra.Command{
	Use:   "archive [blueprint-id]",
	Short: "Archive a blueprint",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("DELETE", fmt.Sprintf("/api/v1/cadreen/blueprints/%s", args[0]), nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 && resp.StatusCode != 204 {
			return handleHTTPError(resp)
		}

		fmt.Println("Blueprint archived.")
		return nil
	},
}

func init() {
	blueprintsListCmd.Flags().String("status", "active", "filter by status (active, archived)")
	blueprintsListCmd.Flags().Int("limit", 50, "max results")

	blueprintsCreateCmd.Flags().String("name", "", "blueprint name (required)")
	blueprintsCreateCmd.Flags().String("description", "", "what this blueprint does")
	blueprintsCreateCmd.Flags().String("param-schema", "", "JSON Schema for parameters")

	blueprintsRunCmd.Flags().StringToString("param", nil, "parameter overrides (key=value)")

	blueprintsRunsCmd.Flags().Int("limit", 10, "max results")

	blueprintsCmd.AddCommand(blueprintsListCmd)
	blueprintsCmd.AddCommand(blueprintsCreateCmd)
	blueprintsCmd.AddCommand(blueprintsRunCmd)
	blueprintsCmd.AddCommand(blueprintsShowCmd)
	blueprintsCmd.AddCommand(blueprintsRunsCmd)
	blueprintsCmd.AddCommand(blueprintsArchiveCmd)
	rootCmd.AddCommand(blueprintsCmd)
}

func printBlueprintsText(blueprints []blueprintSummary) {
	if len(blueprints) == 0 {
		fmt.Println("No blueprints found.")
		return
	}

	fmt.Printf("Blueprints (%d):\n\n", len(blueprints))
	for _, b := range blueprints {
		fmt.Printf("  %s  %s\n", b.ID, b.Name)
		fmt.Printf("    Status: %s | Version: %d\n", b.Status, b.Version)
		fmt.Println()
	}
}

func printBlueprintDetail(b *blueprintSummary) {
	fmt.Printf("ID:       %s\n", b.ID)
	fmt.Printf("Name:     %s\n", b.Name)
	fmt.Printf("Status:   %s\n", b.Status)
	fmt.Printf("Version:  %d\n", b.Version)
	if b.Intent != "" {
		fmt.Printf("Intent:   %s\n", b.Intent)
	}
	if b.Description != "" {
		fmt.Printf("Description: %s\n", b.Description)
	}
	fmt.Println()
}

func printBlueprintRunsText(runs []blueprintRun) {
	if len(runs) == 0 {
		fmt.Println("No runs found.")
		return
	}

	fmt.Printf("Runs (%d):\n\n", len(runs))
	for _, r := range runs {
		fmt.Printf("  %s  %s\n", r.ID, r.Status)
		if r.ResultSummary != "" {
			fmt.Printf("    %s\n", r.ResultSummary)
		}
		fmt.Println()
	}
}
