package cmd

import (
	"fmt"
	"runtime"

	"github.com/instanttls/cli/internal/api"
	"github.com/instanttls/cli/internal/cert"
	"github.com/instanttls/cli/internal/config"
	"github.com/instanttls/cli/internal/trust"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose InstantTLS setup",
	Long: `Run diagnostics to verify your InstantTLS setup is working correctly.

This checks:
  - Login status
  - CA certificate existence
  - Trust store installation
  - Generated certificates

Example:
  instanttls doctor`,
	Run: runDoctor,
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

func runDoctor(cmd *cobra.Command, args []string) {
	pterm.Println()
	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).
		WithTextStyle(pterm.NewStyle(pterm.FgBlack)).
		Println("🩺 InstantTLS Doctor")
	pterm.Println()

	var issues []string

	// Check 1: Login status
	cfg, err := config.Load()
	if err != nil || cfg == nil || cfg.Token == "" {
		pterm.Error.Println("Logged in: ❌")
		issues = append(issues, "Not logged in. Run 'instanttls login'")
	} else {
		pterm.Success.Println(fmt.Sprintf("Logged in: ✅ (%s - %s)", cfg.Email, cfg.Plan))

		// Validate token with API
		client := api.NewClient(cfg.APIBaseURL, cfg.Token)
		if _, err := client.Me(); err != nil {
			pterm.Warning.Println("Token validation: ⚠️ Could not validate token")
		}
	}

	// Check 2: CA exists
	if cert.CAExists() {
		pterm.Success.Println("CA certificate: ✅")
	} else {
		pterm.Error.Println("CA certificate: ❌")
		issues = append(issues, "CA not found. Run 'instanttls init'")
	}

	// Check 3: Trust store
	if trust.IsTrusted() {
		pterm.Success.Println("Trust store: ✅")
	} else {
		pterm.Warning.Println("Trust store: ⚠️ (may not be installed)")
		issues = append(issues, "Trust may not be installed. Run 'instanttls trust'")
	}

	// Check 4: Certificates
	certs, err := cert.ListCerts()
	if err != nil {
		pterm.Error.Println("Certificates: ❌ (could not list)")
	} else if len(certs) == 0 {
		pterm.Info.Println("Certificates: 0 generated")
	} else {
		pterm.Success.Println(fmt.Sprintf("Certificates: ✅ (%d found)", len(certs)))

		pterm.Println()
		pterm.DefaultSection.Println("Generated Certificates")

		tableData := pterm.TableData{
			{"Domain", "Valid Until", "Path"},
		}

		for _, c := range certs {
			tableData = append(tableData, []string{
				c.Domain,
				c.NotAfter.Format("2006-01-02"),
				c.Path,
			})
		}

		pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	}

	// Check 5: Firefox warning
	pterm.Println()
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		pterm.DefaultBox.WithTitle("⚠️ Firefox Note").
			WithBoxStyle(pterm.NewStyle(pterm.FgYellow)).
			Println(`
Firefox uses its own certificate store and may not trust
your system CA by default.

To fix this in Firefox:
1. Open about:config
2. Set security.enterprise_roots.enabled to true
3. Restart Firefox

Or manually import the CA certificate:
1. Open Firefox Settings → Privacy & Security
2. Click "View Certificates" → "Authorities"
3. Import: ` + config.GetCADir() + `/ca.crt
`)
	}

	// Summary
	pterm.Println()
	if len(issues) == 0 {
		pterm.DefaultBox.WithTitle("✅ All Good!").
			WithTitleTopCenter().
			WithBoxStyle(pterm.NewStyle(pterm.FgGreen)).
			Println(`
Your InstantTLS setup is working correctly.
Generate certificates with: instanttls cert "*.local.test"
`)
	} else {
		pterm.DefaultBox.WithTitle("⚠️ Issues Found").
			WithTitleTopCenter().
			WithBoxStyle(pterm.NewStyle(pterm.FgYellow)).
			Println(fmt.Sprintf(`
Found %d issue(s) that need attention.
`, len(issues)))

		pterm.Info.Println("To fix:")
		for i, issue := range issues {
			pterm.Println(fmt.Sprintf("  %d. %s", i+1, issue))
		}
		pterm.Println()
	}
}
