package update

import (
	"fmt"
	"github.com/gocolly/colly"
	"main/database"
	"main/models"
	"main/pkg/logger"
	"strings"
	"time"
)

// Step2 从数据库中获取所有道具，对每个道具进行数据更新
func Step2() {

	fmt.Println("逐条采集道具信息中...")

	var itemsMap = getItemsMapFromDB()
	fetchItemsInfo(itemsMap)
	var items []models.Item
	for _, item := range itemsMap {
		items = append(items, item)
	}

	CreateOrUpdateItems([]string{"name", "image", "wiki_short", "updated_at"}, items)
	// CreateOrUpdateItems([]string{"group", "type"}, items)

	logger.Infof("Step2. 道具基本信息更新完成")
}

func getItemsMapFromDB() map[string]models.Item {
	var items []models.Item
	var itemsMap = make(map[string]models.Item)
	_ = database.DB.Find(&items)
	for _, item := range items {
		itemsMap[item.Name] = item
	}
	return itemsMap
}

func fetchItemsInfo(itemsMap map[string]models.Item) {
	// 需要访问每个道具的页面，所以开启异步
	c := colly.NewCollector(colly.Async(true))
	c.SetRequestTimeout(5 * time.Second)
	c.OnError(func(r *colly.Response, err error) {
		_ = r.Request.Retry()
	})

	// 限制异步频率
	_ = c.Limit(&colly.LimitRule{
		DomainGlob:  "*prts.wiki*",         // limit 规则只会在指定的 domain 生效，所以必须配置
		Parallelism: 5,                     // 同时抓取的协程数
		Delay:       50 * time.Millisecond, // 抓取时延
	})

	// 爬取每个道具页面
	completed := 0
	c.OnHTML("body", func(body *colly.HTMLElement) {
		// 页面中包含的信息有：道具名称、icon 图片链接、wiki 短链接，以及最后编辑时间（视为上次更新时间）
		var name, image, wikishort string
		var updatedAt time.Time

		// 需要更新的字段
		updateColums := []string{"name", "image", "wiki_short"}

		// 1. 获取道具名称
		body.ForEach("#firstHeading", func(i int, e *colly.HTMLElement) {
			// 道具页面的名称与分类页面的名称有误差，需要做修正
			name = strings.ReplaceAll(e.Text, "α", "Α")
			name = strings.ReplaceAll(name, "β", "Β")
			name = strings.ReplaceAll(name, "γ", "Γ")
		})

		// 2. 获取图片链接
		body.ForEach("td.nomobile", func(i int, td *colly.HTMLElement) {
			if i == 1 {
				td.ForEach("a.image > img", func(_ int, img *colly.HTMLElement) {
					if img.Attr("data-src") == "" {
						image = Link(img.Attr("src"))
					} else {
						image = Link(img.Attr("data-src"))
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
			updatedAt = parseTime(e.Text)
			if !updatedAt.IsZero() {
				updateColums = append(updateColums, "updated_at")
			}
		})

		itemsMap[name] = models.Item{
			Name:      itemsMap[name].Name,
			Image:     image,
			WikiShort: wikishort,
			UpdatedAt: updatedAt,
		}

		completed++
		fmt.Printf("=")
		if completed%100 == 0 {
			fmt.Println()
		}
	})

	for _, item := range itemsMap {
		_ = c.Visit(item.Wiki)
	}
	c.Wait()
	fmt.Println()
}

// parseTime 解析源码中的时间字符串
// 原文本为 "此页面最后编辑于2022年5月22日 (星期日) 12:32。"
// 无法直接作为 go 时间解析规则，所以尝试替换 7 次，总有一种能够解析成功
func parseTime(timeStr string) time.Time {
	weeks := []string{"日", "一", "二", "三", "四", "五", "六"}
	var err error
	var t time.Time
	for i := 0; i <= 7; i++ {
		layout := fmt.Sprintf("此页面最后编辑于2006年1月2日 (星期%s) 15:04 -0700", weeks[i])
		t, err = time.Parse(layout, strings.ReplaceAll(strings.TrimSpace(timeStr), "。", " +0800"))
		if err == nil {
			break
		}
	}
	return t
}

func updateItemInColumns(cols []string, item *models.Item) bool {
	res := database.DB.Model(&models.Item{}).Select(cols).Where("name", item.Name).Updates(&item)
	return res.RowsAffected > 0
}
