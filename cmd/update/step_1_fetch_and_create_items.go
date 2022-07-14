package update

import (
	"github.com/gocolly/colly"
	"main/logger"
	"main/models"
	repo "main/repositories"
	"main/wiki"
	"strings"
	"time"
)

// Step1 抓取一次页面，获取所有的道具，如果有不存在的道具，则创建这个道具
func Step1() []models.Item {
	items := fetchAllItems()
	logger.Infof("Step1. 共获取到 %d 条道具数据", len(items))
	repo.CreateOrUpdateItems([]string{"name", "wiki"}, items)
	return items
}

// fetchAllItems 获取所有道具
// 从 [1] 页面中获取前 200 条数据，然后找到下一页的地址，访问下一次获取接下来 200 条，并以此类推。
// 目前共 500+ 数据，所以需要访问三次，若增加到 600+ 条，则需要修改代码，增加一轮访问流程。
//
// [1] https://prts.wiki/w/%E5%88%86%E7%B1%BB:%E9%81%93%E5%85%B7
func fetchAllItems() []models.Item {
	var items []models.Item

	c := colly.NewCollector()
	c.SetRequestTimeout(5 * time.Second)
	c.OnError(func(r *colly.Response, err error) {
		_ = r.Request.Retry()
	})

	// 获取第一轮的 200 条数据
	c.OnHTML(".mw-category-group ul li a", func(a *colly.HTMLElement) {
		items = append(items, models.Item{
			Name: strings.TrimSpace(a.Text),
			Wiki: wiki.Link(a.Attr("href")),
		})
	})

	// 获取第二页链接
	c.OnHTML(`a[title="分类:道具"]:last-of-type`, func(e *colly.HTMLElement) {
		cc := colly.NewCollector()
		cc.SetRequestTimeout(5 * time.Second)
		cc.OnError(func(r *colly.Response, err error) {
			_ = r.Request.Retry()
		})
		// 获取第二轮的 200 条数据
		cc.OnHTML(".mw-category-group ul li a", func(a *colly.HTMLElement) {
			items = append(items, models.Item{
				Name: strings.TrimSpace(a.Text),
				Wiki: wiki.Link(a.Attr("href")),
			})
		})
		// 获取第三页链接
		cc.OnHTML(`a[title="分类:道具"]:last-of-type`, func(ee *colly.HTMLElement) {
			ccc := colly.NewCollector()
			ccc.SetRequestTimeout(5 * time.Second)
			ccc.OnError(func(r *colly.Response, err error) {
				_ = r.Request.Retry()
			})
			// 获取第三页的 200 条数据
			ccc.OnHTML(".mw-category-group ul li a", func(a *colly.HTMLElement) {
				items = append(items, models.Item{
					Name: strings.TrimSpace(a.Text),
					Wiki: wiki.Link(a.Attr("href")),
				})
			})
			_ = ccc.Visit(wiki.Link(ee.Attr("href")))
		})
		_ = cc.Visit(wiki.Link(e.Attr("href")))
	})
	_ = c.Visit("https://prts.wiki/w/分类:道具")
	return items
}
