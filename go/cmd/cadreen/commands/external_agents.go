package commands

import (
	"encoding/json"
	"fmt"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	"github.com/spf13/cobra"
)

var externalAgentsCmd = &cobra.Command{
	Use:   "external-agents",
	Short: "Manage external A2A agent connections",
	Long: `Connect to agents from other systems using the A2A protocol.

External agents are agents built on other platforms (LangChain, CrewAI, etc.)
that can send tasks to your Cadreen agents.`,
}

var externalAgentsConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to an external agent",
	Long: `Connect to an external A2A agent by providing its Agent Card URL.

The connection starts in pending_approval status and must be approved
before it becomes active.

Examples:
  cadreen external-agents connect --agent-id agent_123 --url https://example.com/.well-known/agent.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		agentID, _ := cmd.Flags().GetString("agent-id")
		url, _ := cmd.Flags().GetString("url")

		if agentID == "" {
			return fmt.Errorf("--agent-id is required")
		}
		if url == "" {
			return fmt.Errorf("--url is required")
		}

		payload := map[string]any{"agentCardUrl": url}
		path := fmt.Sprintf("/api/v1/cadreen/agents/%s/external", agentID)

		resp, err := rawDo("POST", path, payload)
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
			fmt.Printf("Connection created: %s\n", result["id"])
			if n, ok := result["agentName"]; ok {
				fmt.Printf("  Agent:  %s\n", n)
			}
			if s, ok := result["status"]; ok {
				fmt.Printf("  Status: %s\n", s)
			}
		}

		return nil
	},
}

var externalAgentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List external agent connections",
	Long: `List all external agent connections for an agent.

Examples:
  cadreen external-agents list --agent-id agent_123
  cadreen external-agents list --agent-id agent_123 --status pending_approval`,
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		agentID, _ := cmd.Flags().GetString("agent-id")
		status, _ := cmd.Flags().GetString("status")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		if agentID == "" {
			return fmt.Errorf("--agent-id is required")
		}

		path := fmt.Sprintf("/api/v1/cadreen/agents/%s/external?limit=%d&offset=%d", agentID, limit, offset)
		if status != "" {
			path += "&status=" + status
		}

		resp, err := rawDo("GET", path, nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result struct {
			Connections []map[string]any `json:"connections"`
			Total       int              `json:"total"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(result, format)
		default:
			if len(result.Connections) == 0 {
				fmt.Println("No external agent connections.")
				return nil
			}
			fmt.Printf("External agent connections (%d):\n\n", result.Total)
			for _, c := range result.Connections {
				fmt.Printf("  %s\n", c["id"])
				if n, ok := c["agentName"]; ok {
					fmt.Printf("    Agent:  %s\n", n)
				}
				if s, ok := c["status"]; ok {
					fmt.Printf("    Status: %s\n", s)
				}
				if h, ok := c["health"]; ok {
					fmt.Printf("    Health: %s\n", h)
				}
				fmt.Println()
			}
		}

		return nil
	},
}

var externalAgentsApproveCmd = &cobra.Command{
	Use:   "approve [connection-id]",
	Short: "Approve a pending external agent connection",
	Long: `Approve a pending external agent connection.

Examples:
  cadreen external-agents approve conn_123 --agent-id agent_456`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		agentID, _ := cmd.Flags().GetString("agent-id")
		connectionID := args[0]

		if agentID == "" {
			return fmt.Errorf("--agent-id is required")
		}

		path := fmt.Sprintf("/api/v1/cadreen/agents/%s/external/%s/approve", agentID, connectionID)

		resp, err := rawDo("POST", path, nil)
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
			fmt.Println("Connection approved.")
		}

		return nil
	},
}

var externalAgentsSuspendCmd = &cobra.Command{
	Use:   "suspend [connection-id]",
	Short: "Suspend an active external agent connection",
	Long: `Suspend an active external agent connection.

Examples:
  cadreen external-agents suspend conn_123 --agent-id agent_456`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		agentID, _ := cmd.Flags().GetString("agent-id")
		connectionID := args[0]

		if agentID == "" {
			return fmt.Errorf("--agent-id is required")
		}

		path := fmt.Sprintf("/api/v1/cadreen/agents/%s/external/%s/suspend", agentID, connectionID)

		resp, err := rawDo("POST", path, nil)
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
			fmt.Println("Connection suspended.")
		}

		return nil
	},
}

var externalAgentsRevokeCmd = &cobra.Command{
	Use:   "revoke [connection-id]",
	Short: "Revoke an external agent connection (permanent)",
	Long: `Revoke an external agent connection. This is permanent.

Examples:
  cadreen external-agents revoke conn_123 --agent-id agent_456`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		agentID, _ := cmd.Flags().GetString("agent-id")
		connectionID := args[0]

		if agentID == "" {
			return fmt.Errorf("--agent-id is required")
		}

		path := fmt.Sprintf("/api/v1/cadreen/agents/%s/external/%s/revoke", agentID, connectionID)

		resp, err := rawDo("POST", path, nil)
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
			fmt.Println("Connection revoked.")
		}

		return nil
	},
}

