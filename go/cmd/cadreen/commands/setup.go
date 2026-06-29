package commands

import (
	"fmt"

	cadreen "github.com/timothy-billingrails/cadreen-sdks/go/cadreen"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up your workspace",
	Long: `Set up and manage your Cadreen workspace.

Use 'setup run' for one-shot setup — describe what you want, and Cadreen
proposes connections, knowledge, and policies.

Use 'setup session' for incremental setup — accumulate resources over time,
then apply them atomically.

Examples:
  cadreen setup run --purpose "Customer support for an e-commerce store"
  cadreen setup run --purpose "Refund processing" --dry-run
  cadreen setup session create --purpose "Customer support"
  cadreen setup session list`,
}

func init() {
	setupCmd.AddCommand(setupRunCmd)
	setupCmd.AddCommand(setupSessionCmd)
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
