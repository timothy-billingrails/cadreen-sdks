package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	"github.com/spf13/cobra"
)

var agentsCmd = &cobra.Command{
	Use:   "agents",
	Short: "Manage agents (create, ask, deploy, govern)",
	Long: `Agents are autonomous workers that follow your rules.

Create an agent, give it knowledge, set governance policies,
then let it handle tasks — or ask it questions in plain language.`,
}

var agentsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an agent",
	Long: `Create a new agent.

Examples:
  cadreen agents create --name "Refund Bot" --description "Handles refund requests"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		if name == "" {
			return fmt.Errorf("--name is required")
		}

		payload := map[string]any{"name": name}
		if description != "" {
			payload["description"] = description
		}

		resp, err := rawDo("POST", "/api/v1/cadreen/agents", payload)
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
			fmt.Printf("Agent created: %s\n", result["id"])
			if n, ok := result["name"]; ok {
				fmt.Printf("  Name:   %s\n", n)
			}
			if s, ok := result["status"]; ok {
				fmt.Printf("  Status: %s\n", s)
			}
		}

		return nil
	},
}

var agentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your agents",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		search, _ := cmd.Flags().GetString("search")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		path := fmt.Sprintf("/api/v1/cadreen/agents?limit=%d&offset=%d", limit, offset)
		if search != "" {
			path += "&search=" + search
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
				fmt.Println("No agents found.")
				return nil
			}
			fmt.Printf("Agents (%d):\n\n", len(result.Agents))
			for _, a := range result.Agents {
				fmt.Printf("  %s  %s\n", a["id"], a["name"])
				if s, ok := a["status"]; ok {
					fmt.Printf("    Status: %s\n", s)
				}
				fmt.Println()
			}
		}

		return nil
	},
}

var agentsGetCmd = &cobra.Command{
	Use:   "get [agent-id]",
	Short: "Get agent details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", fmt.Sprintf("/api/v1/cadreen/agents/%s", args[0]), nil)
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
			printAgentDetail(result)
		}

		return nil
	},
}

var agentsUpdateCmd = &cobra.Command{
	Use:   "update [agent-id]",
	Short: "Update an agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		status, _ := cmd.Flags().GetString("status")

		payload := map[string]any{}
		if name != "" {
			payload["name"] = name
		}
		if description != "" {
			payload["description"] = description
		}
		if status != "" {
			payload["status"] = status
		}

		if len(payload) == 0 {
			return fmt.Errorf("provide at least one flag: --name, --description, --status")
		}

		resp, err := rawDo("PATCH", fmt.Sprintf("/api/v1/cadreen/agents/%s", args[0]), payload)
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
			fmt.Println("Agent updated.")
		}

		return nil
	},
}

var agentsDeleteCmd = &cobra.Command{
	Use:   "delete [agent-id]",
	Short: "Delete an agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("DELETE", fmt.Sprintf("/api/v1/cadreen/agents/%s", args[0]), nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 && resp.StatusCode != 204 {
			return handleHTTPError(resp)
		}

		fmt.Println("Agent deleted.")
		return nil
	},
}

var agentsDeployCmd = &cobra.Command{
	Use:   "deploy [agent-id]",
	Short: "Deploy an agent",
	Long:  `Deploy an agent so it can start handling tasks.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("POST", fmt.Sprintf("/api/v1/cadreen/agents/%s/deploy", args[0]), nil)
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
			fmt.Println("Agent deployed.")
		}

		return nil
	},
}

var agentsCapabilitiesCmd = &cobra.Command{
	Use:   "capabilities [agent-id]",
	Short: "See what an agent can do",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", fmt.Sprintf("/api/v1/cadreen/agents/%s/capabilities", args[0]), nil)
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
			if caps, ok := result["capabilities"].([]any); ok {
				fmt.Printf("This agent can do %d things:\n\n", len(caps))
				for _, c := range caps {
					if m, ok := c.(map[string]any); ok {
						fmt.Printf("  - %s\n", m["name"])
						if d, ok := m["description"]; ok && d != "" {
							fmt.Printf("    %s\n", d)
						}
					}
				}
			} else {
				output.Print(result, output.FormatJSON)
			}
		}

		return nil
	},
}

