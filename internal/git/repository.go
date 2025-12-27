package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// ValidateRepository checks if the current directory is a valid git repository.
// Returns an error if not in a git repository or if git is not installed.
func ValidateRepository() error {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	output, err := cmd.CombinedOutput()

	if err != nil {
		// Check if git command is not found
		if _, lookErr := exec.LookPath("git"); lookErr != nil {
			return fmt.Errorf("git command not found. Please install git and ensure it's in your PATH")
		}

		// Not a git repository
		outputStr := strings.TrimSpace(string(output))
		if strings.Contains(outputStr, "not a git repository") {
			return fmt.Errorf("not a git repository. Run gelete from within a git repository")
		}

		// Other git error
		return fmt.Errorf("git error: %s", outputStr)
	}

	return nil
}

// GetCurrentBranch returns the name of the currently checked-out branch.
// Returns "HEAD" if in detached HEAD state.
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	branch := strings.TrimSpace(string(output))

	// Handle detached HEAD state (empty output)
	if branch == "" {
		return "HEAD", nil
	}

	return branch, nil
}
