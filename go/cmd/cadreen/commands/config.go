package commands

import (
	"fmt"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/config"
	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage local configuration",
	Long: `Manage your local Cadreen configuration.

View and edit settings stored in ~/.cadreen/config.json.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return configListCmd.RunE(configListCmd, args)
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(cfg, format)
			return nil
		}

		fmt.Println("Current configuration:")
		fmt.Println()
		fmt.Printf("  API Key:      %s\n", maskedOr("(not set)", cfg.APIKey))
		fmt.Printf("  Base URL:     %s\n", cfg.BaseURLResolved())
		fmt.Printf("  Workspace:    %s\n", orDefault(cfg.WorkspaceID, "(not set)"))
		fmt.Printf("  Profile:      %s\n", orDefault(cfg.Profile, "full"))
		fmt.Printf("  Output:       %s\n", orDefault(cfg.Output, "text"))
		fmt.Println()
		fmt.Printf("  Config file:  %s\n", config.ConfigPath())
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		val := cfg.Get(args[0])
		if val == "" {
			fmt.Printf("%s: (not set)\n", args[0])
		} else if args[0] == "api_key" {
			fmt.Printf("%s: %s\n", args[0], output.MaskKey(val))
		} else {
			fmt.Printf("%s: %s\n", args[0], val)
		}
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cfg.Set(args[0], args[1]); err != nil {
			return err
		}
		if args[0] == "api_key" {
			fmt.Printf("Set %s = %s\n", args[0], output.MaskKey(args[1]))
		} else {
			fmt.Printf("Set %s = %s\n", args[0], args[1])
		}
		return nil
	},
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset to defaults",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg = &config.Config{
			BaseURL: config.DefaultBaseURL,
			Profile: config.DefaultProfile,
			Output:  config.DefaultOutput,
		}
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Println("Configuration reset to defaults.")
		return nil
	},
}

func init() {
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configResetCmd)
	rootCmd.AddCommand(configCmd)
}

func maskedOr(def, val string) string {
	if val == "" {
		return def
	}
	return output.MaskKey(val)
}

func orDefault(val, def string) string {
	if val == "" {
		return def
	}
	return val
}
