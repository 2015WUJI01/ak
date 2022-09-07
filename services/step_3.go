package services

import (
	"ak/models"
	"ak/pkg/progressbar"
	"sync"
)

func FetchStep3() []models.Item {
	// data channels
	dataChan := make(chan models.Item, 300)

	bar := progressbar.New("Step3. fetching and update items group and type....", 300) // 假设只有 300 个

	var newItems []models.Item
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		total := 0
		for item := range dataChan {
			newItems = append(newItems, item)
			total++
			_ = bar.Add(1)
		}
		bar.ChangeMax(total)
		_ = bar.Finish()
		wg.Done()
	}()
	Step3(dataChan)
	wg.Wait()
	return newItems
}
