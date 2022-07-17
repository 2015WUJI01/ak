package update

import (
	"main/logger"
	"main/models"
	repo "main/repositories"
	"main/wiki"
)

// Step1 抓取一次页面，获取所有的道具，如果有不存在的道具，则创建这个道具
func Step1() {
	// 方法一：获取指针后，遍历出 item 值
	var items []models.Item
	// for _, item := range models.FreshItemWiki() {
	// 	items = append(items, *item)
	// }
	// logger.Infof("Step1. 共获取到 %d 条道具数据", len(items))
	// repo.CreateOrUpdateItems([]string{"name", "wiki"}, items)

	// 方法二：直接传 items 指针
	// items := models.FreshItemWiki()
	// logger.Infof("Step1. 共获取到 %d 条道具数据", len(items))
	// repo.CreateOrUpdateItemsp([]string{"name", "wiki"}, items)

	for n, w := range wiki.FetchAllWiki() {
		items = append(items, models.Item{Name: n, Wiki: w})
	}
	logger.Infof("Step1. 共获取到 %d 条道具数据", len(items))
	repo.CreateOrUpdateItems([]string{"name", "wiki"}, items)
}
