package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	cadreen "github.com/timothy-billingrails/cadreen-sdks/go/cadreen"

	"github.com/spf13/cobra"
)

var documentsCmd = &cobra.Command{
	Use:   "documents",
	Short: "Manage documents (list, get, upload)",
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

var documentsUploadCmd = &cobra.Command{
	Use:   "upload [filepath]",
	Short: "Upload a document",
	Long:  `Upload a file as a document to Cadreen.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		filePath := args[0]
		f, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("open file: %w", err)
		}
		defer f.Close()

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		part, err := writer.CreateFormFile("document", filepath.Base(filePath))
		if err != nil {
			return fmt.Errorf("create form file: %w", err)
		}
		if _, err := io.Copy(part, f); err != nil {
			return fmt.Errorf("copy file: %w", err)
		}
		if err := writer.Close(); err != nil {
			return fmt.Errorf("close writer: %w", err)
		}

		baseURL := cfg.BaseURLResolved()
		url := baseURL + "/api/v1/cadreen/documents/upload"
		req, err := http.NewRequest("POST", url, &body)
		if err != nil {
			return fmt.Errorf("create request: %w", err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+cfg.APIKeyResolved())

		resp, err := httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("upload request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return handleHTTPError(resp)
		}

		var uploadResp cadreen.UploadDocumentResponse
		if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(uploadResp, format)
		default:
			fmt.Printf("Uploaded: %s\n", uploadResp.Name)
			fmt.Printf("ID:       %s\n", uploadResp.ID)
			if uploadResp.Status != "" {
				fmt.Printf("Status:   %s\n", uploadResp.Status)
			}
			if uploadResp.ContentType != "" {
				fmt.Printf("Type:     %s\n", uploadResp.ContentType)
			}
			if uploadResp.Size > 0 {
				fmt.Printf("Size:     %d bytes\n", uploadResp.Size)
			}
		}

		return nil
	},
}

func init() {
	documentsCmd.AddCommand(documentsListCmd)
	documentsCmd.AddCommand(documentsGetCmd)
	documentsCmd.AddCommand(documentsUploadCmd)
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
