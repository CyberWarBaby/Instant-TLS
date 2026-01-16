package commands

import (
	"fmt"

	"github.com/CyberWarBaby/Instant-TLS/cli/internal/cert"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var renewCmd = &cobra.Command{
	Use:   "renew",
	Short: "Renew certificates expiring within 30 days",
	Long: `Check all certificates and renew any that are expiring within 30 days.

Example:
  instanttls renew`,
	Run: runRenew,
}

func init() {
	rootCmd.AddCommand(renewCmd)
}

func runRenew(cmd *cobra.Command, args []string) {
	pterm.Println()
	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgGreen)).
		WithTextStyle(pterm.NewStyle(pterm.FgWhite)).
		Println("🔄 Renew Certificates")
	pterm.Println()

	if !cert.CAExists() {
		printError("CA not found. Run 'instanttls init' first.")
		return
	}

	spinner, _ := pterm.DefaultSpinner.Start("Checking certificates...")

	renewed, err := cert.RenewExpiring(30)
	if err != nil {
		spinner.Fail("Renewal failed")
		printError(err.Error())
		return
	}

	if len(renewed) == 0 {
		spinner.Success("All certificates are valid")
		pterm.Println()
		pterm.Info.Println("No certificates need renewal.")
		return
	}

	spinner.Success(fmt.Sprintf("Renewed %d certificate(s)", len(renewed)))
	pterm.Println()

	pterm.Info.Println("Renewed certificates:")
	for _, domain := range renewed {
		pterm.Println("  ✅ " + domain)
	}
	pterm.Println()
}
