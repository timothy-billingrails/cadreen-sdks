package commands

import (
	"encoding/json"
	"fmt"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	"github.com/spf13/cobra"
)

var learningCmd = &cobra.Command{
	Use:   "learning",
	Short: "Learning insights (patterns, episodes, suggestions)",
	Long: `View detected patterns, learning episodes, and improvement suggestions.

Examples:
  cadreen learning patterns
  cadreen learning episodes
  cadreen learning suggestions`,
}

var learningPatternsCmd = &cobra.Command{
	Use:   "patterns",
	Short: "Detected patterns",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/learning/patterns", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result struct {
			Patterns []struct {
				ID          string   `json:"id"`
				Pattern     string   `json:"pattern"`
				Confidence  float64  `json:"confidence"`
				Occurrences int      `json:"occurrences"`
				Domain      string   `json:"domain"`
			} `json:"patterns"`
			Count int `json:"count"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(result, format)
		default:
			if len(result.Patterns) == 0 {
				fmt.Println("No patterns detected yet.")
				return nil
			}
			fmt.Printf("Detected patterns (%d):\n\n", result.Count)
			for _, p := range result.Patterns {
				fmt.Printf("  %s\n", p.Pattern)
				fmt.Printf("    Confidence: %.0f%%", p.Confidence*100)
				if p.Domain != "" {
					fmt.Printf("  Domain: %s", p.Domain)
				}
				fmt.Println()
			}
		}

		return nil
	},
}

var learningEpisodesCmd = &cobra.Command{
	Use:   "episodes",
	Short: "Learning episodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/learning/episodes", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result struct {
			Episodes []struct {
				ID          string `json:"id"`
				Description string `json:"description"`
				Outcome     string `json:"outcome"`
				Domain      string `json:"domain"`
			} `json:"episodes"`
			Count int `json:"count"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(result, format)
		default:
			if len(result.Episodes) == 0 {
				fmt.Println("No learning episodes yet.")
				return nil
			}
			fmt.Printf("Learning episodes (%d):\n\n", result.Count)
			for _, e := range result.Episodes {
				fmt.Printf("  %s\n", e.Description)
				if e.Outcome != "" {
					fmt.Printf("    Outcome: %s\n", e.Outcome)
				}
				fmt.Println()
			}
		}

		return nil
	},
}

var learningSuggestionsCmd = &cobra.Command{
	Use:   "suggestions",
	Short: "Improvement suggestions",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		resp, err := rawDo("GET", "/api/v1/cadreen/learning/suggestions", nil)
		if err != nil {
			return handleAPIError(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return handleHTTPError(resp)
		}

		var result struct {
			Suggestions []struct {
				ID          string `json:"id"`
				Type        string `json:"type"`
				Description string `json:"description"`
				Impact      string `json:"impact"`
			} `json:"suggestions"`
			Count int `json:"count"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return err
		}

		format := outputFormat()
		switch format {
		case output.FormatJSON:
			output.Print(result, format)
		default:
			if len(result.Suggestions) == 0 {
				fmt.Println("No suggestions at this time.")
				return nil
			}
			fmt.Printf("Suggestions (%d):\n\n", result.Count)
			for _, s := range result.Suggestions {
				fmt.Printf("  [%s] %s\n", s.Type, s.Description)
				if s.Impact != "" {
					fmt.Printf("    Impact: %s\n", s.Impact)
				}
				fmt.Println()
			}
		}

		return nil
	},
}

func init() {
	learningCmd.AddCommand(learningPatternsCmd)
	learningCmd.AddCommand(learningEpisodesCmd)
	learningCmd.AddCommand(learningSuggestionsCmd)
	rootCmd.AddCommand(learningCmd)
}
