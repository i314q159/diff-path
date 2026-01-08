package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	verbose bool
)

var RootCmd = &cobra.Command{
	Use:   "diff-path",
	Short: "AOSP目录对比工具",
	Long:  "AOSP目录对比工具",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("请使用子命令，或使用 --help 查看帮助")
	},
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "详细输出")

	// 在这里添加子命令
	RootCmd.AddCommand(DiffCmd)
}

// Execute 执行根命令
func Execute() error {
	return RootCmd.Execute()
}
