// Package scanner provides core logic for scanning, processing, and cleaning
// heavy dependency folders across multiple package managers.
package scanner //nolint:revive // internal package, not stdlib conflict

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"pumu/internal/pkg"
	"pumu/internal/ui"

	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

// TargetFolder holds the path and calculated size of a detected heavy dependency folder.
type TargetFolder struct {
	Path string
	Size int64
}

// ignoredPaths contains directories that pumu should never descend into.
var ignoredPaths = map[string]bool{
	".Trash": true, ".cache": true, ".npm": true, ".yarn": true,
	".cargo": true, ".rustup": true, "Library": true, "AppData": true,
	"Local": true, "Roaming": true, ".vscode": true, ".idea": true,
}

// deletableTargets contains known heavy dependency/build folders.
var deletableTargets = map[string]bool{
	"node_modules": true, "target": true, ".next": true,
	".svelte-kit": true, ".venv": true, "dist": true, "build": true,
}

func isIgnoredPath(name string) bool     { return ignoredPaths[name] }
func isDeletableTarget(name string) bool { return deletableTargets[name] }

func getTargetFolder(pm pkg.PackageManager) string {
	switch pm {
	case pkg.Bun, pkg.Pnpm, pkg.Yarn, pkg.Npm:
		return "node_modules"
	case pkg.Cargo:
		return "target"
	case pkg.Pip:
		return ".venv"
	}
	return "node_modules"
}

// RefreshCurrentDir detects the package manager in the current directory,
// removes the dependency folder, and reinstalls dependencies.
func RefreshCurrentDir() error {
	dir := "."
	pm := pkg.DetectManager(dir)
	if pm == pkg.Unknown {
		return fmt.Errorf("could not detect package manager in current directory")
	}

	fmt.Printf("🔍 Detected package manager: %s\n", pm)

	targetFolder := getTargetFolder(pm)
	targetPath := filepath.Join(dir, targetFolder)

	if pkg.FileExists(targetPath) {
		fmt.Printf("🗑️  Removing %s...\n", targetFolder)
		duration, err := pkg.RemoveDirectory(targetPath)
		if err != nil {
			return fmt.Errorf("failed to remove %s: %v", targetFolder, err)
		}
		fmt.Printf("✅ Removed in %v\n", duration)
	} else {
		fmt.Printf("ℹ️  No %s found, skipping deletion.\n", targetFolder)
	}

	err := pkg.InstallDependencies(dir, pm, false)
	if err != nil {
		return fmt.Errorf("failed to install dependencies: %v", err)
	}

	fmt.Println("🎉 Refresh complete!")
	return nil
}

// SweepDir scans root for heavy dependency folders and deletes them.
// Pass dryRun=true for list-only mode, reinstall=true to reinstall after deletion,
// and noSelect=true to skip interactive selection.
func SweepDir(root string, dryRun bool, reinstall bool, noSelect bool) error {
	printScanMessage(dryRun, root)

	targets, err := findTargetFolders(root)
	if err != nil {
		return err
	}

	if len(targets) == 0 {
		fmt.Println(ui.SuccessStyle.Render("✨ No heavy folders found!"))
		return nil
	}

	folders := calculateFolderSizes(targets)

	// Interactive selection for deletion
	if !dryRun && !noSelect {
		selected, err := selectFolders(folders, "🗑️  Select folders to delete:")
		if err != nil {
			return fmt.Errorf("selection failed: %w", err)
		}
		if selected == nil {
			fmt.Println(ui.WarningStyle.Render("\n⚠️  Operation canceled."))
			return nil
		}
		folders = selected
	}

	if len(folders) == 0 {
		fmt.Println(ui.SuccessStyle.Render("\n✨ No folders selected for deletion."))
		return nil
	}

	totalFreed, totalDeleted := processFolders(folders, dryRun)
	printSummary(dryRun, folders, totalFreed, totalDeleted)

	if !dryRun && reinstall {
		reinstallDependencies(folders, noSelect)
	}

	return nil
}

