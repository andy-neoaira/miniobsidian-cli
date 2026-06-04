package obsidian

import (
	"encoding/json"
	"errors"
	"github.com/andy-neoaira/miniobsidian-cli/pkg/config"
	"os"
	"path/filepath"
	"strings"
)

// ObsidianConfigFile 和 RunningInWSL 是包级变量，分别指向 config 包中的对应函数。
// 使用变量而非直接调用函数，方便在测试中替换为 Mock 实现。
var ObsidianConfigFile = config.ObsidianFile
var RunningInWSL = config.RunningInWSL

// Path 返回当前 vault 的绝对路径。
// 如果 Vault.Name 已经是绝对路径，直接返回（支持不依赖 Obsidian 配置文件的独立模式）。
// 否则从 Obsidian 的 vault 注册表中查找对应路径。
func (v *Vault) Path() (string, error) {
	if filepath.IsAbs(v.Name) {
		return v.Name, nil
	}

	obsidianConfigFile, err := ObsidianConfigFile()
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(obsidianConfigFile)
	if err != nil {
		return "", errors.New(ObsidianConfigReadError)
	}

	path, err := getPathForVault(content, v.Name)
	if err != nil {
		return "", err
	}

	// 如果在 WSL 环境中运行，自动转换 Windows 路径为 WSL 挂载路径
	if RunningInWSL() {
		return adjustForWslMount(path), nil
	}
	return path, nil
}

// adjustForWslMount 将 Windows 绝对路径转换为 WSL 的 /mnt/ 挂载路径。
// 例如 "C:\Users\name" 会转换为 "/mnt/c/Users/name"。
func adjustForWslMount(dir string) string {
	// 检测 Windows 盘符模式（如 C:, D:, E:）
	if len(dir) >= 2 && dir[1] == ':' && ((dir[0] >= 'A' && dir[0] <= 'Z') || (dir[0] >= 'a' && dir[0] <= 'z')) {
		driveLetter := strings.ToLower(string(dir[0]))
		mnted := "/mnt/" + driveLetter + dir[2:]
		return strings.ReplaceAll(mnted, "\\", "/")
	}

	return dir
}

// getPathForVault 从 Obsidian 配置内容中查找指定名称的 vault 路径。
func getPathForVault(content []byte, name string) (string, error) {
	vaultsContent := ObsidianVaultConfig{}
	if json.Unmarshal(content, &vaultsContent) != nil {
		return "", errors.New(ObsidianConfigParseError)
	}

	for _, element := range vaultsContent.Vaults {
		if element.Path == name ||
			strings.HasSuffix(element.Path, "/"+name) ||
			strings.HasSuffix(element.Path, "\\"+name) {
			return element.Path, nil
		}
	}

	return "", errors.New(ObsidianConfigVaultNotFoundError)
}
