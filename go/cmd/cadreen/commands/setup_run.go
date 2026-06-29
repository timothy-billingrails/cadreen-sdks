package commands

import (
	"context"
	"fmt"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	cadreen "github.com/timothy-billingrails/cadreen-sdks/go/cadreen"

	"github.com/spf13/cobra"
)

var setupRunCmd = &cobra.Command{
	Use:   "run",
	Short: "One-shot setup — describe what you want, Cadreen proposes and applies",
	Long: `Describe what you want Cadreen to do. It proposes connections, knowledge, and policies.

Use --dry-run to preview without creating anything.

Examples:
  cadreen setup run --purpose "Customer support for an e-commerce store"
  cadreen setup run --purpose "Refund processing" --dry-run
  cadreen setup run --memory "Our refund window is 30 days" --memory "Orders over $500 need manager approval"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		purpose, _ := cmd.Flags().GetString("purpose")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		memoryItems, _ := cmd.Flags().GetStringSlice("memory")
		policyRules, _ := cmd.Flags().GetStringSlice("policy")

		if purpose == "" && len(memoryItems) == 0 && len(policyRules) == 0 {
			return fmt.Errorf("provide at least one of: --purpose, --memory, --policy")
		}

		req := cadreen.SetupRequest{
			Purpose: purpose,
			DryRun:  dryRun,
		}

		for _, m := range memoryItems {
			req.Memory = append(req.Memory, cadreen.SetupMemory{
				Type: "reference",
				Content: map[string]any{
					"text": m,
				},
			})
		}

		for _, p := range policyRules {
			req.Policies = append(req.Policies, cadreen.SetupPolicy{
				Name: p,
				Rule: p,
			})
		}

		client := newClient()
		resp, err := client.Setup(context.Background(), req)
		if err != nil {
			return handleAPIError(err)
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(resp, format)
			return nil
		}

		if dryRun {
			fmt.Println("Preview (nothing was created):")
		} else {
			fmt.Println("Setup complete.")
		}
		fmt.Println()

		if len(resp.Proposals) > 0 {
			fmt.Printf("Proposals (%d):\n", len(resp.Proposals))
			for _, p := range resp.Proposals {
				fmt.Printf("  [%s] %s\n", p.Type, p.Description)
				if p.Detail != "" {
					fmt.Printf("    %s\n", p.Detail)
				}
			}
			fmt.Println()
		}

		if resp.Notice != "" {
			fmt.Printf("Notice: %s\n\n", resp.Notice)
		}

		printSetupResults(resp)

		return nil
	},
}

func init() {
	setupRunCmd.Flags().String("purpose", "", "describe what you want Cadreen to do")
	setupRunCmd.Flags().Bool("dry-run", false, "preview without creating anything")
	setupRunCmd.Flags().StringSlice("memory", nil, "knowledge items to add (can specify multiple)")
	setupRunCmd.Flags().StringSlice("policy", nil, "policy rules to create (can specify multiple)")
}
