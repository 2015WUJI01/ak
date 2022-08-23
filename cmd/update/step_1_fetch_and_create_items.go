package update

import (
	"ak/models"
	"ak/pkg/progressbar"
	repo "ak/repositories"
	"ak/wiki/pages"
	"github.com/gocolly/colly"
	"sync"
)

// Step1 抓取一次页面，获取所有的道具，如果有不存在的道具，则创建这个道具
func Step1() {
	c := colly.NewCollector(colly.Async(false), colly.CacheDir("tmp/cache"))
	p := pages.NewItemsPage(c,
		pages.SetItemsPageProgressBar(progressbar.New("[Step.1] 获取 items 数据", 0)),
	).Fetch()

	wg := &sync.WaitGroup{}
	var total int
	var names, wikis []string
	p.ReceiveData(wg, &total, &names, &wikis)
	p.Visit()
	wg.Wait()

	var items []models.Item
	for i := 0; i < len(names); i++ {
		items = append(items, models.Item{Name: names[i], Wiki: wikis[i]})
	}
	repo.CreateOrUpdateItems([]string{"name", "wiki"}, items)
}
