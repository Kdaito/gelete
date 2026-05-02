package ui

import (
	"fmt"
	"strings"
)

// View renders the UI based on the current model state
func (m AppModel) View() string {
	switch m.State {
	case StateSelection:
		return m.renderSelection()
	case StateConfirmation:
		return m.renderConfirmation()
	case StateForceConfirmation:
		return m.renderForceConfirmation()
	case StateDeleting:
		return m.renderDeleting()
	case StateDone:
		return m.renderDone()
	}
	return ""
}

func (m AppModel) renderSelection() string {
	var b strings.Builder

	b.WriteString(TitleStyle.Render("gelete - Interactive Branch Deletion"))
	b.WriteString("\n\n")

	if len(m.Branches) == 0 {
		b.WriteString(HelpStyle.Render("No branches to delete."))
		b.WriteString("\n\n")
		b.WriteString(HelpStyle.Render("Press q to quit."))
		return b.String()
	}

	for i, branch := range m.Branches {
		cursor := "  "
		if i == m.CursorIndex {
			cursor = CursorStyle.Render("> ")
		}

		checkbox := "[ ]"
		style := UnselectedItemStyle
		if m.Selected[branch] {
			checkbox = "[✓]"
			style = SelectedItemStyle
		}

		branchDisplay := branch
		if _, hasWorktree := m.BranchWorktrees[branch]; hasWorktree {
			branchDisplay = branch + " " + WarningStyle.Render("[worktree]")
		}

		fmt.Fprintf(&b, "%s%s %s\n", cursor, checkbox, style.Render(branchDisplay))
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("↑/k: up • ↓/j: down • space/enter: toggle • d: delete selected • q: quit"))
	return b.String()
}

func (m AppModel) renderConfirmation() string {
	var b strings.Builder

	b.WriteString(ConfirmationStyle.Render("Are you sure you want to delete these branches?"))
	b.WriteString("\n\n")

	selectedCount := 0
	for _, branch := range m.Branches {
		if m.Selected[branch] {
			b.WriteString(WarningStyle.Render(fmt.Sprintf("  • %s", branch)))
			b.WriteString("\n")
			selectedCount++
		}
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render(fmt.Sprintf("Total: %d branch(es)", selectedCount)))
	b.WriteString("\n\n")
	b.WriteString(HelpStyle.Render("y: confirm • n: cancel"))
	return b.String()
}

func (m AppModel) renderForceConfirmation() string {
	var b strings.Builder

	b.WriteString(ErrorStyle.Render("⚠ Warning: Unmerged Branches Detected"))
	b.WriteString("\n\n")
	b.WriteString("The following branches have unmerged changes:\n\n")

	for branch, errMsg := range m.UnmergedBranches {
		b.WriteString(WarningStyle.Render(fmt.Sprintf("  • %s", branch)))
		b.WriteString("\n")
		b.WriteString(HelpStyle.Render(fmt.Sprintf("    %s", errMsg)))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(WarningStyle.Render(fmt.Sprintf("Force delete will permanently remove %d unmerged branch(es).", len(m.UnmergedBranches))))
	b.WriteString("\n")
	b.WriteString(ErrorStyle.Render("This action cannot be undone!"))
	b.WriteString("\n\n")
	b.WriteString(HelpStyle.Render("y: force delete • n: cancel and skip these branches"))
	return b.String()
}

func (m AppModel) renderDeleting() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render("Deleting branches..."))
	b.WriteString("\n\n")
	b.WriteString("Please wait...")
	return b.String()
}

func (m AppModel) renderDone() string {
	var b strings.Builder

	b.WriteString(TitleStyle.Render("Deletion Complete"))
	b.WriteString("\n\n")

	if m.DeletedCount > 0 {
		b.WriteString(SuccessStyle.Render(fmt.Sprintf("✓ Successfully deleted %d branch(es)", m.DeletedCount)))
		b.WriteString("\n")
	}

	if len(m.FailedBranches) > 0 {
		b.WriteString("\n")
		b.WriteString(ErrorStyle.Render(fmt.Sprintf("✗ Failed to delete %d branch(es):", len(m.FailedBranches))))
		b.WriteString("\n")
		for branch, err := range m.FailedBranches {
			b.WriteString(ErrorStyle.Render(fmt.Sprintf("  • %s: %s", branch, err)))
			b.WriteString("\n")
		}
	}

	if m.ErrorMsg != "" {
		b.WriteString("\n")
		b.WriteString(ErrorStyle.Render(fmt.Sprintf("Error: %s", m.ErrorMsg)))
	}

	b.WriteString("\n\n")
	b.WriteString(HelpStyle.Render("Press any key to exit."))
	return b.String()
}
