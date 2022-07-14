package tests

import (
	"github.com/stretchr/testify/assert"
	_ "main/logger"
	"main/models"
	"testing"
)

func TestItemFetchWiki(t *testing.T) {
	var tests = []struct {
		item models.Item
		wiki string
	}{
		{models.Item{Name: "龙门币"}, "https://prts.wiki/w/%E9%BE%99%E9%97%A8%E5%B8%81"},
		{models.Item{Name: "令的信物"}, "https://prts.wiki/w/%E4%BB%A4%E7%9A%84%E4%BF%A1%E7%89%A9"},
		{models.Item{Name: "信用"}, "https://prts.wiki/w/%E4%BF%A1%E7%94%A8"},
		{models.Item{Name: "异铁块"}, "https://prts.wiki/w/%E5%BC%82%E9%93%81%E5%9D%97"},
	}
	for _, tt := range tests {
		t.Run(tt.item.Name, func(t *testing.T) {
			tt.item.FreshWiki()
			assert.Equal(t, tt.wiki, tt.item.Wiki)
		})
	}
}

func TestItemFetchGroup(t *testing.T) {
	var tests = []struct {
		item  models.Item
		group string
	}{
		{models.Item{Name: "龙门币"}, "消耗道具"},
		{models.Item{Name: "先锋芯片组"}, "芯片"},
		{models.Item{Name: "初级作战记录"}, "作战记录"},
		{models.Item{Name: "全新装置"}, "材料"},
	}
	for _, tt := range tests {
		t.Run(tt.item.Name, func(t *testing.T) {
			tt.item.FreshGroup()
			assert.Equal(t, tt.group, tt.item.Group)
		})
	}
}
