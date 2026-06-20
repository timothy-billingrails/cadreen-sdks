package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	cadreen "github.com/timothy-billingrails/cadreen-sdks/go/cadreen"

	"github.com/spf13/cobra"
)

var policiesCmd = &cobra.Command{
	Use:   "policies",
	Short: "Manage governance policies (list, evaluate)",
	Long: `Manage the rules Cadreen follows.

List active policies or test what would happen with a given action.`,
}

var policiesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List active rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/policies", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var policies cadreen.ListPoliciesResponse
		if err := json.NewDecoder(resp.Body).Decode(&policies); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(policies, format)
		case output.FormatVerbose:
			printPoliciesVerbose(policies.Policies)
		default:
			printPoliciesText(policies.Policies)
		}

		return nil
	},
}

var policiesEvaluateCmd = &cobra.Command{
	Use:   "evaluate [action]",
	Short: "Test what would happen",
	Long: `Test what would happen if Cadreen tried to do something.

Examples:
  cadreen policies evaluate "Refund $700 to customer 123"
  cadreen policies evaluate "Send email to john@example.com"`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		domain, _ := cmd.Flags().GetString("domain")

		client := newClient()
		resp, err := client.Evaluate(context.Background(), args[0], domain)
		if err != nil {
			return handleAPIError(err)
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(resp, format)
		case output.FormatVerbose:
			printEvaluateVerbose(resp)
		default:
			printEvaluateText(resp)
		}

		return nil
	},
}

func init() {
	policiesEvaluateCmd.Flags().String("domain", "", "policy domain to evaluate against")

	policiesCmd.AddCommand(policiesListCmd)
	policiesCmd.AddCommand(policiesEvaluateCmd)
	rootCmd.AddCommand(policiesCmd)
}

func printPoliciesText(policies []cadreen.Policy) {
	if len(policies) == 0 {
		fmt.Println("No active rules.")
		return
	}

	fmt.Printf("Cadreen follows %d rules:\n\n", len(policies))

	for _, p := range policies {
		fmt.Printf("  %s\n", p.Name)
		if p.Rationale != "" {
			fmt.Printf("  %s\n", p.Rationale)
		}
		fmt.Println()
	}
}

func printPoliciesVerbose(policies []cadreen.Policy) {
	for _, p := range policies {
		fmt.Printf("ID:       %s\n", p.ID)
		fmt.Printf("Name:     %s\n", p.Name)
		fmt.Printf("Domain:   %s\n", p.Domain)
		fmt.Printf("Priority: %d\n", p.Priority)
		if p.Rationale != "" {
			fmt.Printf("Rationale: %s\n", p.Rationale)
		}
		fmt.Println()
	}
}

func printEvaluateText(resp *cadreen.EvaluatePolicyResponse) {
	decision := resp.GovernanceResult.Type
	if decision == "auto" {
		fmt.Println("Cadreen would proceed on its own.")
	} else if decision == "handoff" {
		fmt.Println("Cadreen would ask you first.")
	} else {
		fmt.Printf("Result: %s\n", decision)
	}

	if resp.GovernanceResult.Reason != "" {
		fmt.Printf("\nBecause: %s\n", resp.GovernanceResult.Reason)
	}
	if resp.GovernanceResult.Confidence > 0 {
		fmt.Printf("Confidence: %.0f%%\n", resp.GovernanceResult.Confidence*100)
	}
}

func printEvaluateVerbose(resp *cadreen.EvaluatePolicyResponse) {
	fmt.Printf("Action:     %s\n", resp.Action)
	fmt.Printf("Domain:     %s\n", resp.Domain)
	fmt.Printf("Decision:   %s\n", resp.GovernanceResult.Type)
	fmt.Printf("Confidence: %.2f\n", resp.GovernanceResult.Confidence)
	fmt.Printf("Reason:     %s\n", resp.GovernanceResult.Reason)
}
