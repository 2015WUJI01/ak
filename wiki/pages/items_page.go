package pages

import (
	"github.com/2015WUJI01/looog"
	"github.com/gocolly/colly"
	"main/pkg/progressbar"
	"main/wiki"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// ItemsPage
// "https://prts.Wiki/w/分类:道具"
type ItemsPage struct {
	// 数据通道
	Names, Wikis chan string
	Total        chan int

	c *colly.Collector

	bar *progressbar.ProgressBar

	Count int
}

type ItemsPageFetcher func(*ItemsPage)

func NewItemsPage(c *colly.Collector, fn ...ItemsPageFetcher) *ItemsPage {
	p := &ItemsPage{
		c:     c,
		Names: make(chan string, 600),
		Wikis: make(chan string, 600),
		Total: make(chan int, 1),
		bar:   progressbar.New("[Step.1] 批量获取 itemspage 数据", 0),
	}
	for _, f := range fn {
		f(p)
	}
	return p
}

func (p *ItemsPage) SetCollector(c *colly.Collector) { p.c = c }

// https://prts.wiki/index.php?title=%E5%88%86%E7%B1%BB:%E9%81%93%E5%85%B7
func (p *ItemsPage) page() string { return "https://prts.Wiki/w/分类:道具" }

func (p *ItemsPage) FetchTotal() {
	once := &sync.Once{}
	p.c.OnResponse(func(r *colly.Response) {
		looog.Debug(r.Request.URL)
		once.Do(func() {
			p.Count++
			arr := regexp.MustCompile(`共(.*)个页面`).FindSubmatch(r.Body)
			t, _ := strconv.Atoi(string(arr[1]))
			p.c.MaxDepth = int(math.Floor(float64(t) / 200))
			// looog.Debugf("[step.1] 获取 itemspage 数据总条数: %d", t)
			p.bar.ChangeMax(t)
			p.Total <- t
		})
	})
}

func (p *ItemsPage) FetchName() {
	p.c.OnHTML(".mw-category-group ul li a", func(a *colly.HTMLElement) {
		p.Names <- strings.TrimSpace(a.Text)
		_ = p.bar.Add(1)
	})
}

func (p *ItemsPage) FetchWiki() {
	p.c.OnHTML(".mw-category-group ul li a", func(a *colly.HTMLElement) {
		p.Wikis <- wiki.Link(a.Attr("href"))
	})
}

func (p *ItemsPage) FetchDeeper() {
	p.c.OnHTML(`a[title="分类:道具"]:last-of-type`, func(e *colly.HTMLElement) {
		if e.Text == "下一页" {
			_ = p.c.Visit(e.Request.AbsoluteURL(e.Attr("href")))
		}
	})
}

type ItemsFlag int

const (
	ITEMS_NONE ItemsFlag = 0
	ITEMS_NAME ItemsFlag = 1 << (iota - 1)
	ITEMS_WIKI
	ITEMS_ALL = ITEMS_NAME | ITEMS_WIKI
)

func (p *ItemsPage) Bind(flag ItemsFlag) *ItemsPage {
	p.FetchTotal()
	p.FetchDeeper()

	if flag == ITEMS_NONE {
		return p
	}

	if flag&ITEMS_NAME == ITEMS_NAME {
		p.FetchName()
	}
	if flag&ITEMS_WIKI == ITEMS_WIKI {
		p.FetchWiki()
	}

	return p
}

// Visit 开始爬取
func (p *ItemsPage) Visit() {
	_ = p.c.Visit(p.page())
}
