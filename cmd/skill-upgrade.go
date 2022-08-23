package cmd

import (
	"ak/database"
	"ak/models"
	"fmt"
	"github.com/spf13/cobra"
	"sort"
	"strings"
)

var (
	squeeze   bool
	rankRange string
)

var skillUpgradeCmd = &cobra.Command{
	Use:     "skill-upgrade",
	Aliases: []string{"skup"},
	Short:   "查询技能升级材料",
	Long:    "查询技能升级材料，也可以查专精材料",
	Example: `ak skill-upgrade -o 1 -r 3 令
--order 表示第几个技能
--rank 表示升到几级
--squeeze 合并显示
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 查询干员
		var opr models.Operator
		database.DB.Where("name", args[0]).First(&opr)
		if opr.Name == "" {
			fmt.Println("查无此干员")
			return
		}

		// 解析 rankRange 参数
		rankRange = strings.ReplaceAll(rankRange, "m1", "a")
		rankRange = strings.ReplaceAll(rankRange, "m2", "b")
		rankRange = strings.ReplaceAll(rankRange, "m3", "c")
		ints, err := parseRank(strings.Split(rankRange, ""))
		if err != nil {
			return
		}
		switch len(ints) {
		case 0:
			fmt.Println("rank 参数不能为空") //
			return
		case 1:
			// 只查一个等级
		case 2:
			// 查询等级范围
			sort.Ints(ints)
			for i := ints[0] + 1; i < ints[1]; i++ {
				ints = append(ints, i)
			}
			sort.Ints(ints)
		default:
			fmt.Println("rank 参数过多") //
			return
		}

		// 查询专精时，要指定哪个技能
		if ints[len(ints)-1] > 7 && order == 0 {
			fmt.Println("当查询专精技能材料时需要同时指定哪个技能，例如使用 --order 1 指定查询一技能")
			return
		}

		switch order {
		case 0:
			// 查询所有技能的 7 级以下的升级材料，且不查询专精材料
			var skillLevelMaterials []models.SkillLevelMaterial
			database.DB.Where("opr_name", opr.Name).
				Where("order", 0).
				Where("to_level", ints).
				Order("`order`, `to_level`").
				Find(&skillLevelMaterials)
			toLevel := 0
			for _, m := range skillLevelMaterials {
				if toLevel != m.ToLevel {
					if toLevel != 0 {
						fmt.Println()
					}
					fmt.Printf("Lv%v → Lv%v: ", m.ToLevel-1, m.ToLevel)
					toLevel = m.ToLevel
				} else {
					fmt.Print(" & ")
				}
				fmt.Printf("%v ×%v", m.Amount, "「"+m.ItemName+"」")
			}
			fmt.Println()
			return
		case 1, 2, 3:
			// 可能要查询专精材料
			var skillLevelMaterials []models.SkillLevelMaterial
			database.DB.Where("opr_name", opr.Name).
				Where("order", []int{0, order}).
				Where("to_level", ints).
				Order("`order`, `to_level`").
				Find(&skillLevelMaterials)
			toLevel := 0
			for _, m := range skillLevelMaterials {
				if toLevel != m.ToLevel {
					if toLevel != 0 {
						fmt.Println()
					}
					fmt.Printf("Lv%v → Lv%v: ", m.ToLevel-1, m.ToLevel)
					toLevel = m.ToLevel
				} else {
					fmt.Print(" & ")
				}
				fmt.Printf("%v ×%v", m.Amount, "「"+m.ItemName+"」")
			}
			fmt.Println()
			return
		default:
			fmt.Println("order 参数不符合要求") //
		}
	},
}

func init() {
	skillUpgradeCmd.Flags().BoolVarP(&squeeze, "squeeze", "s", false, "合并显示（暂时没做）")
	skillUpgradeCmd.Flags().IntVarP(&order, "order", "o", 0, "指定干员的第 n 个技能，默认为所有技能")
	skillUpgradeCmd.Flags().StringVarP(&rankRange, "rank", "r", "17", "技能等级或等级范围")
	rootCmd.AddCommand(skillUpgradeCmd)
}
