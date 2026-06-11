package actions_test

import (
	"errors"
	"os"
	"testing"

	"github.com/andy-neoaira/obs-cli/mocks"
	"github.com/andy-neoaira/obs-cli/pkg/actions"
	"github.com/stretchr/testify/assert"
)

func TestMoveNote(t *testing.T) {
	t.Run("Successful move note", func(t *testing.T) {
		// Arrange
		vault := mocks.MockVaultOperator{Name: "myVault"}
		uri := mocks.MockUriManager{}
		note := mocks.MockNoteManager{}
		linkRewriter := mocks.MockLinkRewriter{}
		// Act
		err := actions.MoveNote(&vault, &note, &linkRewriter, &uri, actions.MoveParams{
			CurrentNoteName: "string",
			NewNoteName:     "string",
			ShouldOpen:      true,
		})
		// Assert
		assert.NoError(t, err, "Expected no error")
	})

	t.Run("vault.DefaultName returns an error", func(t *testing.T) {
		// Arrange
		vault := mocks.MockVaultOperator{
			DefaultNameErr: errors.New("Failed to get vault name"),
		}
		// Act
		err := actions.MoveNote(&vault, &mocks.MockNoteManager{}, &mocks.MockLinkRewriter{}, &mocks.MockUriManager{}, actions.MoveParams{
			CurrentNoteName: "string",
			NewNoteName:     "string",
			ShouldOpen:      true,
		})
		// Assert
		assert.Equal(t, err, vault.DefaultNameErr)
	})

	t.Run("vault.Path returns an error", func(t *testing.T) {
		// Arrange
		vaultOp := &mocks.MockVaultOperator{
			PathError: errors.New("Failed to get vault path"),
		}
		// Act
		err := actions.MoveNote(vaultOp, &mocks.MockNoteManager{}, &mocks.MockLinkRewriter{}, &mocks.MockUriManager{}, actions.MoveParams{
			CurrentNoteName: "string",
			NewNoteName:     "string",
			ShouldOpen:      false,
		})
		// Assert
		assert.Equal(t, err, vaultOp.PathError)
	})

	t.Run("note.Move returns an error", func(t *testing.T) {
		// Arrange
		note := mocks.MockNoteManager{
			MoveErr: errors.New("Failed to execute URI"),
		}
		// Act
		err := actions.MoveNote(&mocks.MockVaultOperator{}, &note, &mocks.MockLinkRewriter{}, &mocks.MockUriManager{}, actions.MoveParams{
			CurrentNoteName: "string",
			NewNoteName:     "string",
			ShouldOpen:      false,
		})
		// Assert
		assert.Equal(t, err, note.MoveErr)
	})

	t.Run("linkRewriter.UpdateLinks returns an error", func(t *testing.T) {
		// Arrange
		linkRewriter := mocks.MockLinkRewriter{
			UpdateLinksError: errors.New("Failed to execute URI"),
		}
		// Act
		err := actions.MoveNote(&mocks.MockVaultOperator{}, &mocks.MockNoteManager{}, &linkRewriter, &mocks.MockUriManager{}, actions.MoveParams{
			CurrentNoteName: "string",
			NewNoteName:     "string",
			ShouldOpen:      false,
		})
		// Assert
		assert.Equal(t, err, linkRewriter.UpdateLinksError)
	})

	t.Run("uri.Execute returns an error", func(t *testing.T) {
		// Arrange
		uriManager := &mocks.MockUriManager{
			ExecuteErr: errors.New("Failed to execute URI"),
		}
		// Act
		err := actions.MoveNote(&mocks.MockVaultOperator{}, &mocks.MockNoteManager{}, &mocks.MockLinkRewriter{}, uriManager, actions.MoveParams{
			CurrentNoteName: "string",
			NewNoteName:     "string",
			ShouldOpen:      true,
			UseEditor:       false,
		})
		// Assert
		assert.Equal(t, err, uriManager.ExecuteErr)
	})

	t.Run("Successful move note with editor flag and open", func(t *testing.T) {
		// Arrange
		vault := mocks.MockVaultOperator{Name: "myVault"}
		uri := mocks.MockUriManager{}
		note := mocks.MockNoteManager{}
		linkRewriter := mocks.MockLinkRewriter{}

		// Set EDITOR to a command that will succeed
		originalEditor := os.Getenv("EDITOR")
		defer os.Setenv("EDITOR", originalEditor)
		os.Setenv("EDITOR", "true")

		// Act
		err := actions.MoveNote(&vault, &note, &linkRewriter, &uri, actions.MoveParams{
			CurrentNoteName: "old.md",
			NewNoteName:     "new.md",
			ShouldOpen:      true,
			UseEditor:       true,
		})

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Move note with editor flag fails when editor fails", func(t *testing.T) {
		// Arrange
		vault := mocks.MockVaultOperator{Name: "myVault"}
		uri := mocks.MockUriManager{}
		note := mocks.MockNoteManager{}
		linkRewriter := mocks.MockLinkRewriter{}

		// Set EDITOR to a command that will fail
		originalEditor := os.Getenv("EDITOR")
		defer os.Setenv("EDITOR", originalEditor)
		os.Setenv("EDITOR", "false")

		// Act
		err := actions.MoveNote(&vault, &note, &linkRewriter, &uri, actions.MoveParams{
			CurrentNoteName: "old.md",
			NewNoteName:     "new.md",
			ShouldOpen:      true,
			UseEditor:       true,
		})

		// Assert
		assert.Error(t, err)
	})

	t.Run("Move note with editor flag without open does not use editor", func(t *testing.T) {
		// Arrange
		vault := mocks.MockVaultOperator{Name: "myVault"}
		uri := mocks.MockUriManager{}
		note := mocks.MockNoteManager{}
		linkRewriter := mocks.MockLinkRewriter{}

		// Act - UseEditor is true but ShouldOpen is false
		err := actions.MoveNote(&vault, &note, &linkRewriter, &uri, actions.MoveParams{
			CurrentNoteName: "old.md",
			NewNoteName:     "new.md",
			ShouldOpen:      false,
			UseEditor:       true,
		})

		// Assert - should succeed without opening
		assert.NoError(t, err)
	})
}
