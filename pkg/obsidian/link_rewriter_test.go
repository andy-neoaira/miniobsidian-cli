package obsidian_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/andy-neoaira/obs-cli/pkg/obsidian"
	"github.com/stretchr/testify/assert"
)

func TestLinkRewriter_GenerateReplacements(t *testing.T) {
	t.Run("Can disable basename wikilink replacement", func(t *testing.T) {
		rewriter := obsidian.LinkRewriter{}

		replacements := rewriter.GenerateReplacements("folder/oldNote", "folder/newNote", obsidian.LinkRewriteOptions{
			IncludeBaseLinks: false,
		})

		assert.NotContains(t, replacements, "[[oldNote]]")
		assert.Equal(t, "[[folder/newNote]]", replacements["[[folder/oldNote]]"])
		assert.Equal(t, "](folder/newNote.md)", replacements["](folder/oldNote.md)"])
	})

	t.Run("Can include basename wikilink replacement", func(t *testing.T) {
		rewriter := obsidian.LinkRewriter{}

		replacements := rewriter.GenerateReplacements("folder/oldNote", "folder/newNote", obsidian.LinkRewriteOptions{
			IncludeBaseLinks: true,
		})

		assert.Equal(t, "[[newNote]]", replacements["[[oldNote]]"])
		assert.Equal(t, "[[folder/newNote]]", replacements["[[folder/oldNote]]"])
	})
}

func TestLinkRewriter_UpdateLinks(t *testing.T) {
	t.Run("Updates links without going through Note", func(t *testing.T) {
		tmpDir := t.TempDir()
		notePath := filepath.Join(tmpDir, "links.md")
		assert.NoError(t, os.WriteFile(notePath, []byte("See [[oldNote]] and [md](oldNote.md)"), 0644))

		rewriter := obsidian.LinkRewriter{}
		err := rewriter.UpdateLinks(tmpDir, "oldNote", "newNote")

		assert.NoError(t, err)
		content, readErr := os.ReadFile(notePath)
		assert.NoError(t, readErr)
		assert.Equal(t, "See [[newNote]] and [md](newNote.md)", string(content))
	})
}
