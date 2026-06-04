package actions

import (
	"fmt"
	"github.com/andy-neoaira/miniobsidian-cli/pkg/obsidian"
	"path/filepath"
)

// SearchNotes 是 "search" 命令的业务核心。
// 流程：获取 vault 下所有笔记 → 使用模糊搜索让用户交互式选择 → 打开选中的笔记。
func SearchNotes(vault obsidian.VaultManager, note obsidian.NoteManager, uri obsidian.UriManager, fuzzyFinder obsidian.FuzzyFinderManager, useEditor bool) error {
	vaultName, err := vault.DefaultName()
	if err != nil {
		return err
	}

	vaultPath, err := vault.Path()
	if err != nil {
		return err
	}

	// 获取 vault 中所有 .md 笔记的相对路径列表
	notes, err := note.GetNotesList(vaultPath)
	if err != nil {
		return err
	}

	// 启动终端模糊搜索界面，让用户实时过滤并选择笔记
	index, err := fuzzyFinder.Find(notes, func(i int) string {
		return notes[i]
	})

	if err != nil {
		return err
	}

	// 根据用户选择用编辑器或 Obsidian 打开
	if useEditor {
		fmt.Printf("Opening note: %s\n", notes[index])
		filePath := filepath.Join(vaultPath, notes[index])
		return obsidian.OpenInEditor(filePath)
	}

	obsidianUri := uri.Construct(ObsOpenUrl, map[string]string{
		"file":  notes[index],
		"vault": vaultName,
	})

	err = uri.Execute(obsidianUri)
	if err != nil {
		return err
	}

	return nil
}
