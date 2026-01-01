package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// Worktree represents a git worktree with its metadata
type Worktree struct {
	// Path is the absolute path to the worktree directory
	Path string

	// Branch is the branch name checked out in this worktree
	Branch string

	// Locked indicates if the worktree is locked
	Locked bool
}

// ListWorktrees returns all git worktrees in the current repository.
// Uses `git worktree list --porcelain` for machine-readable output.
func ListWorktrees() ([]Worktree, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil, fmt.Errorf("failed to list worktrees: %s", strings.TrimSpace(string(output)))
	}

	return parseWorktrees(string(output)), nil
}

// parseWorktrees parses the porcelain format output from `git worktree list --porcelain`
// Format:
//
//	worktree /path/to/worktree
//	HEAD <commit-hash>
//	branch refs/heads/branch-name
//	<blank line>
func parseWorktrees(output string) []Worktree {
	var worktrees []Worktree
	lines := strings.Split(output, "\n")

	var currentWorktree *Worktree

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			// Empty line marks end of a worktree entry
			if currentWorktree != nil {
				worktrees = append(worktrees, *currentWorktree)
				currentWorktree = nil
			}
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		switch key {
		case "worktree":
			// Resolve symlinks to get canonical path
			canonicalPath, err := filepath.EvalSymlinks(value)
			if err != nil {
				canonicalPath = value // Fallback to original if resolution fails
			}
			currentWorktree = &Worktree{
				Path:   canonicalPath,
				Locked: false,
			}
		case "branch":
			if currentWorktree != nil {
				// branch format: refs/heads/branch-name
				branch := strings.TrimPrefix(value, "refs/heads/")
				currentWorktree.Branch = branch
			}
		case "locked":
			if currentWorktree != nil {
				currentWorktree.Locked = true
			}
		}
	}

	// Handle last entry if file doesn't end with blank line
	if currentWorktree != nil {
		worktrees = append(worktrees, *currentWorktree)
	}

	return worktrees
}

// RemoveWorktree removes the specified worktree using `git worktree remove`.
// Returns an error if the worktree is locked or doesn't exist.
func RemoveWorktree(worktreePath string) error {
	cmd := exec.Command("git", "worktree", "remove", worktreePath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		outputStr := strings.TrimSpace(string(output))
		return fmt.Errorf("failed to remove worktree '%s': %s", worktreePath, outputStr)
	}

	return nil
}

// ForceRemoveWorktree forcefully removes the specified worktree using `git worktree remove --force --force`.
// This bypasses safety checks and will remove locked worktrees.
// Note: Double --force is required to remove locked worktrees.
func ForceRemoveWorktree(worktreePath string) error {
	cmd := exec.Command("git", "worktree", "remove", "--force", "--force", worktreePath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		outputStr := strings.TrimSpace(string(output))
		return fmt.Errorf("failed to force remove worktree '%s': %s", worktreePath, outputStr)
	}

	return nil
}

// GetWorktreeForBranch returns the worktree associated with a branch, if any.
// Returns nil if the branch is not checked out in any worktree.
func GetWorktreeForBranch(branchName string) (*Worktree, error) {
	worktrees, err := ListWorktrees()
	if err != nil {
		return nil, err
	}

	for _, wt := range worktrees {
		if wt.Branch == branchName {
			return &wt, nil
		}
	}

	return nil, nil
}
