package contract

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRepo creates a temporary git repository for contract testing.
func setupTestRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	exec.Command("git", "init", dir).Run()
	exec.Command("git", "-C", dir, "config", "user.name", "Test User").Run()
	exec.Command("git", "-C", dir, "config", "user.email", "test@example.com").Run()
	exec.Command("git", "-C", dir, "commit", "--allow-empty", "-m", "Initial commit").Run()

	return dir
}

// TestContract_RepositoryValidation tests Contract 1: Repository validation
// Given: User runs gelete outside a git repository
// Then: Display error message and exit with code 1
func TestContract_RepositoryValidation(t *testing.T) {
	// Create a non-git directory
	dir := t.TempDir()

	// Build the gelete binary
	buildCmd := exec.Command("go", "build", "-o", "gelete-test", ".")
	buildCmd.Dir = getProjectRoot(t)
	err := buildCmd.Run()
	require.NoError(t, err, "Failed to build gelete")

	// Run gelete in non-git directory
	cmd := exec.Command("./gelete-test")
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()

	// Should fail with exit code 1
	assert.Error(t, err)
	exitErr, ok := err.(*exec.ExitError)
	require.True(t, ok)
	assert.Equal(t, 1, exitErr.ExitCode(), "Should exit with code 1")

	// Should display error message
	stderrStr := stderr.String()
	assert.Contains(t, stderrStr, "not a git repository", "Error should mention 'not a git repository'")
}

// TestContract_HelpFlag tests Contract 12: Help flag
// Given: User runs `gelete --help`
// Then: Display help text and exit with code 0
func TestContract_HelpFlag(t *testing.T) {
	// Build the gelete binary
	buildCmd := exec.Command("go", "build", "-o", "gelete-test", ".")
	buildCmd.Dir = getProjectRoot(t)
	err := buildCmd.Run()
	require.NoError(t, err, "Failed to build gelete")

	// Run gelete with --help
	cmd := exec.Command("./gelete-test", "--help")
	cmd.Dir = getProjectRoot(t)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err = cmd.Run()

	// Should succeed with exit code 0
	assert.NoError(t, err, "Should exit with code 0")

	// Should display help text
	stdoutStr := stdout.String()
	assert.Contains(t, stdoutStr, "gelete", "Help should mention command name")
	assert.Contains(t, stdoutStr, "Usage", "Help should include usage section")
	assert.Contains(t, stdoutStr, "help", "Help should mention help flag")
}

// TestContract_VersionFlag tests Contract 13: Version flag
// Given: User runs `gelete --version`
// Then: Display version and exit with code 0
func TestContract_VersionFlag(t *testing.T) {
	// Build the gelete binary
	buildCmd := exec.Command("go", "build", "-o", "gelete-test", ".")
	buildCmd.Dir = getProjectRoot(t)
	err := buildCmd.Run()
	require.NoError(t, err, "Failed to build gelete")

	// Run gelete with --version
	cmd := exec.Command("./gelete-test", "--version")
	cmd.Dir = getProjectRoot(t)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err = cmd.Run()

	// Should succeed with exit code 0
	assert.NoError(t, err, "Should exit with code 0")

	// Should display version
	stdoutStr := stdout.String()
	assert.Contains(t, stdoutStr, "gelete", "Version output should mention command name")
}

// TestContract_NoDeletableBranches tests Contract 11: No deletable branches
// Given: Repository has only 1 branch (current branch)
// Then: Display "No branches available for deletion" and exit with code 0
func TestContract_NoDeletableBranches(t *testing.T) {
	repo := setupTestRepo(t)

	// Build the gelete binary
	buildCmd := exec.Command("go", "build", "-o", "gelete-test", ".")
	buildCmd.Dir = getProjectRoot(t)
	err := buildCmd.Run()
	require.NoError(t, err, "Failed to build gelete")

	// Run gelete in repo with only one branch
	cmd := exec.Command("./gelete-test")
	cmd.Dir = repo
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err = cmd.Run()

	// Should succeed with exit code 0
	assert.NoError(t, err, "Should exit with code 0 when no branches to delete")

	// Should display appropriate message
	stdoutStr := stdout.String()
	assert.Contains(t, stdoutStr, "No branches available", "Should indicate no branches available")
}

// getProjectRoot returns the path to the project root directory.
func getProjectRoot(t *testing.T) string {
	t.Helper()

	// Get current working directory
	cwd, err := os.Getwd()
	require.NoError(t, err)

	// Navigate up to project root (from tests/contract/ to project root)
	// This assumes we're in tests/contract/
	root := cwd
	for i := 0; i < 2; i++ {
		parent := strings.TrimSuffix(root, "/tests/contract")
		parent = strings.TrimSuffix(parent, "/tests")
		if parent != root {
			root = parent
			break
		}
		// If not in expected path, try going up one level
		parts := strings.Split(root, "/")
		if len(parts) > 1 {
			root = strings.Join(parts[:len(parts)-1], "/")
		}
	}

	return root
}
