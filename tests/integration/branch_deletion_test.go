package integration

import (
	"os"
	"os/exec"
	"testing"

	"github.com/Kdaito/gelete/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRepo creates a temporary git repository for integration testing.
func setupTestRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	exec.Command("git", "init", dir).Run()
	exec.Command("git", "-C", dir, "config", "user.name", "Test User").Run()
	exec.Command("git", "-C", dir, "config", "user.email", "test@example.com").Run()
	exec.Command("git", "-C", dir, "commit", "--allow-empty", "-m", "Initial commit").Run()

	return dir
}

// TestBranchDeletion_BasicScenario tests basic branch deletion workflow.
func TestBranchDeletion_BasicScenario(t *testing.T) {
	repo := setupTestRepo(t)

	// Change to repository directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// Create test branches
	exec.Command("git", "branch", "feature-a").Run()
	exec.Command("git", "branch", "feature-b").Run()
	exec.Command("git", "branch", "bugfix-1").Run()

	// List branches (should exclude current branch)
	branches, err := git.ListBranches()
	require.NoError(t, err, "ListBranches should succeed")
	assert.Len(t, branches, 3, "Should have 3 deletable branches")
	assert.Contains(t, branches, "feature-a")
	assert.Contains(t, branches, "feature-b")
	assert.Contains(t, branches, "bugfix-1")

	// Delete feature-a
	err = git.DeleteBranch("feature-a")
	assert.NoError(t, err, "DeleteBranch should succeed for merged branch")

	// Verify branch is deleted
	branches, err = git.ListBranches()
	require.NoError(t, err)
	assert.Len(t, branches, 2, "Should have 2 branches remaining")
	assert.NotContains(t, branches, "feature-a", "feature-a should be deleted")

	// Delete feature-b
	err = git.DeleteBranch("feature-b")
	assert.NoError(t, err, "DeleteBranch should succeed")

	// Verify only bugfix-1 remains
	branches, err = git.ListBranches()
	require.NoError(t, err)
	assert.Len(t, branches, 1)
	assert.Contains(t, branches, "bugfix-1")
}

// TestBranchDeletion_MultipleBranches tests deleting multiple branches.
func TestBranchDeletion_MultipleBranches(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// Create 5 test branches
	branchNames := []string{"test-1", "test-2", "test-3", "test-4", "test-5"}
	for _, name := range branchNames {
		exec.Command("git", "branch", name).Run()
	}

	// Verify all branches exist
	branches, err := git.ListBranches()
	require.NoError(t, err)
	assert.Len(t, branches, 5)

	// Delete all test branches
	for _, name := range branchNames {
		err = git.DeleteBranch(name)
		assert.NoError(t, err, "Should delete %s", name)
	}

	// Verify all deleted
	branches, err = git.ListBranches()
	require.NoError(t, err)
	assert.Len(t, branches, 0, "All branches should be deleted")
}

// TestBranchDeletion_CurrentBranchExcluded tests that current branch is excluded from list.
func TestBranchDeletion_CurrentBranchExcluded(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// Get current branch
	currentBranch, err := git.GetCurrentBranch()
	require.NoError(t, err)

	// Create other branches
	exec.Command("git", "branch", "other-1").Run()
	exec.Command("git", "branch", "other-2").Run()

	// List branches
	branches, err := git.ListBranches()
	require.NoError(t, err)

	// Current branch should NOT be in the list
	assert.NotContains(t, branches, currentBranch, "Current branch should be excluded from deletable list")
	assert.Len(t, branches, 2, "Should only show deletable branches")
}

// TestBranchDeletion_EmptyRepository tests behavior when no deletable branches exist.
func TestBranchDeletion_EmptyRepository(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// Only current branch exists, no other branches
	branches, err := git.ListBranches()
	require.NoError(t, err)
	assert.Len(t, branches, 0, "Should have no deletable branches")
}
