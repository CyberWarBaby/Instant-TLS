package cmd

import (
	"os"
	"runtime"

	"github.com/CyberWarBaby/Instant-TLS/cli/internal/api"
	"github.com/CyberWarBaby/Instant-TLS/cli/internal/cert"
	"github.com/CyberWarBaby/Instant-TLS/cli/internal/config"
	"github.com/CyberWarBaby/Instant-TLS/cli/internal/trust"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize local CA and install in trust store",
	Long: `Generate a new local Certificate Authority and install it in your OS trust store.

This command will:
  1. Generate a new CA certificate and private key
  2. Install the CA in your system's trust store
  3. Send a machine ping to the API

After running this, browsers will trust certificates signed by your local CA.

Example:
  instanttls init`,
	Run: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil || cfg == nil || cfg.Token == "" {
		printError("Not logged in. Run 'instanttls login' first.")
		return
	}

	pterm.Println()
	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgMagenta)).
		WithTextStyle(pterm.NewStyle(pterm.FgWhite)).
		Println("🔧 InstantTLS Init")
	pterm.Println()

	// Step 1: Check if CA already exists
	if cert.CAExists() {
		result, _ := pterm.DefaultInteractiveConfirm.
			WithDefaultValue(false).
			Show("CA already exists. Regenerate?")

		if !result {
			pterm.Info.Println("Using existing CA")
			pterm.Println()

			// Offer to reinstall trust
			reinstall, _ := pterm.DefaultInteractiveConfirm.
				WithDefaultValue(true).
				Show("Would you like to reinstall CA in trust store?")

			if reinstall {
				if err := trust.InstallCA(); err != nil {
					printError(err.Error())
					return
				}
			}

			printSuccessBox()
			return
		}
	}

	// Step 2: Generate CA
	spinner, _ := pterm.DefaultSpinner.Start("Generating local CA...")

	if err := cert.GenerateCA(); err != nil {
		spinner.Fail("Failed to generate CA")
		printError(err.Error())
		return
	}

	spinner.Success("CA generated successfully!")
	pterm.Println()

	// Step 3: Install in trust store
	if err := trust.InstallCA(); err != nil {
		printError(err.Error())
		printWarning("You can try again later with 'instanttls trust'")
		return
	}

	pterm.Println()
	pterm.Success.Println("CA installed in trust store!")
	pterm.Println()

	// Step 4: Ping machine
	hostname, _ := os.Hostname()
	client := api.NewClient(cfg.APIBaseURL, cfg.Token)
	_ = client.MachinePing(hostname, runtime.GOOS, runtime.GOARCH)

	printSuccessBox()
}

func printSuccessBox() {
	caDir := config.GetCADir()

	pterm.DefaultBox.WithTitle("🎉 Success: Green Lock Enabled!").
		WithTitleTopCenter().
		WithBoxStyle(pterm.NewStyle(pterm.FgGreen)).
		Println(`
Your local CA has been created and trusted.
Browsers will now trust certificates signed by this CA.
`)

	pterm.Println()
	pterm.Info.Println("Files created:")
	pterm.Println("  CA Certificate: " + caDir + "/ca.crt")
	pterm.Println("  CA Private Key: " + caDir + "/ca.key")
	pterm.Println()
	pterm.Info.Println("Next steps:")
	pterm.Println("  1. Generate a certificate: instanttls cert \"*.local.test\"")
	pterm.Println("  2. Check your setup:       instanttls doctor")
	pterm.Println()
}
