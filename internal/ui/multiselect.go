package ui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
)

// Item represents a selectable item in the multi-select list.
type Item struct {
	Label    string
	Detail   string // e.g. formatted size
	Selected bool
}

// Result holds the outcome of the multi-select interaction.
type Result struct {
	Items    []Item
	Canceled bool
}

// RunMultiSelect launches an interactive multi-select prompt and returns the result.
// All items are pre-selected by default.
func RunMultiSelect(title string, items []Item) (Result, error) {
	if len(items) == 0 {
		return Result{}, errors.New("no items provided")
	}

	options := make([]huh.Option[int], len(items))
	for i, item := range items {
		labelLen := len(item.Label)
		padding := 0
		if labelLen < 55 {
			padding = 55 - labelLen
		}

		display := fmt.Sprintf("%s%s %s", item.Label, strings.Repeat(" ", padding), InfoStyle.Render(item.Detail))

		options[i] = huh.NewOption(display, i).Selected(item.Selected)
	}

	var selectedIndices []int

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[int]().
				Title(PrimaryStyle.Render(title)).
				Options(options...).
				Value(&selectedIndices),
		),
	).WithTheme(huh.ThemeCatppuccin())

	err := form.Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return Result{Canceled: true}, nil
		}
		return Result{}, fmt.Errorf("failed to run multi-select: %w", err)
	}

	// Reconstruct the result array based on selected indices
	selectedMap := make(map[int]bool)
	for _, idx := range selectedIndices {
		selectedMap[idx] = true
	}

	// Update items Selection state
	for i := range items {
		items[i].Selected = selectedMap[i]
	}

	return Result{
		Items:    items,
		Canceled: false,
	}, nil
}
