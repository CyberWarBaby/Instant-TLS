package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/CyberWarBaby/Instant-TLS/cli/internal/api"
	"github.com/CyberWarBaby/Instant-TLS/cli/internal/cert"
	"github.com/CyberWarBaby/Instant-TLS/cli/internal/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "One-time setup: login, create CA, and trust in all browsers",
	Long: `Complete InstantTLS setup in one command.

This will:
  1. Prompt for your API token (from the web dashboard)
  2. Create a local Certificate Authority
  3. Install CA in system trust store
  4. Install CA in Chrome/Chromium (NSS database)
  5. Install CA in Firefox (if installed)

Run with sudo for automatic trust store installation:
  sudo instanttls setup

After setup, just run:
  instanttls cert myapp.local`,
	Run: runSetup,
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func runSetup(cmd *cobra.Command, args []string) {
	pterm.Println()
	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgMagenta)).
		WithTextStyle(pterm.NewStyle(pterm.FgWhite)).
		Println("🚀 InstantTLS Setup")
	pterm.Println()

	// Check if already set up
	if config.Exists() && cert.CAExists() {
		pterm.Warning.Println("InstantTLS is already set up!")
		pterm.Println()
		pterm.Info.Println("To generate certificates, run:")
		pterm.Println("  instanttls cert myapp.local")
		pterm.Println()
		pterm.Info.Println("To re-run setup, first remove the config:")
		pterm.Println("  rm -rf ~/.instanttls")
		return
	}

	// Step 1: Login
	pterm.DefaultSection.Println("Step 1/3: Authentication")

	reader := bufio.NewReader(os.Stdin)

	pterm.Info.Println("Enter API base URL")
	pterm.FgGray.Print("  (default: http://localhost:8081): ")
	apiBaseURL, _ := reader.ReadString('\n')
	apiBaseURL = strings.TrimSpace(apiBaseURL)
	if apiBaseURL == "" {
		apiBaseURL = "http://localhost:8081"
	}

	pterm.Println()
	pterm.Info.Println("Enter your Personal Access Token")
	pterm.FgGray.Println("  (Get one from the web dashboard → Tokens)")
	pterm.FgGray.Print("  Token: ")

	tokenBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		token, _ := reader.ReadString('\n')
		tokenBytes = []byte(strings.TrimSpace(token))
	}
	token := strings.TrimSpace(string(tokenBytes))
	pterm.Println()

	if token == "" {
		printError("Token cannot be empty")
		return
	}

	spinner, _ := pterm.DefaultSpinner.Start("Validating token...")
	client := api.NewClient(apiBaseURL, token)
	user, err := client.Me()
	if err != nil {
		spinner.Fail("Authentication failed")
		printError(fmt.Sprintf("Failed to authenticate: %v", err))
		return
	}
	spinner.Success(fmt.Sprintf("Logged in as %s", user.Email))

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

	// Step 2: Generate CA
	pterm.Println()
	pterm.DefaultSection.Println("Step 2/3: Create Certificate Authority")

	spinner, _ = pterm.DefaultSpinner.Start("Generating CA certificate...")
	if err := cert.GenerateCA(); err != nil {
		spinner.Fail("Failed to generate CA")
		printError(err.Error())
		return
	}
	spinner.Success("CA certificate created")

	// Step 3: Install in trust stores
	pterm.Println()
	pterm.DefaultSection.Println("Step 3/4: Install in Trust Stores")

	caPath := filepath.Join(config.GetCADir(), "ca.crt")

	// System trust store
	spinner, _ = pterm.DefaultSpinner.Start("Installing in system trust store...")
	if err := installSystemTrust(caPath); err != nil {
		spinner.Warning(fmt.Sprintf("System trust: %v", err))
	} else {
		spinner.Success("Installed in system trust store")
	}

	// Chrome/Chromium (NSS)
	spinner, _ = pterm.DefaultSpinner.Start("Installing in Chrome/Chromium...")
	if err := installChromeTrust(caPath); err != nil {
		spinner.Warning(fmt.Sprintf("Chrome: %v", err))
	} else {
		spinner.Success("Installed in Chrome/Chromium")
	}

	// Firefox
	spinner, _ = pterm.DefaultSpinner.Start("Installing in Firefox...")
	if err := installFirefoxTrust(caPath); err != nil {
		spinner.Warning(fmt.Sprintf("Firefox: %v", err))
	} else {
		spinner.Success("Installed in Firefox")
	}

	// Step 4: Add to PATH
	pterm.Println()
	pterm.DefaultSection.Println("Step 4/4: Make CLI Available System-Wide")

	spinner, _ = pterm.DefaultSpinner.Start("Adding instanttls to system PATH...")
	if err := installToPath(); err != nil {
		spinner.Warning(fmt.Sprintf("PATH: %v (you may need to add Go bin to PATH manually)", err))
	} else {
		spinner.Success("instanttls is now available system-wide")
	}

	// Done!
	pterm.Println()
	pterm.DefaultBox.WithTitle("🎉 Setup Complete!").
		WithTitleTopCenter().
		WithBoxStyle(pterm.NewStyle(pterm.FgGreen)).
		Println(`
Your local CA is now trusted by all browsers.

Generate certificates with:
  instanttls cert myapp.local

Don't forget to add your domain to /etc/hosts:
  echo "127.0.0.1 myapp.local" | sudo tee -a /etc/hosts

Then use in your project:
  - Certificate: ~/.instanttls/certs/myapp.local/cert.pem
  - Private Key: ~/.instanttls/certs/myapp.local/key.pem
`)
}

