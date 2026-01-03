package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGoReleaserConfig_Valid tests that .goreleaser.yml is valid and can be checked.
// This verifies FR-021 to FR-025: Multi-platform binary distribution.
func TestGoReleaserConfig_Valid(t *testing.T) {
	// Get project root
	projectRoot, err := filepath.Abs("../..")
	require.NoError(t, err)

	// Check if .goreleaser.yml exists
	goreleaserPath := filepath.Join(projectRoot, ".goreleaser.yml")
	_, err = os.Stat(goreleaserPath)
	require.NoError(t, err, ".goreleaser.yml should exist")

	// Check goreleaser config validity
	cmd := exec.Command("goreleaser", "check")
	cmd.Dir = projectRoot
	output, err := cmd.CombinedOutput()

	if err != nil {
		// If goreleaser is not installed, skip this test
		if strings.Contains(string(output), "executable file not found") ||
			strings.Contains(err.Error(), "executable file not found") {
			t.Skip("goreleaser not installed, skipping validation test")
		}
		t.Logf("goreleaser check output: %s", string(output))
	}

	// If goreleaser is installed, config should be valid
	if err == nil || !strings.Contains(string(output), "not found") {
		assert.NoError(t, err, "goreleaser config should be valid")
	}
}

// TestGoReleaserConfig_HasRequiredBuilds tests that the config includes all required platforms.
// This verifies FR-021, FR-022, FR-023: Builds for linux, darwin, windows on amd64 and arm64.
func TestGoReleaserConfig_HasRequiredBuilds(t *testing.T) {
	projectRoot, err := filepath.Abs("../..")
	require.NoError(t, err)

	goreleaserPath := filepath.Join(projectRoot, ".goreleaser.yml")
	content, err := os.ReadFile(goreleaserPath)
	require.NoError(t, err, ".goreleaser.yml should exist and be readable")

	configStr := string(content)

	// Check for required OS platforms (FR-021, FR-022, FR-023)
	assert.Contains(t, configStr, "linux", "Should build for linux")
	assert.Contains(t, configStr, "darwin", "Should build for darwin")
	assert.Contains(t, configStr, "windows", "Should build for windows")

	// Check for required architectures
	assert.Contains(t, configStr, "amd64", "Should build for amd64")
	assert.Contains(t, configStr, "arm64", "Should build for arm64")
}

// TestGoReleaserConfig_HasChecksums tests that checksums are configured.
// This verifies that SHA256 checksums are generated for security (FR-024 implied).
func TestGoReleaserConfig_HasChecksums(t *testing.T) {
	projectRoot, err := filepath.Abs("../..")
	require.NoError(t, err)

	goreleaserPath := filepath.Join(projectRoot, ".goreleaser.yml")
	content, err := os.ReadFile(goreleaserPath)
	require.NoError(t, err)

	configStr := string(content)

	// Check for checksum configuration
	assert.Contains(t, configStr, "checksum", "Should have checksum configuration")
}

// TestGoReleaserConfig_HasHomebrew tests that Homebrew formula generation is configured.
// This verifies FR-025: Installation via Homebrew.
func TestGoReleaserConfig_HasHomebrew(t *testing.T) {
	projectRoot, err := filepath.Abs("../..")
	require.NoError(t, err)

	goreleaserPath := filepath.Join(projectRoot, ".goreleaser.yml")
	content, err := os.ReadFile(goreleaserPath)
	require.NoError(t, err)

	configStr := string(content)

	// Check for Homebrew configuration
	assert.Contains(t, configStr, "brew", "Should have Homebrew configuration")
}

// TestBuildMatrix_AllPlatforms tests that CI build matrix covers all required platforms.
// This validates that the CI pipeline builds for all 6 platform/arch combinations.
func TestBuildMatrix_AllPlatforms(t *testing.T) {
	projectRoot, err := filepath.Abs("../..")
	require.NoError(t, err)

	ciPath := filepath.Join(projectRoot, ".github/workflows/ci.yml")
	content, err := os.ReadFile(ciPath)
	require.NoError(t, err, "ci.yml should exist")

	ciStr := string(content)

	// Check build matrix includes all required platforms
	assert.Contains(t, ciStr, "goos: [linux, darwin, windows]", "Should build for all OS platforms")
	assert.Contains(t, ciStr, "goarch: [amd64, arm64]", "Should build for all architectures")

	// This matrix creates 3 * 2 = 6 combinations:
	// linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64, windows/arm64
}
