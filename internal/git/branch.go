package git

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

// ListBranches returns a list of all local git branches, excluding the current branch.
// Branches are returned in alphabetical order.
func ListBranches() ([]string, error) {
	// Get current branch to exclude it
	currentBranch, err := GetCurrentBranch()
	if err != nil {
		return nil, fmt.Errorf("failed to get current branch: %w", err)
	}

	// List all branches using git branch --format
	cmd := exec.Command("git", "branch", "--format=%(refname:short)")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	// Parse output (one branch per line)
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var branches []string

	for _, line := range lines {
		branch := strings.TrimSpace(line)
		// Skip empty lines and current branch
		if branch != "" && branch != currentBranch {
			branches = append(branches, branch)
		}
	}

	// Sort alphabetically for consistent output
	sort.Strings(branches)

	return branches, nil
}

// DeleteBranch deletes the specified git branch using safe deletion (git branch -d).
// Returns an error if the branch cannot be deleted (e.g., unmerged changes, doesn't exist).
func DeleteBranch(branchName string) error {
	cmd := exec.Command("git", "branch", "-d", branchName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		outputStr := strings.TrimSpace(string(output))
		return fmt.Errorf("failed to delete branch '%s': %s", branchName, outputStr)
	}

	return nil
}

// ForceDeleteBranch forcefully deletes the specified git branch (git branch -D).
// This bypasses safety checks and will delete branches with unmerged changes.
// Use with caution. Returns an error if the branch doesn't exist.
func ForceDeleteBranch(branchName string) error {
	cmd := exec.Command("git", "branch", "-D", branchName)
	output, err := cmd.CombinedOutput()

	if err != nil {
		outputStr := strings.TrimSpace(string(output))
		return fmt.Errorf("failed to force delete branch '%s': %s", branchName, outputStr)
	}

	return nil
}
