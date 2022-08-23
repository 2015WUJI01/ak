package pages

import (
	"ak/pkg/progressbar"
	"ak/wiki"
	"github.com/gocolly/colly"
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

	// 采集器
	c *colly.Collector

	// 进度条
	bar *progressbar.ProgressBar

	// Total 的值，存一份在内部使用
	count int
}

type ItemsPageSetting func(*ItemsPage)

func NewItemsPage(c *colly.Collector, fn ...ItemsPageSetting) *ItemsPage {
	p := &ItemsPage{
		c:     c,
		Names: make(chan string, 600),
		Wikis: make(chan string, 600),
		Total: make(chan int, 1),
		// bar:   progressbar.New("[Step.1] 获取 items 数据", 0),
	}
	for _, f := range fn {
		f(p)
	}
	return p
}

func SetItemsPageProgressBar(bar *progressbar.ProgressBar) ItemsPageSetting {
	return func(p *ItemsPage) { p.bar = bar }
}

// https://prts.wiki/index.php?title=%E5%88%86%E7%B1%BB:%E9%81%93%E5%85%B7
func (p *ItemsPage) page() string { return "https://prts.Wiki/w/分类:道具" }

func (p *ItemsPage) FetchTotal() {
	once := &sync.Once{}
	p.c.OnResponse(func(r *colly.Response) {
		once.Do(func() {
			// 获取数据总条数
			arr := regexp.MustCompile(`共(.*)个页面`).FindSubmatch(r.Body)
			t, _ := strconv.Atoi(string(arr[1]))
			// 根据数据条数，自动判断需要遍历多少次
			p.count = t
			p.c.MaxDepth = int(math.Floor(float64(t) / 200))
			if p.bar != nil {
				p.bar.ChangeMax(t) // 重新复制进度条的长度
			}
			p.Total <- t
		})
	})
}

func (p *ItemsPage) FetchName() {
	count := 0
	p.c.OnHTML(".mw-category-group ul li a", func(a *colly.HTMLElement) {
		count++
		p.Names <- strings.TrimSpace(a.Text)
		if p.bar != nil {
			_ = p.bar.Add(1)
		}
		if p.count != 0 && p.count == count {
			close(p.Names)
		}
	})
}

func (p *ItemsPage) FetchWiki() {
	count := 0
	p.c.OnHTML(".mw-category-group ul li a", func(a *colly.HTMLElement) {
		count++
		p.Wikis <- wiki.Link(a.Attr("href"))
		if p.count != 0 && p.count == count {
			close(p.Wikis)
		}
	})
}

func (p *ItemsPage) FetchNextPage() {
	p.c.OnHTML(`a[title="分类:道具"]:last-of-type`, func(e *colly.HTMLElement) {
		if e.Text == "下一页" {
			_ = p.c.Visit(e.Request.AbsoluteURL(e.Attr("href")))
		}
	})
}

// Fetch 可以排除一些不想要抓取的内容
// 不过基本上都是要抓的，没什么好排除的
func (p *ItemsPage) Fetch(exception ...string) *ItemsPage {
	p.FetchTotal()
	p.FetchNextPage()
	p.FetchName()
	p.FetchWiki()
	return p
}

// Visit 开始爬取
func (p *ItemsPage) Visit() {
	_ = p.c.Visit(p.page())
}

func (p *ItemsPage) ReceiveData(wg *sync.WaitGroup, total *int, names, wikis *[]string) {
	go func() {
		wg.Add(1)
		defer wg.Done()
		*total = <-p.Total
	}()
	go func() {
		wg.Add(1)
		defer wg.Done()
		for {
			v, ok := <-p.Names
			if !ok {
				break
			}
			*names = append(*names, v)
		}
	}()
	go func() {
		wg.Add(1)
		defer wg.Done()
		for {
			v, ok := <-p.Wikis
			if !ok {
				break
			}
			*wikis = append(*wikis, v)
		}
	}()
}
