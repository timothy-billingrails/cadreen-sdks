package commands

import (
	"encoding/json"
	"fmt"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	"github.com/spf13/cobra"
)

type scheduleSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	BlueprintID string `json:"blueprint_id"`
	Status      string `json:"status"`
	NextRunAt   string `json:"next_run_at,omitempty"`
	LastRunAt   string `json:"last_run_at,omitempty"`
	PauseReason string `json:"pause_reason,omitempty"`
	CreatedAt   string `json:"created_at"`
}

var schedulesCmd = &cobra.Command{
	Use:   "schedules",
	Short: "Manage scheduled blueprints (list, create, pause, resume)",
	Long: `Schedules run blueprints automatically on a timer.

Attach a schedule to a blueprint and it runs daily, weekly, or on a cron expression.`,
}

var schedulesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List schedules",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/schedules", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result struct {
			Schedules []scheduleSummary `json:"schedules"`
			Count     int               `json:"count"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(result, format)
		default:
			printSchedulesText(result.Schedules)
		}

		return nil
	},
}

var schedulesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a schedule",
	Long: `Attach a schedule to a blueprint.

Examples:
  cadreen schedules create --blueprint bp_01JX --name "Daily report" --cron "0 9 * * 1-5"
  cadreen schedules create --blueprint bp_01JX --name "Weekly digest" --frequency weekly --time 09:00`,
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		blueprintID, _ := cmd.Flags().GetString("blueprint")
		name, _ := cmd.Flags().GetString("name")
		cronExpr, _ := cmd.Flags().GetString("cron")
		frequency, _ := cmd.Flags().GetString("frequency")
		time, _ := cmd.Flags().GetString("time")
		timezone, _ := cmd.Flags().GetString("timezone")

		if blueprintID == "" {
			return fmt.Errorf("--blueprint is required")
		}
		if name == "" {
			return fmt.Errorf("--name is required")
		}

		trigger := map[string]any{}
		if cronExpr != "" {
			trigger["type"] = "cron"
			trigger["expression"] = cronExpr
		} else if frequency != "" {
			trigger["type"] = "recurring"
			trigger["frequency"] = frequency
			if time != "" {
				trigger["time"] = time
			}
		} else {
			return fmt.Errorf("provide either --cron or --frequency")
		}

		payload := map[string]any{
			"blueprint_id": blueprintID,
			"name":         name,
			"trigger":      trigger,
		}
		if timezone != "" {
			payload["timezone"] = timezone
		}

		resp, err := rawDo("POST", "/api/v1/cadreen/schedules", payload)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 && resp.StatusCode != 201 {
			return handleHTTPError(resp)
		}

		var result scheduleSummary
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
		} else {
			fmt.Printf("Schedule created: %s\n", result.ID)
			fmt.Printf("  Name:      %s\n", result.Name)
			fmt.Printf("  Status:    %s\n", result.Status)
			if result.NextRunAt != "" {
				fmt.Printf("  Next run:  %s\n", result.NextRunAt)
			}
		}

		return nil
	},
}

var schedulesPauseCmd = &cobra.Command{
	Use:   "pause [schedule-id]",
	Short: "Pause a schedule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		reason, _ := cmd.Flags().GetString("reason")

		payload := map[string]any{}
		if reason != "" {
			payload["reason"] = reason
		}

		resp, err := rawDo("POST", fmt.Sprintf("/api/v1/cadreen/schedules/%s/pause", args[0]), payload)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		fmt.Printf("Schedule paused: %s\n", result.ID)
		return nil
	},
}

var schedulesResumeCmd = &cobra.Command{
	Use:   "resume [schedule-id]",
	Short: "Resume a paused schedule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("POST", fmt.Sprintf("/api/v1/cadreen/schedules/%s/resume", args[0]), nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result struct {
			ID        string `json:"id"`
			Status    string `json:"status"`
			NextRunAt string `json:"next_run_at"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		fmt.Printf("Schedule resumed: %s\n", result.ID)
		if result.NextRunAt != "" {
			fmt.Printf("  Next run: %s\n", result.NextRunAt)
		}
		return nil
	},
}

var schedulesShowCmd = &cobra.Command{
	Use:   "show [schedule-id]",
	Short: "View schedule details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", fmt.Sprintf("/api/v1/cadreen/schedules/%s", args[0]), nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result scheduleSummary
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
		} else {
			printScheduleDetail(&result)
		}

		return nil
	},
}

func init() {
	schedulesCreateCmd.Flags().String("blueprint", "", "blueprint ID to schedule (required)")
	schedulesCreateCmd.Flags().String("name", "", "schedule name (required)")
	schedulesCreateCmd.Flags().String("cron", "", "cron expression (e.g. '0 9 * * 1-5')")
	schedulesCreateCmd.Flags().String("frequency", "", "frequency: daily, weekly, monthly, weekdays")
	schedulesCreateCmd.Flags().String("time", "", "time for recurring (e.g. '09:00')")
	schedulesCreateCmd.Flags().String("timezone", "UTC", "timezone (e.g. 'America/New_York')")

	schedulesPauseCmd.Flags().String("reason", "", "why you're pausing")

	schedulesCmd.AddCommand(schedulesListCmd)
	schedulesCmd.AddCommand(schedulesCreateCmd)
	schedulesCmd.AddCommand(schedulesPauseCmd)
	schedulesCmd.AddCommand(schedulesResumeCmd)
	schedulesCmd.AddCommand(schedulesShowCmd)
	rootCmd.AddCommand(schedulesCmd)
}

func printSchedulesText(schedules []scheduleSummary) {
	if len(schedules) == 0 {
		fmt.Println("No schedules found.")
		return
	}

	fmt.Printf("Schedules (%d):\n\n", len(schedules))
	for _, s := range schedules {
		fmt.Printf("  %s  %s\n", s.ID, s.Name)
		fmt.Printf("    Status: %s", s.Status)
		if s.NextRunAt != "" {
			fmt.Printf(" | Next run: %s", s.NextRunAt)
		}
		fmt.Println()
		fmt.Println()
	}
}

func printScheduleDetail(s *scheduleSummary) {
	fmt.Printf("ID:         %s\n", s.ID)
	fmt.Printf("Name:       %s\n", s.Name)
	fmt.Printf("Status:     %s\n", s.Status)
	fmt.Printf("Blueprint:  %s\n", s.BlueprintID)
	if s.NextRunAt != "" {
		fmt.Printf("Next run:   %s\n", s.NextRunAt)
	}
	if s.LastRunAt != "" {
		fmt.Printf("Last run:   %s\n", s.LastRunAt)
	}
	if s.PauseReason != "" {
		fmt.Printf("Pause reason: %s\n", s.PauseReason)
	}
	fmt.Println()
}
