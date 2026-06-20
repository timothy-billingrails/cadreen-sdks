package commands

import (
	"context"
	"fmt"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	cadreen "github.com/timothy-billingrails/cadreen-sdks/go/cadreen"

	"github.com/spf13/cobra"
)

var memoryCmd = &cobra.Command{
	Use:   "memory",
	Short: "Manage knowledge (list, add, search)",
	Long: `Manage what Cadreen knows.

List memories, add new knowledge, or search existing knowledge.`,
}

var memoryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent knowledge",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		limit, _ := cmd.Flags().GetInt("limit")
		domain, _ := cmd.Flags().GetString("domain")

		client := newClient()
		opts := &cadreen.SearchMemoryOptions{
			Limit: limit,
		}
		if domain != "" {
			opts.Domain = domain
		}

		resp, err := client.Search(context.Background(), "", opts)
		if err != nil {
			return handleAPIError(err)
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(resp, format)
		case output.FormatVerbose:
			printMemoryVerbose(resp)
		default:
			printMemoryText(resp)
		}

		return nil
	},
}

var memoryAddCmd = &cobra.Command{
	Use:   "add [text]",
	Short: "Teach Cadreen something new",
	Long: `Add knowledge to Cadreen's memory.

Examples:
  cadreen memory add "Our refund policy is 30 days"
  cadreen memory add "User prefers Python over JavaScript"`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		domain, _ := cmd.Flags().GetString("domain")
		atomType, _ := cmd.Flags().GetString("type")
		authority, _ := cmd.Flags().GetInt("authority")

		text := args[0]

		client := newClient()
		resp, err := client.Teach(context.Background(), cadreen.RememberRequest{
			Type: atomType,
			Content: map[string]any{
				"text": text,
			},
			Domain:    domain,
			Authority: authority,
		})
		if err != nil {
			return handleAPIError(err)
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(resp, format)
		} else {
			fmt.Printf("Knowledge stored: %s\n", resp.ID)
		}

		return nil
	},
}

func init() {
	memoryListCmd.Flags().Int("limit", 10, "max results")
	memoryListCmd.Flags().String("domain", "", "filter by domain")

	memoryAddCmd.Flags().String("domain", "general", "knowledge domain")
	memoryAddCmd.Flags().String("type", "reference", "atom type (reference, preference, episode, precedent, note)")
	memoryAddCmd.Flags().Int("authority", 5, "authority level (1-10)")

	memoryCmd.AddCommand(memoryListCmd)
	memoryCmd.AddCommand(memoryAddCmd)
	rootCmd.AddCommand(memoryCmd)
}

func printMemoryText(resp *cadreen.SearchMemoryResponse) {
	if resp == nil || len(resp.Results) == 0 {
		fmt.Println("No memories found.")
		return
	}

	fmt.Printf("Cadreen remembers things. Here are the recent ones:\n\n")

	for _, r := range resp.Results {
		fmt.Printf("  %s\n", humanAtomType(r.Type))
		if r.Content.Text != "" {
			fmt.Printf("  \"%s\"\n", r.Content.Text)
		}
		if r.Domain != "" {
			fmt.Printf("  Area: %s\n", r.Domain)
		}
		fmt.Println()
	}
}

func printMemoryVerbose(resp *cadreen.SearchMemoryResponse) {
	if resp == nil || len(resp.Results) == 0 {
		fmt.Println("No memories found.")
		return
	}

	for _, r := range resp.Results {
		fmt.Printf("ID:       %s\n", r.ID)
		fmt.Printf("Type:     %s\n", r.Type)
		fmt.Printf("Domain:   %s\n", r.Domain)
		if r.Content.Text != "" {
			fmt.Printf("Content:  %s\n", r.Content.Text)
		}
		fmt.Printf("Authority: %d\n", r.Authority)
		fmt.Println()
	}
}

func humanAtomType(t string) string {
	switch t {
	case "fact", "reference":
		return "A fact Cadreen learned"
	case "preference":
		return "A preference Cadreen picked up"
	case "episode":
		return "Something that happened"
	case "precedent":
		return "Something Cadreen fixed before"
	case "procedure":
		return "How to do something"
	case "note":
		return "A note"
	default:
		return t
	}
}
