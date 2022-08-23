package cmd

import (
	"ak/cmd/update"
	"ak/models"
	repo "ak/repositories"
	"fmt"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "更新数据",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			update.Step1()
			update.Step2()
			update.Step3()
			update.Step4()
			update.Step5()
		} else {
			m := make(map[string]struct{})
			for _, arg := range args {
				m[arg] = struct{}{}
			}
			if _, ok := m["1"]; ok {
				update.Step1()
			}
			if _, ok := m["2"]; ok {
				update.Step2()
			}
			if _, ok := m["3"]; ok {
				update.Step3()
			}
			if _, ok := m["4"]; ok {
				update.Step4()
			}
			if _, ok := m["5"]; ok {
				update.Step5()
			}
		}
	},
}

func init() {
	var wikiFlag, groupFlag bool
	updateItemCmd := &cobra.Command{
		Use:     "item",
		Short:   "更新指定的 item",
		Example: "ak update item -wiki 龙门币",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, itemname := range args {
				item := models.Item{Name: itemname}
				if wikiFlag {
					models.FreshItemWiki(&item)
					repo.CreateOrUpdateItem(item, "name", "wiki")
					fmt.Printf("更新 %v wiki 成功：%v\n", item.Name, item.Wiki)
				}
				if groupFlag {
					item.FreshGroup()
					fmt.Printf("更新 %v group 成功：%v\n", item.Name, item.Wiki)
				}
			}
		},
	}
	updateItemCmd.Flags().BoolVarP(&wikiFlag, "wiki", "w", false, "更新 wiki 链接")
	updateItemCmd.Flags().BoolVarP(&groupFlag, "group", "g", false, "更新 wiki 链接")
	updateCmd.AddCommand(updateItemCmd)
	rootCmd.AddCommand(updateCmd)
}