var agentsAskCmd = &cobra.Command{
	Use:   "ask [agent-id] [message]",
	Short: "Ask an agent something (plain language)",
	Long: `Send a message to an agent in plain language.

Examples:
  cadreen agents ask ag_01JX "What do you know about customer 123?"
  cadreen agents ask ag_01JX "Process a refund for order 456"`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		agentID := args[0]
		message := strings.Join(args[1:], " ")

		payload := map[string]any{
			"message": message,
		}

		resp, err := rawDo("POST", fmt.Sprintf("/api/v1/cadreen/agents/%s/send", agentID), payload)
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
			if content, ok := result["content"]; ok {
				fmt.Println(content)
			} else if resp2, ok := result["response"]; ok {
				fmt.Println(resp2)
			} else {
				output.Print(result, output.FormatJSON)
			}
		}

		return nil
	},
}

var agentsMessagesCmd = &cobra.Command{
	Use:   "messages [agent-id]",
	Short: "See conversation history",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", fmt.Sprintf("/api/v1/cadreen/agents/%s/messages", args[0]), nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result struct {
			Messages []map[string]any `json:"messages"`
			Count    int              `json:"count"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(result, format)
		default:
			if len(result.Messages) == 0 {
				fmt.Println("No messages yet.")
				return nil
			}
			fmt.Printf("Messages (%d):\n\n", len(result.Messages))
			for _, m := range result.Messages {
				role := m["role"]
				content := m["content"]
				fmt.Printf("  [%s] %s\n", role, content)
				fmt.Println()
			}
		}

		return nil
	},
}

var agentsExecutionsCmd = &cobra.Command{
	Use:   "executions [agent-id]",
	Short: "See what an agent has done",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", fmt.Sprintf("/api/v1/cadreen/agents/%s/executions", args[0]), nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result struct {
			Executions []map[string]any `json:"executions"`
			Count      int              `json:"count"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(result, format)
		default:
			if len(result.Executions) == 0 {
				fmt.Println("No executions yet.")
				return nil
			}
			fmt.Printf("Executions (%d):\n\n", len(result.Executions))
			for _, e := range result.Executions {
				fmt.Printf("  %s  %s\n", e["id"], e["status"])
				if s, ok := e["summary"]; ok && s != "" {
					fmt.Printf("    %s\n", s)
				}
				fmt.Println()
			}
		}

		return nil
	},
}

