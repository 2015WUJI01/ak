package main

import (
	"main/cmd"
)

func main() {
	_ = cmd.Execute()
	// t := time.Now()
	// var names, wikis []string
	// c := colly.NewCollector(colly.Async(true), colly.CacheDir("tmp/cache"))
	// p := pages.NewItemsPage(c).Bind(pages.ITEMS_ALL)
	// done := make(chan bool)
	// go func() {
	// 	bar := progressbar.New("[Step.1] [2/2] 将 itemspage 数据导入数据库", 0)
	// 	const (
	// 		close_name = 1 << iota
	// 		close_wiki
	// 	)
	// 	closed := 0
	// 	for {
	// 		if closed == close_name^close_wiki {
	// 			var items []models.Item
	// 			for i := 0; i < len(names); i++ {
	// 				items = append(items, models.Item{Name: names[i], Wiki: wikis[i]})
	// 				// repo.CreateOrUpdateItem(models.Item{Name: names[i], Wiki: wikis[i]}, []string{"name", "wiki"}...)
	// 			}
	// 			logger.Infof("Step1. 共获取到 %d 条道具数据", len(items))
	// 			repo.CreateOrUpdateItems([]string{"name", "wiki"}, items)
	// 			done <- true
	// 			return
	// 		}
	// 		select {
	// 		case t := <-p.Total:
	// 			bar.ChangeMax(t)
	// 		case n, ok := <-p.Names:
	// 			if ok {
	// 				names = append(names, n)
	// 				bar.Add(1)
	// 			} else {
	// 				closed = closed ^ close_name
	// 			}
	// 		case w, ok := <-p.Wikis:
	// 			if ok {
	// 				wikis = append(wikis, w)
	// 			} else {
	// 				closed = closed ^ close_wiki
	// 			}
	// 		}
	// 	}
	// }()
	// p.Visit()
	// c.Wait()
	// close(p.Names)
	// close(p.Wikis)
	// <-done
	// logger.Infof("cost: %s", time.Now().Sub(t).String())
	// logger.Debug(names, wikis)
}
