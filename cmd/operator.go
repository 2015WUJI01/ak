package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var oprCmd = &cobra.Command{
	Use:               "opr [flags] <args>",
	Short:             "查询单个干员信息",
	Long:              "查询与输入信息最匹配的干员信息\n默认 flags 使用 -in 参数，表示支持使用干员 ID 和干员名称进行查询",
	Example:           "  ak opr 令\n  ak opr --alias 斯卡蒂",
	ValidArgs:         nil,
	ValidArgsFunction: nil,
	Args:              cobra.MinimumNArgs(1), // 至少一个参数
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print("afasd")
		return nil
	},
	Hidden:                false,
	DisableAutoGenTag:     false,
	DisableFlagsInUseLine: true,
}

func init() {
	// rootCmd.AddCommand(oprCmd)
}
