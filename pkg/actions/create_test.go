package actions_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/andy-neoaira/obs-cli/mocks"
	"github.com/andy-neoaira/obs-cli/pkg/actions"
	"github.com/stretchr/testify/assert"
)

// TestCreateNote 测试 CreateNote 业务函数的各种场景。
func TestCreateNote(t *testing.T) {
	t.Run("Successful create note", func(t *testing.T) {
		// Arrange: 准备临时目录和 Mock 对象
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}
		// Act: 调用创建笔记
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName: "note",
		})
		// Assert: 验证文件已创建
		assert.NoError(t, err)
		assert.FileExists(t, filepath.Join(tmpDir, "note.md"))
	})

	t.Run("Successful create note with content", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}
		// Act
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName: "note",
			Content:  "hello world",
		})
		// Assert: 验证文件内容和预期一致
		assert.NoError(t, err)
		content, _ := os.ReadFile(filepath.Join(tmpDir, "note.md"))
		assert.Equal(t, "hello world", string(content))
	})

	t.Run("Successful create note with nested path", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}
		// Act
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName: "folder/note",
		})
		// Assert: 验证子目录和文件都已创建
		assert.NoError(t, err)
		assert.FileExists(t, filepath.Join(tmpDir, "folder", "note.md"))
	})

	t.Run("Existing file is left unchanged without overwrite or append", func(t *testing.T) {
		// Arrange: 先创建一个已有内容的笔记
		tmpDir := t.TempDir()
		notePath := filepath.Join(tmpDir, "note.md")
		if err := os.WriteFile(notePath, []byte("original"), 0644); err != nil {
			t.Fatal(err)
		}
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}
		// Act: 不带 overwrite/append 再次创建同名笔记
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName: "note",
			Content:  "new content",
		})
		// Assert: 原内容应保持不变
		assert.NoError(t, err)
		content, _ := os.ReadFile(notePath)
		assert.Equal(t, "original", string(content))
	})

	t.Run("Successful create note with overwrite", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		notePath := filepath.Join(tmpDir, "note.md")
		if err := os.WriteFile(notePath, []byte("original"), 0644); err != nil {
			t.Fatal(err)
		}
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}
		// Act: 带 overwrite 标志创建
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName:        "note",
			Content:         "overwritten",
			ShouldOverwrite: true,
		})
		// Assert: 内容应被覆盖
		assert.NoError(t, err)
		content, _ := os.ReadFile(notePath)
		assert.Equal(t, "overwritten", string(content))
	})

	t.Run("Successful create note with append", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		notePath := filepath.Join(tmpDir, "note.md")
		if err := os.WriteFile(notePath, []byte("original"), 0644); err != nil {
			t.Fatal(err)
		}
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}
		// Act: 带 append 标志创建
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName:     "note",
			Content:      " appended",
			ShouldAppend: true,
		})
		// Assert: 内容应被追加到末尾
		assert.NoError(t, err)
		content, _ := os.ReadFile(notePath)
		assert.Equal(t, "original appended", string(content))
	})

	t.Run("Successful create note with open in Obsidian", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}
		// Act: 创建并打开（Obsidian URI 模式）
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName:   "note",
			ShouldOpen: true,
			UseEditor:  false,
		})
		// Assert
		assert.NoError(t, err)
	})

	t.Run("Successful create note with open in editor", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}

		// 临时设置 EDITOR 为 "true"（Unix 中总是返回成功的命令）
		originalEditor := os.Getenv("EDITOR")
		defer os.Setenv("EDITOR", originalEditor)
		os.Setenv("EDITOR", "true")

		// Act: 创建并使用编辑器打开
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName:   "note",
			ShouldOpen: true,
			UseEditor:  true,
		})
		// Assert
		assert.NoError(t, err)
	})

	t.Run("vault.DefaultName returns an error", func(t *testing.T) {
		// Arrange: 模拟 DefaultName 出错
		vault := mocks.MockVaultOperator{
			DefaultNameErr: errors.New("Failed to get vault name"),
		}
		// Act
		err := actions.CreateNote(&vault, &mocks.MockUriManager{}, actions.CreateParams{
			NoteName: "note-name",
		})
		// Assert: 错误应原样传递
		assert.Equal(t, vault.DefaultNameErr, err)
	})

	t.Run("vault.Path returns an error", func(t *testing.T) {
		// Arrange: 模拟 Path 出错
		vault := mocks.MockVaultOperator{
			Name:      "myVault",
			PathError: errors.New("Failed to get vault path"),
		}
		// Act
		err := actions.CreateNote(&vault, &mocks.MockUriManager{}, actions.CreateParams{
			NoteName: "note-name",
		})
		// Assert: 错误应原样传递
		assert.Equal(t, vault.PathError, err)
	})

	t.Run("uri.Execute returns an error when opening in Obsidian", func(t *testing.T) {
		// Arrange: 模拟 URI 执行失败
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{
			ExecuteErr: errors.New("Failed to execute URI"),
		}
		// Act
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName:   "note-name",
			ShouldOpen: true,
			UseEditor:  false,
		})
		// Assert
		assert.Equal(t, uri.ExecuteErr, err)
	})

	t.Run("Open in editor fails when editor command fails", func(t *testing.T) {
		// Arrange: 设置 EDITOR 为 "false"（Unix 中总是返回失败的命令）
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}

		originalEditor := os.Getenv("EDITOR")
		defer os.Setenv("EDITOR", originalEditor)
		os.Setenv("EDITOR", "false")

		// Act
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName:   "note",
			ShouldOpen: true,
			UseEditor:  true,
		})
		// Assert: 编辑器命令失败应返回错误
		assert.Error(t, err)
	})

	t.Run("Uses default folder from Obsidian config", func(t *testing.T) {
		// Arrange: 创建模拟的 Obsidian 配置，指定默认新建文件夹为 Inbox
		tmpDir := t.TempDir()
		obsDir := filepath.Join(tmpDir, ".obsidian")
		if err := os.MkdirAll(obsDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(obsDir, "app.json"), []byte(`{
				"newFileLocation": "folder",
				"newFileFolderPath": "Inbox"
			}`), 0644); err != nil {
			t.Fatal(err)
		}

		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}
		// Act
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName: "note",
			Content:  "hello",
		})
		// Assert: 笔记应自动放在 Inbox 文件夹下
		assert.NoError(t, err)
		assert.FileExists(t, filepath.Join(tmpDir, "Inbox", "note.md"))
	})

	t.Run("Explicit path ignores default folder config", func(t *testing.T) {
		// Arrange: 同样配置默认文件夹，但用户显式指定了路径
		tmpDir := t.TempDir()
		obsDir := filepath.Join(tmpDir, ".obsidian")
		if err := os.MkdirAll(obsDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(obsDir, "app.json"), []byte(`{
				"newFileLocation": "folder",
				"newFileFolderPath": "Inbox"
			}`), 0644); err != nil {
			t.Fatal(err)
		}

		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}
		// Act: 用户显式输入 "sub/note" 包含路径分隔符
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName: "sub/note",
			Content:  "hello",
		})
		// Assert: 应放在用户指定的 sub/ 下，而不是 Inbox/sub/ 下
		assert.NoError(t, err)
		assert.FileExists(t, filepath.Join(tmpDir, "sub", "note.md"))
		assert.NoFileExists(t, filepath.Join(tmpDir, "Inbox", "sub", "note.md"))
	})

	t.Run("UseEditor without open does not use editor", func(t *testing.T) {
		// Arrange
		tmpDir := t.TempDir()
		vault := mocks.MockVaultOperator{Name: "myVault", PathValue: tmpDir}
		uri := mocks.MockUriManager{}

		// Act: UseEditor 为 true 但 ShouldOpen 为 false
		err := actions.CreateNote(&vault, &uri, actions.CreateParams{
			NoteName:   "note",
			ShouldOpen: false,
			UseEditor:  true,
		})
		// Assert: 文件应被创建，但编辑器不应被调用
		assert.NoError(t, err)
		assert.FileExists(t, filepath.Join(tmpDir, "note.md"))
	})
}

