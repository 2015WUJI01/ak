package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"main/database"
	"main/logger"
	"main/models"
	"strconv"
	"strings"
)

var (
	order int
	rank  []string
)

var skillCmd = &cobra.Command{
	Use:   "skill",
	Short: "查询干员技能信息",
	Long:  "查询干员技能信息",
	Example: `    ak skill XXX // 查询某干员的所有技能，参数为干员 ID 或 name 与 opr 相同，默认 -in
    ak skill --id 42 // 查询 ID 为 42 的干员的所有技能
    ak skill --name 史尔特尔 // 查询 史尔特尔 的所有技能
    ak skill --alias 42 // 查询 史尔特尔 的所有技能
    ak skill --order X XXX // 查询某干员的第 X 个技能，会列出该技能所有等级（R1-7、M1-3）的描述
    ak skill --order 3 史尔特尔 // 查询 史尔特尔 的第 3 个技能
    ak skill --rank X XXX // 查询干员指定技能等级的所有技能描述，未专精用数字 1-7 表示 1-7 级，专精用 m1-m3 表示，或用 a,b,c 表示，无视大小写
    ak skill --rank 7 令 // 查询令所有技能等级 7 级的技能描述
    ak skill --rank m1 令 // 查询令所有技能等级专精 1 级的技能描述
    ak skill --rank a 令 // 查询令所有技能等级专精 1 级的技能描述
    ak skill --order 3 --rank C 令 // 查询令三技能等级专精 3 级的所有技能描述`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		var opr models.Operator
		var key = args[0]

		// 优先查询干员 ID
		if id {
			database.DB.Where("id", key).First(&opr)
		}

		// 查询干员名称
		if name && opr == (models.Operator{}) {
			database.DB.Where("name", key).First(&opr)
		}

		// 再查询别名
		if alias {
			fmt.Printf("暂不支持别名查询\n")
			// database.DB.Where("id", key).Find(&oprs)
			// for _, opr := range oprs {
			// 	fmt.Printf("%+v\n", opr)
			// }
		}

		// 找不到干员
		if opr == (models.Operator{}) {
			fmt.Printf("查询干员无果\n")
			return
		}

		type SkInfo struct {
			skill  models.Skill
			levels []models.SkillLevel
		}

		// 校验 rank 参数
		var valids []int
		var err error
		if valids, err = parseRank(rank); err != nil {
			fmt.Println(err)
			return
		} else if len(valids) == 0 {
			valids = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		}

		// 要查几个技能
		// 查询多个技能
		if order == 0 {
			var skills []models.Skill
			database.DB.Where("opr_name", opr.Name).Find(&skills)
			var sks []SkInfo
			for _, skill := range skills {
				var sklvs []models.SkillLevel
				database.DB.Where("opr_name", opr.Name).
					Where("order", skill.Order).
					Where("level", valids).
					Find(&sklvs)
				sks = append(sks, SkInfo{
					skill:  skill,
					levels: sklvs,
				})
			}
			logger.Debug(sks)
		} else if order <= 3 {
			// 查询单个技能
			var skill models.Skill
			database.DB.Where("opr_name", opr.Name).First(&skill)

			var sk SkInfo
			var sklvs []models.SkillLevel
			database.DB.Where("opr_name", opr.Name).
				Where("order", skill.Order).
				Where("level", valids).
				Find(&sklvs)
			sk = SkInfo{
				skill:  skill,
				levels: sklvs,
			}

			logger.Debug(sk)
		}

	},
}

func init() {

	// 指定干员
	skillCmd.Flags().BoolVarP(&id, "id", "i", false, "使用「干员 ID」进行查询")
	skillCmd.Flags().BoolVarP(&name, "name", "n", true, "使用「干员名称」进行查询")
	skillCmd.Flags().BoolVarP(&alias, "alias", "a", false, "使用「干员别名」进行查询")

	// 指定技能
	skillCmd.Flags().IntVarP(&order, "order", "o", 0, "指定干员的第 n 个技能，默认为所有技能")
	skillCmd.Flags().StringArrayVarP(&rank, "rank", "r", []string{}, "指定干员技能的等级，默认为所有等级")

	rootCmd.AddCommand(skillCmd)
}

func parseRank(rank []string) ([]int, error) {
	var valid []int
	for _, s := range rank {
		switch strings.ToLower(s) {
		case "1", "2", "3", "4", "5", "6", "7", "8", "9", "10":
		case "m1", "a":
			s = "8"
		case "m2", "b":
			s = "9"
		case "m3", "c":
			s = "10"
		default:
			return nil, errors.New(fmt.Sprintf("请检查 --rank 参数是否符合要求，rank 必须为 1-7 或 m1-m3 或 a-c 中任意一种，而不是: %s", s))
		}
		i, _ := strconv.Atoi(s)
		valid = append(valid, i)
	}
	return valid, nil
}
