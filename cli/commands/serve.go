package commands

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/CyberWarBaby/Instant-TLS/cli/internal/cert"
	"github.com/CyberWarBaby/Instant-TLS/cli/internal/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	servePort   int
	serveTarget string
)

var serveCmd = &cobra.Command{
	Use:   "serve <domain> --to <target>",
	Short: "Start HTTPS reverse proxy to your local server",
	Long: `Start an HTTPS reverse proxy that forwards requests to your local HTTP server.

This eliminates the need to configure TLS in your application - just run your
app on HTTP and let InstantTLS handle the HTTPS.

Examples:
  instanttls serve myapp.local --to localhost:3000
  instanttls serve myapp.local --to 127.0.0.1:8080 --port 443
  instanttls serve api.local --to localhost:4000 --port 8443

The command will:
  1. Generate a certificate for the domain (if not exists)
  2. Start an HTTPS server on the specified port (default: 443)
  3. Proxy all requests to your HTTP backend

Note: Port 443 requires sudo. Use --port 8443 to avoid sudo.`,
	Args: cobra.ExactArgs(1),
	Run:  runServe,
}

func init() {
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 443, "HTTPS port to listen on")
	serveCmd.Flags().StringVarP(&serveTarget, "to", "t", "", "Target HTTP server (e.g., localhost:3000)")
	serveCmd.MarkFlagRequired("to")
	rootCmd.AddCommand(serveCmd)
}

func runServe(cmd *cobra.Command, args []string) {
	domain := args[0]

	// Validate target
	if serveTarget == "" {
		printError("Target is required. Use --to localhost:3000")
		return
	}

	// Parse target to ensure it's valid
	if !strings.Contains(serveTarget, ":") {
		printError("Target must include port (e.g., localhost:3000)")
		return
	}

	// Check if port 443 and not root
	if servePort == 443 && os.Geteuid() != 0 {
		pterm.Warning.Println("Port 443 requires root access. Try:")
		pterm.Println("  sudo instanttls serve " + domain + " --to " + serveTarget)
		pterm.Println("  or use: --port 8443")
		return
	}

	pterm.Println()
	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).
		WithTextStyle(pterm.NewStyle(pterm.FgWhite)).
		Println("🔒 HTTPS Reverse Proxy")
	pterm.Println()

	// Check for existing cert or generate new one
	certDir := filepath.Join(config.GetCertsDir(), domain)
	certPath := filepath.Join(certDir, "cert.pem")
	keyPath := filepath.Join(certDir, "key.pem")

	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		pterm.Info.Println("Generating certificate for " + domain + "...")

		if !cert.CAExists() {
			printError("CA not found. Run 'instanttls setup' first.")
			return
		}

		domains := []string{domain, "localhost", "127.0.0.1"}
		if _, err := cert.GenerateCertMulti(domains); err != nil {
			printError(fmt.Sprintf("Failed to generate certificate: %v", err))
			return
		}
		pterm.Success.Println("Certificate generated")
	} else {
		pterm.Success.Println("Using existing certificate for " + domain)
	}

	// Parse target URL
	targetURL, err := url.Parse("http://" + serveTarget)
	if err != nil {
		printError(fmt.Sprintf("Invalid target URL: %v", err))
		return
	}

	// Create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Custom error handler
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		pterm.Error.Printf("Proxy error: %v\n", err)
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(fmt.Sprintf("Backend unavailable: %v", err)))
	}

	// Load TLS certificate
	tlsCert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		printError(fmt.Sprintf("Failed to load certificate: %v", err))
		return
	}

	// Create HTTPS server
	server := &http.Server{
		Addr: ":" + strconv.Itoa(servePort),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Log request
			pterm.FgGray.Printf("  → %s %s\n", r.Method, r.URL.Path)
			proxy.ServeHTTP(w, r)
		}),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
		},
	}

	// Print startup info
	pterm.Println()
	pterm.DefaultBox.WithTitle("🚀 Proxy Started").
		WithTitleTopCenter().
		WithBoxStyle(pterm.NewStyle(pterm.FgGreen)).
		Printf(`
  HTTPS: https://%s:%d
     ↓
  HTTP:  http://%s

  Press Ctrl+C to stop
`, domain, servePort, serveTarget)
	pterm.Println()

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		pterm.Println()
		pterm.Info.Println("Shutting down proxy...")
		server.Close()
	}()

	// Start server (TLS is already configured)
	pterm.Info.Printf("Listening on https://%s:%d\n", domain, servePort)
	pterm.FgGray.Println("Proxying to http://" + serveTarget)
	pterm.Println()

	if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
		printError(fmt.Sprintf("Server error: %v", err))
	}
}
