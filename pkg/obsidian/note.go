package obsidian

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Note 是 NoteManager 的具体实现，提供笔记的增删改查和链接管理能力。
// 使用空结构体是因为方法不依赖内部状态，仅作为方法接收者的载体。
type Note struct{}

// isHiddenDir 判断目录条目是否是隐藏目录（以 . 开头但不是当前目录 .）。
func isHiddenDir(d fs.DirEntry) bool {
	return d.IsDir() && d.Name() != "." && strings.HasPrefix(d.Name(), ".")
}

// NoteMatch 保存单次搜索结果的信息。
type NoteMatch struct {
	FilePath   string // 匹配所在的笔记相对路径
	LineNumber int    // 匹配行号（0 表示文件名匹配）
	MatchLine  string // 匹配行的内容摘要
}

// NoteManager 定义了与笔记交互的接口。
// 通过接口解耦，方便在测试中注入 Mock 对象。
type NoteManager interface {
	Move(string, string) error
	Delete(string) error
	UpdateLinks(string, string, string) error
	GetContents(string, string) (string, error)
	SetContents(string, string, string) error
	GetNotesList(string) ([]string, error)
	SearchNotesWithSnippets(string, string) ([]NoteMatch, error)
	FindBacklinks(string, string) ([]NoteMatch, error)
}

// Move 将笔记从原路径移动到新路径，会自动补充 .md 后缀。
func (m *Note) Move(originalPath string, newPath string) error {
	o := AddMdSuffix(originalPath)
	n := AddMdSuffix(newPath)

	err := os.Rename(o, n)

	if err != nil {
		return errors.New(NoteDoesNotExistError)
	}

	message := fmt.Sprintf(`Moved note
from %s
to %s`, o, n)

	fmt.Println(message)
	return nil
}

// Delete 删除指定路径的笔记，会自动补充 .md 后缀。
func (m *Note) Delete(path string) error {
	note := AddMdSuffix(path)
	err := os.Remove(note)
	if err != nil {
		return errors.New(NoteDoesNotExistError)
	}
	fmt.Println("Deleted note: ", note)
	return nil
}

