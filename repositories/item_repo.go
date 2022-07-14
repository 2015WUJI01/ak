package repo

import (
	"gorm.io/gorm/clause"
	"main/database"
	"main/models"
)

// CreateOrUpdateItems 批量创建或更新数据
func CreateOrUpdateItems(cols []string, items []models.Item) {
	database.DB.Select(cols).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns(cols),
	}).Create(&items)
}
