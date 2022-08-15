package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	repo "main/repositories"
)

var itemCmd = &cobra.Command{
	Use:     "item",
	Short:   "查询道具信息",
	Long:    "查询道具信息",
	Example: `    ak item 龙门币`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if item, ok := repo.FindItemByName(args[0]); ok {
			fmt.Println(item.Info(repo.FindItemAlias(item.Name)))
			return
		} else if item, ok = repo.FindItemByAlias(args[0]); ok {
			fmt.Println(item.Info(repo.FindItemAlias(item.Name)))
			return
		}
		fmt.Println("查询无果")
	},
}

func init() {
	rootCmd.AddCommand(itemCmd)
}
