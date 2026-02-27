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
	Example: `  pumu list            # scan current directory
  pumu list -p ~/dev   # scan a custom path`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := cmd.Root().PersistentFlags().GetString("path")
		if err != nil {
			return err
		}
		return scanner.SweepDir(path, true, false, true)
	},
}
