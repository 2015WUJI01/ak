package ri

import (
	"fmt"
	"github.com/gocolly/colly"
	"main/pkg/logger"
	"time"
)

type ModelFetcher struct {
	f *Fetcher `gorm:"-:all"`
}

// Fetcher 单例 fetcher
func (m *ModelFetcher) Fetcher() *Fetcher {
	if m.f == nil {
		m.f = NewFetcher().AutoRetry(0)
	}
	return m.f
}

type Fetcher struct {
	colly        *colly.Collector
	RetryCounter int
}

func NewFetcher() *Fetcher {
	return &Fetcher{
		colly: colly.NewCollector(),
	}
}

// Fetcher 配置相关

// Colly 返回 colly 对象
func (f *Fetcher) Colly() *colly.Collector {
	return f.colly
}

// AutoRetry 配置自动重试的时间
func (f *Fetcher) AutoRetry(t time.Duration) *Fetcher {
	f.colly.OnError(func(r *colly.Response, err error) {
		time.Sleep(t)
		f.RetryCounter++
		logger.Debugf(fmt.Sprintf("访问 %s 失败，正在第 %d 次重试", r.Request.URL, f.RetryCounter))
		_ = r.Request.Retry()
	})
	return f
}

// 覆写 Colly 相关函数

func (f *Fetcher) OnHTML(selector string, callback colly.HTMLCallback) *Fetcher {
	f.colly.OnHTML(selector, callback)
	return f
}

func (f *Fetcher) OnResponse(callback colly.ResponseCallback) *Fetcher {
	f.colly.OnResponse(callback)
	return f
}

func (f *Fetcher) OnError(callback colly.ErrorCallback) *Fetcher {
	f.colly.OnError(callback)
	return f
}

func (f *Fetcher) Visit(url string) error {
	return f.colly.Visit(url)
}
