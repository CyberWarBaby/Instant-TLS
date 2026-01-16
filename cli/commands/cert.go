package commands

import (
	"fmt"
	"strings"

	"github.com/CyberWarBaby/Instant-TLS/cli/internal/api"
	"github.com/CyberWarBaby/Instant-TLS/cli/internal/cert"
	"github.com/CyberWarBaby/Instant-TLS/cli/internal/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var certCmd = &cobra.Command{
	Use:   "cert <domain> [additional-domains...]",
	Short: "Generate a certificate for one or more domains",
	Long: `Generate a TLS certificate for one or more domains.

The certificate will be signed by your local CA. Make sure you have run
'instanttls init' first.

localhost and 127.0.0.1 are automatically included as Subject Alternative Names (SANs).

Examples:
  instanttls cert myapp.local                    # Single domain + localhost
  instanttls cert myapp.local api.local          # Multiple domains
  instanttls cert "*.local.dev"                  # Wildcard certificate
  instanttls cert myapp.local localhost:3000     # With port (port ignored in cert)`,
	Args: cobra.MinimumNArgs(1),
	Run:  runCert,
}

func init() {
	rootCmd.AddCommand(certCmd)
}

func runCert(cmd *cobra.Command, args []string) {
	domains := args
	primaryDomain := domains[0]

	cfg, err := config.Load()
	if err != nil || cfg == nil || cfg.Token == "" {
		printError("Not logged in. Run 'instanttls login' first.")
		return
	}

	pterm.Println()
	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgBlue)).
		WithTextStyle(pterm.NewStyle(pterm.FgWhite)).
		Println("📜 Generate Certificate")
	pterm.Println()

	// Check CA exists
	if !cert.CAExists() {
		printError("CA not found. Run 'instanttls init' first.")
		return
	}

	// Check plan limits
	isWildcard := strings.HasPrefix(primaryDomain, "*.")

	if isWildcard && cfg.Plan == "free" {
		// Check license from API
		client := api.NewClient(cfg.APIBaseURL, cfg.Token)
		license, err := client.License()
		if err != nil {
			printWarning("Could not verify license, proceeding with local check")
		} else {
			maxCerts := license.Limits["max_wildcard_certs"]
			currentCount := cert.CountWildcardCerts()

			if maxCerts > 0 && currentCount >= maxCerts {
				printError(fmt.Sprintf(
					"Free plan limit reached (%d/%d wildcard certs). Upgrade to Pro for unlimited certs.",
					currentCount, maxCerts,
				))
				pterm.Println()
				pterm.Info.Println("Upgrade at: https://instanttls.dev/pricing")
				return
			}
		}
	}

	// Generate certificate with multiple domains
	domainsDisplay := strings.Join(domains, ", ")
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Generating certificate for %s...", domainsDisplay))

	certDir, err := cert.GenerateCertMulti(domains)
	if err != nil {
		spinner.Fail("Failed to generate certificate")
		printError(err.Error())
		return
	}

	spinner.Success(fmt.Sprintf("Certificate generated for %s", domainsDisplay))
	pterm.Println()

	pterm.DefaultBox.WithTitle("📁 Certificate Files").
		WithTitleTopCenter().
		Println(fmt.Sprintf(`
  Certificate: %s/cert.pem
  Private Key: %s/key.pem
  
  Domains: %s + localhost + 127.0.0.1
`, certDir, certDir, domainsDisplay))

	pterm.Println()
	pterm.DefaultSection.Println("Usage Examples")
	pterm.Println()

	// Node.js example
	pterm.FgCyan.Println("Node.js:")
	pterm.DefaultBox.Println(fmt.Sprintf(`const https = require('https');
const fs = require('fs');

const options = {
  key: fs.readFileSync('%s/key.pem'),
  cert: fs.readFileSync('%s/cert.pem')
};

https.createServer(options, (req, res) => {
  res.writeHead(200);
  res.end('Hello HTTPS!');
}).listen(443);`, certDir, certDir))

	pterm.Println()

	// Nginx example
	pterm.FgCyan.Println("Nginx:")
	pterm.DefaultBox.Println(fmt.Sprintf(`server {
    listen 443 ssl;
    server_name %s;

    ssl_certificate     %s/cert.pem;
    ssl_certificate_key %s/key.pem;

    location / {
        # your config here
    }
}`, primaryDomain, certDir, certDir))

	pterm.Println()

	// Caddy example
	pterm.FgCyan.Println("Caddy:")
	pterm.DefaultBox.Println(fmt.Sprintf(`%s {
    tls %s/cert.pem %s/key.pem

    respond "Hello HTTPS!"
}`, strings.TrimPrefix(primaryDomain, "*."), certDir, certDir))

	pterm.Println()
}
