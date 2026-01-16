package cmd

import (
	"fmt"

	"github.com/instanttls/cli/internal/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Display current user and plan information",
	Long: `Show information about the currently logged-in user.

Example:
  instanttls whoami`,
	Run: runWhoami,
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}

func runWhoami(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil || cfg == nil || cfg.Token == "" {
		printError("Not logged in. Run 'instanttls login' first.")
		return
	}

	pterm.Println()
	pterm.DefaultBox.WithTitle("👤 Current User").
		WithTitleTopCenter().
		Println(fmt.Sprintf(`
  Email:        %s
  Plan:         %s
  API URL:      %s
  Token Prefix: %s...
`, cfg.Email, planBadge(cfg.Plan), cfg.APIBaseURL, cfg.TokenPrefix))
	pterm.Println()
}
