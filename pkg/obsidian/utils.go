package obsidian

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// AddMdSuffix 为没有 .md 后缀的字符串追加 .md。
// 这是为了保证用户输入的笔记名最终都对应到 Markdown 文件。
func AddMdSuffix(str string) string {
	if !strings.HasSuffix(str, ".md") {
		return str + ".md"
	}
	return str
}

// RemoveMdSuffix 移除字符串末尾的 .md 后缀（如果存在）。
func RemoveMdSuffix(str string) string {
	if strings.HasSuffix(str, ".md") {
		return strings.TrimSuffix(str, ".md")
	}
	return str
}

// normalizePathSeparators 将反斜杠转换为正斜杠，保证跨平台一致性。
// Obsidian 在所有操作系统中都使用正斜杠作为链接分隔符。
func normalizePathSeparators(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

// wikiLinkPatterns 返回一个笔记名称对应的三种 wikilink 模式：
// [[name]]、[[name|、[[name#
func wikiLinkPatterns(name string) [3]string {
	return [3]string{
		"[[" + name + "]]",
		"[[" + name + "|",
		"[[" + name + "#",
	}
}

// GenerateNoteLinkTexts 生成指向某篇笔记的 wikilink 模式数组。
// 它会先取笔记名的 basename 并去掉 .md 后缀。
func GenerateNoteLinkTexts(noteName string) [3]string {
	noteName = filepath.Base(noteName)
	noteName = RemoveMdSuffix(noteName)
	return wikiLinkPatterns(noteName)
}

// GenerateBacklinkSearchPatterns 创建用于查找反向链接的搜索模式列表。
// 包括基于 basename 的 wikilink、基于完整路径的 wikilink，以及多种 Markdown 链接格式。
func GenerateBacklinkSearchPatterns(notePath string) []string {
	normalized := normalizePathSeparators(notePath)
	pathNoExt := RemoveMdSuffix(normalized)
	baseName := RemoveMdSuffix(path.Base(normalized))

	// 1. 基于 basename 的 wikilink（最常用）
	basePatterns := wikiLinkPatterns(baseName)
	patterns := append([]string{}, basePatterns[:]...)

	// 2. 基于完整路径的 wikilink（仅在路径与 basename 不同时）
	if pathNoExt != baseName {
		pathPatterns := wikiLinkPatterns(pathNoExt)
		patterns = append(patterns, pathPatterns[:]...)
	}

	// 3. Markdown 标准链接和相对链接
	mdPath := AddMdSuffix(normalized)
	patterns = append(patterns,
		"]("+mdPath+")",
		"]("+pathNoExt+")",
		"](./"+mdPath+")",
		"](./"+pathNoExt+")",
	)

	return patterns
}

// GenerateLinkReplacements 创建移动笔记时需要替换的链接映射表。
// 默认包含 basename wikilink 替换，保持旧版本行为。
func GenerateLinkReplacements(oldNotePath, newNotePath string) map[string]string {
	return GenerateLinkReplacementsWithOptions(oldNotePath, newNotePath, true)
}

// GenerateLinkReplacementsWithOptions 创建移动笔记时需要替换的链接映射表。
// 它会处理各种 Obsidian 链接格式：简单 wikilink、路径 wikilink、Markdown 链接。
// 所有路径都会被归一化为正斜杠，以保证跨平台一致。
//
// includeBaseLinks 用于控制是否替换 [[basename]] 这类不带目录的链接。
// 当 vault 中存在多个同名笔记时，basename 链接无法唯一确定目标，调用方应传 false，
// 只更新路径明确的链接，避免把无关笔记的引用误改掉。
func GenerateLinkReplacementsWithOptions(oldNotePath, newNotePath string, includeBaseLinks bool) map[string]string {
	replacements := make(map[string]string)

	// 将路径归一化为正斜杠，确保匹配一致
	oldNormalized := normalizePathSeparators(oldNotePath)
	newNormalized := normalizePathSeparators(newNotePath)

	// 取 basename（不含扩展名）
	oldBase := RemoveMdSuffix(path.Base(oldNormalized))
	newBase := RemoveMdSuffix(path.Base(newNormalized))

	// 取完整路径（不含扩展名）
	oldPathNoExt := RemoveMdSuffix(oldNormalized)
	newPathNoExt := RemoveMdSuffix(newNormalized)

	// 1. 简单 wikilink（仅 basename）——仅在调用方确认无歧义时替换
	if includeBaseLinks {
		replacements["[["+oldBase+"]]"] = "[[" + newBase + "]]"
		replacements["[["+oldBase+"|"] = "[[" + newBase + "|"
		replacements["[["+oldBase+"#"] = "[[" + newBase + "#"
	}

	// 2. 基于路径的 wikilink（路径与 basename 不同时）
	if oldPathNoExt != oldBase {
		replacements["[["+oldPathNoExt+"]]"] = "[[" + newPathNoExt + "]]"
		replacements["[["+oldPathNoExt+"|"] = "[[" + newPathNoExt + "|"
		replacements["[["+oldPathNoExt+"#"] = "[[" + newPathNoExt + "#"
	}

	// 3. Markdown 链接（多种格式）
	oldMd := AddMdSuffix(oldNormalized)
	newMd := AddMdSuffix(newNormalized)

	// 标准 Markdown 链接: [text](folder/note.md)
	replacements["]("+oldMd+")"] = "](" + newMd + ")"
	replacements["]("+oldPathNoExt+")"] = "](" + newPathNoExt + ")"

	// 相对 Markdown 链接: [text](./folder/note.md)
	replacements["](./"+oldMd+")"] = "](./" + newMd + ")"
	replacements["](./"+oldPathNoExt+")"] = "](./" + newPathNoExt + ")"

	return replacements
}

// ReplaceContent 批量替换 content 中的字符串，使用 replacements map 中的键值对。
func ReplaceContent(content []byte, replacements map[string]string) []byte {
	for o, n := range replacements {
		content = bytes.ReplaceAll(content, []byte(o), []byte(n))
	}
	return content
}

// ReplaceContentSkippingFencedCode 批量替换 Markdown 内容，但跳过 ``` 或 ~~~ 包裹的代码块。
// 这不能替代完整 Markdown parser，但能避免 move 命令修改代码示例中的链接文本，
// 对当前 CLI 的依赖体积和行为稳定性是更务实的折中。
func ReplaceContentSkippingFencedCode(content []byte, replacements map[string]string) []byte {
	lines := bytes.SplitAfter(content, []byte("\n"))
	inFence := false
	for i, line := range lines {
		trimmed := bytes.TrimSpace(line)
		if bytes.HasPrefix(trimmed, []byte("```")) || bytes.HasPrefix(trimmed, []byte("~~~")) {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		lines[i] = ReplaceContent(line, replacements)
	}
	return bytes.Join(lines, nil)
}

// IsExcluded 判断 relPath（相对于 vault 根目录的正斜杠路径）是否匹配任何排除规则。
// 支持的规则类型：
//   - 纯路径前缀："Archive"、"Templates/"
//   - 通配符："*.pdf"
//   - 双星号：**/drafts（匹配任意深度）
func IsExcluded(relPath string, filters []string) bool {
	normalized := filepath.ToSlash(relPath)
	for _, filter := range filters {
		if matchFilter(normalized, filter) {
			return true
		}
	}
	return false
}

// matchFilter 是 IsExcluded 的底层匹配逻辑。
func matchFilter(normalizedPath, filter string) bool {
	filter = strings.TrimRight(filter, "/")

	// 纯路径：精确匹配或前缀匹配（以 "/" 结尾表示目录前缀）
	if !strings.ContainsAny(filter, "*?[") {
		return normalizedPath == filter || strings.HasPrefix(normalizedPath, filter+"/")
	}

	// "**/" 前缀：匹配剩余部分在任意深度出现
	if strings.HasPrefix(filter, "**/") {
		return matchPathOrSegments(normalizedPath, filter[3:])
	}

	// 简单 glob：同时匹配完整路径和每个路径段
	return matchPathOrSegments(normalizedPath, filter)
}

// matchPathOrSegments 先用 filepath.Match 匹配完整路径，如果不匹配则逐个匹配路径段。
// 这样 "*.pdf" 也能匹配 "sub/file.pdf" 中的 file.pdf 段。
func matchPathOrSegments(path, pattern string) bool {
	if matched, _ := filepath.Match(pattern, path); matched {
		return true
	}
	for _, segment := range strings.Split(path, "/") {
		if matched, _ := filepath.Match(pattern, segment); matched {
			return true
		}
	}
	return false
}

// ShouldSkipDirectoryOrFile 判断文件/目录是否应该被跳过（不处理）。
// 跳过条件：是目录、是隐藏文件（以 . 开头）、不是 .md 文件。
func ShouldSkipDirectoryOrFile(info os.FileInfo) bool {
	isDirectory := info.IsDir()
	isHidden := info.Name()[0] == '.'
	isNonMarkdownFile := filepath.Ext(info.Name()) != ".md"
	if isDirectory || isHidden || isNonMarkdownFile {
		return true
	}
	return false
}

// OpenInEditor 使用用户偏好的编辑器打开指定文件。
// 支持常见 GUI 编辑器（VS Code、Sublime、Atom 等）的自动 --wait 等待标志，
// 也支持 EEDITOR 环境变量中包含参数的情况（如 "code -w"）。
func OpenInEditor(filePath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim" // 如果未设置 EDITOR，默认回退到 vim
	}

	// 将 EDITOR 拆分为命令和参数（处理 "code -w" 这类带参数的情况）
	parts := strings.Fields(editor)
	editorBin := parts[0]
	userArgs := parts[1:]

	// 构建参数列表：用户参数 + 自动等待标志 + 文件路径
	var args []string
	args = append(args, userArgs...)

	// 检测常见 GUI 编辑器，自动添加 --wait 防止命令立即返回
	editorLower := strings.ToLower(filepath.Base(editorBin))
	needsWait := strings.Contains(editorLower, "code") ||
		strings.Contains(editorLower, "vscode") ||
		strings.Contains(editorLower, "subl") ||
		strings.Contains(editorLower, "atom") ||
		strings.Contains(editorLower, "mate")

	// 如果用户已经在参数中传了 --wait，则不再重复添加
	if needsWait && !containsWaitFlag(userArgs) {
		args = append(args, "--wait")
	}

	args = append(args, filePath)

	cmd := exec.Command(editorBin, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open file in editor '%s': %w", editor, err)
	}

	return nil
}

// containsWaitFlag 检查参数列表中是否已经包含等待类 flag。
func containsWaitFlag(args []string) bool {
	for _, a := range args {
		if a == "--wait" || a == "-w" {
			return true
		}
	}
	return false
}
