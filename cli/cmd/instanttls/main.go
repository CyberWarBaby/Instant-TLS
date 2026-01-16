package main

import (
	"os"

	"github.com/CyberWarBaby/Instant-TLS/cli/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
