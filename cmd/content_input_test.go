package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveContentInput(t *testing.T) {
	t.Run("Uses inline content when no file is provided", func(t *testing.T) {
		content, err := resolveContentInput("hello", "")

		assert.NoError(t, err)
		assert.Equal(t, "hello", content)
	})

	t.Run("Reads content from file", func(t *testing.T) {
		tmpDir := t.TempDir()
		contentPath := filepath.Join(tmpDir, "content.md")
		assert.NoError(t, os.WriteFile(contentPath, []byte("# Title\nbody"), 0644))

		content, err := resolveContentInput("", contentPath)

		assert.NoError(t, err)
		assert.Equal(t, "# Title\nbody", content)
	})

	t.Run("Rejects inline content and file content together", func(t *testing.T) {
		content, err := resolveContentInput("inline", "content.md")

		assert.Error(t, err)
		assert.Empty(t, content)
		assert.Contains(t, err.Error(), "cannot be used together")
	})
}
