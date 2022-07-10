package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"main/database"
	"main/models"
)

// opr 命令的参数
var (
	id    bool
	name  bool
	alias bool
)

var oprCmd = &cobra.Command{
	Use: `opr [flags] <args>
  默认仅支持的「干员名称」查询，即 [flags] 默认仅使用 --name。当存在其他 [flags] 时，默认的 --name 同样生效，此时，多 [flags] 情况下，结果会按照「干员 ID」、「干员名称」、「干员别名」递减的优先级进行排序，返回优先级最高的结果。`,
	Short: "查询单个干员信息",
	Long:  "查询最符合要求的首位干员的所有数据",
	Example: `  ak opr 斯卡蒂           // 查询结果：斯卡蒂。精准查找，不会查询出浊心斯卡蒂或其他干员

  ak opr --id 123         // 查询结果：「干员 ID」为 123 的干员。
  ak opr --alias 42       // 查询结果：史尔特尔。「干员别名」为 42 的只有史尔特尔
  ak opr --alias 小车     // 查询结果：最匹配的所有小车其中之一。若需要查询所有匹配的干员，需要使用 oprs 命令

  ak opr --id --alias 42  // 查询结果：「干员 ID」为 123 的干员。会查找出「干员 ID」为 42 的干员和「干员别名」为 42 的史尔特尔，但因为「干员别名」优先级较低而不会显示，若想同时显示这两位干员需要使用 oprs 命令
  ak opr -ia 42           // 效果与上条命令等同。[flags] 可以使用缩写，缩写时使用单个减号，并且 [flags] 可以连起来一起写`,
	ValidArgsFunction:     nil,
	Args:                  cobra.MinimumNArgs(1), // 至少一个参数
	Hidden:                false,
	DisableAutoGenTag:     false,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		var opr models.Operator
		var key = args[0]

		// 优先查询干员 ID
		if id {
			database.DB.Where("id", key).First(&opr)
			if (models.Operator{}) != opr {
				fmt.Printf("%+v\n", opr)
				return
			}
		}

		// 查询干员名称
		if name {
			database.DB.Where("name", key).First(&opr)
			if (models.Operator{}) != opr {
				fmt.Printf("%+v\n", opr)
				return
			}
		}

		// 再查询别名
		if alias {
			fmt.Printf("暂不支持别名查询\n")
			// database.DB.Where("id", key).Find(&oprs)
			// for _, opr := range oprs {
			// 	fmt.Printf("%+v\n", opr)
			// }
		}
		fmt.Printf("查询无果\n")
	},
}

func init() {
	oprCmd.Flags().BoolVarP(&id, "id", "i", false, "使用「干员 ID」进行查询")
	oprCmd.Flags().BoolVarP(&name, "name", "n", true, "使用「干员名称」进行查询")
	oprCmd.Flags().BoolVarP(&alias, "alias", "a", false, "使用「干员别名」进行查询")
	oprCmd.Flags().BoolP("auto", "A", false, "自动模式，即使用 -ian 进行模糊查询")

	rootCmd.AddCommand(oprCmd)
}
