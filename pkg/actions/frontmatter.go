package actions

import (
	"errors"
	"fmt"

	"github.com/andy-neoaira/obs-cli/pkg/frontmatter"
	"github.com/andy-neoaira/obs-cli/pkg/obsidian"
)

// FrontmatterParams 定义了 frontmatter 命令所需的业务参数。
type FrontmatterParams struct {
	NoteName string // 笔记名称或路径
	Print    bool   // 是否打印 frontmatter
	Edit     bool   // 是否编辑某个 key
	Delete   bool   // 是否删除某个 key
	Key      string // 要操作的 key
	Value    string // 编辑时的新值
}

// Frontmatter 是 "frontmatter" 命令的业务核心。
// 根据 Print/Edit/Delete 标志分别处理笔记的 YAML frontmatter。
func Frontmatter(vault obsidian.VaultManager, note obsidian.NoteManager, params FrontmatterParams) (string, error) {
	_, err := vault.DefaultName()
	if err != nil {
		return "", err
	}

	vaultPath, err := vault.Path()
	if err != nil {
		return "", err
	}

	// 读取笔记当前内容
	contents, err := note.GetContents(vaultPath, params.NoteName)
	if err != nil {
		return "", err
	}

	// 分发到具体的操作处理函数
	if params.Print {
		return handlePrint(contents)
	}

	if params.Edit {
		return handleEdit(note, vaultPath, params.NoteName, contents, params.Key, params.Value)
	}

	if params.Delete {
		return handleDelete(note, vaultPath, params.NoteName, contents, params.Key)
	}

	return "", errors.New("no operation specified: use --print, --edit, or --delete")
}

// handlePrint 处理 --print 操作：解析并格式化输出 frontmatter。
func handlePrint(contents string) (string, error) {
	if !frontmatter.HasFrontmatter(contents) {
		return "", nil // 没有 frontmatter 的笔记返回空字符串
	}

	fm, _, err := frontmatter.Parse(contents)
	if err != nil {
		return "", err
	}

	formatted, err := frontmatter.Format(fm)
	if err != nil {
		return "", err
	}

	return formatted, nil
}

// handleEdit 处理 --edit 操作：修改或新增指定 key，然后写回文件。
func handleEdit(note obsidian.NoteManager, vaultPath, noteName, contents, key, value string) (string, error) {
	if key == "" {
		return "", errors.New("--key is required for edit operation")
	}
	if value == "" {
		return "", errors.New("--value is required for edit operation")
	}

	updatedContent, err := frontmatter.SetKey(contents, key, value)
	if err != nil {
		return "", err
	}

	err = note.SetContents(vaultPath, noteName, updatedContent)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Updated frontmatter key '%s' in %s", key, noteName), nil
}

// handleDelete 处理 --delete 操作：删除指定 key，然后写回文件。
func handleDelete(note obsidian.NoteManager, vaultPath, noteName, contents, key string) (string, error) {
	if key == "" {
		return "", errors.New("--key is required for delete operation")
	}

	updatedContent, err := frontmatter.DeleteKey(contents, key)
	if err != nil {
		return "", err
	}

	err = note.SetContents(vaultPath, noteName, updatedContent)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Deleted frontmatter key '%s' from %s", key, noteName), nil
}