func printScanMessage(dryRun bool, root string) {
	if dryRun {
		fmt.Println(ui.InfoStyle.Render(fmt.Sprintf("🔎 Listing heavy dependency folders in '%s'...", root)))
	} else {
		fmt.Println(ui.InfoStyle.Render(fmt.Sprintf("🔎 Scanning for heavy dependency folders in '%s'...", root)))
	}
}

func findTargetFolders(root string) ([]string, error) {
	var targets []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			if d.Name() == ".git" || isIgnoredPath(d.Name()) {
				return filepath.SkipDir
			}

			if isDeletableTarget(d.Name()) {
				targets = append(targets, path)
				return filepath.SkipDir
			}
		}

		return nil
	})

	return targets, err
}

func calculateFolderSizes(targets []string) []TargetFolder {
	var folders []TargetFolder
	action := func() {
		var wg sync.WaitGroup
		var mu sync.Mutex

		sem := make(chan struct{}, 20)

		for _, tPath := range targets {
			wg.Add(1)
			go func(p string) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()

				size, err := dirSize(p)
				if err != nil {
					size = 0
				}

				mu.Lock()
				folders = append(folders, TargetFolder{Path: p, Size: size})
				mu.Unlock()
			}(tPath)
		}

		wg.Wait()
	}

	err := spinner.New().
		Title(ui.WarningStyle.Render(fmt.Sprintf("⏱️  Calculating sizes for %d folders...", len(targets)))).
		Action(action).
		Run()
	if err != nil {
		return nil
	}

	sort.Slice(folders, func(i, j int) bool {
		return folders[i].Size > folders[j].Size
	})

	return folders
}

func processFolders(folders []TargetFolder, dryRun bool) (int64, int64) {
	var totalFreed int64
	var deletedWg sync.WaitGroup
	var totalDeleted int64

	fmt.Println()
	fmt.Println(ui.RenderRow([]string{"Folder Path", "Size"}, []int{80, 15}))
	fmt.Println(lipgloss.NewStyle().Foreground(ui.ColorSubtext).Render(strings.Repeat("-", 100)))

	for _, folder := range folders {
		printFolderInfo(folder)
		totalFreed += folder.Size
	}

	if !dryRun {
		tickChan := make(chan struct{})
		sem := make(chan struct{}, 20)

		go func() {
			for _, folder := range folders {
				deletedWg.Add(1)
				go func(p string, s int64) {
					defer deletedWg.Done()
					sem <- struct{}{}
					defer func() { <-sem }()
					_, err := pkg.RemoveDirectory(p)
					if err == nil {
						atomic.AddInt64(&totalDeleted, s)
					}
					tickChan <- struct{}{}
				}(folder.Path, folder.Size)
			}
			deletedWg.Wait()
			close(tickChan)
		}()

		fmt.Println()
		if err := ui.TrackProgress("🗑️  Deleting concurrent folders", len(folders), tickChan); err != nil {
			// If progress bar fails, we just continue (or could return error)
			// for now let's just log it or ignore formally if it's non-critical
			// but to satisfy errcheck we handle it.
			return totalFreed, totalDeleted
		}
	}

	return totalFreed, totalDeleted
}

func printFolderInfo(folder TargetFolder) {
	sizeMB := float64(folder.Size) / 1024 / 1024
	formattedSize := formatSize(folder.Size)

	var sizeStr string
	if sizeMB > 1000 {
		sizeStr = ui.AlertStyle.Render(fmt.Sprintf("%10s 🚨", formattedSize))
	} else if sizeMB > 100 {
		sizeStr = ui.WarningStyle.Render(fmt.Sprintf("%10s ⚠️", formattedSize))
	} else {
		sizeStr = ui.SuccessStyle.Render(fmt.Sprintf("%10s", formattedSize))
	}

	displayPath := folder.Path
	if len(displayPath) > 80 {
		displayPath = "..." + displayPath[len(displayPath)-77:]
	}

	fmt.Println(ui.RenderRow([]string{displayPath, sizeStr}, []int{80, 15}))
}

