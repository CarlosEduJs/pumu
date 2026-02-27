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
	RunE: func(cmd *cobra.Command, args []string) error {
		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			return err
		}
		return scanner.RepairDir(".", verbose)
	},
}
