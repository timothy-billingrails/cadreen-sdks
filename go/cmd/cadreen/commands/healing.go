package commands

import (
	"encoding/json"
	"fmt"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	cadreen "github.com/timothy-billingrails/cadreen-sdks/go/cadreen"

	"github.com/spf13/cobra"
)

var healingCmd = &cobra.Command{
	Use:   "healing",
	Short: "Self-healing stats and diagnostics",
	Long:  `View healing statistics, precedents, and diagnose failures.`,
}

var healingStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Healing statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/healing/stats", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var stats cadreen.HealingStatsResponse
		if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(stats, format)
		default:
			printHealingStats(stats)
		}

		return nil
	},
}

var healingPrecedentsCmd = &cobra.Command{
	Use:   "precedents",
	Short: "List healing precedents",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/healing/precedents", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result struct {
			Precedents []cadreen.HealingPrecedent `json:"precedents"`
			Count      int                        `json:"count"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(result, format)
		case output.FormatVerbose:
			printPrecedentsVerbose(result.Precedents)
		default:
			printPrecedentsText(result.Precedents)
		}

		return nil
	},
}

var healingDiagnoseCmd = &cobra.Command{
	Use:   "diagnose [error]",
	Short: "Diagnose a failure",
	Long: `Diagnose why something failed and get recovery suggestions.

Examples:
  cadreen healing diagnose "connection timeout to stripe API"
  cadreen healing diagnose "403 forbidden" --tool stripe_refund`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		toolName, _ := cmd.Flags().GetString("tool")

		body := map[string]string{"error_message": args[0]}
		if toolName != "" {
			body["tool_name"] = toolName
		}

		resp, err := rawDo("POST", "/api/v1/cadreen/healing/diagnose", body)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var diag cadreen.HealingDiagnosis
		if err := json.NewDecoder(resp.Body).Decode(&diag); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(diag, format)
		default:
			printDiagnosis(diag)
		}

		return nil
	},
}

func init() {
	healingDiagnoseCmd.Flags().String("tool", "", "tool name that failed")

	healingCmd.AddCommand(healingStatsCmd)
	healingCmd.AddCommand(healingPrecedentsCmd)
	healingCmd.AddCommand(healingDiagnoseCmd)
	rootCmd.AddCommand(healingCmd)
}

func printHealingStats(stats cadreen.HealingStatsResponse) {
	fmt.Println("Healing Statistics")
	fmt.Println()
	if stats.TotalPrecedents > 0 {
		fmt.Printf("  Total precedents: %d\n", stats.TotalPrecedents)
	}
	if stats.SuccessfulRecoveries > 0 {
		fmt.Printf("  Successful: %d\n", stats.SuccessfulRecoveries)
	}
	if stats.FailedRecoveries > 0 {
		fmt.Printf("  Failed: %d\n", stats.FailedRecoveries)
	}
	if stats.SuccessRate > 0 {
		fmt.Printf("  Success rate: %.0f%%\n", stats.SuccessRate*100)
	}
}

func printPrecedentsText(precedents []cadreen.HealingPrecedent) {
	if len(precedents) == 0 {
		fmt.Println("No healing precedents.")
		return
	}

	fmt.Printf("Healing precedents (%d):\n\n", len(precedents))
	for _, p := range precedents {
		status := "✗"
		if p.Success {
			status = "✓"
		}
		fmt.Printf("  %s %s\n", status, p.ErrorType)
		if p.ToolName != "" {
			fmt.Printf("    Tool: %s\n", p.ToolName)
		}
		if p.RecoveryStrategy != "" {
			fmt.Printf("    Strategy: %s\n", p.RecoveryStrategy)
		}
		fmt.Println()
	}
}

func printPrecedentsVerbose(precedents []cadreen.HealingPrecedent) {
	for _, p := range precedents {
		fmt.Printf("ID:        %s\n", p.ID)
		fmt.Printf("Error:     %s\n", p.ErrorType)
		fmt.Printf("Success:   %v\n", p.Success)
		fmt.Printf("Attempts:  %d\n", p.Attempts)
		if p.ToolName != "" {
			fmt.Printf("Tool:      %s\n", p.ToolName)
		}
		if p.RecoveryStrategy != "" {
			fmt.Printf("Strategy:  %s\n", p.RecoveryStrategy)
		}
		if p.WhatWorked != "" {
			fmt.Printf("Worked:    %s\n", p.WhatWorked)
		}
		fmt.Println()
	}
}

func printDiagnosis(diag cadreen.HealingDiagnosis) {
	fmt.Println("Diagnosis")
	fmt.Println()
	if diag.ErrorCategory != "" {
		fmt.Printf("  Category:    %s\n", diag.ErrorCategory)
	}
	if diag.RootCause != "" {
		fmt.Printf("  Root cause:  %s\n", diag.RootCause)
	}
	if diag.SemanticReason != "" {
		fmt.Printf("  Reason:      %s\n", diag.SemanticReason)
	}
	if diag.CanRetry {
		fmt.Println("  Can retry:   yes")
	}
	if diag.NeedsHuman {
		fmt.Println("  Needs human: yes")
	}
	if diag.Confidence > 0 {
		fmt.Printf("  Confidence:  %.0f%%\n", diag.Confidence*100)
	}
}
