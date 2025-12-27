package unit

import (
	"os"
	"os/exec"
	"testing"

	"github.com/Kdaito/gelete/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRepo creates a temporary git repository for testing.
// Returns the absolute path to the repository directory.
// The repository is automatically cleaned up when the test finishes.
func setupTestRepo(t *testing.T) string {
	t.Helper()

	// Create temporary directory
	dir := t.TempDir()

	// Initialize git repository
	cmd := exec.Command("git", "init", dir)
	err := cmd.Run()
	require.NoError(t, err, "Failed to initialize git repository")

	// Configure git user (required for commits)
	exec.Command("git", "-C", dir, "config", "user.name", "Test User").Run()
	exec.Command("git", "-C", dir, "config", "user.email", "test@example.com").Run()

	// Create initial commit (required for branches to exist)
	exec.Command("git", "-C", dir, "commit", "--allow-empty", "-m", "Initial commit").Run()

	return dir
}

// TestValidateRepository_ValidRepo tests that ValidateRepository succeeds in a valid git repository.
func TestValidateRepository_ValidRepo(t *testing.T) {
	repo := setupTestRepo(t)

	// Change to repository directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(repo)
	require.NoError(t, err)

	// Test validation
	err = git.ValidateRepository()
	assert.NoError(t, err, "ValidateRepository should succeed in a valid git repository")
}

// TestValidateRepository_NotARepo tests that ValidateRepository fails when not in a git repository.
func TestValidateRepository_NotARepo(t *testing.T) {
	// Create a non-git directory
	dir := t.TempDir()

	// Change to non-git directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(dir)
	require.NoError(t, err)

	// Test validation should fail
	err = git.ValidateRepository()
	assert.Error(t, err, "ValidateRepository should fail in a non-git directory")
	assert.Contains(t, err.Error(), "not a git repository", "Error message should mention 'not a git repository'")
}

// TestGetCurrentBranch_MainBranch tests getting the current branch name.
func TestGetCurrentBranch_MainBranch(t *testing.T) {
	repo := setupTestRepo(t)

	// Change to repository directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(repo)
	require.NoError(t, err)

	// Get current branch (should be "main" or "master" depending on git config)
	branch, err := git.GetCurrentBranch()
	assert.NoError(t, err, "GetCurrentBranch should succeed")
	assert.NotEmpty(t, branch, "Branch name should not be empty")
	// Common default branches
	assert.Contains(t, []string{"main", "master"}, branch, "Should be on default branch")
}

// TestGetCurrentBranch_CustomBranch tests getting current branch after switching.
func TestGetCurrentBranch_CustomBranch(t *testing.T) {
	repo := setupTestRepo(t)

	// Change to repository directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(repo)
	require.NoError(t, err)

	// Create and switch to a new branch
	exec.Command("git", "checkout", "-b", "test-branch").Run()

	// Get current branch
	branch, err := git.GetCurrentBranch()
	assert.NoError(t, err, "GetCurrentBranch should succeed")
	assert.Equal(t, "test-branch", branch, "Should be on test-branch")
}

// TestGetCurrentBranch_DetachedHEAD tests getting current branch in detached HEAD state.
func TestGetCurrentBranch_DetachedHEAD(t *testing.T) {
	repo := setupTestRepo(t)

	// Change to repository directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err := os.Chdir(repo)
	require.NoError(t, err)

	// Create a commit to checkout
	exec.Command("git", "commit", "--allow-empty", "-m", "Second commit").Run()

	// Get the commit hash
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, _ := cmd.Output()
	commitHash := string(output)[:7] // First 7 chars

	// Checkout the commit (detached HEAD)
	exec.Command("git", "checkout", commitHash).Run()

	// Get current branch
	branch, err := git.GetCurrentBranch()
	assert.NoError(t, err, "GetCurrentBranch should succeed even in detached HEAD")
	assert.Equal(t, "HEAD", branch, "Should return 'HEAD' in detached state")
}

// TestListBranches_MultipleBranches tests listing branches with multiple branches present.
func TestListBranches_MultipleBranches(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// Create test branches
	exec.Command("git", "branch", "feature-a").Run()
	exec.Command("git", "branch", "feature-b").Run()
	exec.Command("git", "branch", "test-1").Run()

	// List branches
	branches, err := git.ListBranches()
	assert.NoError(t, err, "ListBranches should succeed")
	assert.Len(t, branches, 3, "Should return 3 branches (excluding current)")
	assert.Contains(t, branches, "feature-a")
	assert.Contains(t, branches, "feature-b")
	assert.Contains(t, branches, "test-1")
}

// TestListBranches_ExcludesCurrentBranch tests that current branch is excluded.
func TestListBranches_ExcludesCurrentBranch(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// Get current branch
	currentBranch, err := git.GetCurrentBranch()
	require.NoError(t, err)

	// Create other branches
	exec.Command("git", "branch", "other").Run()

	// List branches
	branches, err := git.ListBranches()
	assert.NoError(t, err)
	assert.NotContains(t, branches, currentBranch, "Current branch should be excluded")
	assert.Contains(t, branches, "other")
}

// TestListBranches_OnlyCurrentBranch tests when only current branch exists.
func TestListBranches_OnlyCurrentBranch(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// No other branches exist
	branches, err := git.ListBranches()
	assert.NoError(t, err, "ListBranches should succeed")
	assert.Len(t, branches, 0, "Should return empty list when only current branch exists")
}

// TestListBranches_Sorted tests that branches are returned in alphabetical order.
func TestListBranches_Sorted(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// Create branches in non-alphabetical order
	exec.Command("git", "branch", "zebra").Run()
	exec.Command("git", "branch", "alpha").Run()
	exec.Command("git", "branch", "beta").Run()

	// List branches
	branches, err := git.ListBranches()
	assert.NoError(t, err)

	// Should be sorted alphabetically
	assert.Equal(t, []string{"alpha", "beta", "zebra"}, branches, "Branches should be sorted alphabetically")
}

// TestDeleteBranch_Success tests successful branch deletion.
func TestDeleteBranch_Success(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// Create a test branch
	exec.Command("git", "branch", "test-delete").Run()

	// Delete the branch
	err = git.DeleteBranch("test-delete")
	assert.NoError(t, err, "DeleteBranch should succeed for merged branch")

	// Verify branch is deleted (it should not appear in branch list)
	branches, _ := git.ListBranches()
	assert.NotContains(t, branches, "test-delete", "Branch should be deleted")
}

// TestDeleteBranch_NonExistent tests deleting a non-existent branch.
func TestDeleteBranch_NonExistent(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// Try to delete non-existent branch
	err = git.DeleteBranch("does-not-exist")
	assert.Error(t, err, "DeleteBranch should fail for non-existent branch")
}