var externalAgentsDeleteCmd = &cobra.Command{
	Use:   "delete [connection-id]",
	Short: "Delete an external agent connection",
	Long: `Delete an external agent connection.

Examples:
  cadreen external-agents delete conn_123 --agent-id agent_456`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		agentID, _ := cmd.Flags().GetString("agent-id")
		connectionID := args[0]

		if agentID == "" {
			return fmt.Errorf("--agent-id is required")
		}

		path := fmt.Sprintf("/api/v1/cadreen/agents/%s/external/%s", agentID, connectionID)

		resp, err := rawDo("DELETE", path, nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 && resp.StatusCode != 204 {
			return handleHTTPError(resp)
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(map[string]string{"status": "deleted"}, format)
		} else {
			fmt.Println("Connection deleted.")
		}

		return nil
	},
}

var externalAgentsListAllCmd = &cobra.Command{
	Use:   "list-all",
	Short: "List all external agent connections in the workspace",
	Long: `List all external agent connections across all agents in the workspace.

Examples:
  cadreen external-agents list-all`,
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/external-agents/connections", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result struct {
			Connections []map[string]any `json:"connections"`
			Total       int              `json:"total"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(result, format)
		default:
			if len(result.Connections) == 0 {
				fmt.Println("No external agent connections.")
				return nil
			}
			fmt.Printf("External agent connections (%d):\n\n", result.Total)
			for _, c := range result.Connections {
				fmt.Printf("  %s\n", c["id"])
				if n, ok := c["agentName"]; ok {
					fmt.Printf("    Agent:  %s\n", n)
				}
				if s, ok := c["status"]; ok {
					fmt.Printf("    Status: %s\n", s)
				}
				if h, ok := c["health"]; ok {
					fmt.Printf("    Health: %s\n", h)
				}
				fmt.Println()
			}
		}

		return nil
	},
}

var externalAgentsSettingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "View or update external agent settings",
	Long: `View or update external agent settings for the workspace.

Examples:
  cadreen external-agents settings
  cadreen external-agents settings --enable
  cadreen external-agents settings --disable`,
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		enable, _ := cmd.Flags().GetBool("enable")
		disable, _ := cmd.Flags().GetBool("disable")

		if enable && disable {
			return fmt.Errorf("cannot use both --enable and --disable")
		}

		if enable || disable {
			payload := map[string]any{"enabled": enable}
			resp, err := rawDo("PUT", "/api/v1/cadreen/external-agents/settings", payload)
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
				if enable {
					fmt.Println("External agents enabled.")
				} else {
					fmt.Println("External agents disabled.")
				}
			}
		} else {
			resp, err := rawDo("GET", "/api/v1/cadreen/external-agents/settings", nil)
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
				enabled, _ := result["enabled"].(bool)
				if enabled {
					fmt.Println("External agents: enabled")
				} else {
					fmt.Println("External agents: disabled")
				}
			}
		}

		return nil
	},
}

func init() {
	externalAgentsConnectCmd.Flags().String("agent-id", "", "agent ID (required)")
	externalAgentsConnectCmd.Flags().String("url", "", "Agent Card URL (required)")

	externalAgentsListCmd.Flags().String("agent-id", "", "agent ID (required)")
	externalAgentsListCmd.Flags().String("status", "", "filter by status (pending_approval, active, suspended, revoked)")
	externalAgentsListCmd.Flags().Int("limit", 50, "max results")
	externalAgentsListCmd.Flags().Int("offset", 0, "skip this many results")

	externalAgentsApproveCmd.Flags().String("agent-id", "", "agent ID (required)")
	externalAgentsSuspendCmd.Flags().String("agent-id", "", "agent ID (required)")
	externalAgentsRevokeCmd.Flags().String("agent-id", "", "agent ID (required)")

	externalAgentsSettingsCmd.Flags().Bool("enable", false, "enable external agents")
	externalAgentsSettingsCmd.Flags().Bool("disable", false, "disable external agents")

	externalAgentsCmd.AddCommand(externalAgentsConnectCmd)
	externalAgentsCmd.AddCommand(externalAgentsListCmd)
	externalAgentsCmd.AddCommand(externalAgentsApproveCmd)
	externalAgentsCmd.AddCommand(externalAgentsSuspendCmd)
	externalAgentsCmd.AddCommand(externalAgentsRevokeCmd)
	externalAgentsCmd.AddCommand(externalAgentsDeleteCmd)
	externalAgentsCmd.AddCommand(externalAgentsListAllCmd)
	externalAgentsCmd.AddCommand(externalAgentsSettingsCmd)
	rootCmd.AddCommand(externalAgentsCmd)
}
