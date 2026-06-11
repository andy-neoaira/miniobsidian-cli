package cmd

import (
	"fmt"
	"io"
	"os"
)

// resolveContentInput 统一处理命令中的文本输入来源。
//
// 约束：
//  1. --content 适合短文本和简单脚本。
//  2. --content-file 适合多行 Markdown 或包含引号的内容，避免 shell 拼接破坏文本。
//  3. --content-file - 从 stdin 读取，方便 agent 和管道安全传入长内容。
//  4. 两种来源同时出现时直接报错，避免用户误以为内容会自动合并。
func resolveContentInput(inlineContent, contentFile string) (string, error) {
	if inlineContent != "" && contentFile != "" {
		return "", fmt.Errorf("--content and --content-file cannot be used together")
	}
	if contentFile == "" {
		return inlineContent, nil
	}
	if contentFile == "-" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("failed to read content from stdin: %w", err)
		}
		return string(data), nil
	}
	data, err := os.ReadFile(contentFile)
	if err != nil {
		return "", fmt.Errorf("failed to read content file %q: %w", contentFile, err)
	}
	return string(data), nil
}
