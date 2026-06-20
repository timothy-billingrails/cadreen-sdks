package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/config"
	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	cadreen "github.com/timothy-billingrails/cadreen-sdks/go/cadreen"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Cadreen",
	Long: `Authenticate with Cadreen.

Paste your API key or press Enter to open the browser.

Your API key is saved to ~/.cadreen/config.json.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenFlag, _ := cmd.Flags().GetBool("token")

		if cfg.IsAuthenticated() && !tokenFlag {
			fmt.Println("Already authenticated.")
			fmt.Printf("  API Key: %s\n", output.MaskKey(cfg.APIKeyResolved()))
			fmt.Println()
			fmt.Println("Run 'cadreen logout' first to switch accounts.")
			return nil
		}

		if tokenFlag || !hasBrowser() {
			return loginWithToken()
		}
		return loginInteractive()
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cfg.Delete(); err != nil {
			return err
		}
		fmt.Println("Done. Credentials removed.")
		return nil
	},
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !cfg.IsAuthenticated() {
			fmt.Println("Not authenticated.")
			fmt.Println("Run 'cadreen init' or 'cadreen login' to get started.")
			return nil
		}

		resp, err := rawDo("GET", "/api/v1/cadreen/health", nil)
		if err != nil {
			output.PrintError("Can't reach the API.")
			output.PrintHint(err.Error())
			output.PrintTry("cadreen login — to re-authenticate")
			return nil
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			output.PrintError("Your API key isn't working.")
			output.PrintHint("It might be expired, revoked, or typed wrong.")
			output.PrintTry("cadreen login — to authenticate again")
			return nil
		}

		var health struct {
			Status      string `json:"status"`
			Workspace   string `json:"workspace"`
			TenantID    string `json:"tenant_id"`
			Connections struct {
				Total   int `json:"total"`
				Healthy int `json:"healthy"`
			} `json:"connections"`
			Governance struct {
				Active int `json:"active"`
			} `json:"governance"`
			Memory struct {
				Status string `json:"status"`
			} `json:"memory"`
		}
		if err := readJSON(resp.Body, &health); err != nil {
			return err
		}

		fmt.Printf("Authenticated as: %s\n", cfg.WorkspaceID)
		fmt.Printf("API Key:         %s\n", output.MaskKey(cfg.APIKeyResolved()))
		fmt.Printf("Base URL:        %s\n", cfg.BaseURLResolved())
		fmt.Printf("Status:          %s\n", health.Status)

		return nil
	},
}

func init() {
	loginCmd.Flags().Bool("token", false, "paste API key directly (no browser)")
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(whoamiCmd)
}

func loginInteractive() error {
	fmt.Println("How would you like to authenticate?")
	fmt.Println()
	fmt.Println("  1) Paste an API key (I already have one)")
	fmt.Println("  2) Open browser (I need to get a key)")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Choose [1]: ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	if choice == "2" {
		return loginWithBrowser()
	}
	return loginWithToken()
}

func loginWithToken() error {
	fmt.Print("Enter your API key: ")
	reader := bufio.NewReader(os.Stdin)
	key, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return fmt.Errorf("reading input: %w", err)
	}
	key = strings.TrimSpace(key)

	if key == "" {
		return fmt.Errorf("no API key provided")
	}

	cfg.APIKey = key
	if err := validateAndSave(); err != nil {
		cfg.APIKey = ""
		return err
	}

	return nil
}

func loginWithBrowser() error {
	fmt.Println("Opening browser for authentication...")
	fmt.Println("If the browser doesn't open, visit:")
	fmt.Printf("  %s/infra/auth/cli\n", cfg.BaseURLResolved())
	fmt.Println()
	fmt.Println("(Browser OAuth not yet implemented. Use paste method for now.)")
	fmt.Println()
	return loginWithToken()
}

func validateAndSave() error {
	resp, err := rawDo("GET", "/api/v1/cadreen/health", nil)
	if err != nil {
		output.PrintError("Can't reach the API.")
		output.PrintHint(err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		output.PrintError("Your API key isn't working.")
		output.PrintHint("It might be expired, revoked, or typed wrong.")
		output.PrintTry("cadreen login — to authenticate again")
		return fmt.Errorf("authentication failed")
	}

	if err := cfg.Save(); err != nil {
		return err
	}

	fmt.Println()
	output.PrintSuccess("Authenticated.")
	fmt.Printf("  API Key: %s\n", output.MaskKey(cfg.APIKeyResolved()))
	fmt.Printf("  Config:  %s\n", config.ConfigPath())
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  cadreen ask \"what can you do?\"")
	fmt.Println("  cadreen chat")
	fmt.Println("  cadreen status")

	return nil
}

func hasBrowser() bool {
	return os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != "" ||
		os.Getenv("TERM_PROGRAM") != "" || os.Getenv("WT_SESSION") != ""
}

func newClient() *cadreen.Client {
	return cadreen.NewClient(cadreen.CadreenConfig{
		APIKey:  cfg.APIKeyResolved(),
		BaseURL: cfg.BaseURLResolved(),
	})
}
