package wiki

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/schollz/progressbar/v3"
	"net/url"
	"strings"
	"time"
)

type ItemPage struct {
	GoalsName []string
	Data      []ItemPageData
}

type ItemPageData struct {
	name      string
	image     string
	wikishort string
	updatedat time.Time
}

func (data ItemPageData) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":      data.name,
		"image":     data.image,
		"wikishort": data.wikishort,
		"updatedat": data.updatedat,
	}
}

func FetchItemInfo(names ...string) map[string]map[string]interface{} {
	var (
		mnames = make(map[string]struct{})
		minfo  = make(map[string]map[string]interface{})
	)
	for _, name := range names {
		mnames[name] = struct{}{}
	}
	res := fetchItemPage(mnames)
	for name, data := range res {
		minfo[name] = data.ToMap()

	}
	return minfo
}

func fetchItemPage(names map[string]struct{}) (res map[string]ItemPageData) {
	res = make(map[string]ItemPageData)

	// 需要访问每个道具的页面，所以开启异步
	c := NewAsyncCollector(4, 50, 5*time.Second)
	// 爬取每个道具页面
	bar := progressbar.NewOptions(len(names),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
		}),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription("逐条采集道具信息中..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: "[green]-[reset]",
			BarStart:      "[",
			BarEnd:        "]",
		}))

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

		if _, ok := names[name]; !ok {
			return
		}

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
		})

		res[name] = ItemPageData{
			name:      name,
			image:     image,
			wikishort: wikishort,
			updatedat: updatedAt,
		}
		_ = bar.Add(1)
	})

	for name := range names {
		_ = c.Visit(Link("/w/" + url.PathEscape(name)))
	}
	c.Wait()
	fmt.Println()
	return res
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
