package cmd

import (
	"fmt"
	"strings"

	"github.com/instanttls/cli/internal/api"
	"github.com/instanttls/cli/internal/cert"
	"github.com/instanttls/cli/internal/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var certCmd = &cobra.Command{
	Use:   "cert <domain>",
	Short: "Generate a certificate for a domain",
	Long: `Generate a TLS certificate for a domain or wildcard pattern.

The certificate will be signed by your local CA. Make sure you have run
'instanttls init' first.

Examples:
  instanttls cert "*.local.test"     # Wildcard certificate
  instanttls cert "myapp.local"      # Single domain
  instanttls cert "localhost"        # Localhost certificate`,
	Args: cobra.ExactArgs(1),
	Run:  runCert,
}

func init() {
	rootCmd.AddCommand(certCmd)
}

func runCert(cmd *cobra.Command, args []string) {
	domain := args[0]

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
	isWildcard := strings.HasPrefix(domain, "*.")

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

	// Generate certificate
	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Generating certificate for %s...", domain))

	certDir, err := cert.GenerateCert(domain)
	if err != nil {
		spinner.Fail("Failed to generate certificate")
		printError(err.Error())
		return
	}

	spinner.Success(fmt.Sprintf("Certificate generated for %s", domain))
	pterm.Println()

	pterm.DefaultBox.WithTitle("📁 Certificate Files").
		WithTitleTopCenter().
		Println(fmt.Sprintf(`
  Certificate: %s/cert.pem
  Private Key: %s/key.pem
`, certDir, certDir))

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
}`, domain, certDir, certDir))

	pterm.Println()

	// Caddy example
	pterm.FgCyan.Println("Caddy:")
	pterm.DefaultBox.Println(fmt.Sprintf(`%s {
    tls %s/cert.pem %s/key.pem

    respond "Hello HTTPS!"
}`, strings.TrimPrefix(domain, "*."), certDir, certDir))

	pterm.Println()
}
