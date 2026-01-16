package cmd

import (
	"github.com/instanttls/cli/internal/trust"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var trustCmd = &cobra.Command{
	Use:   "trust",
	Short: "Reinstall CA certificate in OS trust store",
	Long: `Reinstall the CA certificate in your operating system's trust store.

Use this if:
  - Trust store was reset
  - You reinstalled your OS
  - Browsers don't trust your local certificates

Example:
  instanttls trust`,
	Run: runTrust,
}

func init() {
	rootCmd.AddCommand(trustCmd)
}

func runTrust(cmd *cobra.Command, args []string) {
	pterm.Println()
	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgYellow)).
		WithTextStyle(pterm.NewStyle(pterm.FgBlack)).
		Println("🔒 Install CA Trust")
	pterm.Println()

	if err := trust.InstallCA(); err != nil {
		printError(err.Error())
		return
	}

	pterm.Println()
	pterm.Success.Println("CA certificate installed in trust store!")
	pterm.Println()
}
