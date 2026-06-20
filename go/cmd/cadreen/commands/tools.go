package commands

import (
	"encoding/json"
	"fmt"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	"github.com/spf13/cobra"
)

var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "List available tools",
	Long: `List the tools Cadreen can use.

Shows connected services and their capabilities.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/tools", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var tools ListToolsResponseData
		if err := json.NewDecoder(resp.Body).Decode(&tools); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(tools, format)
		case output.FormatVerbose:
			printToolsVerbose(tools)
		default:
			printToolsText(tools)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(toolsCmd)
}

func printToolsText(t ListToolsResponseData) {
	if len(t.Data) == 0 {
		fmt.Println("No tools available.")
		return
	}

	fmt.Printf("Cadreen can do %d things:\n\n", len(t.Data))

	for _, tool := range t.Data {
		name := tool.Function.Name
		if tool.Function.Description != "" {
			fmt.Printf("  %s\n", name)
			fmt.Printf("    %s\n", tool.Function.Description)
		} else {
			fmt.Printf("  %s\n", name)
		}
		fmt.Println()
	}
}

func printToolsVerbose(t ListToolsResponseData) {
	for _, tool := range t.Data {
		fmt.Printf("Name:        %s\n", tool.Function.Name)
		fmt.Printf("Description: %s\n", tool.Function.Description)
		fmt.Println()
	}
}
