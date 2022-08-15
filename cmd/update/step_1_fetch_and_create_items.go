package update

import (
	"github.com/gocolly/colly"
	"main/logger"
	"main/models"
	repo "main/repositories"
	"main/wiki/pages"
	"time"
)

// Step1 抓取一次页面，获取所有的道具，如果有不存在的道具，则创建这个道具
func Step1() {
	t := time.Now()
	// 方法一：获取指针后，遍历出 item 值
	// var items []models.Item

	// var names, wikis []string
	c := colly.NewCollector(colly.Async(false), colly.CacheDir("tmp/cache"))
	p := pages.NewItemsPage(c).Bind(pages.ITEMS_ALL)
	done := make(chan bool)
	go CreateOrUpdateItems(p, done)
	logger.Infof("repare: %s", time.Now().Sub(t).String())
	t = time.Now()
	p.Visit()
	c.Wait()
	close(p.Names)
	close(p.Wikis)
	<-done

	// for n, w := range wiki.FetchAllWiki() {
	// 	items = append(items, models.Item{Name: n, Wiki: w})
	// }
	// logger.Infof("Step1. 共获取到 %d 条道具数据", len(items))
	// repo.CreateOrUpdateItems([]string{"name", "wiki"}, items)
	logger.Infof("cost: %s", time.Now().Sub(t).String())
}

func CreateOrUpdateItems(p *pages.ItemsPage, done chan bool) {
	var names, wikis []string

	const (
		close_name = 1 << iota
		close_wiki
	)
	closed := 0
	for {
		if closed == close_name^close_wiki {
			break
		}
		select {
		case <-p.Total:
		case n, ok := <-p.Names:
			if ok {
				names = append(names, n)
			} else {
				closed = closed ^ close_name
			}
		case w, ok := <-p.Wikis:
			if ok {
				wikis = append(wikis, w)
			} else {
				closed = closed ^ close_wiki
			}
		}
	}
	var items []models.Item
	for i := 0; i < len(names); i++ {
		items = append(items, models.Item{Name: names[i], Wiki: wikis[i]})
		// repo.CreateOrUpdateItem(models.Item{Name: names[i], Wiki: wikis[i]}, []string{"name", "wiki"}...)
	}
	// logger.Infof("Step1. 共获取到 %d 条道具数据", len(items))
	repo.CreateOrUpdateItems([]string{"name", "wiki"}, items)
	// logger.Debug(names, wikis)
	done <- true
}
