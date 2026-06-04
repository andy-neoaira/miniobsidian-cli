package mocks

// MockUriManager 是 UriManager 接口的 Mock 实现，用于单元测试。
// 可以记录 Construct 的调用参数，并控制 Execute 的返回值。
type MockUriManager struct {
	ConstructedURI string            // Construct() 返回的固定 URI
	LastBase       string            // 记录上次 Construct() 的 base 参数
	LastParams     map[string]string // 记录上次 Construct() 的 params 参数
	ExecuteErr     error             // 模拟 Execute() 返回的错误
	ExecuteCalls   int               // 记录 Execute() 被调用的次数
}

func (m *MockUriManager) Construct(base string, params map[string]string) string {
	m.LastBase = base
	m.LastParams = params
	return m.ConstructedURI
}

func (m *MockUriManager) Execute(uri string) error {
	m.ExecuteCalls++
	return m.ExecuteErr
}
