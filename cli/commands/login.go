package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/CyberWarBaby/Instant-TLS/cli/internal/api"
	"github.com/CyberWarBaby/Instant-TLS/cli/internal/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with your Personal Access Token",
	Long: `Login to InstantTLS with your Personal Access Token (PAT).

Get your token from the InstantTLS dashboard at https://instanttls.dev/app/tokens

Example:
  instanttls login`,
	Run: runLogin,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func runLogin(cmd *cobra.Command, args []string) {
	pterm.Println()
	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).
		WithTextStyle(pterm.NewStyle(pterm.FgBlack)).
		Println("🔐 InstantTLS Login")
	pterm.Println()

	reader := bufio.NewReader(os.Stdin)

	// Prompt for API base URL
	pterm.Info.Println("Enter API base URL")
	pterm.FgGray.Print("  (default: http://localhost:8081): ")

	apiBaseURL, _ := reader.ReadString('\n')
	apiBaseURL = strings.TrimSpace(apiBaseURL)
	if apiBaseURL == "" {
		apiBaseURL = "http://localhost:8081"
	}

	// Prompt for token
	pterm.Println()
	pterm.Info.Println("Enter your Personal Access Token")
	pterm.FgGray.Print("  Token: ")

	// Read token with hidden input
	tokenBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		// Fallback to regular input if term.ReadPassword fails
		token, _ := reader.ReadString('\n')
		tokenBytes = []byte(strings.TrimSpace(token))
	}
	token := strings.TrimSpace(string(tokenBytes))
	pterm.Println()
	pterm.Println()

	if token == "" {
		printError("Token cannot be empty")
		return
	}

	// Validate token
	spinner, _ := pterm.DefaultSpinner.Start("Validating token...")

	client := api.NewClient(apiBaseURL, token)
	user, err := client.Me()
	if err != nil {
		spinner.Fail("Authentication failed")
		pterm.Println()
		printError(fmt.Sprintf("Failed to authenticate: %v", err))
		return
	}

	spinner.Success("Token validated!")

	// Save config
	tokenPrefix := token
	if len(tokenPrefix) > 12 {
		tokenPrefix = tokenPrefix[:12]
	}

	cfg := &config.Config{
		APIBaseURL:  apiBaseURL,
		Token:       token,
		TokenPrefix: tokenPrefix,
		Email:       user.Email,
		Plan:        user.Plan,
	}

	if err := config.Save(cfg); err != nil {
		printError(fmt.Sprintf("Failed to save config: %v", err))
		return
	}

	pterm.Println()
	pterm.DefaultBox.WithTitle("✅ Login Successful").
		WithTitleTopCenter().
		WithBoxStyle(pterm.NewStyle(pterm.FgGreen)).
		Println(fmt.Sprintf(`
  Email: %s
  Plan:  %s
  API:   %s
`, user.Email, planBadge(user.Plan), apiBaseURL))

	pterm.Println()
	pterm.Info.Println("Next steps:")
	pterm.Println("  1. Run 'instanttls init' to create your local CA")
	pterm.Println("  2. Run 'instanttls cert \"*.local.test\"' to generate a certificate")
	pterm.Println("  3. Run 'instanttls doctor' to verify your setup")
	pterm.Println()
}

func planBadge(plan string) string {
	switch plan {
	case "pro":
		return pterm.FgMagenta.Sprint("⭐ Pro")
	case "team":
		return pterm.FgCyan.Sprint("🏢 Team")
	default:
		return pterm.FgGray.Sprint("Free")
	}
}
