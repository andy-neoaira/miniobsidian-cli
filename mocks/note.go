package mocks

import "github.com/andy-neoaira/miniobsidian-cli/pkg/obsidian"

// MockNoteManager 是 NoteManager 接口的 Mock 实现，用于单元测试。
// 通过设置各个字段的值，可以控制 Mock 对象的行为（返回指定值或模拟错误）。
type MockNoteManager struct {
	DeleteErr           error                // 模拟 Delete() 返回的错误
	MoveErr             error                // 模拟 Move() 返回的错误
	UpdateLinksError    error                // 模拟 UpdateLinks() 返回的错误
	GetContentsError    error                // 模拟 GetContents() 返回的错误
	SetContentsError    error                // 模拟 SetContents() 返回的错误
	FindBacklinksErr    error                // 模拟 FindBacklinks() 返回的错误
	FindBacklinksResult []obsidian.NoteMatch // FindBacklinks() 返回的自定义结果
	NoMatches           bool                 // 是否返回空结果
	Contents            string               // GetContents() 返回的自定义内容
}

func (m *MockNoteManager) Delete(string) error {
	return m.DeleteErr
}

func (m *MockNoteManager) Move(string, string) error {
	return m.MoveErr
}

func (m *MockNoteManager) UpdateLinks(string, string, string) error {
	return m.UpdateLinksError
}

func (m *MockNoteManager) GetContents(string, string) (string, error) {
	if m.Contents != "" {
		return m.Contents, m.GetContentsError
	}
	return "example contents", m.GetContentsError
}

func (m *MockNoteManager) SetContents(string, string, string) error {
	return m.SetContentsError
}

func (m *MockNoteManager) GetNotesList(string) ([]string, error) {
	return []string{"note1", "note2", "note3"}, m.GetContentsError
}

func (m *MockNoteManager) SearchNotesWithSnippets(string, string) ([]obsidian.NoteMatch, error) {
	if m.GetContentsError != nil {
		return nil, m.GetContentsError
	}
	if m.NoMatches {
		return []obsidian.NoteMatch{}, nil
	}
	return []obsidian.NoteMatch{
		{FilePath: "note1.md", LineNumber: 5, MatchLine: "example match line"},
		{FilePath: "note2.md", LineNumber: 10, MatchLine: "another match"},
	}, nil
}

func (m *MockNoteManager) FindBacklinks(string, string) ([]obsidian.NoteMatch, error) {
	if m.FindBacklinksErr != nil {
		return nil, m.FindBacklinksErr
	}
	if m.FindBacklinksResult != nil {
		return m.FindBacklinksResult, nil
	}
	if m.NoMatches {
		return []obsidian.NoteMatch{}, nil
	}
	return []obsidian.NoteMatch{
		{FilePath: "linking-note.md", LineNumber: 5, MatchLine: "This links to [[target]]"},
		{FilePath: "another-note.md", LineNumber: 10, MatchLine: "Also references [[target]]"},
	}, nil
}
