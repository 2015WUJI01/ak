package repo

import (
	"gorm.io/gorm/clause"
	"main/database"
	"main/models"
)

func GetItemsMap() map[string]models.Item {
	var mitems = make(map[string]models.Item)
	var items []models.Item
	_ = database.DB.Find(&items)
	for _, item := range items {
		mitems[item.Name] = item
	}
	return mitems
}

// CreateOrUpdateItems 批量创建或更新数据
func CreateOrUpdateItems(cols []string, items []models.Item) {
	database.DB.Select(cols).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns(cols),
	}).Create(&items)
}

func CreateOrUpdateItemsp(cols []string, items []*models.Item) {
	database.DB.Select(cols).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns(cols),
	}).Create(&items)
}

// CreateOrUpdateItem 批量创建或更新数据
// 传 cols 的时候，需要把 primaryKey 字段同时传入
func CreateOrUpdateItem(item models.Item, cols ...string) {
	database.DB.Select(cols).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns(cols),
	}).Create(&item)
}
