package scanner

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"pumu/internal/pkg"
	"pumu/internal/ui"

	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

// PruneDir scans for dependency folders and intelligently prunes based on safety score.
func PruneDir(root string, threshold int, dryRun bool) error {
	if dryRun {
		fmt.Println(ui.InfoStyle.Render(fmt.Sprintf("🌿 Analyzing safely deletable folders in '%s' (dry-run)...", root)))
	} else {
		fmt.Println(ui.InfoStyle.Render(fmt.Sprintf("🌿 Pruning safely deletable folders in '%s'...", root)))
	}

	targets, err := findTargetFolders(root)
	if err != nil {
		return fmt.Errorf("failed to scan: %w", err)
	}

	if len(targets) == 0 {
		fmt.Println(ui.SuccessStyle.Render("✨ No heavy folders found!"))
		return nil
	}

	folders := calculateFolderSizes(targets)

	// Analyze each folder
	var results []pkg.PruneResult
	err = spinner.New().
		Title(ui.WarningStyle.Render(fmt.Sprintf("🧐 Analyzing %d folders...", len(folders)))).
		Action(func() {
			results = analyzeAllFolders(folders)
		}).
		Run()
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	// Sort by score descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	fmt.Println()
	fmt.Println(ui.RenderRow([]string{"Folder Path", "Size", "Score", "Reason"}, []int{55, 12, 8, 30}))
	fmt.Println(lipgloss.NewStyle().Foreground(ui.ColorSubtext).Render(strings.Repeat("-", 110)))

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
		fmt.Println(ui.SuccessStyle.Render(fmt.Sprintf("✨ No folders meet the prune threshold (score ≥ %d).", threshold)))
		fmt.Println(ui.InfoStyle.Render(fmt.Sprintf("🤓 Total found: %s across %d folders", formatSize(totalSize), len(results))))
		return nil
	}

	if dryRun {
		fmt.Println(ui.SuccessStyle.Render(fmt.Sprintf("🌿 Analysis complete! %d/%d folders can be pruned (score ≥ %d).",
			prunableCount, len(results), threshold)))
		fmt.Println(ui.InfoStyle.Render(fmt.Sprintf("🤓 Space that can be freed: %s (of %s total found)",
			formatSize(prunableSize), formatSize(totalSize))))
		return nil
	}

	// Actually delete prunable folders
	tickChan := make(chan struct{})
	var deletedWg sync.WaitGroup
	var totalDeleted int64
	sem := make(chan struct{}, 20)

	go func() {
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
				tickChan <- struct{}{}
			}(r.Path, r.Size)
		}
		deletedWg.Wait()
		close(tickChan)
	}()

	fmt.Println()
	if err := ui.TrackProgress("🗑️  Deleting concurrent folders", prunableCount, tickChan); err != nil {
		return fmt.Errorf("cleanup failed: %w", err)
	}

	fmt.Println(ui.SuccessStyle.Render(fmt.Sprintf("\n🌿 Prune complete! Removed %d folders (score ≥ %d).", prunableCount, threshold)))
	fmt.Println(ui.InfoStyle.Render(fmt.Sprintf("💾 Space freed: %s (of %s total found)",
		formatSize(totalDeleted), formatSize(totalSize))))

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
		scoreStr = ui.AlertStyle.Render(fmt.Sprintf("%5d", r.Score))
	} else if r.Score >= 50 {
		scoreStr = ui.WarningStyle.Render(fmt.Sprintf("%5d", r.Score))
	} else {
		scoreStr = lipgloss.NewStyle().Foreground(ui.ColorSubtext).Render(fmt.Sprintf("%5d", r.Score))
	}

	// Dim the row if below threshold
	if r.Score < threshold {
		fmt.Println(ui.SubtextStyle.Render(ui.RenderRow([]string{displayPath, sizeStr, scoreStr, r.Reason}, []int{55, 12, 8, 30})))
	} else {
		fmt.Println(ui.RenderRow([]string{displayPath, sizeStr, scoreStr, r.Reason}, []int{55, 12, 8, 30}))
	}
}
