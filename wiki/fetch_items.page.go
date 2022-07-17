package wiki

import (
	"github.com/gocolly/colly"
	"strings"
	"time"
)

type itemsPageData struct {
	Name string
	Wiki string
}

// 根据名称获取 wiki
// 若不传参数，名称为空，则默认获取所有的 wiki
// 返回名称和 wiki 的键值对
func FetchAllWiki(names ...string) map[string]string {
	var (
		mnames = make(map[string]struct{})
		mwikis = make(map[string]string)
	)
	for _, name := range names {
		mnames[name] = struct{}{}
	}
	res := fetchItemsPage(mnames)
	for name, data := range res {
		mwikis[name] = data.Wiki
	}
	return mwikis
}

func fetchItemsPage(names map[string]struct{}) (res map[string]itemsPageData) {
	fetchAll := len(names) == 0
	res = make(map[string]itemsPageData)

	c := colly.NewCollector()
	c.SetRequestTimeout(5 * time.Second)
	c.OnError(func(r *colly.Response, err error) {
		_ = r.Request.Retry()
	})

	// 获取第一轮的 200 条数据
	c.OnHTML(".mw-category-group ul li a", func(a *colly.HTMLElement) {
		name := strings.TrimSpace(a.Text)
		if _, ok := names[name]; fetchAll || ok {
			res[name] = itemsPageData{
				Name: name,
				Wiki: Link(a.Attr("href")),
			}
		}
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
			name := strings.TrimSpace(a.Text)
			if _, ok := names[name]; fetchAll || ok {
				res[name] = itemsPageData{
					Name: name,
					Wiki: Link(a.Attr("href")),
				}
			}
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
				name := strings.TrimSpace(a.Text)
				if _, ok := names[name]; fetchAll || ok {
					res[name] = itemsPageData{
						Name: name,
						Wiki: Link(a.Attr("href")),
					}
				}
			})
			_ = ccc.Visit(Link(ee.Attr("href")))
		})
		_ = cc.Visit(Link(e.Attr("href")))
	})
	_ = c.Visit("https://prts.Wiki/w/分类:道具")
	return res
}
