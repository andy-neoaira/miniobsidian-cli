package mocks

// MockFuzzyFinder 是 FuzzyFinderManager 接口的 Mock 实现，用于单元测试。
// 可以控制返回的选中索引或模拟错误。
type MockFuzzyFinder struct {
	SelectedIndex int   // Find() 返回的固定索引
	FindErr       error // 模拟 Find() 返回的错误
}

func (f *MockFuzzyFinder) Find(slice interface{}, itemFunc func(i int) string, opts ...interface{}) (int, error) {
	if f.FindErr != nil {
		return -1, f.FindErr
	}
	return f.SelectedIndex, nil
}
