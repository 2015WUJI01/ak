package services

import (
	"ak/models"
	"ak/pkg/progressbar"
	"sync"
)

func FetchStep1(items *[]models.Item) {
	// data channels
	dataChan := make(chan ItemsPageData, 600)
	totalChan := make(chan int, 1)

	bar := progressbar.New("Step1. fetching and store items name and wiki...", 0)

	wg := sync.WaitGroup{}
	go func() {
		for {
			select {
			case t := <-totalChan:
				bar.ChangeMax(t)
				wg.Add(t)
			case data, ok := <-dataChan:
				if !ok {
					break
				}
				*items = append(*items, models.Item{Name: data.Name, Wiki: data.Wiki})
				bar.Add(1)
				wg.Done()
			}
		}
	}()

	Step1(dataChan, totalChan)
	wg.Wait()
}
