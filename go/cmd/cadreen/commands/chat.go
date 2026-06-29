package commands

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	cadreen "github.com/timothy-billingrails/cadreen-sdks/go/cadreen"

	"github.com/spf13/cobra"
)

type chatState int

const (
	chatStateIdle chatState = iota
	chatStateStreaming
	chatStateAwaitingConfirmation
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start an interactive chat session",
	Long: `Start an interactive chat session with Cadreen.

Your conversation is persisted automatically. Governance prompts
appear inline when Cadreen needs your permission.

Type 'exit' or 'quit' to leave. Type 'clear' to reset.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		memoryOff, _ := cmd.Flags().GetBool("no-memory")
		convID, _ := cmd.Flags().GetString("conversation-id")
		userID, _ := cmd.Flags().GetString("user-id")

		reader := bufio.NewReader(os.Stdin)
		state := chatStateIdle
		var pendingActionsLocal []pendingAction

		fmt.Println("Cadreen Chat (type 'exit' to quit, 'clear' to reset)")
		fmt.Println()

		for {
			switch state {
			case chatStateIdle:
				fmt.Print("You: ")
				input, err := reader.ReadString('\n')
				if err != nil {
					return nil
				}
				input = strings.TrimSpace(input)

				if input == "" {
					continue
				}
				if input == "exit" || input == "quit" {
					fmt.Println("Goodbye.")
					return nil
				}
				if input == "clear" {
					convID = ""
					fmt.Println("Conversation cleared.")
					fmt.Println()
					continue
				}

				req := cadreen.ChatCompletionRequest{
					Messages: []cadreen.ChatMessage{
						{Role: "user", Content: input},
					},
				}
				if memoryOff {
					req.Context = map[string]any{"memory": false}
				}
				if convID != "" {
					req.ConversationID = convID
				}
				if userID != "" {
					req.UserID = userID
				}

				state = chatStateStreaming
				pendingActions, newConvID, err := chatStreamResponse(req)
				if err != nil {
					state = chatStateIdle
					continue
				}
				if newConvID != "" {
					convID = newConvID
				}

				if len(pendingActions) > 0 {
					state = chatStateAwaitingConfirmation
					pendingActionsLocal = pendingActions
				} else {
					state = chatStateIdle
					fmt.Println()
				}

			case chatStateAwaitingConfirmation:
				fmt.Println()
				showConfirmationPrompt(pendingActionsLocal)
				fmt.Print("> ")

				input, err := reader.ReadString('\n')
				if err != nil {
					return nil
				}
				input = strings.TrimSpace(strings.ToLower(input))

				if input == "yes" || input == "y" {
					req := cadreen.ChatCompletionRequest{
						Messages: []cadreen.ChatMessage{
							{Role: "user", Content: "yes"},
						},
					}
					if convID != "" {
						req.ConversationID = convID
					}

					state = chatStateStreaming
					_, newConvID, err := chatStreamResponse(req)
					if err != nil {
						state = chatStateIdle
						continue
					}
					if newConvID != "" {
						convID = newConvID
					}
					state = chatStateIdle
					fmt.Println()
				} else if input == "no" || input == "n" {
					req := cadreen.ChatCompletionRequest{
						Messages: []cadreen.ChatMessage{
							{Role: "user", Content: "no"},
						},
					}
					if convID != "" {
						req.ConversationID = convID
					}

					state = chatStateStreaming
					_, _, err := chatStreamResponse(req)
					if err != nil {
						state = chatStateIdle
						continue
					}
					state = chatStateIdle
					fmt.Println()
				} else {
					fmt.Println("Cancelled.")
					state = chatStateIdle
					fmt.Println()
				}
			}
		}
	},
}

func init() {
	chatCmd.Flags().Bool("no-memory", false, "disable cross-conversation memory")
	chatCmd.Flags().String("conversation-id", "", "resume a specific conversation")
	chatCmd.Flags().String("user-id", "", "end-user ID for permission scoping")
	rootCmd.AddCommand(chatCmd)
}

func chatStreamResponse(req cadreen.ChatCompletionRequest) ([]pendingAction, string, error) {
	client := newClient()
	ch, err := client.ChatCompletionsStream(context.Background(), req)
	if err != nil {
		handleAPIError(err)
		return nil, "", err
	}

	fmt.Print("Cadreen: ")

	var convID string
	var pendingActions []pendingAction

	for event := range ch {
		if event.Error != nil {
			fmt.Fprintln(os.Stderr)
			handleAPIError(event.Error)
			return nil, "", event.Error
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
			var extra streamExtra
			if err := json.Unmarshal(event.RawJSON, &extra); err == nil {
				if extra.ConversationID != "" {
					convID = extra.ConversationID
				}
				if len(extra.PendingActions) > 0 {
					pendingActions = extra.PendingActions
				}
			}
		}
	}

	return pendingActions, convID, nil
}

type streamExtra struct {
	ConversationID string          `json:"conversation_id"`
	PendingActions []pendingAction `json:"pending_actions"`
}

type pendingAction struct {
	Tool        string `json:"tool"`
	ConfirmWith string `json:"confirm_with"`
	ExpiresAt   int64  `json:"expires_at"`
}

func showConfirmationPrompt(actions []pendingAction) {
	fmt.Println()
	fmt.Println("Before I can do this, I need your permission.")
	fmt.Println()

	if len(actions) == 1 {
		fmt.Printf("I want to: %s\n", humanToolName(actions[0].Tool))
	} else {
		fmt.Printf("I need your permission for %d things:\n\n", len(actions))
		for i, a := range actions {
			fmt.Printf("  %d. %s\n", i+1, humanToolName(a.Tool))
		}
	}

	fmt.Println()
	fmt.Println("Proceed? (yes/no)")
}

func humanToolName(tool string) string {
	return strings.ReplaceAll(tool, "_", " ")
}
