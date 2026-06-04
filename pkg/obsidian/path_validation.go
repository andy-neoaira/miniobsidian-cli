package obsidian

import (
	"errors"
	"path/filepath"
	"strings"
)

// ErrPathTraversal 当路径试图逃逸出基础目录时返回此错误。
var ErrPathTraversal = errors.New("path traversal detected: path must remain within vault directory")

// ValidatePath 确保 relativePath 与 basePath 拼接后仍位于 basePath 内部。
// 返回清理后的绝对路径；如果路径试图逃逸出基础目录，返回 ErrPathTraversal。
//
// 安全检查步骤：
//  1. 拒绝绝对路径（防止直接传入 /etc/passwd）
//  2. 清理并拼接路径
//  3. 验证结果路径是否以 basePath 为前缀
func ValidatePath(basePath, relativePath string) (string, error) {
	// 拒绝绝对路径输入
	if filepath.IsAbs(relativePath) {
		return "", ErrPathTraversal
	}

	// 将基础路径转换为清理后的绝对路径
	absBase, err := filepath.Abs(filepath.Clean(basePath))
	if err != nil {
		return "", err
	}

	// 清理相对路径并与基础路径拼接
	cleanRelative := filepath.Clean(relativePath)
	joinedPath := filepath.Join(absBase, cleanRelative)

	// 获取拼接后的绝对路径
	absJoined, err := filepath.Abs(joinedPath)
	if err != nil {
		return "", err
	}

	// 验证拼接后的路径是否以基础路径为前缀。
	// 在 basePath 末尾加上路径分隔符，防止部分匹配（如 /vault-backup 匹配 /vault）。
	if !strings.HasPrefix(absJoined, absBase+string(filepath.Separator)) && absJoined != absBase {
		return "", ErrPathTraversal
	}

	return absJoined, nil
}
