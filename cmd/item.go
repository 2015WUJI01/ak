package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"main/database"
	"main/models"
)

var itemCmd = &cobra.Command{
	Use:     "item",
	Short:   "查询道具信息",
	Long:    "查询道具信息",
	Example: `    ak item 龙门币`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var item models.Item
		database.DB.Where("name", args[0]).First(&item)
		fmt.Printf("%+v", item)
	},
}

func init() {
	rootCmd.AddCommand(itemCmd)
}
