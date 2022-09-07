package cmd

import (
	"ak/cmd/update"
	"ak/database"
	"ak/models"
	"ak/services"
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

			var items []models.Item
			if _, ok := m["1"]; ok {
				services.FetchStep1(&items)
			} else {
				database.DB.Where(&models.Item{}).Find(&items)
			}
			if _, ok := m["2"]; ok {
				services.FetchStep2(items)
			}
			if _, ok := m["3"]; ok {
				services.FetchStep3()
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
	rootCmd.AddCommand(updateCmd)
}
