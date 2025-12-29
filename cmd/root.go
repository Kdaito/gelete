package cmd

import (
	"fmt"

	"github.com/Kdaito/gelete/internal/git"
	"github.com/Kdaito/gelete/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	// Version is set by goreleaser during build
	Version = "dev"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:     "gelete",
	Short:   "Interactive git branch deletion tool",
	Long:    `gelete provides an interactive terminal UI for selecting and deleting local git branches.`,
	Version: Version,
	RunE:    run,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

// run is the main command execution function
func run(cmd *cobra.Command, args []string) error {
	// Validate we're in a git repository
	if err := git.ValidateRepository(); err != nil {
		return fmt.Errorf("not a git repository: %w", err)
	}

	// Get list of deletable branches
	branches, err := git.ListBranches()
	if err != nil {
		return fmt.Errorf("failed to list branches: %w", err)
	}

	// Check if there are any branches to delete
	if len(branches) == 0 {
		fmt.Println("No branches to delete.")
		fmt.Println("(Current branch is excluded from the list)")
		return nil
	}

	// Initialize the UI model
	model := ui.AppModel{
		Branches:       branches,
		Selected:       make(map[string]bool),
		CursorIndex:    0,
		State:          ui.StateSelection,
		FailedBranches: make(map[string]string),
	}

	// Start the bubbletea program
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running UI: %w", err)
	}

	return nil
}

func init() {
	rootCmd.Flags().BoolP("version", "v", false, "Print version information")
}
