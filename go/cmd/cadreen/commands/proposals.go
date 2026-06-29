package commands

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	cadreen "github.com/timothy-billingrails/cadreen-sdks/go/cadreen"

	"github.com/spf13/cobra"
)

var proposalsCmd = &cobra.Command{
	Use:   "proposals",
	Short: "Task proposals",
	Long: `List, view, accept, dismiss, and get stats on task proposals.

Proposals are system-detected opportunities — Cadreen proposes, you decide.

Examples:
  cadreen proposals list
  cadreen proposals list --status proposed
  cadreen proposals accept prop_01abc
  cadreen proposals dismiss prop_01abc --reason "Not relevant"
  cadreen proposals stats`,
}

var proposalsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List proposals",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		status, _ := cmd.Flags().GetString("status")
		limit, _ := cmd.Flags().GetInt("limit")

		path := "/api/v1/cadreen/proposals"
		params := url.Values{}
		if status != "" {
			params.Set("status", status)
		}
		if limit > 0 {
			params.Set("limit", strconv.Itoa(limit))
		}
		if len(params) > 0 {
			path += "?" + params.Encode()
		}

		resp, err := rawDo("GET", path, nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result cadreen.ListProposalsResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(result, format)
		default:
			printProposalsList(result.Proposals)
		}

		return nil
	},
}

var proposalsGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get a proposal",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/proposals/"+args[0], nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var proposal cadreen.TaskProposal
		if err := json.NewDecoder(resp.Body).Decode(&proposal); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(proposal, format)
		default:
			printProposal(proposal)
		}

		return nil
	},
}

var proposalsAcceptCmd = &cobra.Command{
	Use:   "accept [id]",
	Short: "Accept a proposal",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("POST", "/api/v1/cadreen/proposals/"+args[0]+"/accept", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result cadreen.AcceptProposalResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(result, format)
		default:
			printAcceptProposal(result)
		}

		return nil
	},
}

var proposalsDismissCmd = &cobra.Command{
	Use:   "dismiss [id]",
	Short: "Dismiss a proposal",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		reason, _ := cmd.Flags().GetString("reason")

		body := map[string]string{}
		if reason != "" {
			body["reason"] = reason
		}

		resp, err := rawDo("POST", "/api/v1/cadreen/proposals/"+args[0]+"/dismiss", body)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result cadreen.DismissProposalResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(result, format)
		default:
			output.PrintSuccess("Proposal dismissed.")
		}

		return nil
	},
}

var proposalsStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Proposal statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/proposals/stats", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var stats cadreen.ProposalStatsResponse
		if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(stats, format)
		default:
			printProposalStats(stats)
		}

		return nil
	},
}

func init() {
	proposalsListCmd.Flags().String("status", "", "filter by status: proposed, accepted, dismissed, expired, all")
	proposalsListCmd.Flags().Int("limit", 0, "max results to return")
	proposalsDismissCmd.Flags().String("reason", "", "reason for dismissing")

	proposalsCmd.AddCommand(proposalsListCmd)
	proposalsCmd.AddCommand(proposalsGetCmd)
	proposalsCmd.AddCommand(proposalsAcceptCmd)
	proposalsCmd.AddCommand(proposalsDismissCmd)
	proposalsCmd.AddCommand(proposalsStatsCmd)
	rootCmd.AddCommand(proposalsCmd)
}

func printProposalsList(proposals []cadreen.TaskProposal) {
	if len(proposals) == 0 {
		fmt.Println("No proposals.")
		return
	}

	fmt.Printf("Proposals (%d):\n\n", len(proposals))
	for _, p := range proposals {
		statusLabel := p.Status
		if statusLabel == "" {
			statusLabel = "unknown"
		}
		fmt.Printf("  [%s] %s\n", statusLabel, p.Title)
		if p.ProposalType != "" {
			fmt.Printf("    Type: %s\n", p.ProposalType)
		}
		if p.Confidence > 0 {
			fmt.Printf("    Confidence: %.0f%%\n", p.Confidence*100)
		}
		fmt.Println()
	}
}

func printProposal(p cadreen.TaskProposal) {
	fmt.Println("Proposal")
	fmt.Println()
	fmt.Printf("  ID:          %s\n", p.ID)
	fmt.Printf("  Title:       %s\n", p.Title)
	if p.Description != "" {
		fmt.Printf("  Description: %s\n", p.Description)
	}
	if p.Intent != "" {
		fmt.Printf("  Intent:      %s\n", p.Intent)
	}
	if p.ProposalType != "" {
		fmt.Printf("  Type:        %s\n", p.ProposalType)
	}
	if p.Status != "" {
		fmt.Printf("  Status:      %s\n", p.Status)
	}
	if p.Confidence > 0 {
		fmt.Printf("  Confidence:  %.0f%%\n", p.Confidence*100)
	}
	if p.Domain != "" {
		fmt.Printf("  Domain:      %s\n", p.Domain)
	}
	if p.TriggerType != "" {
		fmt.Printf("  Trigger:     %s\n", p.TriggerType)
	}
	if p.TriggerSource != "" {
		fmt.Printf("  Source:      %s\n", p.TriggerSource)
	}
	if len(p.Evidence) > 0 {
		fmt.Printf("  Evidence:    %d item(s)\n", len(p.Evidence))
	}
	if p.CreatedAt != "" {
		fmt.Printf("  Created:     %s\n", p.CreatedAt)
	}
	if p.ExpiresAt != "" {
		fmt.Printf("  Expires:     %s\n", p.ExpiresAt)
	}
}

func printAcceptProposal(r cadreen.AcceptProposalResponse) {
	fmt.Println("Proposal Accepted")
	fmt.Println()
	fmt.Printf("  Status:   %s\n", r.Status)
	if r.Action != "" {
		fmt.Printf("  Action:   %s\n", r.Action)
	}
	if r.Intent != "" {
		fmt.Printf("  Intent:   %s\n", r.Intent)
	}
	if r.ExecutionID != "" {
		fmt.Printf("  Execution: %s\n", r.ExecutionID)
	}
	if r.NextStep != "" {
		fmt.Printf("  Next:     %s\n", r.NextStep)
	}
	if r.AutoApproved != nil && *r.AutoApproved {
		fmt.Println("  Auto:     yes")
	}
}

func printProposalStats(stats cadreen.ProposalStatsResponse) {
	fmt.Println("Proposal Statistics")
	fmt.Println()
	total := stats.Proposed + stats.Accepted + stats.Dismissed + stats.Expired
	if total > 0 {
		fmt.Printf("  Total:     %d\n", total)
	}
	if stats.Proposed > 0 {
		fmt.Printf("  Proposed:  %d\n", stats.Proposed)
	}
	if stats.Accepted > 0 {
		fmt.Printf("  Accepted:  %d\n", stats.Accepted)
	}
	if stats.Dismissed > 0 {
		fmt.Printf("  Dismissed: %d\n", stats.Dismissed)
	}
	if stats.Expired > 0 {
		fmt.Printf("  Expired:   %d\n", stats.Expired)
	}
	if total == 0 {
		fmt.Println("  No proposals yet.")
	}
}


