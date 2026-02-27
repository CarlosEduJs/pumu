package cmd

import (
	"pumu/internal/scanner"

	"github.com/spf13/cobra"
)

func init() {
	sweepCmd.Flags().Bool("reinstall", false, "Reinstall packages after removing their folders")
	sweepCmd.Flags().Bool("no-select", false, "Skip interactive selection (delete/reinstall all found folders)")
	rootCmd.AddCommand(sweepCmd)
}

var sweepCmd = &cobra.Command{
	Use:   "sweep",
	Short: "Sweep (delete) heavy dependency folders",
	Long: `Scans for heavy dependency folders and removes them.
Use --reinstall to automatically reinstall packages after deletion,
and --no-select to skip the interactive selection prompt.`,
	Example: `  pumu sweep                              # interactive selection
  pumu sweep --no-select                  # delete all without prompting
  pumu sweep --reinstall                  # delete and reinstall
  pumu sweep --no-select --reinstall ~/projects`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := cmd.Root().PersistentFlags().GetString("path")
		if err != nil {
			return err
		}
		reinstall, err := cmd.Flags().GetBool("reinstall")
		if err != nil {
			return err
		}
		noSelect, err := cmd.Flags().GetBool("no-select")
		if err != nil {
			return err
		}
		return scanner.SweepDir(path, false, reinstall, noSelect)
	},
}
