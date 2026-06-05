package cmd

import (
	"bytes"
	"testing"

	"github.com/andy-neoaira/obs-cli/pkg/obsidian"
	"github.com/stretchr/testify/assert"
)

// TestFormatVaultsTable 测试 vault 列表的表格格式化输出。
func TestFormatVaultsTable(t *testing.T) {
	t.Run("Aligns columns with varying name lengths", func(t *testing.T) {
		// 准备不同名称长度的 vault 数据
		vaults := []obsidian.VaultInfo{
			{Name: "Notes", Path: "/home/user/Notes"},
			{Name: "LongVaultName", Path: "/home/user/LongVaultName"},
			{Name: "Work", Path: "/home/user/Work"},
		}

		var buf bytes.Buffer
		formatVaultsTable(&buf, vaults, "")
		output := buf.String()

		// 验证输出了 3 行
		lines := bytes.Split(bytes.TrimSpace([]byte(output)), []byte("\n"))
		assert.Len(t, lines, 3)

		// 每行都应包含名称和路径
		assert.Contains(t, output, "Notes")
		assert.Contains(t, output, "/home/user/Notes")
		assert.Contains(t, output, "LongVaultName")
		assert.Contains(t, output, "/home/user/LongVaultName")

		// tabwriter 应保证路径列对齐：找到每行路径的起始偏移量并比较
		pathOffsets := make([]int, len(lines))
		for i, line := range lines {
			pathOffsets[i] = bytes.Index(line, []byte("/home"))
		}
		assert.Equal(t, pathOffsets[0], pathOffsets[1], "path columns should be aligned")
		assert.Equal(t, pathOffsets[1], pathOffsets[2], "path columns should be aligned")
	})

	t.Run("Single vault produces output", func(t *testing.T) {
		vaults := []obsidian.VaultInfo{
			{Name: "MyVault", Path: "/tmp/MyVault"},
		}

		var buf bytes.Buffer
		formatVaultsTable(&buf, vaults, "")
		output := buf.String()

		assert.Contains(t, output, "MyVault")
		assert.Contains(t, output, "/tmp/MyVault")
	})

	t.Run("Marks default vault", func(t *testing.T) {
		vaults := []obsidian.VaultInfo{
			{Name: "Notes", Path: "/home/user/Notes"},
			{Name: "Work", Path: "/home/user/Work"},
		}

		var buf bytes.Buffer
		formatVaultsTable(&buf, vaults, "Work")
		output := buf.String()

		assert.Contains(t, output, "Work")
		assert.Contains(t, output, "(default)")
		// Notes 行不应包含 (default)
		lines := bytes.Split(bytes.TrimSpace([]byte(output)), []byte("\n"))
		assert.Len(t, lines, 2)
		assert.NotContains(t, string(lines[0]), "(default)")
		assert.Contains(t, string(lines[1]), "(default)")
	})

	t.Run("Empty vault list produces no output", func(t *testing.T) {
		var buf bytes.Buffer
		formatVaultsTable(&buf, []obsidian.VaultInfo{}, "")

		assert.Empty(t, buf.String())
	})
}
