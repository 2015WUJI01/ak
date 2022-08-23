package spider

import (
	"github.com/gocolly/colly"
	"time"
)

type Spider struct {
	c *colly.Collector
}

func Default() *Spider {
	return New(
		CacheDir("/tmp/cache"),
		AutoRetry(),
	)
}

func New(options ...func(*Spider)) *Spider {
	s := &Spider{
		c: colly.NewCollector(),
	}

	for _, f := range options {
		f(s)
	}
	return s
}

/******** Option Functions ********/

func AutoRetry(beforeRetry ...func(r *colly.Response, err error)) func(*Spider) {
	return func(s *Spider) {
		s.c.OnError(func(r *colly.Response, err error) {
			for _, fn := range beforeRetry {
				fn(r, err)
			}
			_ = r.Request.Retry()
		})
	}
}

func Async(domain string, delay, randomdelay time.Duration, parallelism int) func(*Spider) {
	limit := Limit(domain, delay, randomdelay, parallelism)
	return func(s *Spider) {
		limit(s)
		s.c.Async = true
	}
}

func Limit(domain string, delay, randomdelay time.Duration, parallelism int) func(*Spider) {
	return func(s *Spider) {
		_ = s.c.Limit(&colly.LimitRule{
			DomainRegexp: domain,
			DomainGlob:   domain,
			Delay:        delay,
			RandomDelay:  randomdelay,
			Parallelism:  parallelism,
		})
	}
}

func CacheDir(path string) func(*Spider) {
	return func(s *Spider) {
		s.c.CacheDir = path
	}
}
