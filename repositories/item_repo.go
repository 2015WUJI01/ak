package repo

import (
	"gorm.io/gorm/clause"
	"main/database"
	"main/models"
	"strings"
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

// FindItemByName 通过 name 精准查询
func FindItemByName(name string) (item models.Item, ok bool) {
	database.DB.Where("lower(name) = ?", strings.ToLower(name)).First(&item)
	if item.Name != "" {
		return item, true
	}
	return models.Item{}, false
}

// FindItemByAlias 通过别名查找最匹配的道具
func FindItemByAlias(alias string) (item models.Item, ok bool) {
	var a models.Alias
	// 完全匹配
	database.DB.Where("lower(alias) = ?", strings.ToLower(alias)).
		Where("type", models.ItemAliasType).First(&a)
	if a.Name != "" {
		return FindItemByName(a.Name)
	}

	// 前缀匹配
	database.DB.Where("lower(alias) like ?", strings.ToLower(alias)+"%").
		Where("type", models.ItemAliasType).First(&a)
	if a.Name != "" {
		return FindItemByName(a.Name)
	}

	// 后缀匹配
	database.DB.Where("lower(alias) like ?", "%"+strings.ToLower(alias)).
		Where("type", models.ItemAliasType).First(&a)
	if a.Name != "" {
		return FindItemByName(a.Name)
	}

	// 中间包含
	database.DB.Where("lower(alias) like ?", "%"+strings.ToLower(alias)+"%").
		Where("type", models.ItemAliasType).First(&a)
	if a.Name != "" {
		return FindItemByName(a.Name)
	}
	return models.Item{}, false
}
