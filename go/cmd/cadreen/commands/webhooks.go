package commands

import (
	"encoding/json"
	"fmt"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	cadreen "github.com/timothy-billingrails/cadreen-sdks/go/cadreen"

	"github.com/spf13/cobra"
)

var webhooksCmd = &cobra.Command{
	Use:   "webhooks",
	Short: "Manage webhooks (list, create, delete)",
	Long:  `Manage webhook subscriptions for event notifications.`,
}

var webhooksListCmd = &cobra.Command{
	Use:   "list",
	Short: "List webhooks",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/webhooks", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var webhooks cadreen.ListWebhooksResponse
		if err := json.NewDecoder(resp.Body).Decode(&webhooks); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(webhooks, format)
		case output.FormatVerbose:
			printWebhooksVerbose(webhooks.Webhooks)
		default:
			printWebhooksText(webhooks.Webhooks)
		}

		return nil
	},
}

var webhooksCreateCmd = &cobra.Command{
	Use:   "create [url]",
	Short: "Create a webhook",
	Long: `Subscribe to events at a URL.

Examples:
  cadreen webhooks create https://example.com/hooks --events escalation.created
  cadreen webhooks create https://example.com/hooks --events execution.completed --secret mysecret`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		events, _ := cmd.Flags().GetStringSlice("events")
		secret, _ := cmd.Flags().GetString("secret")

		client := newClient()
		webhook, err := client.CreateWebhook(cmd.Context(), cadreen.CreateWebhookRequest{
			URL:    args[0],
			Events: events,
			Secret: secret,
		})
		if err != nil {
			return handleAPIError(err)
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(webhook, format)
		} else {
			fmt.Printf("Webhook created: %s\n", webhook.ID)
			fmt.Printf("URL: %s\n", webhook.URL)
		}

		return nil
	},
}

var webhooksDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a webhook",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		client := newClient()
		if err := client.DeleteWebhook(cmd.Context(), args[0]); err != nil {
			return handleAPIError(err)
		}

		fmt.Printf("Webhook %s deleted.\n", args[0])
		return nil
	},
}

func init() {
	webhooksCreateCmd.Flags().StringSlice("events", nil, "events to subscribe to")
	webhooksCreateCmd.Flags().String("secret", "", "HMAC secret for signature verification")

	webhooksCmd.AddCommand(webhooksListCmd)
	webhooksCmd.AddCommand(webhooksCreateCmd)
	webhooksCmd.AddCommand(webhooksDeleteCmd)
	rootCmd.AddCommand(webhooksCmd)
}

func printWebhooksText(webhooks []cadreen.Webhook) {
	if len(webhooks) == 0 {
		fmt.Println("No webhooks configured.")
		return
	}

	fmt.Printf("Webhooks (%d):\n\n", len(webhooks))
	for _, w := range webhooks {
		fmt.Printf("  %s\n", w.URL)
		if len(w.Events) > 0 {
			fmt.Printf("    Events: %v\n", w.Events)
		}
		fmt.Println()
	}
}

func printWebhooksVerbose(webhooks []cadreen.Webhook) {
	for _, w := range webhooks {
		fmt.Printf("ID:       %s\n", w.ID)
		fmt.Printf("URL:      %s\n", w.URL)
		fmt.Printf("Active:   %v\n", w.IsActive)
		if len(w.Events) > 0 {
			fmt.Printf("Events:   %v\n", w.Events)
		}
		fmt.Println()
	}
}
