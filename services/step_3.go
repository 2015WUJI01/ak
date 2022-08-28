package services

import (
	"ak/models"
	"ak/pkg/progressbar"
	"sync"
)

func FetchStep3() []models.Item {
	dataChan := make(chan models.Item, 300) // data channels
	bar := progressbar.New("Step3. fetching and filling in item group and type....", -1)

	var newItems []models.Item
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for item := range dataChan {
			newItems = append(newItems, item)
			_ = bar.Add(1)
		}
		_ = bar.Finish()
		wg.Done()
	}()
	Step3(dataChan)
	wg.Wait()
	return newItems
}
