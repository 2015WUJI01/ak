package services

import (
	"github.com/gocolly/colly"
	"main/pkg/progressbar"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ItemsPageData struct {
	Name string
	Wiki string
}

func Step1(datach chan ItemsPageData) {
	defer close(datach)

	bar := progressbar.New("123", 0)

	c := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(3),
		colly.CacheDir("tmp/cache"),
	)
	c.SetRequestTimeout(3 * time.Second)

	c.OnError(func(r *colly.Response, err error) {
		_ = r.Request.Retry()
	})

	c.OnResponse(func(r *colly.Response) {
		if bar.GetMax() == 0 {
			arr := regexp.MustCompile(`共(.*)个页面`).FindSubmatch(r.Body)
			t, _ := strconv.Atoi(string(arr[1]))
			bar.ChangeMax(t)
		}
	})

	// 一次获取 200 条数据
	c.OnHTML(".mw-category-group ul li a", func(e *colly.HTMLElement) {
		_ = bar.Add(1)
		datach <- ItemsPageData{
			Name: strings.TrimSpace(e.Text),
			Wiki: e.Request.AbsoluteURL(e.Attr("href")),
		}
	})

	// 获取下一页链接
	c.OnHTML(`a[title="分类:道具"]:last-of-type`, func(e *colly.HTMLElement) {
		url := e.Request.AbsoluteURL(e.Attr("href"))
		if e.Text == "下一页" {
			_ = c.Visit(url)
		}
	})

	_ = c.Visit("https://prts.Wiki/w/分类:道具")
	c.Wait()
}
