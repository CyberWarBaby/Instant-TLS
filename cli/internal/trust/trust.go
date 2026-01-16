package trust

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/CyberWarBaby/Instant-TLS/cli/internal/config"
	"github.com/pterm/pterm"
)

// InstallCA installs the CA certificate into the OS trust store
func InstallCA() error {
	caDir := config.GetCADir()
	certPath := filepath.Join(caDir, "ca.crt")

	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return fmt.Errorf("CA certificate not found. Run 'instanttls init' first")
	}

	switch runtime.GOOS {
	case "darwin":
		return installDarwin(certPath)
	case "linux":
		return installLinux(certPath)
	case "windows":
		return installWindows(certPath)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func installDarwin(certPath string) error {
	cmd := fmt.Sprintf("sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain %s", certPath)

	pterm.Info.Println("Installing CA certificate into macOS Keychain...")
	pterm.Println()
	pterm.DefaultBox.WithTitle("Command to run").Println(cmd)
	pterm.Println()

	result, _ := pterm.DefaultInteractiveConfirm.
		WithDefaultValue(true).
		Show("This requires sudo. Continue?")

	if !result {
		return fmt.Errorf("installation cancelled by user")
	}

	cmdExec := exec.Command("sudo", "security", "add-trusted-cert", "-d", "-r", "trustRoot", "-k", "/Library/Keychains/System.keychain", certPath)
	cmdExec.Stdin = os.Stdin
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr

	if err := cmdExec.Run(); err != nil {
		return fmt.Errorf("failed to install CA: %w", err)
	}

	return nil
}

func installLinux(certPath string) error {
	// Check for common Linux distributions
	destPath := "/usr/local/share/ca-certificates/instanttls.crt"

	cmd := fmt.Sprintf("sudo cp %s %s && sudo update-ca-certificates", certPath, destPath)

	pterm.Info.Println("Installing CA certificate into Linux trust store...")
	pterm.Println()
	pterm.DefaultBox.WithTitle("Commands to run").Println(cmd)
	pterm.Println()

	// Check if update-ca-certificates exists
	if _, err := exec.LookPath("update-ca-certificates"); err != nil {
		pterm.Warning.Println("update-ca-certificates not found.")
		pterm.Println()
		pterm.DefaultBox.WithTitle("Manual Installation").Println(`
For Fedora/RHEL/CentOS:
  sudo cp ` + certPath + ` /etc/pki/ca-trust/source/anchors/
  sudo update-ca-trust

For Arch Linux:
  sudo trust anchor --store ` + certPath + `

For other distributions, consult your documentation.
`)
		return fmt.Errorf("automatic installation not supported on this distribution")
	}

	result, _ := pterm.DefaultInteractiveConfirm.
		WithDefaultValue(true).
		Show("This requires sudo. Continue?")

	if !result {
		return fmt.Errorf("installation cancelled by user")
	}

	// Copy certificate
	cpCmd := exec.Command("sudo", "cp", certPath, destPath)
	cpCmd.Stdin = os.Stdin
	cpCmd.Stdout = os.Stdout
	cpCmd.Stderr = os.Stderr

	if err := cpCmd.Run(); err != nil {
		return fmt.Errorf("failed to copy CA certificate: %w", err)
	}

	// Update CA certificates
	updateCmd := exec.Command("sudo", "update-ca-certificates")
	updateCmd.Stdin = os.Stdin
	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr

	if err := updateCmd.Run(); err != nil {
		return fmt.Errorf("failed to update CA certificates: %w", err)
	}

	// Also install to Chrome/Chromium NSS database
	installChromeNSS(certPath)

	return nil
}

// installChromeNSS installs the CA to Chrome/Chromium's NSS database on Linux
func installChromeNSS(certPath string) {
	// Check if certutil is available
	if _, err := exec.LookPath("certutil"); err != nil {
		pterm.Println()
		pterm.Warning.Println("Chrome/Chromium uses its own certificate store.")
		pterm.Info.Println("Install libnss3-tools and run:")
		pterm.Println()
		pterm.DefaultBox.Println(`sudo apt install libnss3-tools
certutil -d sql:$HOME/.pki/nssdb -A -t "C,," -n "InstantTLS Local CA" -i ` + certPath)
		return
	}

	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	nssDB := filepath.Join(homeDir, ".pki", "nssdb")

	// Check if NSS database exists
	if _, err := os.Stat(nssDB); os.IsNotExist(err) {
		// Create the directory
		os.MkdirAll(nssDB, 0700)
	}

	pterm.Println()
	pterm.Info.Println("Adding CA to Chrome/Chromium certificate store...")

	cmd := exec.Command("certutil", "-d", "sql:"+nssDB, "-A", "-t", "C,,", "-n", "InstantTLS Local CA", "-i", certPath)
	if err := cmd.Run(); err != nil {
		pterm.Warning.Printf("Failed to add to Chrome store: %v\n", err)
		return
	}

	pterm.Success.Println("CA added to Chrome/Chromium!")
	pterm.FgYellow.Println("  ⚠️  Restart Chrome for changes to take effect")
}

func installWindows(certPath string) error {
	cmd := fmt.Sprintf("certutil -addstore Root \"%s\"", certPath)

	pterm.Info.Println("Installing CA certificate into Windows trust store...")
	pterm.Println()
	pterm.DefaultBox.WithTitle("Command to run").Println(cmd)
	pterm.Println()

	result, _ := pterm.DefaultInteractiveConfirm.
		WithDefaultValue(true).
		Show("This requires administrator privileges. Continue?")

	if !result {
		return fmt.Errorf("installation cancelled by user")
	}

	cmdExec := exec.Command("certutil", "-addstore", "Root", certPath)
	cmdExec.Stdin = os.Stdin
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr

	if err := cmdExec.Run(); err != nil {
		return fmt.Errorf("failed to install CA: %w. Try running as Administrator", err)
	}

	return nil
}

// IsTrusted checks if the CA is installed in the trust store
func IsTrusted() bool {
	caDir := config.GetCADir()
	certPath := filepath.Join(caDir, "ca.crt")

	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return false
	}

	switch runtime.GOOS {
	case "darwin":
		return isTrustedDarwin(certPath)
	case "linux":
		return isTrustedLinux()
	case "windows":
		return isTrustedWindows()
	default:
		return false
	}
}

func isTrustedDarwin(certPath string) bool {
	cmd := exec.Command("security", "find-certificate", "-c", "InstantTLS Local Development CA", "/Library/Keychains/System.keychain")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func isTrustedLinux() bool {
	destPath := "/usr/local/share/ca-certificates/instanttls.crt"
	_, err := os.Stat(destPath)
	return err == nil
}

func isTrustedWindows() bool {
	cmd := exec.Command("certutil", "-verify", "-urlfetch")
	if err := cmd.Run(); err != nil {
		// This is a rough check, Windows cert verification is complex
		return false
	}
	return true
}
