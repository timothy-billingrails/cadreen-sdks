package commands

import (
	"encoding/json"
	"fmt"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	cadreen "github.com/timothy-billingrails/cadreen-sdks/go/cadreen"

	"github.com/spf13/cobra"
)

var escalationsCmd = &cobra.Command{
	Use:   "escalations",
	Short: "Manage escalations (list, get, resolve)",
	Long:  `View and resolve escalations that need human attention.`,
}

var escalationsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List escalations",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/escalations", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var escalations cadreen.ListEscalationsResponse
		if err := json.NewDecoder(resp.Body).Decode(&escalations); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(escalations, format)
		case output.FormatVerbose:
			printEscalationsVerbose(escalations.Escalations)
		default:
			printEscalationsText(escalations.Escalations)
		}

		return nil
	},
}

var escalationsGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get an escalation by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/escalations/"+args[0], nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var esc cadreen.Escalation
		if err := json.NewDecoder(resp.Body).Decode(&esc); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(esc, format)
		default:
			printEscalationDetail(esc)
		}

		return nil
	},
}

var escalationsResolveCmd = &cobra.Command{
	Use:   "resolve [id] [resolution]",
	Short: "Resolve an escalation",
	Long: `Resolve an escalation with a resolution message.

Examples:
  cadreen escalations resolve esc_123 "Approved refund of $50"
  cadreen escalations resolve esc_456 "Denied — exceeds policy limit"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		client := newClient()
		esc, err := client.ResolveEscalation(cmd.Context(), args[0], args[1])
		if err != nil {
			return handleAPIError(err)
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(esc, format)
		} else {
			fmt.Printf("Escalation %s resolved.\n", esc.ID)
		}

		return nil
	},
}

func init() {
	escalationsCmd.AddCommand(escalationsListCmd)
	escalationsCmd.AddCommand(escalationsGetCmd)
	escalationsCmd.AddCommand(escalationsResolveCmd)
	rootCmd.AddCommand(escalationsCmd)
}

func printEscalationsText(escs []cadreen.Escalation) {
	if len(escs) == 0 {
		fmt.Println("No escalations.")
		return
	}

	fmt.Printf("Escalations (%d):\n\n", len(escs))
	for _, e := range escs {
		statusIcon := "!"
		if e.Status == "resolved" {
			statusIcon = "✓"
		}
		fmt.Printf("  %s %s\n", statusIcon, e.ID)
		if e.Intent != "" {
			fmt.Printf("    %s\n", e.Intent)
		}
		if e.Severity != "" {
			fmt.Printf("    Severity: %s\n", e.Severity)
		}
		fmt.Println()
	}
}

func printEscalationsVerbose(escs []cadreen.Escalation) {
	for _, e := range escs {
		printEscalationDetail(e)
	}
}

func printEscalationDetail(e cadreen.Escalation) {
	fmt.Printf("ID:       %s\n", e.ID)
	fmt.Printf("Status:   %s\n", e.Status)
	if e.Intent != "" {
		fmt.Printf("Intent:   %s\n", e.Intent)
	}
	if e.Severity != "" {
		fmt.Printf("Severity: %s\n", e.Severity)
	}
	if e.ErrorMessage != "" {
		fmt.Printf("Error:    %s\n", e.ErrorMessage)
	}
	if e.HumanPrompt != "" {
		fmt.Printf("Prompt:   %s\n", e.HumanPrompt)
	}
	if e.Resolution != "" {
		fmt.Printf("Resolution: %s\n", e.Resolution)
	}
	fmt.Println()
}
