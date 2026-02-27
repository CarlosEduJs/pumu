package cmd

import (
	"pumu/internal/scanner"

	"github.com/spf13/cobra"
)

func init() {
	// --threshold is prune-specific (local flag)
	pruneCmd.Flags().Int("threshold", 50, "Minimum staleness score to prune (0-100)")
	// --dry-run is also exposed on root as a persistent flag so other cmds can share it,
	// but prune registers its own local copy to avoid double-registration.
	pruneCmd.Flags().Bool("dry-run", false, "Only analyze and list, don't delete")
	rootCmd.AddCommand(pruneCmd)
}

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Prune dependency folders by staleness score",
	Long: `Analyzes dependency folders and removes those whose staleness score
is above the given threshold. Use --dry-run to preview without deleting.`,
	Example: `  pumu prune                          # prune with default threshold (50)
  pumu prune --threshold 70           # only prune if score >= 70
  pumu prune --dry-run                # preview without deleting
  pumu prune --dry-run -p ~/projects  # preview on a custom path`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := cmd.Root().PersistentFlags().GetString("path")
		if err != nil {
			return err
		}
		threshold, err := cmd.Flags().GetInt("threshold")
		if err != nil {
			return err
		}
		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			return err
		}
		return scanner.PruneDir(path, threshold, dryRun)
	},
}
