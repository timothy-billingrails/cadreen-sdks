package commands

import (
	"encoding/json"
	"fmt"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	cadreen "github.com/timothy-billingrails/cadreen-sdks/go/cadreen"

	"github.com/spf13/cobra"
)

var credentialsCmd = &cobra.Command{
	Use:   "credentials",
	Short: "Manage credentials (list, create, delete)",
	Long:  `Manage API keys and credentials for connected services.`,
}

var credentialsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/credentials", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var creds cadreen.ListCredentialsResponse
		if err := json.NewDecoder(resp.Body).Decode(&creds); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(creds, format)
		case output.FormatVerbose:
			printCredentialsVerbose(creds.Credentials)
		default:
			printCredentialsText(creds.Credentials)
		}

		return nil
	},
}

var credentialsDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a credential",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		client := newClient()
		if err := client.DeleteCredential(cmd.Context(), args[0]); err != nil {
			return handleAPIError(err)
		}

		fmt.Printf("Credential %s deleted.\n", args[0])
		return nil
	},
}

func init() {
	credentialsCmd.AddCommand(credentialsListCmd)
	credentialsCmd.AddCommand(credentialsDeleteCmd)
	rootCmd.AddCommand(credentialsCmd)
}

func printCredentialsText(creds []cadreen.CredentialMetadata) {
	if len(creds) == 0 {
		fmt.Println("No credentials configured.")
		return
	}

	fmt.Printf("Credentials (%d):\n\n", len(creds))
	for _, c := range creds {
		status := "inactive"
		if c.IsActive {
			status = "active"
		}
		fmt.Printf("  %s (%s) — %s\n", c.CredentialName, c.Provider, status)
		fmt.Println()
	}
}

func printCredentialsVerbose(creds []cadreen.CredentialMetadata) {
	for _, c := range creds {
		fmt.Printf("ID:       %s\n", c.ID)
		fmt.Printf("Provider: %s\n", c.Provider)
		fmt.Printf("Name:     %s\n", c.CredentialName)
		fmt.Printf("Active:   %v\n", c.IsActive)
		fmt.Println()
	}
}
