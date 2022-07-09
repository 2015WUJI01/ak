package cmd

import (
	"github.com/spf13/cobra"
	"main/cmd/update"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "更新数据",
	RunE: func(cmd *cobra.Command, args []string) error {
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
		return nil
	},
}
