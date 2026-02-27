// Package cmd defines the CLI commands for pumu.
package cmd

import (
	"fmt"
	"os"

	"pumu/internal/scanner"

	"github.com/spf13/cobra"
)

const version = "v1.2.0-beta.0"

var rootCmd = &cobra.Command{
	Use:   "pumu",
	Short: "pumu – clean heavy dependency folders from your projects",
	Long: `pumu scans your filesystem for heavy dependency folders
(node_modules, target, .venv, etc.) and lets you sweep, list,
repair or prune them with ease.

Running pumu with no subcommand refreshes the current directory.`,
	Version: version,
	// Running bare `pumu` refreshes the current dir.
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Running refresh in current directory...")
		return scanner.RefreshCurrentDir()
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
