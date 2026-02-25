package scanner

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"pumu/internal/pkg"

	"github.com/fatih/color"
)

// PruneDir scans for dependency folders and intelligently prunes based on safety score.
func PruneDir(root string, threshold int, dryRun bool) error {
	if dryRun {
		color.Cyan("🌿 Analyzing safely deletable folders in '%s' (dry-run)...\n", root)
	} else {
		color.Cyan("🌿 Pruning safely deletable folders in '%s'...\n", root)
	}

	targets, err := findTargetFolders(root)
	if err != nil {
		return fmt.Errorf("failed to scan: %w", err)
	}

	if len(targets) == 0 {
		color.Green("✨ No heavy folders found!\n")
		return nil
	}

	folders := calculateFolderSizes(targets)

	// Analyze each folder
	color.Yellow("🧐 Analyzing %d folders...\n", len(folders))

	results := analyzeAllFolders(folders)

	// Sort by score descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// Print table header
	fmt.Println()
	color.Set(color.FgWhite, color.Underline)
	fmt.Printf("%-55s | %10s | %5s | %s\n", "Folder Path", "Size", "Score", "Reason")
	color.Unset()

	var prunableCount int
	var prunableSize int64
	var totalSize int64

	for _, r := range results {
		totalSize += r.Size
		printPruneRow(r, threshold)
		if r.Score >= threshold {
			prunableCount++
			prunableSize += r.Size
		}
	}

	// Summary
	fmt.Println(strings.Repeat("-", 110))

	if prunableCount == 0 {
		color.Green("✨ No folders meet the prune threshold (score ≥ %d).", threshold)
		color.Cyan("🤓 Total found: %s across %d folders\n", formatSize(totalSize), len(results))
		return nil
	}

	if dryRun {
		color.Green("🌿 Analysis complete! %d/%d folders can be pruned (score ≥ %d).",
			prunableCount, len(results), threshold)
		color.Cyan("🤓 Space that can be freed: %s (of %s total found)\n",
			formatSize(prunableSize), formatSize(totalSize))
		return nil
	}

	// Actually delete prunable folders
	color.Yellow("\n🗑️  Deleting %d folders concurrently...", prunableCount)

	var deletedWg sync.WaitGroup
	var totalDeleted int64
	sem := make(chan struct{}, 20)

	for _, r := range results {
		if r.Score < threshold {
			continue
		}

		deletedWg.Add(1)
		go func(path string, size int64) {
			defer deletedWg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			_, err := pkg.RemoveDirectory(path)
			if err == nil {
				atomic.AddInt64(&totalDeleted, size)
			}
		}(r.Path, r.Size)
	}

	deletedWg.Wait()

	color.Green("\n🌿 Prune complete! Removed %d folders (score ≥ %d).", prunableCount, threshold)
	color.Cyan("💾 Space freed: %s (of %s total found)\n",
		formatSize(totalDeleted), formatSize(totalSize))

	return nil
}

// analyzeAllFolders runs AnalyzeFolder concurrently on all found folders.
func analyzeAllFolders(folders []TargetFolder) []pkg.PruneResult {
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make([]pkg.PruneResult, 0, len(folders))
	sem := make(chan struct{}, 20)

	for _, f := range folders {
		wg.Add(1)
		go func(folder TargetFolder) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			result := pkg.AnalyzeFolder(folder.Path, folder.Size)

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(f)
	}

	wg.Wait()
	return results
}

// printPruneRow prints a single row in the prune analysis table.
func printPruneRow(r pkg.PruneResult, threshold int) {
	sizeStr := formatSize(r.Size)

	displayPath := r.Path
	if len(displayPath) > 55 {
		displayPath = "..." + displayPath[len(displayPath)-52:]
	}

	// Color the score based on value
	var scoreStr string
	if r.Score >= 80 {
		scoreStr = color.RedString("%5d", r.Score)
	} else if r.Score >= 50 {
		scoreStr = color.YellowString("%5d", r.Score)
	} else {
		scoreStr = color.HiBlackString("%5d", r.Score)
	}

	// Dim the row if below threshold
	if r.Score < threshold {
		fmt.Printf("%-55s | %10s | %s | %s\n",
			color.HiBlackString(displayPath),
			color.HiBlackString(sizeStr),
			scoreStr,
			color.HiBlackString(r.Reason),
		)
	} else {
		fmt.Printf("%-55s | %10s | %s | %s\n",
			displayPath,
			sizeStr,
			scoreStr,
			r.Reason,
		)
	}
}
