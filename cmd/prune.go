package cmd

import (
	"pumu/internal/scanner"

	"github.com/spf13/cobra"
)

func init() {
	pruneCmd.Flags().Int("threshold", 50, "Minimum score to prune (0-100)")
	pruneCmd.Flags().Bool("dry-run", false, "Only analyze and list, don't delete")
	rootCmd.AddCommand(pruneCmd)
}

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Prune dependency folders by staleness score",
	Long: `Analyzes dependency folders and removes those whose staleness score
is above the given threshold. Use --dry-run to preview without deleting.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		threshold, err := cmd.Flags().GetInt("threshold")
		if err != nil {
			return err
		}
		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			return err
		}
		return scanner.PruneDir(".", threshold, dryRun)
	},
}
