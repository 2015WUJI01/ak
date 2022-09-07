package services

import (
	"ak/models"
	"ak/pkg/progressbar"
	"sync"
)

func FetchStep2(items []models.Item) []models.Item {
	// data channels
	dataChan := make(chan models.Item, 200)

	bar := progressbar.New("Step2. fetching and update items basic info...", len(items))

	var newItems []models.Item
	wg := sync.WaitGroup{}
	wg.Add(len(items))
	go func() {
		for {
			select {
			case item, ok := <-dataChan:
				if !ok {
					break
				}
				newItems = append(newItems, item)
				bar.Add(1)
				wg.Done()
			}
		}
	}()
	Step2(items, dataChan)
	wg.Wait()
	return items
}
