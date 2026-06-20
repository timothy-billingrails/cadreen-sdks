package commands

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/config"
	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	"github.com/spf13/cobra"
)

var (
	doctorOK   = color.New(color.FgGreen).SprintFunc()
	doctorFail = color.New(color.FgRed).SprintFunc()
	doctorTodo = color.New(color.FgYellow).SprintFunc()
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check if your setup is ready",
	Long: `Check if your Cadreen setup is ready to use.

Runs a series of checks: authentication, API connectivity,
memory, rules, and connections. Shows what's working and
what needs attention.`,
	RunE: runDoctor,
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

func runDoctor(cmd *cobra.Command, args []string) error {
	if cfg == nil {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return err
		}
	}

	fmt.Println("Checking your setup...")
	fmt.Println()

	authStatus, authOK := checkAuth()
	if authOK {
		fmt.Printf("  %s %-20s %s\n", doctorOK("✓"), "Authenticated", authStatus)
	} else {
		fmt.Printf("  %s %-20s %s\n", doctorFail("✗"), "Authenticated", authStatus)
		fmt.Printf("  %s %-20s %s\n", doctorTodo("○"), "Test request", "Send: cadreen ask \"hello\"")
		fmt.Printf("  %s %-20s %s\n", doctorTodo("○"), "Self-healing", "Not tested yet")
		fmt.Println()
		fmt.Printf("You're 0/4 ready. Fix the issues above first.\n")
		return nil
	}

	health, err := fetchHealth()
	if err != nil {
		fmt.Printf("  %s %-20s %s\n", doctorFail("✗"), "API reachable", err.Error())
		fmt.Printf("  %s %-20s %s\n", doctorTodo("○"), "Test request", "Send: cadreen ask \"hello\"")
		fmt.Printf("  %s %-20s %s\n", doctorTodo("○"), "Self-healing", "Not tested yet")
		fmt.Println()
		fmt.Printf("You're 1/4 ready. Fix the issues above first.\n")
		return nil
	}

	fmt.Printf("  %s %-20s %s\n", doctorOK("✓"), "API reachable", "Everything is working")

	passed := 1

	memStatus, memOK := checkMemoryFromHealth(health)
	if memOK {
		fmt.Printf("  %s %-20s %s\n", doctorOK("✓"), "Memory", memStatus)
		passed++
	} else {
		fmt.Printf("  %s %-20s %s\n", doctorFail("✗"), "Memory", memStatus)
	}

	govStatus, govOK := checkGovernanceFromHealth(health)
	if govOK {
		fmt.Printf("  %s %-20s %s\n", doctorOK("✓"), "Rules", govStatus)
		passed++
	} else {
		fmt.Printf("  %s %-20s %s\n", doctorFail("✗"), "Rules", govStatus)
	}

	connStatus, connOK := checkConnectionsFromHealth(health)
	if connOK {
		fmt.Printf("  %s %-20s %s\n", doctorOK("✓"), "Connections", connStatus)
		passed++
	} else {
		fmt.Printf("  %s %-20s %s\n", doctorFail("✗"), "Connections", connStatus)
	}

	fmt.Printf("  %s %-20s %s\n", doctorTodo("○"), "Test request", "Send: cadreen ask \"hello\"")
	fmt.Printf("  %s %-20s %s\n", doctorTodo("○"), "Self-healing", "Not tested yet")
	fmt.Println()

	total := 4 // API + memory + governance + connections
	if passed == total {
		fmt.Printf("You're %d/%d ready. Run 'cadreen ask \"hello\"' to test.\n", passed, total)
	} else {
		fmt.Printf("You're %d/%d ready. Fix the issues above first.\n", passed, total)
	}

	return nil
}

func checkAuth() (string, bool) {
	if !cfg.IsAuthenticated() {
		return "No API key configured", false
	}
	return output.MaskKey(cfg.APIKeyResolved()), true
}

func fetchHealth() (*healthResponse, error) {
	resp, err := rawDo("GET", "/api/v1/cadreen/health", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var health healthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return nil, err
	}

	return &health, nil
}

func checkMemoryFromHealth(h *healthResponse) (string, bool) {
	switch h.Memory.Status {
	case "healthy":
		return "Working well", true
	case "degraded":
		return "Needs attention", false
	case "not_configured":
		return "Not configured", false
	default:
		return h.Memory.Status, false
	}
}

func checkGovernanceFromHealth(h *healthResponse) (string, bool) {
	if h.Governance.Active > 0 {
		return fmt.Sprintf("%d active", h.Governance.Active), true
	}
	return "No active rules", false
}

func checkConnectionsFromHealth(h *healthResponse) (string, bool) {
	if h.Connections.Total == 0 {
		return "No connections", false
	}
	if h.Connections.Healthy == h.Connections.Total {
		return fmt.Sprintf("%d services, all healthy", h.Connections.Total), true
	}
	return fmt.Sprintf("%d/%d healthy", h.Connections.Healthy, h.Connections.Total), false
}
