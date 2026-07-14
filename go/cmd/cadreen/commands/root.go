package commands

import (
	"fmt"
	"os"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/config"
	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	Version = "0.4.0"
	Commit  = "unknown"
	Date    = "unknown"

	cfgFile string
	apiKey  string
	format  string
	verbose bool
	cfg     *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "cadreen",
	Short: "Intelligence infrastructure for developers",
	Long: `Cadreen — intelligence infrastructure for developers.

Cadreen remembers things, follows rules, connects to services,
and heals itself when things go wrong. All governed. All observable.

Run 'cadreen init' to get started.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if cfgFile != "" {
			cfg, err = config.LoadFrom(cfgFile)
		} else {
			cfg, err = config.Load()
		}
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
		if apiKey != "" {
			cfg.APIKey = apiKey
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if !cfg.IsAuthenticated() {
			fmt.Println("Welcome to Cadreen. Let's get you set up.")
			fmt.Println()
			fmt.Println("  cadreen init    — Set up your account")
			fmt.Println("  cadreen login   — Authenticate with an API key")
			fmt.Println()
			fmt.Println("Run 'cadreen init' to get started.")
			return nil
		}

		memoryOff, _ := cmd.Flags().GetBool("no-memory")
		model := initialTUIModel(memoryOff)
		p := tea.NewProgram(model, tea.WithAltScreen())
		_, err := p.Run()
		return err
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default ~/.cadreen/config.json)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key (overrides config and env)")
	rootCmd.PersistentFlags().StringVar(&format, "format", "", "output format: text, json, verbose")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "show detailed output")
	rootCmd.Flags().Bool("no-memory", false, "disable cross-conversation memory (TUI mode)")

	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cadreen %s (commit: %s, built: %s)\n", Version, Commit, Date)
	},
}

func outputFormat() output.Format {
	if format != "" {
		return output.ResolveFormat(format)
	}
	if verbose {
		return output.FormatVerbose
	}
	if cfg != nil && cfg.Output == "json" {
		return output.FormatJSON
	}
	return output.FormatText
}

func requireAuth() {
	if cfg == nil || !cfg.IsAuthenticated() {
		output.PrintError("Not authenticated.")
		output.PrintHint("Run 'cadreen init' to set up your account.")
		output.PrintTry("cadreen init")
		os.Exit(1)
	}
}
