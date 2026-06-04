package obsidian

import (
	"errors"
	"github.com/ktr0731/go-fuzzyfinder"
)

// FuzzyFinder 是 FuzzyFinderManager 的具体实现，封装了 go-fuzzyfinder 库。
type FuzzyFinder struct{}

// FuzzyFinderManager 定义了模糊搜索交互界面的接口。
type FuzzyFinderManager interface {
	Find(slice interface{}, itemFunc func(i int) string, opts ...interface{}) (int, error)
}

// Find 启动终端模糊搜索界面。
// 接收字符串切片和显示函数，返回用户选中项的索引。
// 如果用户取消选择或传入类型不正确，返回错误。
func (f *FuzzyFinder) Find(slice interface{}, itemFunc func(i int) string, opts ...interface{}) (int, error) {
	items, ok := slice.([]string)
	if !ok {
		return -1, errors.New("invalid slice type, expected []string")
	}

	index, err := fuzzyfinder.Find(items, func(i int) string {
		return itemFunc(i)
	})
	if err != nil {
		return -1, errors.New(NoteDoesNotExistError)
	}
	return index, nil
}
