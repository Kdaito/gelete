package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// TitleStyle is used for the application title
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginBottom(1)

	// SelectedItemStyle is used for selected branches in the list
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#04B575")).
				Bold(true)

	// UnselectedItemStyle is used for unselected branches in the list
	UnselectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF"))

	// CursorStyle is used for the cursor indicator
	CursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF69B4"))

	// HelpStyle is used for help text
	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)

	// ErrorStyle is used for error messages
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	// SuccessStyle is used for success messages
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	// WarningStyle is used for warning messages
	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Bold(true)

	// ConfirmationStyle is used for confirmation prompts
	ConfirmationStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFA500")).
				Bold(true).
				MarginTop(1)
)
