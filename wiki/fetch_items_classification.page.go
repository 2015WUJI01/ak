package wiki

import "github.com/gocolly/colly"

type itemClassificationPageData struct {
	name  string // 道具名称
	group string // 道具分组（大类），例如：芯片；材料；
	typee string // 道具类型（小类），例如：芯片、芯片组、双芯片；T1、T2、T3、T4、T5；
}

func (data itemClassificationPageData) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":  data.name,
		"group": data.group,
		"typee": data.typee,
	}
}

// FetchItemsClassification
// 对 [1] 页面进行抓取，分析最后的表格，从中梳理出每个道具的分类情况
// [1] 理智信息页面 https://prts.wiki/w/%E7%90%86%E6%99%BA
func FetchItemsClassification(names ...string) map[string]map[string]interface{} {
	var (
		mnames = make(map[string]struct{})
		micpds = make(map[string]map[string]interface{})
	)
	for _, name := range names {
		mnames[name] = struct{}{}
	}

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
					if _, ok := mnames[e.Text]; ok {
						micpds[e.Text] = itemClassificationPageData{
							name:  e.Text,
							group: group,
							typee: "",
						}.ToMap()
					}
				})
		} else {
			// 二级元素
			tr.ForEach("td:has(table) table tr:has(th)", func(i int, tr *colly.HTMLElement) {
				tr.ForEach("td li .smw-value", func(_ int, e *colly.HTMLElement) {
					if _, ok := mnames[e.Text]; ok {
						micpds[e.Text] = itemClassificationPageData{
							name:  e.Text,
							group: group,
							typee: types[i],
						}.ToMap()
					}
				})
			})
		}
	})
	_ = c.Visit("https://prts.wiki/w/%E7%90%86%E6%99%BA")
	return micpds
}
