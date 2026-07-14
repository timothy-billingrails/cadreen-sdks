package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	"github.com/spf13/cobra"
)

var responsesCmd = &cobra.Command{
	Use:   "responses",
	Short: "Manage responses (create, get)",
	Long: `Interact with the Responses API.

Create a new response or retrieve an existing one by ID.`,
}

var responsesCreateCmd = &cobra.Command{
	Use:   "create [input]",
	Short: "Create a response",
	Long: `Create a new response from the Responses API.

The model field is required by the API format but ignored — Cadreen decides which model to use.

Examples:
  cadreen responses create "What tools do I have?"
  cadreen responses create "What tools do I have?" --stream
  cadreen responses create "Refund order 12345" --tools '[{"type":"function","name":"process_refund","parameters":{}}]'
  
  # Pipe from stdin
  echo "summarize this" | cadreen responses create
  cat prompt.txt | cadreen responses create --stream`,
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		model, _ := cmd.Flags().GetString("model")
		stream, _ := cmd.Flags().GetBool("stream")

		// Get input from args or stdin
		var input string
		if len(args) > 0 {
			input = strings.Join(args, " ")
		} else if stat, _ := os.Stdin.Stat(); stat != nil && (stat.Mode()&os.ModeCharDevice) == 0 {
			limited := io.LimitReader(os.Stdin, 1024*1024) // 1MB cap
			data, err := io.ReadAll(limited)
			if err != nil {
				return fmt.Errorf("reading stdin: %w", err)
			}
			input = string(data)
		} else {
			return cmd.Help()
		}

		if strings.TrimSpace(input) == "" {
			return cmd.Help()
		}

		// Model is optional — Cadreen decides which model to use
		if model == "" {
			model = "cadreen"
		}

		payload := map[string]any{
			"model": model,
			"input": input,
		}
		if stream {
			payload["stream"] = true
		}

		resp, err := rawDo("POST", "/api/v1/cadreen/responses", payload)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if stream {
			return responsesStream(resp)
		}

		if resp.StatusCode != 200 && resp.StatusCode != 201 {
			return handleHTTPError(resp)
		}

		var result map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(result, format)
		case output.FormatVerbose:
			printResponseVerbose(result)
		default:
			printResponseText(result)
		}

		return nil
	},
}

var responsesGetCmd = &cobra.Command{
	Use:   "get [response-id]",
	Short: "Get a response by ID",
	Long: `Retrieve an existing response by its ID.

Examples:
  cadreen responses get resp_abc123`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		responseID := args[0]
		path := fmt.Sprintf("/api/v1/cadreen/responses/%s", responseID)

		resp, err := rawDo("GET", path, nil)
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
		switch format {
		case output.FormatJSON:
			output.Print(result, format)
		case output.FormatVerbose:
			printResponseVerbose(result)
		default:
			printResponseText(result)
		}

		return nil
	},
}

func init() {
	responsesCreateCmd.Flags().String("model", "cadreen", "model to use (default: cadreen, ignored by Cadreen)")
	responsesCreateCmd.Flags().Bool("stream", false, "stream the response")

	responsesCmd.AddCommand(responsesCreateCmd)
	responsesCmd.AddCommand(responsesGetCmd)
	rootCmd.AddCommand(responsesCmd)
}

func responsesStream(resp *http.Response) error {
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return handleHTTPError(resp)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var event map[string]any
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(event, format)
			continue
		}

		printStreamEvent(event)
	}

	fmt.Println()
	return nil
}

func printStreamEvent(event map[string]any) {
	// Server format: {type, sequence, delta (string), response, item}
	// delta is a plain string, not an object
	if delta, ok := event["delta"].(string); ok {
		fmt.Print(delta)
		return
	}

	// Fallback: check for nested delta.content (OpenAI SDK format)
	if delta, ok := event["delta"].(map[string]any); ok {
		if content, ok := delta["content"].(string); ok {
			fmt.Print(content)
			return
		}
	}

	// Fallback: check for output_text in response object
	if resp, ok := event["response"].(map[string]any); ok {
		if outputText, ok := resp["output_text"].(string); ok {
			fmt.Print(outputText)
			return
		}
	}

	if output, ok := event["output"].(string); ok {
		fmt.Print(output)
	}

	if text, ok := event["text"].(string); ok {
		fmt.Print(text)
	}
}

func printResponseText(result map[string]any) {
	if id, ok := result["id"]; ok {
		fmt.Printf("Response: %v\n", id)
	}

	if output, ok := result["output"]; ok {
		switch o := output.(type) {
		case string:
			fmt.Println(o)
		case []any:
			for _, item := range o {
				if m, ok := item.(map[string]any); ok {
					if content, ok := m["content"]; ok {
						switch c := content.(type) {
						case string:
							fmt.Println(c)
						case []any:
							for _, part := range c {
								if pm, ok := part.(map[string]any); ok {
									if text, ok := pm["text"].(string); ok {
										fmt.Println(text)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	if content, ok := result["content"]; ok {
		switch c := content.(type) {
		case string:
			fmt.Println(c)
		case []any:
			for _, item := range c {
				if m, ok := item.(map[string]any); ok {
					if text, ok := m["text"].(string); ok {
						fmt.Println(text)
					}
				}
			}
		}
	}
}

func printResponseVerbose(result map[string]any) {
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("Error formatting response: %s\n", err)
		return
	}
	fmt.Println(string(b))
}
