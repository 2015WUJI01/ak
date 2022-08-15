package pages

import (
	"github.com/gocolly/colly"
	"log"
	"main/wiki"
	"strings"
	"time"
)

type ItemPage struct {
	Url  string
	data struct {
		Name      string
		Image     string
		ShortWiki string
		Updatedat time.Time
	}

	c           *colly.Collector
	onCompleted func(p *ItemPage)
}

type ItemPageOption func(p *ItemPage)

func NewItemPage(opts ...ItemPageOption) *ItemPage {
	p := &ItemPage{}
	for _, fn := range opts {
		fn(p)
	}
	return p
}

func (p *ItemPage) SetCollector(c *colly.Collector) {
	p.c = c
}

func (p *ItemPage) SetUrl(url string) { p.Url = url }

func (p *ItemPage) Data() struct {
	Name      string
	Image     string
	ShortWiki string
	Updatedat time.Time
} {
	return p.data
}

// FetchImage 获取图片链接
func (p *ItemPage) FetchImage(c *colly.Collector) {
	c.OnHTML("body td.nomobile", func(e *colly.HTMLElement) {
		if e.Index == 1 {
			e.ForEach("a.image > img", func(_ int, img *colly.HTMLElement) {
				if img.Attr("data-src") == "" {
					p.data.Image = wiki.Link(img.Attr("src"))
				} else {
					p.data.Image = wiki.Link(img.Attr("data-src"))
				}
			})
		}
	})
}

// FetchShortWiki 获取 wiki 短链接
func (p *ItemPage) FetchShortWiki(c *colly.Collector) {
	c.OnHTML("body", func(body *colly.HTMLElement) {
		body.ForEach(".copyUrl", func(i int, e *colly.HTMLElement) {
			p.data.ShortWiki = e.Attr("data-clipboard-text")
		})
	})
}

// FetchUpdatedAt 获取 wiki 短链接
func (p *ItemPage) FetchUpdatedAt(c *colly.Collector) {
	c.OnHTML("body", func(body *colly.HTMLElement) {
		// 4. 获取最后编辑时间
		body.ForEach("#footer-info-lastmod", func(_ int, e *colly.HTMLElement) {
			p.data.Updatedat = wiki.UpdatedAtStr(e.Text).AsTime()
		})
	})
}

// FetchName 获取道具名称
func (p *ItemPage) FetchName(c *colly.Collector) {
	c.OnHTML("#firstHeading", func(e *colly.HTMLElement) {
		p.data.Name = e.Text
		for old, now := range map[string]string{"α": "Α", "β": "Β", "γ": "Γ"} {
			p.data.Name = strings.ReplaceAll(p.data.Name, old, now)
		}
	})
}

func (p *ItemPage) OnCompleted() {}
func (p *ItemPage) SetOnCompleted(f func(p *ItemPage)) {
	p.onCompleted = f
}

func (p *ItemPage) Visit() {
	// 初始化 collector
	if p.c == nil {
		p.c = p.defaultCollector()
	}

	// 注册需要爬取的对象
	p.FetchName(p.c)
	p.FetchImage(p.c)
	p.FetchShortWiki(p.c)
	p.FetchUpdatedAt(p.c)

	// 开始爬取
	log.Println(p.Url)
	// 处理数据
	_ = p.c.Visit(p.Url)
	if p.onCompleted != nil {
		p.onCompleted(p)
	} else {
		p.OnCompleted()
	}
}

func (p *ItemPage) Wait() {
	p.c.Wait()
}
func (p *ItemPage) defaultCollector() *colly.Collector {
	c := colly.NewCollector()
	c.SetRequestTimeout(5 * time.Second)
	c.OnError(func(r *colly.Response, err error) {
		_ = r.Request.Retry()
	})
	return c
}
