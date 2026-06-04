package actions

import (
	"github.com/andy-neoaira/miniobsidian-cli/pkg/obsidian"
)

// ListParams 定义了 list 命令所需的业务参数。
type ListParams struct {
	Path string // 要列出的目标路径（相对于 vault 根目录），空字符串表示根目录
}

// ListEntries 是 "list" 命令的业务核心。
// 获取指定 vault 下某路径中的所有文件和文件夹列表。
func ListEntries(vault obsidian.VaultManager, params ListParams) ([]string, error) {
	// 触发默认 vault 校验（如果没有设置会报错）
	_, err := vault.DefaultName()
	if err != nil {
		return nil, err
	}

	vaultPath, err := vault.Path()
	if err != nil {
		return nil, err
	}

	// 委托给 obsidian 包执行实际的文件系统遍历
	return obsidian.ListEntries(vaultPath, params.Path)
}
