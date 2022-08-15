package tests

import (
	"github.com/stretchr/testify/assert"
	_ "main/logger"
	"main/models"
	"testing"
)

func TestFreshAllWiki(t *testing.T) {
	var tests = []struct {
		// in
		items []*models.Item
		// out
		wiki map[string]string
	}{
		// 单个
		{
			[]*models.Item{{Name: "龙门币"}},
			map[string]string{"龙门币": "https://prts.wiki/w/%E9%BE%99%E9%97%A8%E5%B8%81"},
		},
		// 多个
		{
			[]*models.Item{{Name: "令的信物"}, {Name: "信用"}},
			map[string]string{
				"令的信物": "https://prts.wiki/w/%E4%BB%A4%E7%9A%84%E4%BF%A1%E7%89%A9",
				"信用":   "https://prts.wiki/w/%E4%BF%A1%E7%94%A8",
			},
		},
		// 测试全部
		{
			[]*models.Item{},
			map[string]string{},
		},
	}
	for i, tt := range tests {
		tt.items = models.FreshItemWiki(tt.items...)
		for ii, item := range tt.items {
			t.Logf("%d.%d. %+v", i, ii, item)
		}
		// t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
		// 	models.FreshItemWiki(tt.items...)
		// 	for _, item := range tt.items {
		// 单条或多条
		// assert.Equal(t, tt.wiki[item.name], item.wiki)
		// 全部
		// }
		// })
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
