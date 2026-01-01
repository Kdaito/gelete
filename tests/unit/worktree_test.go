package unit

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/Kdaito/gelete/internal/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestListWorktrees_NoWorktrees tests listing worktrees when none exist.
func TestListWorktrees_NoWorktrees(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// List worktrees (should include main worktree only)
	worktrees, err := git.ListWorktrees()
	assert.NoError(t, err, "ListWorktrees should succeed")
	// Main repository is also a worktree
	assert.GreaterOrEqual(t, len(worktrees), 1, "Should have at least main worktree")
}

// TestListWorktrees_WithWorktrees tests listing worktrees when they exist.
func TestListWorktrees_WithWorktrees(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// Create a worktree
	exec.Command("git", "branch", "test-wt").Run()
	worktreePath := t.TempDir()
	exec.Command("git", "worktree", "add", worktreePath, "test-wt").Run()

	// Resolve symlinks in expected path for comparison
	expectedPath, _ := filepath.EvalSymlinks(worktreePath)

	// List worktrees
	worktrees, err := git.ListWorktrees()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(worktrees), 2, "Should have main + test worktree")

	// Find the test worktree
	found := false
	for _, wt := range worktrees {
		if wt.Branch == "test-wt" {
			found = true
			assert.Equal(t, expectedPath, wt.Path)
			assert.False(t, wt.Locked, "New worktree should not be locked")
			break
		}
	}
	assert.True(t, found, "test-wt should be in worktree list")

	// Cleanup
	exec.Command("git", "worktree", "remove", worktreePath).Run()
}

// TestRemoveWorktree_Success tests successful worktree removal.
func TestRemoveWorktree_Success(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// Create a worktree
	exec.Command("git", "branch", "test-rm").Run()
	worktreePath := t.TempDir()
	exec.Command("git", "worktree", "add", worktreePath, "test-rm").Run()

	// Remove worktree
	err = git.RemoveWorktree(worktreePath)
	assert.NoError(t, err, "RemoveWorktree should succeed")

	// Verify removal
	worktrees, _ := git.ListWorktrees()
	for _, wt := range worktrees {
		assert.NotEqual(t, "test-rm", wt.Branch, "Removed worktree should not be listed")
	}
}

// TestRemoveWorktree_LockedWorktree tests that removing a locked worktree fails.
func TestRemoveWorktree_LockedWorktree(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// Create and lock a worktree
	exec.Command("git", "branch", "test-locked").Run()
	worktreePath := t.TempDir()
	exec.Command("git", "worktree", "add", worktreePath, "test-locked").Run()
	exec.Command("git", "worktree", "lock", worktreePath).Run()

	// Attempt to remove locked worktree
	err = git.RemoveWorktree(worktreePath)
	assert.Error(t, err, "RemoveWorktree should fail for locked worktree")

	// Cleanup with force
	exec.Command("git", "worktree", "remove", "--force", worktreePath).Run()
}

// TestForceRemoveWorktree_Success tests force removal of worktree.
func TestForceRemoveWorktree_Success(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// Create a worktree
	exec.Command("git", "branch", "test-force").Run()
	worktreePath := t.TempDir()
	exec.Command("git", "worktree", "add", worktreePath, "test-force").Run()

	// Force remove worktree
	err = git.ForceRemoveWorktree(worktreePath)
	assert.NoError(t, err, "ForceRemoveWorktree should succeed")

	// Verify removal
	worktrees, _ := git.ListWorktrees()
	for _, wt := range worktrees {
		assert.NotEqual(t, "test-force", wt.Branch, "Removed worktree should not be listed")
	}
}

// TestForceRemoveWorktree_LockedWorktree tests force removing a locked worktree.
func TestForceRemoveWorktree_LockedWorktree(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// Create and lock a worktree
	exec.Command("git", "branch", "test-force-locked").Run()
	worktreePath := t.TempDir()
	exec.Command("git", "worktree", "add", worktreePath, "test-force-locked").Run()
	exec.Command("git", "worktree", "lock", worktreePath).Run()

	// Force remove should succeed even if locked
	err = git.ForceRemoveWorktree(worktreePath)
	assert.NoError(t, err, "ForceRemoveWorktree should succeed for locked worktree")

	// Verify removal
	worktrees, _ := git.ListWorktrees()
	for _, wt := range worktrees {
		assert.NotEqual(t, "test-force-locked", wt.Branch, "Removed worktree should not be listed")
	}
}

// TestRemoveWorktree_NonExistent tests removing a non-existent worktree.
func TestRemoveWorktree_NonExistent(t *testing.T) {
	repo := setupTestRepo(t)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	err := os.Chdir(repo)
	require.NoError(t, err)

	// Attempt to remove non-existent worktree
	err = git.RemoveWorktree("/path/does/not/exist")
	assert.Error(t, err, "RemoveWorktree should fail for non-existent worktree")
}
