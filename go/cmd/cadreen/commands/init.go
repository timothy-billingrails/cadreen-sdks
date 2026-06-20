package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Set up Cadreen in the current directory",
	Long: `Set up Cadreen.

Checks for an API key, validates it, and saves your config.
Run this once to get started.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.IsAuthenticated() {
			fmt.Println("Already authenticated.")
			fmt.Printf("  API Key: %s\n", output.MaskKey(cfg.APIKeyResolved()))
			fmt.Println()
			runDoctor(cmd, args)
			return nil
		}

		envKey := os.Getenv("CADREEN_API_KEY")
		if envKey != "" {
			cfg.APIKey = envKey
			if err := validateAndSave(); err != nil {
				cfg.APIKey = ""
				return err
			}
			runDoctor(cmd, args)
			return nil
		}

		fmt.Println("Welcome to Cadreen. Let's get you set up.")
		fmt.Println()

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter your API key (or press Enter to open browser): ")
		key, _ := reader.ReadString('\n')
		key = strings.TrimSpace(key)

		if key == "" {
			return loginWithBrowser()
		}

		cfg.APIKey = key
		if err := validateAndSave(); err != nil {
			cfg.APIKey = ""
			return err
		}

		runDoctor(cmd, args)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
