package assistant

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserPromptForTests(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("creates new test prompt when no file exists", func(t *testing.T) {
		filename := filepath.Join(tempDir, "test.go")
		content := []byte("package main\n\nfunc main() {}")

		prompt, err := UserPromptForTests(content, filename)
		require.NoError(t, err)
		assert.False(t, prompt.AppendToFile)
		assert.Contains(t, prompt.TestPrompt, string(content))
		assert.Contains(t, prompt.TestPrompt, "Return full test file as go code")
		assert.Equal(t, filepath.Join(tempDir, "test_test.go"), prompt.TestFilename)
	})

	t.Run("creates append prompt when file exists", func(t *testing.T) {
		filename := filepath.Join(tempDir, "existing.go")
		testFilename := filepath.Join(tempDir, "existing_test.go")
		content := []byte("package main\n\nfunc main() {}")
		testContent := []byte("package main\n\nfunc TestMain(t *testing.T) {}")

		require.NoError(t, os.WriteFile(testFilename, testContent, 0644))

		prompt, err := UserPromptForTests(content, filename)
		require.NoError(t, err)
		assert.True(t, prompt.AppendToFile)
		assert.Contains(t, prompt.TestPrompt, string(content))
		assert.Contains(t, prompt.TestPrompt, string(testContent))
		assert.Contains(t, prompt.TestPrompt, "Add only new test functions")
	})

	t.Run("returns error when test file exists but unreadable", func(t *testing.T) {
		filename := filepath.Join(tempDir, "unreadable.go")
		testFilename := filepath.Join(tempDir, "unreadable_test.go")
		content := []byte("package main")

		require.NoError(t, os.WriteFile(testFilename, []byte{}, 0000)) // No permissions

		_, err := UserPromptForTests(content, filename)
		assert.Error(t, err)
	})
}

func TestTestFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple case", "file.go", "file_test.go"},
		{"with path", "path/to/file.go", "path/to/file_test.go"},
		{"already test", "file_test.go", "file_test_test.go"},
		{"no extension", "file", "file_test.go"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, testFilename(tt.input))
		})
	}
}
