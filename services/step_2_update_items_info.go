package services

import (
	"ak/common"
	"ak/models"
	"ak/wiki"
	"github.com/2015WUJI01/looog"
	"github.com/gocolly/colly"
	"strings"
	"time"
)

// Step2 从数据库中获取所有道具，对每个道具进行数据更新
func Step2(items []models.Item, dataChan chan models.Item) {
	defer close(dataChan)
	// 需要访问每个道具的页面，所以开启异步
	c := colly.NewCollector(
		colly.Async(true),
		// colly.CacheDir("tmp/cache"),
	)
	c.SetRequestTimeout(3 * time.Second)
	_ = c.Limit(&colly.LimitRule{
		DomainGlob:  "*",                    // limit 规则只会在指定的 domain 生效，所以必须配置
		Parallelism: 2,                      // 同时抓取的协程数
		Delay:       200 * time.Millisecond, // 抓取时延
	})

	c.OnError(func(r *colly.Response, err error) {
		looog.Warn("retrying")
		_ = r.Request.Retry()
	})

	// 爬取每个道具页面
	c.OnHTML("body", func(body *colly.HTMLElement) {
		// 页面中包含的信息有：道具名称、icon 图片链接、wiki 短链接，以及最后编辑时间（视为上次更新时间）
		var name, image, wikishort string
		var updatedAt time.Time

		// 1. 获取道具名称
		body.ForEach("#firstHeading", func(i int, e *colly.HTMLElement) {
			// 道具页面的名称与分类页面的名称有误差，需要做修正
			name = e.Text
			for old, n := range map[string]string{"α": "Α", "β": "Β", "γ": "Γ"} {
				name = strings.ReplaceAll(name, old, n)
			}
		})

		// 2. 获取图片链接
		body.ForEach("td.nomobile", func(i int, td *colly.HTMLElement) {
			if i == 1 {
				td.ForEach("a.image > img", func(_ int, img *colly.HTMLElement) {
					if img.Attr("data-src") == "" {
						image = wiki.Link(img.Attr("src"))
					} else {
						image = wiki.Link(img.Attr("data-src"))
					}
				})
			}
		})

		// 3. 获取 wiki 短链接
		body.ForEach(".copyUrl", func(i int, e *colly.HTMLElement) {
			wikishort = e.Attr("data-clipboard-text")
		})

		// 4. 获取最后编辑时间
		body.ForEach("#footer-info-lastmod", func(_ int, e *colly.HTMLElement) {
			updatedAt = common.LastMod(e.Text).Time()
		})

		i := models.Item{
			Name:      name,
			Image:     image,
			WikiShort: wikishort,
			UpdatedAt: updatedAt,
		}
		dataChan <- i
	})

	for _, item := range items {
		_ = c.Visit(item.Wiki)
	}
	c.Wait()
}
