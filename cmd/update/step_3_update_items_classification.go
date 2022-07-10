package update

import (
	"github.com/gocolly/colly"
	"main/logger"
	"main/models"
)

// 道具的分类
type classification struct {
	Name  string // 道具名称
	Group string // 道具分组（大类），例如：芯片；材料；
	Type  string // 道具类型（小类），例如：芯片、芯片组、双芯片；T1、T2、T3、T4、T5；
}

// Step3 抓取所有道具分类情况，并依次更新
func Step3() {

	var itemsMap = getItemsMapFromDB()
	var items []models.Item
	for _, cf := range fetchClassifications() {
		item, ok := itemsMap[cf.Name]
		if !ok {
			logger.Warn(cf.Name)
			continue
		}
		items = append(items, models.Item{
			Name:  item.Name,
			Group: cf.Group,
			Type:  cf.Type,
		})
	}
	CreateOrUpdateItems([]string{"name", "group", "type"}, items)
	logger.Infof("Step3. 道具分类信息更新完成")
}

// fetchClassifications
// 对 [1] 页面进行抓取，分析最后的表格，从中梳理出每个道具的分类情况
// [1] 理智信息页面 https://prts.wiki/w/%E7%90%86%E6%99%BA
func fetchClassifications() []classification {
	var cfs []classification

	c := colly.NewCollector()
	c.OnError(func(r *colly.Response, err error) {
		_ = r.Request.Retry()
	})
	// 表格中一行就是一个 group
	c.OnHTML("table.uncollapsed > tbody > tr:has(th:not(.navbox-title))", func(tr *colly.HTMLElement) {
		var group string
		var types []string

		// 获取 Group 名称
		tr.ForEach(`th.navbox-group:not([style="text-align:center; width:5%"])`,
			func(_ int, th *colly.HTMLElement) { group = th.Text })

		// 获取 Type 列表，有可能当前 group 没有 type
		tr.ForEach(`th.navbox-group[style="text-align:center; width:5%"]`,
			func(_ int, th *colly.HTMLElement) { types = append(types, th.Text) })

		if len(types) == 0 {
			// 当 group 没有 type 时，此时的一行即
			tr.ForEach(`td:has(div:not(:empty)[style="padding:0em 0.25em"]) .smw-value`,
				func(_ int, e *colly.HTMLElement) {
					cfs = append(cfs, classification{Name: e.Text, Group: group})
				})
		} else {
			// 二级元素
			tr.ForEach("td:has(table) table tr:has(th)", func(i int, tr *colly.HTMLElement) {
				tr.ForEach("td li .smw-value", func(_ int, e *colly.HTMLElement) {
					cfs = append(cfs, classification{Name: e.Text, Group: group, Type: types[i]})
				})
			})
		}
	})
	_ = c.Visit("https://prts.wiki/w/%E7%90%86%E6%99%BA")
	return cfs
}
