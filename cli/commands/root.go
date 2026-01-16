package commands

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	// Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED")).
			PaddingBottom(1)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3B82F6"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B"))

	mutedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280"))
)

var rootCmd = &cobra.Command{
	Use:   "instanttls",
	Short: "InstantTLS - Trusted HTTPS locally with zero browser warnings",
	Long: `
  ██╗███╗   ██╗███████╗████████╗ █████╗ ███╗   ██╗████████╗████████╗██╗     ███████╗
  ██║████╗  ██║██╔════╝╚══██╔══╝██╔══██╗████╗  ██║╚══██╔══╝╚══██╔══╝██║     ██╔════╝
  ██║██╔██╗ ██║███████╗   ██║   ███████║██╔██╗ ██║   ██║      ██║   ██║     ███████╗
  ██║██║╚██╗██║╚════██║   ██║   ██╔══██║██║╚██╗██║   ██║      ██║   ██║     ╚════██║
  ██║██║ ╚████║███████║   ██║   ██║  ██║██║ ╚████║   ██║      ██║   ███████╗███████║
  ╚═╝╚═╝  ╚═══╝╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═══╝   ╚═╝      ╚═╝   ╚══════╝╚══════╝

InstantTLS gives developers trusted HTTPS locally with zero browser warnings.

Generate a local Certificate Authority, install it in your OS trust store,
and create wildcard certificates for all your local development domains.

Get started:
  1. Login with your Personal Access Token:  instanttls login
  2. Initialize local CA:                     instanttls init
  3. Generate a certificate:                  instanttls cert "*.local.test"
  4. Check your setup:                        instanttls doctor

Learn more at https://instanttls.dev
`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Customize help template
	rootCmd.SetHelpTemplate(getHelpTemplate())

	// Disable default completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func getHelpTemplate() string {
	return `{{.Long}}

{{with (or .Commands .HasAvailableSubCommands)}}
` + titleStyle.Render("Commands:") + `
{{range .}}
  {{rpad .Name .NamePadding }} {{.Short}}
{{end}}{{end}}{{if .HasAvailableLocalFlags}}

` + titleStyle.Render("Flags:") + `
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

` + titleStyle.Render("Global Flags:") + `
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasExample}}

` + titleStyle.Render("Examples:") + `
{{.Example}}{{end}}

Use "{{.CommandPath}} [command] --help" for more information about a command.
`
}

func printSuccess(message string) {
	pterm.Success.Println(message)
}

func printError(message string) {
	pterm.Error.Println(message)
}

func printInfo(message string) {
	pterm.Info.Println(message)
}

func printWarning(message string) {
	pterm.Warning.Println(message)
}

func printStep(step, total int, message string) {
	pterm.DefaultSection.WithLevel(2).Println(
		fmt.Sprintf("[%d/%d] %s", step, total, message),
	)
}

func exitWithError(message string) {
	printError(message)
	os.Exit(1)
}