func installSystemTrust(caPath string) error {
	// Check if running as root
	if os.Geteuid() != 0 {
		// Try with sudo
		cmd := exec.Command("sudo", "cp", caPath, "/usr/local/share/ca-certificates/instanttls.crt")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("need sudo access")
		}
		cmd = exec.Command("sudo", "update-ca-certificates")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// Running as root
	if err := exec.Command("cp", caPath, "/usr/local/share/ca-certificates/instanttls.crt").Run(); err != nil {
		return err
	}
	return exec.Command("update-ca-certificates").Run()
}

func installChromeTrust(caPath string) error {
	// Get actual user's home (not root's home when using sudo)
	home := os.Getenv("HOME")
	if sudoUser := os.Getenv("SUDO_USER"); sudoUser != "" {
		home = "/home/" + sudoUser
	}

	nssDB := filepath.Join(home, ".pki", "nssdb")

	// Create NSS database directory if it doesn't exist
	if _, err := os.Stat(nssDB); os.IsNotExist(err) {
		os.MkdirAll(nssDB, 0755)
	}

	// Remove existing cert if any
	exec.Command("certutil", "-d", "sql:"+nssDB, "-D", "-n", "InstantTLS").Run()

	// Add the CA certificate
	cmd := exec.Command("certutil", "-d", "sql:"+nssDB, "-A", "-t", "C,,", "-n", "InstantTLS", "-i", caPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}

	// Fix ownership if running as sudo
	if sudoUser := os.Getenv("SUDO_USER"); sudoUser != "" {
		exec.Command("chown", "-R", sudoUser+":"+sudoUser, nssDB).Run()
	}

	return nil
}

func installFirefoxTrust(caPath string) error {
	// Get actual user's home
	home := os.Getenv("HOME")
	if sudoUser := os.Getenv("SUDO_USER"); sudoUser != "" {
		home = "/home/" + sudoUser
	}

	firefoxDir := filepath.Join(home, ".mozilla", "firefox")
	if _, err := os.Stat(firefoxDir); os.IsNotExist(err) {
		return fmt.Errorf("Firefox not installed")
	}

	// Find all Firefox profiles and add cert to each
	profiles, err := filepath.Glob(filepath.Join(firefoxDir, "*.default*"))
	if err != nil || len(profiles) == 0 {
		return fmt.Errorf("no Firefox profiles found")
	}

	for _, profile := range profiles {
		// Remove existing cert if any
		exec.Command("certutil", "-d", "sql:"+profile, "-D", "-n", "InstantTLS").Run()

		// Add the CA certificate
		cmd := exec.Command("certutil", "-d", "sql:"+profile, "-A", "-t", "C,,", "-n", "InstantTLS", "-i", caPath)
		cmd.Run() // Ignore errors for individual profiles
	}

	// Fix ownership if running as sudo
	if sudoUser := os.Getenv("SUDO_USER"); sudoUser != "" {
		exec.Command("chown", "-R", sudoUser+":"+sudoUser, firefoxDir).Run()
	}

	return nil
}

func installToPath() error {
	// Get the path to the current executable
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot find executable path: %w", err)
	}

	// Resolve symlinks to get the real path
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("cannot resolve executable path: %w", err)
	}

	// Target path in /usr/local/bin
	targetPath := "/usr/local/bin/instanttls"

	// Remove existing symlink if it exists
	os.Remove(targetPath)

	// Create symlink
	if err := os.Symlink(execPath, targetPath); err != nil {
		// If symlink fails, try copying
		input, err := os.ReadFile(execPath)
		if err != nil {
			return fmt.Errorf("cannot read executable: %w", err)
		}
		if err := os.WriteFile(targetPath, input, 0755); err != nil {
			return fmt.Errorf("cannot copy to /usr/local/bin: %w", err)
		}
	}

	return nil
}
