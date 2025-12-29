package integration

import (
	"os"
	"os/exec"
	"testing"

	"github.com/Kdaito/gelete/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUnmergedBranch_SafeDeleteFails tests that safe deletion fails for unmerged branches.
// This verifies FR-008: System MUST detect unmerged branches and prevent deletion.
func TestUnmergedBranch_SafeDeleteFails(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// Get current branch name
	currentBranch, err := git.GetCurrentBranch()
	require.NoError(t, err)

	// Create a branch with unmerged changes
	exec.Command("git", "checkout", "-b", "experimental").Run()
	exec.Command("git", "commit", "--allow-empty", "-m", "Unmerged commit").Run()

	// Switch back to original branch
	exec.Command("git", "checkout", currentBranch).Run()

	// Attempt to delete the unmerged branch with safe delete
	err = git.DeleteBranch("experimental")

	// Should fail because branch has unmerged changes
	assert.Error(t, err, "DeleteBranch should fail for unmerged branch")
	assert.Contains(t, err.Error(), "experimental", "Error should mention branch name")
}

// TestUnmergedBranch_ForceDeleteSucceeds tests that force deletion succeeds for unmerged branches.
// This verifies FR-009: System MUST offer force deletion option.
func TestUnmergedBranch_ForceDeleteSucceeds(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// Get current branch name
	currentBranch, err := git.GetCurrentBranch()
	require.NoError(t, err)

	// Create a branch with unmerged changes
	exec.Command("git", "checkout", "-b", "experimental").Run()
	exec.Command("git", "commit", "--allow-empty", "-m", "Unmerged commit").Run()

	// Switch back to original branch
	exec.Command("git", "checkout", currentBranch).Run()

	// Force delete should succeed
	err = git.ForceDeleteBranch("experimental")
	assert.NoError(t, err, "ForceDeleteBranch should succeed for unmerged branch")

	// Verify branch is deleted
	branches, _ := git.ListBranches()
	assert.NotContains(t, branches, "experimental", "Branch should be deleted")
}

// TestUnmergedBranch_MultipleScenarios tests handling of mixed merged and unmerged branches.
func TestUnmergedBranch_MultipleScenarios(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// Get current branch name
	currentBranch, err := git.GetCurrentBranch()
	require.NoError(t, err)

	// Create merged branch (no extra commits)
	exec.Command("git", "branch", "merged-branch").Run()

	// Create unmerged branch with extra commit
	exec.Command("git", "checkout", "-b", "unmerged-branch").Run()
	exec.Command("git", "commit", "--allow-empty", "-m", "Unmerged commit").Run()

	// Switch back to original branch
	exec.Command("git", "checkout", currentBranch).Run()

	// Safe delete of merged branch should succeed
	err = git.DeleteBranch("merged-branch")
	assert.NoError(t, err, "DeleteBranch should succeed for merged branch")

	// Safe delete of unmerged branch should fail
	err = git.DeleteBranch("unmerged-branch")
	assert.Error(t, err, "DeleteBranch should fail for unmerged branch")

	// Force delete of unmerged branch should succeed
	err = git.ForceDeleteBranch("unmerged-branch")
	assert.NoError(t, err, "ForceDeleteBranch should succeed for unmerged branch")

	// Verify both branches are deleted
	branches, _ := git.ListBranches()
	assert.NotContains(t, branches, "merged-branch")
	assert.NotContains(t, branches, "unmerged-branch")
}

// TestUnmergedBranch_ErrorMessageFormat tests that error messages are clear and helpful.
func TestUnmergedBranch_ErrorMessageFormat(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// Get current branch name
	currentBranch, err := git.GetCurrentBranch()
	require.NoError(t, err)

	// Create a branch with unmerged changes
	exec.Command("git", "checkout", "-b", "experimental").Run()
	exec.Command("git", "commit", "--allow-empty", "-m", "Unmerged commit").Run()

	// Switch back to original branch
	exec.Command("git", "checkout", currentBranch).Run()

	// Attempt to delete
	err = git.DeleteBranch("experimental")

	// Error message should be clear and helpful
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "experimental", "Error should mention branch name")
	// Git's error message typically includes "not fully merged" or similar
}
