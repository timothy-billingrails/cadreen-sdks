package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show system health",
	Long: `Check if your Cadreen system is working.

Shows health status, memory, rules, and connections.
Use --verbose for full details. Use --json for scripting.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/health", nil)
		if err != nil {
			output.PrintError("Can't reach the API.")
			output.PrintHint(err.Error())
			output.PrintTry("cadreen status — to try again")
			return nil
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var health healthResponse
		if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(health, format)
		case output.FormatVerbose:
			printStatusVerbose(health)
		default:
			printStatusText(health)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

type healthResponse struct {
	Status      string          `json:"status"`
	Timestamp   string          `json:"timestamp"`
	Memory      memoryStatus    `json:"memory"`
	Connections connectionStatus `json:"connections"`
	Governance  governanceStatus `json:"governance"`
	Workspace   string          `json:"workspace"`
	TenantID    string          `json:"tenant_id"`
}

type memoryStatus struct {
	Status string `json:"status"`
}

type connectionStatus struct {
	Total   int `json:"total"`
	Healthy int `json:"healthy"`
}

type governanceStatus struct {
	Active int `json:"active"`
}

func printStatusText(h healthResponse) {
	if h.Status == "healthy" {
		fmt.Println("Everything is working.")
	} else {
		fmt.Printf("Status: %s\n", h.Status)
	}
	fmt.Println()

	details := []string{}
	if h.Memory.Status != "" {
		details = append(details, fmt.Sprintf("Memory: %s", humanMemoryStatus(h.Memory.Status)))
	}
	if h.Governance.Active > 0 {
		details = append(details, fmt.Sprintf("Rules: %d active", h.Governance.Active))
	}
	if h.Connections.Total > 0 {
		details = append(details, fmt.Sprintf("Connections: %d total, %d healthy", h.Connections.Total, h.Connections.Healthy))
	}

	for _, d := range details {
		fmt.Printf("  %s\n", d)
	}
}

func printStatusVerbose(h healthResponse) {
	fmt.Printf("Status:      %s\n", h.Status)
	fmt.Printf("Workspace:   %s\n", h.Workspace)
	fmt.Printf("Memory:      %s\n", h.Memory.Status)
	fmt.Printf("Governance:  %d active\n", h.Governance.Active)
	fmt.Printf("Connections: %d total, %d healthy\n", h.Connections.Total, h.Connections.Healthy)
	fmt.Printf("Timestamp:   %s\n", h.Timestamp)
}

func humanMemoryStatus(s string) string {
	switch s {
	case "healthy":
		return "working well"
	case "degraded":
		return "needs attention"
	case "not_configured":
		return "not configured"
	default:
		return s
	}
}

func handleHTTPError(resp *http.Response) error {
	limited := io.LimitReader(resp.Body, 64*1024) // 64KB cap
	body, err := io.ReadAll(limited)
	if err != nil {
		return fmt.Errorf("HTTP %d (failed to read body)", resp.StatusCode)
	}
	var apiErr struct {
		Message string `json:"message"`
		Hint    string `json:"hint"`
	}
	json.Unmarshal(body, &apiErr)

	if apiErr.Message != "" {
		output.PrintError(apiErr.Message)
		output.PrintHint(apiErr.Hint)
	} else {
		output.PrintError(fmt.Sprintf("HTTP %d", resp.StatusCode))
	}

	return fmt.Errorf("HTTP %d", resp.StatusCode)
}
