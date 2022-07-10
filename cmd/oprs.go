package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"main/database"
	"main/models"
	"main/pkg/logger"
)

// opr 命令的参数 与 opr 相同

var oprsCmd = &cobra.Command{
	Use: `oprs [flags] <args>
  与 opr 命令相同，默认仅支持的「干员名称」查询，即 [flags] 默认仅使用 --name。当存在其他 [flags] 时，默认的 --name 同样生效，此时，多 [flags] 情况下，结果会按照「干员 ID」、「干员名称」、「干员别名」递减的优先级进行排序，与 opr 命令不同的是，opr 返回优先级最高的结果，而 oprs 返回所有排序后的结果。`,
	Short: "查询多个干员信息",
	Long:  "查询符合要求的所有干员的所有数据",
	Example: `  ak oprs 斯卡蒂         // 查询结果：斯卡蒂、浊心斯卡蒂
  ak oprs --alias 小车   // 查询结果：所有小车
  ak oprs --id --alias 42  // 查询结果：「干员 ID」为 123 的干员，以及「干员别名」为 42 的史尔特尔
  ak oprs -ia 42           // 效果与上条命令等同。[flags] 可以使用缩写，缩写时使用单个减号，并且 [flags] 可以连起来一起写`,
	ValidArgsFunction:     nil,
	Args:                  cobra.MinimumNArgs(1), // 至少一个参数
	Hidden:                false,
	DisableAutoGenTag:     false,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		var oprs []models.Operator
		// var oprsmap = make(map[string]int) // map[opr.Name]idx
		var keys = args

		// 优先查询干员 ID
		if id {
			var res []models.Operator
			database.DB.Where("id", keys).Find(&res)
			logger.Debug(len(res))
			oprs = append(oprs, res...)
			logger.Debug(oprs, len(oprs))
		}

		// 查询干员名称
		if name {
			var res []models.Operator
			stmt := database.DB.Where("name like ?", "%"+keys[0]+"%")
			for i := 1; i < len(keys); i++ {
				stmt.Or("name like ?", "%"+keys[i]+"%")
			}
			stmt.Find(&res)
			oprs = append(oprs, res...)
		}

		// 再查询别名
		if alias {
			fmt.Printf("暂不支持别名查询\n")
			// database.DB.Where("id", key).Find(&oprs)
			// for _, opr := range oprs {
			// 	fmt.Printf("%+v\n", opr)
			// }
		}

		// oprs 去重
		var idxmap = make(map[int]struct{})
		var res []models.Operator
		for i := 0; i < len(oprs); i++ {
			if _, ok := idxmap[oprs[i].ID]; !ok {
				idxmap[oprs[i].ID] = struct{}{}
				res = append(res, oprs[i])
			}
		}
		for _, opr := range res {
			fmt.Printf("%+v\n", opr)
		}
	},
}

func init() {
	oprsCmd.Flags().BoolVarP(&id, "id", "i", false, "使用「干员 ID」进行查询")
	oprsCmd.Flags().BoolVarP(&name, "name", "n", true, "使用「干员名称」进行查询")
	oprsCmd.Flags().BoolVarP(&alias, "alias", "a", false, "使用「干员别名」进行查询")

	rootCmd.AddCommand(oprsCmd)
}
