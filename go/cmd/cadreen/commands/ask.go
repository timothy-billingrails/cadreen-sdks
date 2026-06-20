package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	cadreen "github.com/timothy-billingrails/cadreen-sdks/go/cadreen"

	"github.com/spf13/cobra"
)

var askCmd = &cobra.Command{
	Use:   "ask [question]",
	Short: "Ask Cadreen a question (one-shot)",
	Long: `Ask Cadreen a question and get a streaming response.

Supports stdin for piping:
  echo "summarize this" | cadreen ask
  cat error.log | cadreen ask "what went wrong?"
  cadreen ask < prompt.txt`,
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		var message string
		if len(args) > 0 {
			message = strings.Join(args, " ")
		} else if stat, _ := os.Stdin.Stat(); stat != nil && (stat.Mode()&os.ModeCharDevice) == 0 {
			limited := io.LimitReader(os.Stdin, 1024*1024) // 1MB cap
			data, err := io.ReadAll(limited)
			if err != nil {
				return fmt.Errorf("reading stdin: %w", err)
			}
			message = string(data)
		} else {
			return cmd.Help()
		}

		if strings.TrimSpace(message) == "" {
			return cmd.Help()
		}

		noStream, _ := cmd.Flags().GetBool("no-stream")
		memoryOff, _ := cmd.Flags().GetBool("no-memory")
		convID, _ := cmd.Flags().GetString("conversation-id")

		req := cadreen.ChatCompletionRequest{
			Messages: []cadreen.ChatMessage{
				{Role: "user", Content: message},
			},
		}
		if memoryOff {
			req.Context = map[string]any{"memory": false}
		}
		if convID != "" {
			req.ConversationID = convID
		}

		if noStream {
			return askNonStreaming(req)
		}
		return askStreaming(req)
	},
}

func init() {
	askCmd.Flags().Bool("no-stream", false, "disable streaming")
	askCmd.Flags().Bool("no-memory", false, "disable cross-conversation memory")
	askCmd.Flags().String("conversation-id", "", "continue an existing conversation")
	rootCmd.AddCommand(askCmd)
}

func askStreaming(req cadreen.ChatCompletionRequest) error {
	client := newClient()
	ch, err := client.ChatCompletionsStream(context.Background(), req)
	if err != nil {
		return handleAPIError(err)
	}

	var convID string
	for event := range ch {
		if event.Error != nil {
			fmt.Fprintln(os.Stderr)
			return handleAPIError(event.Error)
		}
		if event.Chunk == nil {
			continue
		}

		if len(event.Chunk.Choices) > 0 {
			delta := event.Chunk.Choices[0].Delta
			if delta.Content != "" {
				fmt.Print(delta.Content)
			}
		}

		if len(event.RawJSON) > 0 {
			var extra struct {
				ConversationID string `json:"conversation_id"`
			}
			if err := json.Unmarshal(event.RawJSON, &extra); err == nil && extra.ConversationID != "" {
				convID = extra.ConversationID
			}
		}
	}

	fmt.Println()
	if convID != "" {
		fmt.Fprintf(os.Stderr, "conversation: %s\n", convID)
	}

	return nil
}

func askNonStreaming(req cadreen.ChatCompletionRequest) error {
	client := newClient()
	resp, err := client.ChatCompletions(context.Background(), req)
	if err != nil {
		return handleAPIError(err)
	}

	if len(resp.Choices) > 0 {
		fmt.Println(resp.Choices[0].Message.Content)
	}

	return nil
}

func handleAPIError(err error) error {
	var apiErr *cadreen.APIError
	if errors.As(err, &apiErr) {
		output.PrintError(apiErr.Message)
		if apiErr.Hint != "" {
			output.PrintHint(apiErr.Hint)
		}
		if apiErr.NextAction != nil {
			output.PrintTry(apiErr.NextAction.Label)
		}
		return err
	}

	output.PrintError(err.Error())
	return err
}

func readJSON(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}
