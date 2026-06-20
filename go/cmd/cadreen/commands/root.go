package commands

import (
	"fmt"
	"os"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/config"
	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	"github.com/spf13/cobra"
)

var (
	Version = "dev"
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

Your AI remembers things, follows rules, connects to services,
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
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default ~/.cadreen/config.json)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key (overrides config and env)")
	rootCmd.PersistentFlags().StringVar(&format, "format", "", "output format: text, json, verbose")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "show detailed output")

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
		output.PrintHint("Run 'cadreen init' or 'cadreen login' to authenticate.")
		os.Exit(1)
	}
}
