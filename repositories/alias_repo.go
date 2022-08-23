package repo

import (
	"ak/database"
	"ak/models"
	"strings"
)

func FindItemAlias(name string) []string {
	var alias []models.Alias
	database.DB.Where("lower(name) = ?", strings.ToLower(name)).
		Where("type", models.ItemAliasType).Find(&alias)
	var arr []string
	for _, a := range alias {
		arr = append(arr, a.Alias)
	}
	return arr
}

func FindOprAlias(name string) []string {
	var alias []models.Alias
	database.DB.Where("lower(name) = ?", strings.ToLower(name)).
		Where("type", models.OprAliasType).Find(&alias)
	var arr []string
	for _, a := range alias {
		arr = append(arr, a.Alias)
	}
	return arr
}
