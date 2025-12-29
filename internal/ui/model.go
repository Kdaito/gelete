package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// AppState represents the current state of the application
type AppState int

const (
	// StateSelection: User is selecting branches to delete
	StateSelection AppState = iota
	// StateConfirmation: User is confirming deletion
	StateConfirmation
	// StateForceConfirmation: User is confirming force deletion of unmerged branches
	StateForceConfirmation
	// StateDeleting: Deletion is in progress
	StateDeleting
	// StateDone: Deletion complete or cancelled
	StateDone
)

// AppModel represents the application state following bubbletea's Elm architecture
type AppModel struct {
	// Branches contains all deletable branches (excludes current branch)
	Branches []string

	// Selected tracks which branches are selected for deletion (branch name -> bool)
	Selected map[string]bool

	// CursorIndex is the current cursor position in the branch list
	CursorIndex int

	// State represents the current application state
	State AppState

	// ErrorMsg holds any error message to display
	ErrorMsg string

	// SuccessMsg holds any success message to display
	SuccessMsg string

	// DeletedCount tracks how many branches were successfully deleted
	DeletedCount int

	// FailedBranches tracks branches that failed to delete with error messages
	FailedBranches map[string]string

	// UnmergedBranches tracks branches that failed due to unmerged changes
	// and are candidates for force deletion
	UnmergedBranches map[string]string
}

// Init initializes the bubbletea model
func (m AppModel) Init() tea.Cmd {
	return nil
}
