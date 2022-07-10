package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"main/database"
	"main/models"
)

// opr 命令的参数
var (
	isOpr  bool
	isItem bool
)

var aliasCmd = &cobra.Command{
	Use:   `alias`,
	Short: "查询或操作干员的别名",
	Long:  "查询或操作干员的别名",
	Example: `
    ak alias XXX // 查询 XXX 这个别名指的是什么东西

    # 添加别名，需要有 类型（可选）、原名、别名三种
	ak alias --type-opr
    ak alias --type-item

    ak alias --set
    ak alias --set 斯卡蒂 42`,
	ValidArgsFunction:     nil,
	Args:                  cobra.MinimumNArgs(1), // 至少一个参数
	Hidden:                false,
	DisableAutoGenTag:     false,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			var aliases []models.Alias
			database.DB.Where("alias", args[0]).Limit(15).Find(&aliases)
			for _, a := range aliases {
				fmt.Println(a.Name)
			}

			if len(aliases) == 0 {
				fmt.Println("查询无果")
			}
		}
	},
}

var aGetCmd = &cobra.Command{Use: "get", Short: "查询一个别名信息",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

	}}

var aRmCmd = &cobra.Command{Use: "rm", Aliases: []string{"del"}, Short: "删除一个别名信息",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

	}}
var aExportCmd = &cobra.Command{Use: "export", Short: "导出所有的别名信息",
	Run: func(cmd *cobra.Command, args []string) {

	}}

func init() {
	aliasCmd.PersistentFlags().BoolVarP(&isOpr, "opr", "o", true, "操作干员的别名")
	aliasCmd.PersistentFlags().BoolVarP(&isItem, "item", "i", false, "操作道具的别名")

	aliasCmd.AddCommand(aGetCmd)
	aliasCmd.AddCommand(aRmCmd)
	aliasCmd.AddCommand(aExportCmd)
	rootCmd.AddCommand(aliasCmd)
}
