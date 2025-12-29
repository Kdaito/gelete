package ui

import (
	"strings"

	"github.com/Kdaito/gelete/internal/git"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model state
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.State {
		case StateSelection:
			return m.handleSelectionInput(msg)
		case StateConfirmation:
			return m.handleConfirmationInput(msg)
		case StateForceConfirmation:
			return m.handleForceConfirmationInput(msg)
		case StateDone:
			return m, tea.Quit
		}
	}

	return m, nil
}

// handleSelectionInput handles keyboard input in the selection state
func (m AppModel) handleSelectionInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "up", "k":
		if m.CursorIndex > 0 {
			m.CursorIndex--
		}

	case "down", "j":
		if m.CursorIndex < len(m.Branches)-1 {
			m.CursorIndex++
		}

	case " ", "enter":
		if len(m.Branches) > 0 {
			branch := m.Branches[m.CursorIndex]
			m.Selected[branch] = !m.Selected[branch]
		}

	case "d":
		// Check if any branches are selected
		hasSelection := false
		for _, selected := range m.Selected {
			if selected {
				hasSelection = true
				break
			}
		}

		if hasSelection {
			m.State = StateConfirmation
		}
	}

	return m, nil
}

// handleConfirmationInput handles keyboard input in the confirmation state
func (m AppModel) handleConfirmationInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		m.State = StateDeleting
		return m, m.deleteBranches

	case "n", "q", "ctrl+c":
		m.State = StateSelection
	}

	return m, nil
}

// handleForceConfirmationInput handles keyboard input in the force confirmation state
func (m AppModel) handleForceConfirmationInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		m.State = StateDeleting
		return m, m.forceDeleteBranches

	case "n", "q", "ctrl+c":
		// Skip unmerged branches and mark as done
		m.State = StateDone
	}

	return m, nil
}

// deleteBranches executes branch deletion and returns a command
// If unmerged branches are detected, transitions to StateForceConfirmation
func (m AppModel) deleteBranches() tea.Msg {
	m.DeletedCount = 0
	m.FailedBranches = make(map[string]string)
	m.UnmergedBranches = make(map[string]string)

	for _, branch := range m.Branches {
		if m.Selected[branch] {
			err := git.DeleteBranch(branch)
			if err != nil {
				// Check if error is due to unmerged changes
				if isUnmergedError(err.Error()) {
					m.UnmergedBranches[branch] = err.Error()
				} else {
					m.FailedBranches[branch] = err.Error()
				}
			} else {
				m.DeletedCount++
			}
		}
	}

	// If there are unmerged branches, prompt for force delete
	if len(m.UnmergedBranches) > 0 {
		m.State = StateForceConfirmation
	} else {
		m.State = StateDone
	}

	return m
}

// forceDeleteBranches executes force deletion of unmerged branches
func (m AppModel) forceDeleteBranches() tea.Msg {
	for branch := range m.UnmergedBranches {
		err := git.ForceDeleteBranch(branch)
		if err != nil {
			m.FailedBranches[branch] = err.Error()
		} else {
			m.DeletedCount++
			delete(m.UnmergedBranches, branch)
		}
	}

	m.State = StateDone
	return m
}

// isUnmergedError checks if an error message indicates unmerged changes
func isUnmergedError(errMsg string) bool {
	// Git typically returns errors containing "not fully merged" for unmerged branches
	return strings.Contains(errMsg, "not fully merged") ||
		strings.Contains(errMsg, "not merged")
}
