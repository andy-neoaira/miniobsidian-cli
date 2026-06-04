package obsidian

import (
	"errors"
	"os"
	"sort"
	"strings"
)

// ListEntries 列出 vault 中指定路径下的所有文件和文件夹。
// 返回的目录名以 "/" 结尾，文件夹在前、文件在后，且各自按字母顺序排序。
func ListEntries(vaultPath, relativePath string) ([]string, error) {
	targetPath := vaultPath
	if strings.TrimSpace(relativePath) != "" {
		validatedPath, err := ValidatePath(vaultPath, relativePath)
		if err != nil {
			return nil, err
		}
		targetPath = validatedPath
	}

	info, err := os.Stat(targetPath)
	if err != nil {
		return nil, errors.New(VaultAccessError)
	}
	if !info.IsDir() {
		return nil, errors.New(VaultAccessError)
	}

	entries, err := os.ReadDir(targetPath)
	if err != nil {
		return nil, errors.New(VaultReadError)
	}

	// 分离目录和文件，便于分别排序后按目录优先的顺序返回
	dirs := make([]string, 0, len(entries))
	files := make([]string, 0, len(entries))

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue // 跳过隐藏文件和目录
		}
		if entry.IsDir() {
			dirs = append(dirs, name+"/")
			continue
		}
		files = append(files, name)
	}

	sort.Strings(dirs)
	sort.Strings(files)

	return append(dirs, files...), nil
}
