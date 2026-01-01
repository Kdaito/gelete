package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/Kdaito/gelete/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWorktree_ListWorktrees tests listing worktrees in a repository.
// This verifies FR-010: System MUST detect branches that are checked out as git worktrees.
func TestWorktree_ListWorktrees(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// Create a branch and a worktree
	exec.Command("git", "branch", "feature-1").Run()
	worktreePath := t.TempDir()
	exec.Command("git", "worktree", "add", worktreePath, "feature-1").Run()

	// Resolve symlinks for comparison
	expectedPath, _ := filepath.EvalSymlinks(worktreePath)

	// List worktrees
	worktrees, err := git.ListWorktrees()
	assert.NoError(t, err, "ListWorktrees should succeed")
	assert.GreaterOrEqual(t, len(worktrees), 1, "Should have at least one worktree")

	// Check if feature-1 is in the worktree list
	found := false
	for _, wt := range worktrees {
		if wt.Branch == "feature-1" {
			found = true
			assert.Equal(t, expectedPath, wt.Path, "Worktree path should match")
			break
		}
	}
	assert.True(t, found, "feature-1 should be listed as a worktree")

	// Cleanup
	exec.Command("git", "worktree", "remove", worktreePath).Run()
}

// TestWorktree_RemoveWorktree tests removing a worktree.
// This verifies FR-013: System MUST remove worktree directory before deleting branch.
func TestWorktree_RemoveWorktree(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// Create a branch and a worktree
	exec.Command("git", "branch", "feature-2").Run()
	worktreePath := t.TempDir()
	exec.Command("git", "worktree", "add", worktreePath, "feature-2").Run()

	// Remove the worktree
	err = git.RemoveWorktree(worktreePath)
	assert.NoError(t, err, "RemoveWorktree should succeed")

	// Verify worktree is removed
	worktrees, _ := git.ListWorktrees()
	for _, wt := range worktrees {
		assert.NotEqual(t, "feature-2", wt.Branch, "feature-2 should not be in worktree list")
	}

	// Now branch can be deleted normally
	err = git.DeleteBranch("feature-2")
	assert.NoError(t, err, "DeleteBranch should succeed after worktree removal")
}

// TestWorktree_ForceRemoveWorktree tests force removing a worktree.
// This verifies FR-014: System MUST handle locked worktrees.
func TestWorktree_ForceRemoveWorktree(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// Create a branch and a worktree
	exec.Command("git", "branch", "feature-3").Run()
	worktreePath := t.TempDir()
	exec.Command("git", "worktree", "add", worktreePath, "feature-3").Run()

	// Lock the worktree
	exec.Command("git", "worktree", "lock", worktreePath).Run()

	// Normal remove should fail
	err = git.RemoveWorktree(worktreePath)
	assert.Error(t, err, "RemoveWorktree should fail for locked worktree")

	// Force remove should succeed
	err = git.ForceRemoveWorktree(worktreePath)
	assert.NoError(t, err, "ForceRemoveWorktree should succeed for locked worktree")

	// Verify worktree is removed
	worktrees, _ := git.ListWorktrees()
	for _, wt := range worktrees {
		assert.NotEqual(t, "feature-3", wt.Branch, "feature-3 should not be in worktree list")
	}
}

// TestWorktree_BranchDeletionWithWorktree tests that branch deletion fails when worktree exists.
func TestWorktree_BranchDeletionWithWorktree(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// Create a branch and a worktree
	exec.Command("git", "branch", "feature-4").Run()
	worktreePath := t.TempDir()
	exec.Command("git", "worktree", "add", worktreePath, "feature-4").Run()

	// Attempting to delete branch with active worktree should fail
	err = git.DeleteBranch("feature-4")
	assert.Error(t, err, "DeleteBranch should fail when worktree exists")

	// After removing worktree, deletion should succeed
	git.RemoveWorktree(worktreePath)
	err = git.DeleteBranch("feature-4")
	assert.NoError(t, err, "DeleteBranch should succeed after worktree removal")
}

// TestWorktree_MultipleWorktrees tests handling multiple worktrees.
func TestWorktree_MultipleWorktrees(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// Create multiple branches and worktrees
	exec.Command("git", "branch", "feature-a").Run()
	exec.Command("git", "branch", "feature-b").Run()
	worktreePathA := t.TempDir()
	worktreePathB := t.TempDir()
	exec.Command("git", "worktree", "add", worktreePathA, "feature-a").Run()
	exec.Command("git", "worktree", "add", worktreePathB, "feature-b").Run()

	// List worktrees
	worktrees, err := git.ListWorktrees()
	assert.NoError(t, err)

	// Count how many of our feature branches are in worktrees
	count := 0
	for _, wt := range worktrees {
		if wt.Branch == "feature-a" || wt.Branch == "feature-b" {
			count++
		}
	}
	assert.Equal(t, 2, count, "Should have 2 feature worktrees")

	// Cleanup
	git.RemoveWorktree(worktreePathA)
	git.RemoveWorktree(worktreePathB)
	git.DeleteBranch("feature-a")
	git.DeleteBranch("feature-b")
}
