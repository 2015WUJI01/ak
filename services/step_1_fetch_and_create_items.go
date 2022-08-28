package services

import (
	"github.com/2015WUJI01/looog"
	"github.com/gocolly/colly"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ItemsPageData struct {
	Name string
	Wiki string
}

func Step1(dataChan chan ItemsPageData, total chan int) {
	defer close(dataChan)

	c := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(3),
		colly.CacheDir("tmp/cache"),
	)
	c.SetRequestTimeout(3 * time.Second)

	c.OnError(func(r *colly.Response, err error) {
		looog.Warnf("retrying: %s", r.Request.URL)
		_ = r.Request.Retry()
	})

	sendOnce := sync.Once{}
	c.OnResponse(func(r *colly.Response) {
		sendOnce.Do(func() {
			arr := regexp.MustCompile(`共(.*)个页面`).FindSubmatch(r.Body)
			t, _ := strconv.Atoi(string(arr[1]))
			total <- t
		})
	})

	// 一次获取 200 条数据
	c.OnHTML(".mw-category-group ul li a", func(e *colly.HTMLElement) {
		dataChan <- ItemsPageData{
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
