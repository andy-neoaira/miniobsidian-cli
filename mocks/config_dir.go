package mocks

import (
	"testing"
)

// CreateMockObsidianConfigFile 创建临时的 Obsidian 配置文件路径，用于测试。
// 使用 t.TempDir() 确保测试结束后自动清理。
func CreateMockObsidianConfigFile(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	return tmpDir + "/obsidian.json"
}

// CreateMockCliConfigDirectories 创建临时的 CLI 配置目录和文件路径，用于测试。
func CreateMockCliConfigDirectories(t *testing.T) (string, string) {
	t.Helper()
	tmpDir := t.TempDir()
	return tmpDir, tmpDir + "/preferences.json"
}
