package ri

import "C"
import (
	"github.com/gocolly/colly"
	"log"
	"main/database"
	"main/pkg/logger"
	"main/wiki"
	"strings"
)

var OC OperatorCenter

type OperatorCenter struct {
	// 所有干员
	Operators []Operator
	// 所有干员的别名（包含原名）
	AliasList []string
	// 别名查真名用的map
	RevAliasMap map[string]string
}

// UpdateOperatorsName 增量更新干员名称
func (oc *OperatorCenter) UpdateOperatorsName() (oprs []Operator) {

	// 第一步：获取所有的干员名称（以及是否是肉鸽限定）
	names, roguelikes := fetchOprNames()

	// 第二步：将得到的干员名称数据存到数据库中
	for i := 0; i < len(names); i++ {
		o := Operator{
			UUID:      i + 1,
			Name:      names[i],
			Roguelike: roguelikes[i],
		}

		res := database.DB.Exec("INSERT INTO `operators`(`uuid`, `name`, `roguelike`) VALUES(?,?,?) ON DUPLICATE KEY UPDATE `name`=?, `roguelike`=?", o.UUID, o.Name, o.Roguelike, o.Name, o.Roguelike)

		if res.RowsAffected > 0 {
			if opr, ok := OC.FindOprByUuid(i + 1); ok {
				oprs = append(oprs, opr)
				logger.Debugf("增量更新干员名称数据：%d %s", opr.UUID, opr.Name)
			}
		}
	}
	// 第三步：从数据库中更新 oprs 数据
	_ = oc.Update()
	return
}

// Update 刷新 OC 的数据
func (oc *OperatorCenter) Update() error {
	// 从数据库中获取 oprs 数据
	database.DB.Model(Operator{}).Find(&oc.Operators)
	// 初始化 RevAliasMap
	oc.InitRevAliasMap()
	return nil
}

// InitRevAliasMap 初始化 RevAliasMap
func (oc *OperatorCenter) InitRevAliasMap() {
	var aliasMap = make(map[string][]string)
	for _, opr := range oc.Operators {
		if len(opr.Alias) != 0 {
			aliasMap[opr.Name] = opr.Alias
		}
	}
	oc.RevAliasMap = reverseAliasMap(&aliasMap)
}

func reverseAliasMap(aliasMap *map[string][]string) map[string]string {
	var revAliasMap = make(map[string]string)
	for name, aliases := range *aliasMap {
		revAliasMap[name] = name
		for _, alias := range aliases {
			revAliasMap[alias] = name
		}
	}
	return revAliasMap
}

// FindOprByName 在数据库中按名称精准查询
func (oc OperatorCenter) FindOprByName(name string) (opr Operator, ok bool) {
	res := database.DB.Model(Operator{}).Where("name = ?", name).First(&opr)
	if res.Error != nil {
		return Operator{}, false
	}
	return opr, true
}

// FindOprByAlias 在数据库中按别名精准查询
func (oc OperatorCenter) FindOprByAlias(alias string) (opr Operator, ok bool) {
	var oprs []Operator
	res := database.DB.Model(&Operator{}).Where("JSON_CONTAINS(alias, ?)", alias).Or("name = ?", alias).Find(&oprs)
	if res.Error != nil {
		log.Println("使用别名查询干员时发生错误，" + res.Error.Error())
		return opr, false
	}
	if len(oprs) == 0 {
		return opr, false
	} else {
		if len(oprs) > 1 {
			log.Println(oprs)
		}
		return oprs[0], true
	}
}

func (oc OperatorCenter) FindOprByUuid(uuid int) (opr Operator, ok bool) {
	res := database.DB.Model(&Operator{}).Where("uuid = ?", uuid).Find(&opr)
	if res.Error != nil {
		log.Println("使用 Uuid 查询干员时发生错误，" + res.Error.Error())
		return opr, false
	}
	return opr, true
}

func fetchOprNames() (names []string, roguelikes []bool) {
	_ = NewFetcher().
		AutoRetry(0).
		OnHTML("#mw-content-text table tbody tr td:nth-of-type(1)", func(a *colly.HTMLElement) {
			names = append(names, strings.TrimSpace(a.Text))
			switch a.Text {
			case "Touch", "Sharp", "Stormeye", "Pith", "预备干员-术师", "预备干员-近战", "预备干员-狙击", "预备干员-后勤", "暮落(集成战略)":
				roguelikes = append(roguelikes, true)
			default:
				roguelikes = append(roguelikes, false)
			}
		}).
		Visit(wiki.Page("/w/干员上线时间一览"))
	return
}