// TestNormalizeContent 测试内容转义字符的还原功能。
func TestNormalizeContent(t *testing.T) {
	t.Run("Replaces escape sequences with actual characters", func(t *testing.T) {
		// Arrange: 输入带转义字符的字符串
		input := "Hello\\nWorld\\tTabbed\\rReturn\\\"Quote\\'SingleQuote\\\\Backslash"
		expected := "Hello\nWorld\tTabbed\rReturn\"Quote'SingleQuote\\Backslash"

		// Act: 调用 NormalizeContent 还原
		result := actions.NormalizeContent(input)

		// Assert: 验证转义序列被正确替换
		assert.Equal(t, expected, result, "The content should have the escape sequences replaced correctly")
	})

	t.Run("Handles empty input", func(t *testing.T) {
		// Arrange
		input := ""
		expected := ""

		// Act
		result := actions.NormalizeContent(input)

		// Assert
		assert.Equal(t, expected, result, "Empty input should return empty output")
	})

	t.Run("No escape sequences in input", func(t *testing.T) {
		// Arrange
		input := "Plain text with no escapes"
		expected := "Plain text with no escapes"

		// Act
		result := actions.NormalizeContent(input)

		// Assert: 无转义序列时应保持原样
		assert.Equal(t, expected, result, "Content without escape sequences should remain unchanged")
	})
}
