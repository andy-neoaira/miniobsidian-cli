package mocks

// MockLinkRewriter 是 LinkRewriteManager 接口的 Mock 实现。
// MoveNote 测试通过它单独模拟“链接重写成功/失败”，不再把该职责塞进 MockNoteManager。
type MockLinkRewriter struct {
	UpdateLinksError error
}

func (m *MockLinkRewriter) UpdateLinks(string, string, string) error {
	return m.UpdateLinksError
}
