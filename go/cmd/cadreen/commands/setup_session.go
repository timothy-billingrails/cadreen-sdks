package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	cadreen "github.com/timothy-billingrails/cadreen-sdks/go/cadreen"

	"github.com/spf13/cobra"
)

var setupSessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Manage setup sessions",
	Long: `Create and manage stateful setup sessions. Accumulate resources incrementally,
then apply them atomically.

Examples:
  cadreen setup session create --purpose "Customer support"
  cadreen setup session list
  cadreen setup session get sess_01abc
  cadreen setup session add sess_01abc --memory "Refund policy: 30 days"
  cadreen setup session apply sess_01abc --confirm`,
}

var setupSessionCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new setup session",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		purpose, _ := cmd.Flags().GetString("purpose")
		constraints, _ := cmd.Flags().GetStringSlice("constraints")

		req := cadreen.SetupSessionCreateRequest{
			Purpose:     purpose,
			Constraints: constraints,
		}

		client := newClient()
		resp, err := client.CreateSetupSession(context.Background(), req)
		if err != nil {
			return handleAPIError(err)
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(resp, format)
			return nil
		}

		fmt.Printf("Session created: %s\n", resp.ID)
		fmt.Printf("  Status: %s\n", resp.Status)
		if resp.Purpose != "" {
			fmt.Printf("  Purpose: %s\n", resp.Purpose)
		}
		fmt.Printf("  Created: %s\n", resp.CreatedAt)
		return nil
	},
}

var setupSessionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List setup sessions",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		client := newClient()
		resp, err := client.ListSetupSessions(context.Background())
		if err != nil {
			return handleAPIError(err)
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(resp, format)
			return nil
		}

		if len(resp.Sessions) == 0 {
			fmt.Println("No setup sessions found.")
			return nil
		}

		fmt.Printf("Setup Sessions (%d):\n\n", len(resp.Sessions))
		for _, s := range resp.Sessions {
			fmt.Printf("  %s  %-10s  %s\n", s.ID, s.Status, s.Purpose)
		}
		return nil
	},
}

var setupSessionGetCmd = &cobra.Command{
	Use:   "get [session-id]",
	Short: "Get a setup session by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		client := newClient()
		resp, err := client.GetSetupSession(context.Background(), args[0])
		if err != nil {
			return handleAPIError(err)
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(resp, format)
			return nil
		}

		printSetupSession(resp)
		return nil
	},
}

var setupSessionAddCmd = &cobra.Command{
	Use:   "add [session-id]",
	Short: "Add resources to a draft session",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		memoryItems, _ := cmd.Flags().GetStringSlice("memory")
		policyRules, _ := cmd.Flags().GetStringSlice("policy")

		if len(memoryItems) == 0 && len(policyRules) == 0 {
			return fmt.Errorf("provide at least one of: --memory, --policy")
		}

		req := cadreen.SetupSessionAddRequest{}

		for _, m := range memoryItems {
			req.Memory = append(req.Memory, cadreen.SetupMemory{
				Type: "reference",
				Content: map[string]any{
					"text": m,
				},
			})
		}

		for _, p := range policyRules {
			parts := strings.SplitN(p, ":", 2)
			name := parts[0]
			rule := p
			if len(parts) == 2 {
				rule = strings.TrimSpace(parts[1])
			}
			req.Policies = append(req.Policies, cadreen.SetupPolicy{
				Name: name,
				Rule: rule,
			})
		}

		client := newClient()
		resp, err := client.AddToSetupSession(context.Background(), args[0], req)
		if err != nil {
			return handleAPIError(err)
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(resp, format)
			return nil
		}

		fmt.Println("Resources added.")
		printSetupSession(resp)
		return nil
	},
}

var setupSessionApplyCmd = &cobra.Command{
	Use:   "apply [session-id]",
	Short: "Apply a draft session atomically",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		confirm, _ := cmd.Flags().GetBool("confirm")
		if !confirm {
			return fmt.Errorf("--confirm is required to apply a session (safety catch)")
		}

		client := newClient()
		resp, err := client.ApplySetupSession(context.Background(), args[0], cadreen.SetupSessionApplyRequest{
			Confirm: true,
		})
		if err != nil {
			return handleAPIError(err)
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(resp, format)
			return nil
		}

		fmt.Printf("Session %s: %s\n", resp.SessionID, resp.Status)
		fmt.Printf("  Applied: %d | Failed: %d\n", resp.Applied, resp.Failed)
		return nil
	},
}

func init() {
	setupSessionCreateCmd.Flags().String("purpose", "", "describe what you want Cadreen to do")
	setupSessionCreateCmd.Flags().StringSlice("constraints", nil, "rules Cadreen must follow")

	setupSessionAddCmd.Flags().StringSlice("memory", nil, "knowledge items to add (can specify multiple)")
	setupSessionAddCmd.Flags().StringSlice("policy", nil, "policy rules to create (format: 'name:rule')")

	setupSessionApplyCmd.Flags().Bool("confirm", false, "confirm application (required)")

	setupSessionCmd.AddCommand(setupSessionCreateCmd)
	setupSessionCmd.AddCommand(setupSessionListCmd)
	setupSessionCmd.AddCommand(setupSessionGetCmd)
	setupSessionCmd.AddCommand(setupSessionAddCmd)
	setupSessionCmd.AddCommand(setupSessionApplyCmd)
}

func printSetupSession(s *cadreen.SetupSession) {
	fmt.Printf("Session: %s\n", s.ID)
	fmt.Printf("  Status:      %s\n", s.Status)
	if s.Purpose != "" {
		fmt.Printf("  Purpose:     %s\n", s.Purpose)
	}
	if len(s.Constraints) > 0 {
		fmt.Printf("  Constraints: %s\n", strings.Join(s.Constraints, ", "))
	}
	fmt.Printf("  Created:     %s\n", s.CreatedAt)
	fmt.Printf("  Updated:     %s\n", s.UpdatedAt)

	total := len(s.Connections) + len(s.Credentials) + len(s.Memory) + len(s.Policies)
	fmt.Printf("\n  Resources (%d):\n", total)
	for _, c := range s.Connections {
		fmt.Printf("    Connection: %s\n", c.Capability)
	}
	for _, c := range s.Credentials {
		fmt.Printf("    Credential: %s/%s\n", c.Provider, c.Name)
	}
	for _, m := range s.Memory {
		fmt.Printf("    Memory: %s\n", m.Type)
	}
	for _, p := range s.Policies {
		fmt.Printf("    Policy: %s\n", p.Name)
	}

	if len(s.Proposals) > 0 {
		fmt.Printf("\n  Proposals (%d):\n", len(s.Proposals))
		for _, p := range s.Proposals {
			fmt.Printf("    [%s] %s\n", p.Type, p.Description)
		}
	}

	if s.AppliedCount > 0 || s.FailedCount > 0 {
		fmt.Printf("\n  Applied: %d | Failed: %d\n", s.AppliedCount, s.FailedCount)
	}
}
