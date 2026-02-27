package cmd

import (
	"pumu/internal/scanner"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List heavy dependency folders (dry-run)",
	Long:  `Scans for heavy dependency folders and lists them without deleting anything.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return scanner.SweepDir(".", true, false, true)
	},
}
