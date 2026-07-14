package commands

import (
	"encoding/json"
	"fmt"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	"github.com/spf13/cobra"
)

var federationCmd = &cobra.Command{
	Use:   "federation",
	Short: "Manage federation links between workspaces",
	Long: `Federation lets workspaces share agents and knowledge.

Create a link, approve it, then share agents across workspace boundaries.`,
}

var federationCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a federation link",
	Long: `Create a federation link to another workspace.

Examples:
  cadreen federation create --target-workspace ws_01JX`,
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		targetWorkspace, _ := cmd.Flags().GetString("target-workspace")
		if targetWorkspace == "" {
			return fmt.Errorf("--target-workspace is required")
		}

		payload := map[string]any{
			"target_workspace_id": targetWorkspace,
		}

		resp, err := rawDo("POST", "/api/v1/cadreen/federation", payload)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 && resp.StatusCode != 201 {
			return handleHTTPError(resp)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
		} else {
			fmt.Printf("Federation link created: %s\n", result["id"])
			if s, ok := result["status"]; ok {
				fmt.Printf("  Status: %s\n", s)
			}
		}

		return nil
	},
}

var federationListCmd = &cobra.Command{
	Use:   "list",
	Short: "List federation links",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/federation", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result struct {
			Links []map[string]any `json:"links"`
			Count int              `json:"count"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(result, format)
		default:
			if len(result.Links) == 0 {
				fmt.Println("No federation links.")
				return nil
			}
			fmt.Printf("Federation links (%d):\n\n", len(result.Links))
			for _, l := range result.Links {
				fmt.Printf("  %s  %s\n", l["id"], l["status"])
				if t, ok := l["targetWorkspaceId"]; ok {
					fmt.Printf("    Workspace: %s\n", t)
				}
				fmt.Println()
			}
		}

		return nil
	},
}

var federationGetCmd = &cobra.Command{
	Use:   "get [federation-id]",
	Short: "Get federation link details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", fmt.Sprintf("/api/v1/cadreen/federation/%s", args[0]), nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
		} else {
			printFederationDetail(result)
		}

		return nil
	},
}

var federationApproveCmd = &cobra.Command{
	Use:   "approve [federation-id]",
	Short: "Approve a federation link",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("POST", fmt.Sprintf("/api/v1/cadreen/federation/%s/approve", args[0]), nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
		} else {
			fmt.Println("Federation link approved.")
		}

		return nil
	},
}

var federationSuspendCmd = &cobra.Command{
	Use:   "suspend [federation-id]",
	Short: "Suspend a federation link",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		reason, _ := cmd.Flags().GetString("reason")

		payload := map[string]any{}
		if reason != "" {
			payload["reason"] = reason
		}

		resp, err := rawDo("POST", fmt.Sprintf("/api/v1/cadreen/federation/%s/suspend", args[0]), payload)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
		} else {
			fmt.Println("Federation link suspended.")
		}

		return nil
	},
}

var federationRevokeCmd = &cobra.Command{
	Use:   "revoke [federation-id]",
	Short: "Revoke a federation link",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		reason, _ := cmd.Flags().GetString("reason")

		payload := map[string]any{}
		if reason != "" {
			payload["reason"] = reason
		}

		resp, err := rawDo("POST", fmt.Sprintf("/api/v1/cadreen/federation/%s/revoke", args[0]), payload)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 && resp.StatusCode != 204 {
			return handleHTTPError(resp)
		}

		fmt.Println("Federation link revoked.")
		return nil
	},
}

var federationPermissionsCmd = &cobra.Command{
	Use:   "permissions [federation-id]",
	Short: "See what's shared",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", fmt.Sprintf("/api/v1/cadreen/federation/%s/permissions", args[0]), nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
		} else {
			if perms, ok := result["permissions"].([]any); ok {
				fmt.Printf("Shared permissions (%d):\n\n", len(perms))
				for _, p := range perms {
					fmt.Printf("  - %s\n", p)
				}
			} else {
				output.Print(result, output.FormatJSON)
			}
		}

		return nil
	},
}

var federationLinkAgentCmd = &cobra.Command{
	Use:   "link-agent [federation-id]",
	Short: "Share an agent with the linked workspace",
	Long: `Link a local agent to a remote agent in the federated workspace.

Examples:
  cadreen federation link-agent fed_01JX --local-agent ag_01KY --remote-agent ag_02LZ`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		localAgent, _ := cmd.Flags().GetString("local-agent")
		remoteAgent, _ := cmd.Flags().GetString("remote-agent")

		if localAgent == "" {
			return fmt.Errorf("--local-agent is required")
		}
		if remoteAgent == "" {
			return fmt.Errorf("--remote-agent is required")
		}

		payload := map[string]any{
			"local_agent_id":  localAgent,
			"remote_agent_id": remoteAgent,
		}

		resp, err := rawDo("POST", fmt.Sprintf("/api/v1/cadreen/federation/%s/agents", args[0]), payload)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 && resp.StatusCode != 201 {
			return handleHTTPError(resp)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(result, format)
		} else {
			fmt.Println("Agent linked.")
		}

		return nil
	},
}

var federationAgentsCmd = &cobra.Command{
	Use:   "agents [federation-id]",
	Short: "See linked agents",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", fmt.Sprintf("/api/v1/cadreen/federation/%s/agents", args[0]), nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result struct {
			Agents []map[string]any `json:"agents"`
			Count  int              `json:"count"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(result, format)
		default:
			if len(result.Agents) == 0 {
				fmt.Println("No linked agents.")
				return nil
			}
			fmt.Printf("Linked agents (%d):\n\n", len(result.Agents))
			for _, a := range result.Agents {
			fmt.Printf("  Local:  %s\n", a["localAgentId"])
			fmt.Printf("  Remote: %s\n", a["remoteAgentId"])
				if s, ok := a["status"]; ok {
					fmt.Printf("  Status: %s\n", s)
				}
				fmt.Println()
			}
		}

		return nil
	},
}

func init() {
	federationCreateCmd.Flags().String("target-workspace", "", "target workspace ID (required)")

	federationSuspendCmd.Flags().String("reason", "", "why you're suspending this link")

	federationRevokeCmd.Flags().String("reason", "", "why you're revoking this link")

	federationLinkAgentCmd.Flags().String("local-agent", "", "local agent ID (required)")
	federationLinkAgentCmd.Flags().String("remote-agent", "", "remote agent ID (required)")

	federationCmd.AddCommand(federationCreateCmd)
	federationCmd.AddCommand(federationListCmd)
	federationCmd.AddCommand(federationGetCmd)
	federationCmd.AddCommand(federationApproveCmd)
	federationCmd.AddCommand(federationSuspendCmd)
	federationCmd.AddCommand(federationRevokeCmd)
	federationCmd.AddCommand(federationPermissionsCmd)
	federationCmd.AddCommand(federationLinkAgentCmd)
	federationCmd.AddCommand(federationAgentsCmd)
	rootCmd.AddCommand(federationCmd)
}

func printFederationDetail(l map[string]any) {
	fmt.Printf("ID:      %s\n", l["id"])
	if s, ok := l["status"]; ok {
		fmt.Printf("Status:  %s\n", s)
	}
	if t, ok := l["targetWorkspaceId"]; ok {
		fmt.Printf("Target:  %s\n", t)
	}
	if c, ok := l["created_at"]; ok {
		fmt.Printf("Created: %s\n", c)
	}
	fmt.Println()
}
