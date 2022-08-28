package services

import (
	"ak/models"
	"github.com/2015WUJI01/looog"
	"github.com/gocolly/colly"
	"time"
)

func Step3(dataChan chan models.Item) {
	defer close(dataChan)
	// 需要访问每个道具的页面，所以开启异步
	c := colly.NewCollector(
		colly.CacheDir("tmp/cache"),
	)
	c.SetRequestTimeout(3 * time.Second)
	c.OnError(func(r *colly.Response, err error) {
		looog.Warnf("retrying: %s", r.Request.URL)
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
					dataChan <- models.Item{
						Name:  e.Text,
						Group: group,
						Type:  "",
					}
				})
		} else {
			// 二级元素
			tr.ForEach("td:has(table) table tr:has(th)", func(i int, tr *colly.HTMLElement) {
				tr.ForEach("td li .smw-value", func(_ int, e *colly.HTMLElement) {
					dataChan <- models.Item{
						Name:  e.Text,
						Group: group,
						Type:  types[i],
					}
				})
			})
		}
	})
	_ = c.Visit("https://prts.wiki/w/理智")
	c.Wait()
}
