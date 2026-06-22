package commands

import (
	"context"
	"fmt"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	cadreen "github.com/timothy-billingrails/cadreen-sdks/go/cadreen"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up your workspace",
	Long: `Describe what you want Cadreen to do. It proposes connections, knowledge, and policies.

Use --dry-run to preview without creating anything.

Examples:
  cadreen setup --purpose "Customer support for an e-commerce store"
  cadreen setup --purpose "Refund processing" --dry-run
  cadreen setup --memory "Our refund window is 30 days" --memory "Orders over $500 need manager approval"`,
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
	setupCmd.Flags().String("purpose", "", "describe what you want Cadreen to do")
	setupCmd.Flags().Bool("dry-run", false, "preview without creating anything")
	setupCmd.Flags().StringSlice("memory", nil, "knowledge items to add (can specify multiple)")
	setupCmd.Flags().StringSlice("policy", nil, "policy rules to create (can specify multiple)")

	rootCmd.AddCommand(setupCmd)
}

func printSetupResults(resp *cadreen.SetupResult) {
	total := resp.Applied + resp.Failed
	if total == 0 {
		return
	}

	if resp.DryRun {
		fmt.Printf("Would apply: %d item(s)\n", resp.Applied)
	} else {
		fmt.Printf("Applied: %d | Failed: %d\n", resp.Applied, resp.Failed)
	}

	for _, c := range resp.Connections {
		status := c.Status
		if status == "would_create" {
			status = "would connect"
		}
		fmt.Printf("  Connection %s: %s\n", c.Capability, status)
	}
	for _, c := range resp.Credentials {
		fmt.Printf("  Credential %s/%s: %s\n", c.Provider, c.Name, c.Status)
	}
	for _, m := range resp.Memory {
		fmt.Printf("  Memory %s: %s\n", m.ID, m.Status)
	}
	for _, p := range resp.Policies {
		fmt.Printf("  Policy %s: %s\n", p.Name, p.Status)
	}
}
