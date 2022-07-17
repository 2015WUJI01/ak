package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var moduleCmd = &cobra.Command{
	Use:     "module",
	Aliases: []string{"mz"},
	Short:   "查询干员模组",
	Long:    "查询干员模组",
	Example: `    ak module 令 // 查询令有哪些模组
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO")
	},
}

func init() {
	rootCmd.AddCommand(moduleCmd)
}
