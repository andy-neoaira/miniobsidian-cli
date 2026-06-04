package actions

import (
	"fmt"
	"strings"

	"github.com/andy-neoaira/miniobsidian-cli/pkg/obsidian"
)

// PrintParams 定义了 print 命令所需的业务参数。
type PrintParams struct {
	NoteName        string // 笔记名称或路径
	IncludeMentions bool   // 是否在末尾附加显示反向链接
}

// PrintNote 是 "print" 命令的业务核心。
// 读取笔记内容，如果指定了 --mentions，还会查找并格式化所有反向链接。
func PrintNote(vault obsidian.VaultManager, note obsidian.NoteManager, params PrintParams) (string, error) {
	_, err := vault.DefaultName()
	if err != nil {
		return "", err
	}

	vaultPath, err := vault.Path()
	if err != nil {
		return "", err
	}

	// 获取笔记的完整文本内容
	contents, err := note.GetContents(vaultPath, params.NoteName)
	if err != nil {
		return "", err
	}

	// 如果用户要求显示反向链接，查找并追加到内容末尾
	if params.IncludeMentions {
		backlinks, err := note.FindBacklinks(vaultPath, params.NoteName)
		if err != nil {
			return "", err
		}

		if len(backlinks) > 0 {
			contents += formatMentions(backlinks)
		}
	}

	return contents, nil
}

// formatMentions 将反向链接列表格式化为 Markdown 文本，按来源笔记分组展示。
func formatMentions(backlinks []obsidian.NoteMatch) string {
	var sb strings.Builder
	sb.WriteString("\n\n## Linked Mentions\n")

	// 按文件路径分组，同时保持首次出现的顺序
	grouped := make(map[string][]obsidian.NoteMatch)
	var order []string

	for _, match := range backlinks {
		if _, exists := grouped[match.FilePath]; !exists {
			order = append(order, match.FilePath)
		}
		grouped[match.FilePath] = append(grouped[match.FilePath], match)
	}

	// 逐组输出：先显示笔记名，再列出该笔记中所有引用位置
	for _, filePath := range order {
		noteName := obsidian.RemoveMdSuffix(filePath)
		fmt.Fprintf(&sb, "\n**[[%s]]**\n", noteName)
		for _, match := range grouped[filePath] {
			sb.WriteString("- ")
			sb.WriteString(match.MatchLine)
			sb.WriteByte('\n')
		}
	}

	return sb.String()
}
