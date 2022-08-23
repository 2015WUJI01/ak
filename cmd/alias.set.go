package cmd

import (
	"ak/database"
	"ak/models"
	"fmt"
	"github.com/spf13/cobra"
	"gorm.io/gorm/clause"
)

const (
	_             = iota
	AliasTypeItem // 道具作为第一个需要抓取的信息，所以放在第一个
	AliasTypeOpr
)

var aSetCmd = &cobra.Command{
	Use:     "set",
	Aliases: []string{"new", "add"},
	Short:   "新增一个别名信息",
	Long:    `新增一个别名信息，先写干员原名，再写别名，可以一次设置多个别名`,
	Example: `  ak alias set <opr_name> <opr_alias>...
  
  ak alias set 史尔特尔 42                 // 设置干员史尔特尔的别名为 42
  ak alias set --item 龙门币 钱            // 默认是设置干员的别名，即 --opr 参数，若设置道具别名需要加上 --item 或 -i 参数
  ak alias set 浊心斯卡蒂 蒂蒂 红蒂 浊蒂   // 可以一次为一位干员设置多个别名
`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var cnt int64
		database.DB.Model(&models.Operator{}).Where("name", args[0]).Count(&cnt)
		if cnt == 0 {
			fmt.Println("干员名称输入错误")
			return
		}

		// 别名类型
		atype := AliasTypeOpr
		if isItem {
			atype = AliasTypeItem
		}

		// 添加别名
		var aliases []models.Alias
		for _, a := range args[1:] {
			aliases = append(aliases, models.Alias{
				Name:  args[0],
				Alias: a,
				Type:  atype,
			})
		}
		database.DB.Select([]string{"name", "type", "alias"}).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}, {Name: "type"}, {Name: "alias"}},
			DoNothing: true,
		}).Create(&aliases)
		fmt.Printf("新增干员 %s 别名 %v 成功\n", args[0], args[1:])
	}}

func init() {
	aliasCmd.AddCommand(aSetCmd)
	aSetCmd.MarkFlagsMutuallyExclusive("opr", "item")
}