func printSummary(dryRun bool, folders []TargetFolder, totalFreed, totalDeleted int64) {
	fmt.Println(lipgloss.NewStyle().Foreground(ui.ColorSubtext).Render(strings.Repeat("-", 100)))
	if dryRun {
		fmt.Println(ui.SuccessStyle.Render(fmt.Sprintf("📋 List complete! Found %d heavy folders.", len(folders))))
		fmt.Println(ui.InfoStyle.Render(fmt.Sprintf("💾 Total space that can be freed: %s", formatSize(totalFreed))))
	} else {
		fmt.Println(ui.SuccessStyle.Render(fmt.Sprintf("🧹 Sweep complete! Processed %d heavy folders.", len(folders))))
		fmt.Println(ui.InfoStyle.Render(fmt.Sprintf("💾 Total space actually freed: %s", formatSize(totalDeleted))))
	}
}

// selectFolders presents an interactive multi-select for choosing folders.
// Returns nil if the user canceled, or the filtered list of selected folders.
func selectFolders(folders []TargetFolder, title string) ([]TargetFolder, error) {
	items := make([]ui.Item, len(folders))
	for i, f := range folders {
		items[i] = ui.Item{
			Label:    f.Path,
			Detail:   formatSize(f.Size),
			Selected: true,
		}
	}

	result, err := ui.RunMultiSelect(title, items)
	if err != nil {
		return nil, err
	}
	if result.Canceled {
		return nil, nil
	}

	var selected []TargetFolder
	for i, item := range result.Items {
		if item.Selected {
			selected = append(selected, folders[i])
		}
	}
	return selected, nil
}

func reinstallDependencies(folders []TargetFolder, noSelect bool) {
	// Build unique project list with their detected package managers
	seen := make(map[string]bool)
	type reinstallTarget struct {
		Dir string
		PM  pkg.PackageManager
	}
	var targets []reinstallTarget

	for _, folder := range folders {
		baseDir := filepath.Dir(folder.Path)
		if seen[baseDir] {
			continue
		}
		seen[baseDir] = true

		pm := pkg.DetectManager(baseDir)
		if pm != pkg.Unknown {
			targets = append(targets, reinstallTarget{Dir: baseDir, PM: pm})
		}
	}

	if len(targets) == 0 {
		fmt.Println(ui.WarningStyle.Render("\n⚠️  No projects with known package managers found for reinstallation."))
		return
	}

	// Interactive selection for reinstallation
	if !noSelect {
		items := make([]ui.Item, len(targets))
		for i, t := range targets {
			items[i] = ui.Item{
				Label:    t.Dir,
				Detail:   string(t.PM),
				Selected: true,
			}
		}

		result, err := ui.RunMultiSelect("📦 Select projects to reinstall:", items)
		if err != nil {
			fmt.Println(ui.AlertStyle.Render(fmt.Sprintf("❌ Selection failed: %v", err)))
			return
		}
		if result.Canceled {
			fmt.Println(ui.WarningStyle.Render("\n⚠️  Reinstallation canceled."))
			return
		}

		// Filter to only selected targets
		var selected []reinstallTarget
		for i, item := range result.Items {
			if item.Selected {
				selected = append(selected, targets[i])
			}
		}
		targets = selected
	}

	if len(targets) == 0 {
		fmt.Println(ui.SuccessStyle.Render("\n✨ No projects selected for reinstallation."))
		return
	}

	fmt.Println(ui.WarningStyle.Render("\n⚙️  Reinstalling dependencies sequentially..."))
	for _, t := range targets {
		fmt.Printf("📦 Reinstalling for %s (%s)...\n", t.Dir, t.PM)
		err := pkg.InstallDependencies(t.Dir, t.PM, true)
		if err != nil {
			fmt.Println(ui.AlertStyle.Render(fmt.Sprintf("❌ Failed to reinstall %s: %v", t.Dir, err)))
		} else {
			fmt.Println(ui.SuccessStyle.Render(fmt.Sprintf("✅ Reinstalled %s", t.Dir)))
		}
	}
	fmt.Println(ui.SuccessStyle.Render("🎉 All target reinstallations complete!"))
}

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.WalkDir(path, func(_ string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err == nil {
				size += info.Size()
			}
		}
		return nil
	})
	return size, err
}

// formatSize converts a byte count into a human-readable string (KB, MB, GB, etc.)
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	// "KMGTPE" represents Kilo, Mega, Giga, Tera, Peta, Exa
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
