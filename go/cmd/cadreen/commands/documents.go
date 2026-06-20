package commands

import (
	"encoding/json"
	"fmt"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	cadreen "github.com/timothy-billingrails/cadreen-sdks/go/cadreen"

	"github.com/spf13/cobra"
)

var documentsCmd = &cobra.Command{
	Use:   "documents",
	Short: "Manage documents (list, get)",
	Long:  `Manage uploaded documents.`,
}

var documentsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List documents",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/documents", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var docs cadreen.ListDocumentsResponse
		if err := json.NewDecoder(resp.Body).Decode(&docs); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(docs, format)
		case output.FormatVerbose:
			printDocumentsVerbose(docs.Documents)
		default:
			printDocumentsText(docs.Documents)
		}

		return nil
	},
}

var documentsGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get a document by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/documents/"+args[0], nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var doc cadreen.Document
		if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(doc, format)
		default:
			fmt.Printf("ID:      %s\n", doc.ID)
			fmt.Printf("Name:    %s\n", doc.Name)
			if doc.ContentType != "" {
				fmt.Printf("Type:    %s\n", doc.ContentType)
			}
			if doc.Size > 0 {
				fmt.Printf("Size:    %d bytes\n", doc.Size)
			}
			if doc.Status != "" {
				fmt.Printf("Status:  %s\n", doc.Status)
			}
		}

		return nil
	},
}

func init() {
	documentsCmd.AddCommand(documentsListCmd)
	documentsCmd.AddCommand(documentsGetCmd)
	rootCmd.AddCommand(documentsCmd)
}

func printDocumentsText(docs []cadreen.Document) {
	if len(docs) == 0 {
		fmt.Println("No documents found.")
		return
	}

	fmt.Printf("Documents (%d):\n\n", len(docs))
	for _, d := range docs {
		fmt.Printf("  %s\n", d.Name)
		if d.ContentType != "" {
			fmt.Printf("    Type: %s\n", d.ContentType)
		}
		fmt.Println()
	}
}

func printDocumentsVerbose(docs []cadreen.Document) {
	for _, d := range docs {
		fmt.Printf("ID:      %s\n", d.ID)
		fmt.Printf("Name:    %s\n", d.Name)
		fmt.Printf("Type:    %s\n", d.ContentType)
		fmt.Printf("Size:    %d\n", d.Size)
		fmt.Printf("Status:  %s\n", d.Status)
		fmt.Println()
	}
}
