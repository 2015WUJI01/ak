package services

import (
	"ak/models"
	"ak/pkg/progressbar"
	"fmt"
	"sync"
)

func FetchStep1() []models.Item {

	// data channels
	dataChan := make(chan ItemsPageData, 600)
	total := make(chan int, 1)

	bar := progressbar.New(fmt.Sprintf("%-25s", "Step1. Fetch items"), 0)

	var items []models.Item
	wg := sync.WaitGroup{}
	go func() {
		for {
			select {
			case t := <-total:
				bar.ChangeMax(t)
				wg.Add(t)
			case data, ok := <-dataChan:
				if !ok {
					break
				}
				items = append(items, models.Item{Name: data.Name, Wiki: data.Wiki})
				bar.Add(1)
				wg.Done()
			}
		}
	}()

	Step1(dataChan, total)
	wg.Wait()
	return items
}
