package cmd

import (
	"pumu/internal/scanner"

	"github.com/spf13/cobra"
)

func init() {
	repairCmd.Flags().Bool("verbose", false, "Show details for all projects, including healthy ones")
	rootCmd.AddCommand(repairCmd)
}

var repairCmd = &cobra.Command{
	Use:   "repair",
	Short: "Repair dependency folders",
	Long:  `Scans for projects with missing or corrupted dependency folders and reinstalls them.`,
	Example: `  pumu repair                   # repair current directory
  pumu repair --verbose         # show details for healthy projects too
  pumu repair -p ~/projects     # repair a custom path`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := cmd.Root().PersistentFlags().GetString("path")
		if err != nil {
			return err
		}
		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			return err
		}
		return scanner.RepairDir(path, verbose)
	},
}
