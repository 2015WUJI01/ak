package wiki

import (
	"github.com/gocolly/colly"
	"log"
	"time"
)

func NewCollector(timeout ...time.Duration) *colly.Collector {
	if len(timeout) == 0 {
		timeout = append(timeout, 5*time.Second)
	}
	c := colly.NewCollector()
	c.SetRequestTimeout(timeout[0])
	c.OnError(func(r *colly.Response, err error) {
		log.Println("err: ", err)
		_ = r.Request.Retry()
	})
	return c
}

func NewAsyncCollector(parallelism int, delay time.Duration, timeout ...time.Duration) *colly.Collector {
	if len(timeout) == 0 {
		timeout = append(timeout, 5*time.Second)
	}
	c := colly.NewCollector(colly.Async(true))
	c.SetRequestTimeout(timeout[0])
	c.OnError(func(r *colly.Response, err error) {
		log.Println("err: ", err)
		_ = r.Request.Retry()
	})
	_ = c.Limit(&colly.LimitRule{
		DomainGlob:  "*prts.wiki*",            // limit 规则只会在指定的 domain 生效，所以必须配置
		Parallelism: parallelism,              // 同时抓取的协程数
		Delay:       delay * time.Millisecond, // 抓取时延
	})
	return c
}

func Link(uri string) string {
	if uri == "" {
		return ""
	}
	return "https://prts.wiki" + uri
}
