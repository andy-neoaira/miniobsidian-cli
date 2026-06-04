package cmd

import (
	"fmt"
	"log"

	"github.com/andy-neoaira/miniobsidian-cli/pkg/obsidian"
	"github.com/spf13/cobra"
)

// printDefaultDeprecatedCmd 是已弃用的 "print-default" 命令。
// 建议用户使用 "list-vaults --default" 替代。
var printDefaultDeprecatedCmd = &cobra.Command{
	Use:        "print-default",
	Aliases:    []string{"pd"},
	Short:      "prints default vault name and path (deprecated: use list-vaults --default)",
	Args:       cobra.ExactArgs(0),
	Deprecated: "use list-vaults --default instead",
	Run: func(cmd *cobra.Command, args []string) {
		pathOnly, _ := cmd.Flags().GetBool("path-only")

		vault := obsidian.Vault{}
		name, err := vault.DefaultName()
		if err != nil {
			log.Fatal(err)
		}
		path, err := vault.Path()
		if err != nil {
			log.Fatal(err)
		}

		// --path-only 模式下只输出路径，方便脚本调用
		if pathOnly {
			fmt.Print(path)
			return
		}

		openType, _ := vault.DefaultOpenType()

		fmt.Println("Default vault name:", name)
		fmt.Println("Default vault path:", path)
		fmt.Println("Default open type:", openType)
	},
}

func init() {
	printDefaultDeprecatedCmd.Flags().Bool("path-only", false, "print only the vault path")
	rootCmd.AddCommand(printDefaultDeprecatedCmd)
}
