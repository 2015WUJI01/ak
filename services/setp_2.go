package services

import (
	"ak/models"
	"ak/pkg/progressbar"
	"fmt"
	"sync"
)

func FetchStep2(items []models.Item) []models.Item {
	// data channels
	bar := progressbar.New(fmt.Sprintf("%-25s", "Step2. Fill in items"), len(items))
	dataChan := make(chan models.Item, 200)

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