// GetContents 读取 vault 中指定笔记的完整文本内容。
// 搜索策略：先尝试完整相对路径匹配，再回退到 basename 匹配（向后兼容）。
func (m *Note) GetContents(vaultPath string, noteName string) (string, error) {
	note := AddMdSuffix(noteName)

	var notePath string
	// 遍历 vault 目录树查找目标笔记
	err := filepath.WalkDir(vaultPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err // 遇到错误继续下一个路径
		}
		if d.IsDir() {
			return nil // 跳过目录
		}

		// 优先完整相对路径匹配
		relPath, err := filepath.Rel(vaultPath, path)
		if err == nil && relPath == note {
			notePath = path
			return filepath.SkipDir // 找到后停止遍历
		}

		// 回退到 basename 匹配（向后兼容旧版本行为）
		if filepath.Base(path) == note {
			notePath = path
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil || notePath == "" {
		return "", errors.New(NoteDoesNotExistError)
	}

	file, err := os.Open(notePath)
	if err != nil {
		return "", errors.New(VaultReadError)
	}
	defer file.Close() // 确保文件句柄在函数返回前关闭

	content, err := io.ReadAll(file)
	if err != nil {
		return "", errors.New(VaultReadError)
	}

	return string(content), nil
}

// SetContents 将内容写入 vault 中的指定笔记。
// 先通过遍历找到笔记的实际路径，然后覆盖写入。
func (m *Note) SetContents(vaultPath string, noteName string, content string) error {
	note := AddMdSuffix(noteName)

	var notePath string
	err := filepath.WalkDir(vaultPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// 优先完整相对路径匹配
		relPath, err := filepath.Rel(vaultPath, path)
		if err == nil && relPath == note {
			notePath = path
			return filepath.SkipDir
		}

		// 回退到 basename 匹配
		if filepath.Base(path) == note {
			notePath = path
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil || notePath == "" {
		return errors.New(NoteDoesNotExistError)
	}

	err = os.WriteFile(notePath, []byte(content), 0644)
	if err != nil {
		return errors.New(VaultWriteError)
	}

	return nil
}

// UpdateLinks 遍历 vault 中所有笔记，将指向旧笔记的链接更新为新笔记的链接。
// 这是 move 命令的重要后续操作，保证笔记间的引用关系不会断裂。
func (m *Note) UpdateLinks(vaultPath string, oldNoteName string, newNoteName string) error {
	err := filepath.Walk(vaultPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.New(VaultAccessError)
		}

		if ShouldSkipDirectoryOrFile(info) {
			return nil
		}

		originalContent, err := os.ReadFile(path)
		if err != nil {
			return errors.New(VaultReadError)
		}

		// 生成所有需要替换的链接模式，并执行替换
		replacements := GenerateLinkReplacements(oldNoteName, newNoteName)
		updatedContent := ReplaceContent(originalContent, replacements)

		// 如果没有实际变化，跳过写入以提高性能
		if bytes.Equal(originalContent, updatedContent) {
			return nil
		}

		err = os.WriteFile(path, updatedContent, info.Mode())
		if err != nil {
			return errors.New(VaultWriteError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// GetNotesList 获取 vault 中所有 .md 笔记的相对路径列表。
// 会自动跳过隐藏目录和用户配置的排除路径。
func (m *Note) GetNotesList(vaultPath string) ([]string, error) {
	excluded := ExcludedPaths(vaultPath)
	var notes []string
	err := filepath.WalkDir(vaultPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if isHiddenDir(d) {
			return filepath.SkipDir // 跳过隐藏目录（如 .git、.obsidian）
		}
		relPath, err := filepath.Rel(vaultPath, path)
		if err != nil {
			return err
		}
		if relPath != "." && IsExcluded(relPath, excluded) {
			if d.IsDir() {
				return filepath.SkipDir // 被排除的目录直接跳过
			}
			return nil
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
			notes = append(notes, relPath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return notes, nil
}

// SearchNotesWithSnippets 在 vault 中搜索包含 query 的笔记。
// 同时搜索文件名和文件内容，返回匹配片段列表。
// 为避免内存问题，大于 10MB 的文件会被跳过。
func (m *Note) SearchNotesWithSnippets(vaultPath string, query string) ([]NoteMatch, error) {
	excluded := ExcludedPaths(vaultPath)
	var matches []NoteMatch
	queryLower := strings.ToLower(query)

	err := filepath.WalkDir(vaultPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if isHiddenDir(d) {
			return filepath.SkipDir
		}
		relPath, relErr := filepath.Rel(vaultPath, path)
		if relErr != nil {
			return relErr
		}
		if relPath != "." && IsExcluded(relPath, excluded) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
			fileNameMatches := strings.Contains(strings.ToLower(relPath), queryLower)
			var hasContentMatch bool

			// 检查文件大小，避免读取超大文件（>10MB）导致内存问题
			if info, err := d.Info(); err == nil && info.Size() < 10*1024*1024 {
				content, err := os.ReadFile(path)
				if err == nil {
					lines := strings.Split(string(content), "\n")
					for lineNum, line := range lines {
						if strings.Contains(strings.ToLower(line), queryLower) {
							hasContentMatch = true
							matchLine := strings.TrimSpace(line)
							if len(matchLine) > 80 {
								// 如果匹配行太长，截取查询词前后各 20 个字符作为摘要
								queryPos := strings.Index(strings.ToLower(matchLine), queryLower)
								if queryPos != -1 {
									start := queryPos - 20
									end := queryPos + len(query) + 20
									if start < 0 {
										start = 0
									}
									if end > len(matchLine) {
										end = len(matchLine)
									}
									if start > 0 {
										matchLine = "..." + matchLine[start:]
									}
									if end < len(strings.TrimSpace(line)) {
										matchLine = matchLine[:end-start] + "..."
									}
								} else {
									matchLine = matchLine[:80] + "..."
								}
							}

							matches = append(matches, NoteMatch{
								FilePath:   relPath,
								LineNumber: lineNum + 1, // 行号从 1 开始
								MatchLine:  matchLine,
							})
						}
					}
				}
			}

			// 只有当内容没有匹配时，才将文件名匹配作为补充结果加入
			if fileNameMatches && !hasContentMatch {
				matches = append(matches, NoteMatch{
					FilePath:   relPath,
					LineNumber: 0, // 0 表示文件名匹配
					MatchLine:  fmt.Sprintf("(filename match: %s)", filepath.Base(relPath)),
				})
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return matches, nil
}

// maxFileSizeBytes 定义了反向链接查找时跳过的文件大小上限（10MB）。
const maxFileSizeBytes = 10 * 1024 * 1024 // 10MB

// containsAnyPattern 快速检查 contentLower 中是否包含任意一个 pattern（不区分大小写）。
// 用于在大文件扫描前先进行粗略过滤，避免逐行解析不必要的文件。
func containsAnyPattern(contentLower []byte, patterns [][]byte) bool {
	for _, pattern := range patterns {
		if bytes.Contains(contentLower, pattern) {
			return true
		}
	}
	return false
}

// findMatchingLines 逐行扫描 content，找出包含任意 pattern 的所有行。
func findMatchingLines(content []byte, patternsLower [][]byte) []NoteMatch {
	var matches []NoteMatch
	lineNum := 0

	for len(content) > 0 {
		lineNum++

		// 找到当前行的结束位置（换行符）
		idx := bytes.IndexByte(content, '\n')
		var line []byte
		if idx == -1 {
			line = content
			content = nil
		} else {
			line = content[:idx]
			content = content[idx+1:]
		}

		// 检查当前行是否匹配任一模式
		lineLower := bytes.ToLower(line)
		for _, pattern := range patternsLower {
			if bytes.Contains(lineLower, pattern) {
				matches = append(matches, NoteMatch{
					LineNumber: lineNum,
					MatchLine:  string(bytes.TrimSpace(line)),
				})
				break // 一行只需记录一次
			}
		}
	}

	return matches
}

// FindBacklinks 查找 vault 中所有引用指定笔记的反向链接。
// 结果按文件修改时间降序排列（最新的在前）。
func (m *Note) FindBacklinks(vaultPath, noteName string) ([]NoteMatch, error) {
	noteName = RemoveMdSuffix(noteName)
	excluded := ExcludedPaths(vaultPath)

	// 生成搜索模式并预先转换为小写字节切片，避免对每个文件重复转换
	patterns := GenerateBacklinkSearchPatterns(noteName)
	patternsLower := make([][]byte, len(patterns))
	for i, p := range patterns {
		patternsLower[i] = []byte(strings.ToLower(p))
	}

	var matches []NoteMatch
	fileModTimes := make(map[string]int64)

	err := filepath.WalkDir(vaultPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if isHiddenDir(d) {
			return filepath.SkipDir
		}

		relPath, err := filepath.Rel(vaultPath, path)
		if err != nil {
			return err
		}
		if relPath != "." && IsExcluded(relPath, excluded) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}

		// 跳过笔记自身（避免将自身的自引用算作反向链接）
		if RemoveMdSuffix(normalizePathSeparators(relPath)) == noteName {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil //nolint:nilerr
		}
		if info.Size() > maxFileSizeBytes {
			fmt.Fprintf(os.Stderr, "Skipping file %s: size %d bytes exceeds limit %d bytes\n", relPath, info.Size(), maxFileSizeBytes)
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil //nolint:nilerr
		}

		// 快速预检：如果文件内容中不包含任何模式，直接跳过
		contentLower := bytes.ToLower(content)
		if !containsAnyPattern(contentLower, patternsLower) {
			return nil
		}

		// 找到所有匹配行
		modTime := info.ModTime().UnixNano()
		fileMatches := findMatchingLines(content, patternsLower)
		for i := range fileMatches {
			fileMatches[i].FilePath = relPath
			fileModTimes[relPath] = modTime
		}
		matches = append(matches, fileMatches...)

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 按文件修改时间降序排序（最新修改的文件排在前面）
	sort.Slice(matches, func(i, j int) bool {
		return fileModTimes[matches[i].FilePath] > fileModTimes[matches[j].FilePath]
	})

	return matches, nil
}