var agentsKnowledgeCmd = &cobra.Command{
	Use:   "knowledge [agent-id]",
	Short: "See what an agent knows",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", fmt.Sprintf("/api/v1/cadreen/agents/%s/knowledge", args[0]), nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result struct {
			Knowledge []map[string]any `json:"knowledge"`
			Count     int              `json:"count"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(result, format)
		default:
			if len(result.Knowledge) == 0 {
				fmt.Println("This agent doesn't know anything yet.")
				return nil
			}
			fmt.Printf("Knowledge (%d):\n\n", len(result.Knowledge))
			for _, k := range result.Knowledge {
				fmt.Printf("  %s\n", k["subject"])
				if o, ok := k["object"]; ok && o != "" {
					fmt.Printf("    %s\n", o)
				}
				fmt.Println()
			}
		}

		return nil
	},
}

var agentsKnowledgeAddCmd = &cobra.Command{
	Use:   "knowledge-add [agent-id]",
	Short: "Teach an agent something",
	Long: `Add a piece of knowledge to an agent's memory.

Examples:
  cadreen agents knowledge-add ag_01JX --subject "Refund policy" --object "Refunds over $500 need manager approval" --type reference`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		subject, _ := cmd.Flags().GetString("subject")
		object, _ := cmd.Flags().GetString("object")
		knowledgeType, _ := cmd.Flags().GetString("type")

		if subject == "" {
			return fmt.Errorf("--subject is required")
		}
		if object == "" {
			return fmt.Errorf("--object is required")
		}

		payload := map[string]any{
			"subject": subject,
			"object":  object,
		}
		if knowledgeType != "" {
			payload["type"] = knowledgeType
		}

		resp, err := rawDo("POST", fmt.Sprintf("/api/v1/cadreen/agents/%s/knowledge", args[0]), payload)
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
			fmt.Println("Knowledge added.")
		}

		return nil
	},
}

var agentsKnowledgeSearchCmd = &cobra.Command{
	Use:   "knowledge-search [agent-id]",
	Short: "Search what an agent knows",
	Long: `Search an agent's knowledge base.

Examples:
  cadreen agents knowledge-search ag_01JX --query "refund policy"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		query, _ := cmd.Flags().GetString("query")
		if query == "" {
			return fmt.Errorf("--query is required")
		}

		resp, err := rawDo("POST", fmt.Sprintf("/api/v1/cadreen/agents/%s/knowledge/search", args[0]), map[string]any{"query": query})
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result struct {
			Results []map[string]any `json:"results"`
			Count   int              `json:"count"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(result, format)
		default:
			if len(result.Results) == 0 {
				fmt.Println("Nothing found.")
				return nil
			}
			fmt.Printf("Found %d results:\n\n", len(result.Results))
			for _, r := range result.Results {
				fmt.Printf("  %s\n", r["subject"])
				if o, ok := r["object"]; ok && o != "" {
					fmt.Printf("    %s\n", o)
				}
				fmt.Println()
			}
		}

		return nil
	},
}

var agentsGovernanceCmd = &cobra.Command{
	Use:   "governance [agent-id]",
	Short: "See governance policies for an agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", fmt.Sprintf("/api/v1/cadreen/agents/%s/governance", args[0]), nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result struct {
			Policies []map[string]any `json:"policies"`
			Count    int              `json:"count"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(result, format)
		default:
			if len(result.Policies) == 0 {
				fmt.Println("No governance policies for this agent.")
				return nil
			}
			fmt.Printf("Governance policies (%d):\n\n", len(result.Policies))
			for _, p := range result.Policies {
				fmt.Printf("  %s\n", p["name"])
				if r, ok := p["rationale"]; ok && r != "" {
					fmt.Printf("    %s\n", r)
				}
				fmt.Println()
			}
		}

		return nil
	},
}

var agentsAuditCmd = &cobra.Command{
	Use:   "audit [agent-id]",
	Short: "See what an agent has done (audit log)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", fmt.Sprintf("/api/v1/cadreen/agents/%s/audit", args[0]), nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result struct {
			Entries []map[string]any `json:"entries"`
			Count   int              `json:"count"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(result, format)
		default:
			if len(result.Entries) == 0 {
				fmt.Println("No audit entries yet.")
				return nil
			}
			fmt.Printf("Audit log (%d):\n\n", len(result.Entries))
			for _, e := range result.Entries {
				fmt.Printf("  %s  %s\n", e["timestamp"], e["action"])
				if d, ok := e["detail"]; ok && d != "" {
					fmt.Printf("    %s\n", d)
				}
				fmt.Println()
			}
		}

		return nil
	},
}

var agentsNegotiateCmd = &cobra.Command{
	Use:   "negotiate [agent-id]",
	Short: "Start a negotiation with another agent",
	Long: `Propose an action to another agent and negotiate terms.

Examples:
  cadreen agents negotiate ag_01JX --to ag_02KY --proposal "Handle refund for order 789"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		to, _ := cmd.Flags().GetString("to")
		proposal, _ := cmd.Flags().GetString("proposal")

		if to == "" {
			return fmt.Errorf("--to is required (target agent ID)")
		}
		if proposal == "" {
			return fmt.Errorf("--proposal is required")
		}

		payload := map[string]any{
			"to":       to,
			"proposal": proposal,
		}

		resp, err := rawDo("POST", fmt.Sprintf("/api/v1/cadreen/agents/%s/negotiate", args[0]), payload)
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
			fmt.Printf("Negotiation started: %s\n", result["id"])
			if s, ok := result["status"]; ok {
				fmt.Printf("  Status: %s\n", s)
			}
		}

		return nil
	},
}

