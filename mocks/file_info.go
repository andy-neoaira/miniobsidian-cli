package mocks

import (
	"os"
	"time"
)

// MockFileInfo 是 os.FileInfo 接口的 Mock 实现，用于单元测试。
// 可以控制文件名、是否是目录等属性。
type MockFileInfo struct {
	FileName    string // 文件名
	IsDirectory bool   // 是否是目录
}

func (fi *MockFileInfo) Name() string {
	return fi.FileName
}

func (fi *MockFileInfo) Size() int64 {
	return 0
}

func (fi *MockFileInfo) Mode() os.FileMode {
	return 0
}

func (fi *MockFileInfo) ModTime() time.Time {
	return time.Now()
}

func (fi *MockFileInfo) IsDir() bool {
	return fi.IsDirectory
}

func (fi *MockFileInfo) Sys() interface{} {
	return nil
}
