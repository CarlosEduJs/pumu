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
	Example: `  pumu                        # refresh current directory
  pumu list                   # list all heavy folders
  pumu sweep --no-select      # delete all without prompting`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			return err
		}
		fmt.Printf("Running refresh in %s...\n", path)
		return scanner.RefreshCurrentDir()
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("path", "p", ".", "Root path to scan")

	rootCmd.SetVersionTemplate("pumu version {{.Version}}\n")

	// Disable the default completion command output on –completion-script-*
	rootCmd.CompletionOptions.DisableDefaultCmd = false
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
