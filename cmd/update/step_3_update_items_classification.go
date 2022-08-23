package update

import (
	"ak/database"
	"ak/logger"
	"ak/models"
	repo "ak/repositories"
	"ak/wiki"
)

// 道具的分类
type classification struct {
	Name  string // 道具名称
	Group string // 道具分组（大类），例如：芯片；材料；
	Type  string // 道具类型（小类），例如：芯片、芯片组、双芯片；T1、T2、T3、T4、T5；
}

// Step3 抓取所有道具分类情况，并依次更新
func Step3() {
	var items []models.Item
	_ = database.DB.Find(&items)
	micpds := wiki.FetchItemsClassification(func() []string {
		var names []string
		for _, item := range items {
			names = append(names, item.Name)
		}
		return names
	}()...)
	for i, item := range items {
		if v, ok := micpds[item.Name]; ok {
			items[i].Group = v["group"].(string)
			items[i].Type = v["typee"].(string)
		}
	}
	repo.CreateOrUpdateItems([]string{"name", "group", "type"}, items)
	logger.Infof("Step3. 道具分类信息更新完成")
}
