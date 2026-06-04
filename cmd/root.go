package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd 是 obs-cli 的根命令定义。
// Cobra 框架会根据这里的定义自动生成帮助信息、版本号和子命令树。
var rootCmd = &cobra.Command{
	Use:     "obs-cli",
	Short:   "Interact with Obsidian vaults from the terminal",
	Version: "v0.3.6",
	Long:    "Interact with Obsidian vaults from the terminal",
}

// Execute 是 CLI 的入口函数，由 main.go 调用。
// 它会解析命令行参数并执行对应的子命令；如果出错则打印错误信息并退出程序。
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