var agentsNegotiationsCmd = &cobra.Command{
	Use:   "negotiations [agent-id]",
	Short: "See ongoing negotiations",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", fmt.Sprintf("/api/v1/cadreen/agents/%s/negotiations", args[0]), nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result struct {
			Negotiations []map[string]any `json:"negotiations"`
			Count        int              `json:"count"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(result, format)
		default:
			if len(result.Negotiations) == 0 {
				fmt.Println("No negotiations.")
				return nil
			}
			fmt.Printf("Negotiations (%d):\n\n", len(result.Negotiations))
			for _, n := range result.Negotiations {
				fmt.Printf("  %s  %s\n", n["id"], n["status"])
				if p, ok := n["proposal"]; ok && p != "" {
					fmt.Printf("    %s\n", p)
				}
				fmt.Println()
			}
		}

		return nil
	},
}

func init() {
	agentsCreateCmd.Flags().String("name", "", "agent name (required)")
	agentsCreateCmd.Flags().String("description", "", "what this agent does")

	agentsListCmd.Flags().String("search", "", "search agents by name")
	agentsListCmd.Flags().Int("limit", 50, "max results")
	agentsListCmd.Flags().Int("offset", 0, "skip this many results")

	agentsUpdateCmd.Flags().String("name", "", "new name")
	agentsUpdateCmd.Flags().String("description", "", "new description")
	agentsUpdateCmd.Flags().String("status", "", "new status (active, paused)")

	agentsKnowledgeAddCmd.Flags().String("subject", "", "what you're teaching (required)")
	agentsKnowledgeAddCmd.Flags().String("object", "", "the knowledge content (required)")
	agentsKnowledgeAddCmd.Flags().String("type", "", "knowledge type (reference, procedure, preference)")

	agentsKnowledgeSearchCmd.Flags().String("query", "", "search query (required)")

	agentsNegotiateCmd.Flags().String("to", "", "target agent ID (required)")
	agentsNegotiateCmd.Flags().String("proposal", "", "what you're proposing (required)")

	agentsCmd.AddCommand(agentsCreateCmd)
	agentsCmd.AddCommand(agentsListCmd)
	agentsCmd.AddCommand(agentsGetCmd)
	agentsCmd.AddCommand(agentsUpdateCmd)
	agentsCmd.AddCommand(agentsDeleteCmd)
	agentsCmd.AddCommand(agentsDeployCmd)
	agentsCmd.AddCommand(agentsCapabilitiesCmd)
	agentsCmd.AddCommand(agentsAskCmd)
	agentsCmd.AddCommand(agentsMessagesCmd)
	agentsCmd.AddCommand(agentsExecutionsCmd)
	agentsCmd.AddCommand(agentsKnowledgeCmd)
	agentsCmd.AddCommand(agentsKnowledgeAddCmd)
	agentsCmd.AddCommand(agentsKnowledgeSearchCmd)
	agentsCmd.AddCommand(agentsGovernanceCmd)
	agentsCmd.AddCommand(agentsAuditCmd)
	agentsCmd.AddCommand(agentsNegotiateCmd)
	agentsCmd.AddCommand(agentsNegotiationsCmd)
	rootCmd.AddCommand(agentsCmd)
}

func printAgentDetail(a map[string]any) {
	fmt.Printf("ID:     %s\n", a["id"])
	fmt.Printf("Name:   %s\n", a["name"])
	if s, ok := a["status"]; ok {
		fmt.Printf("Status: %s\n", s)
	}
	if d, ok := a["description"]; ok && d != "" {
		fmt.Printf("Description: %s\n", d)
	}
	if c, ok := a["created_at"]; ok {
		fmt.Printf("Created: %s\n", c)
	}
	fmt.Println()
}
