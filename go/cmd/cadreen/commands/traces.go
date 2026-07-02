package commands

import (
	"encoding/json"
	"fmt"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	cadreen "github.com/timothy-billingrails/cadreen-sdks/go/cadreen"

	"github.com/spf13/cobra"
)

var tracesCmd = &cobra.Command{
	Use:   "traces",
	Short: "Browse what happened",
	Long: `Browse recent activity.

Shows what Cadreen did, whether it handled things on its own
or asked for permission.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		limit, _ := cmd.Flags().GetInt("limit")
		filterType, _ := cmd.Flags().GetString("type")

		path := fmt.Sprintf("/api/v1/cadreen/intelligence?limit=%d", limit)
		if filterType != "" {
			path += "&type=" + filterType
		}

		resp, err := rawDo("GET", path, nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var traces cadreen.ListIntelligenceResponse
		if err := json.NewDecoder(resp.Body).Decode(&traces); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(traces, format)
		case output.FormatVerbose:
			printTracesVerbose(traces)
		default:
			printTracesText(traces)
		}

		return nil
	},
}

func init() {
	tracesCmd.Flags().Int("limit", 10, "max results")
	tracesCmd.Flags().String("type", "", "filter by type")
	rootCmd.AddCommand(tracesCmd)
}

func printTracesText(t cadreen.ListIntelligenceResponse) {
	if len(t.Traces) == 0 {
		fmt.Println("No activity yet.")
		return
	}

	fmt.Printf("What happened (%d requests):\n\n", len(t.Traces))

	for _, tr := range t.Traces {
		statusIcon := humanTraceStatus(tr.Meta.Governance.Decision)
		fmt.Printf("  ● %s\n", statusIcon)
		if tr.Meta.Summary != "" {
			fmt.Printf("    \"%s\"\n", tr.Meta.Summary)
		}
		if tr.Meta.Governance.Confidence > 0 {
			fmt.Printf("    Confidence: %.0f%%\n", tr.Meta.Governance.Confidence*100)
		}
		if tr.Meta.Process.DurationMs > 0 {
			fmt.Printf("    Duration: %.1fs\n", float64(tr.Meta.Process.DurationMs)/1000)
		}
		fmt.Println()
		fmt.Println()
	}
}

func printTracesVerbose(t cadreen.ListIntelligenceResponse) {
	for _, tr := range t.Traces {
		fmt.Printf("ID:         %s\n", tr.ID)
		fmt.Printf("Domain:     %s\n", tr.Domain)
		fmt.Printf("Decision:   %s\n", tr.Meta.Governance.Decision)
		fmt.Printf("Summary:    %s\n", tr.Meta.Summary)
		fmt.Printf("Confidence: %.2f\n", tr.Meta.Governance.Confidence)
		fmt.Printf("Duration:   %dms\n", tr.Meta.Process.DurationMs)
		fmt.Printf("Created:    %s\n", tr.CreatedAt)
		fmt.Println()
	}
}

func humanTraceStatus(decision string) string {
	switch decision {
	case "auto", "auto_complete":
		return "Handled on its own"
	case "handoff":
		return "Asked for permission"
	case "escalate":
		return "Needs a human"
	case "abstain":
		return "Declined to answer"
	default:
		if decision == "" {
			return "Handled on its own"
		}
		return decision
	}
}
