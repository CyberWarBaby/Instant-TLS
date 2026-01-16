package main

import (
	"os"

	"github.com/CyberWarBaby/Instant-TLS/cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
