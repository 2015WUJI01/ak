package repo

import (
	"ak/database"
	"ak/models"
	"strings"
)

func FindOprByID(id string) (opr models.Operator, ok bool) {
	database.DB.Where("id", id).First(&opr)
	if opr.ID != 0 {
		return opr, true
	}
	return models.Operator{}, false
}

// FindOprByName 通过 name 精准查询
func FindOprByName(name string) (opr models.Operator, ok bool) {
	database.DB.Where("lower(name) = ?", strings.ToLower(name)).First(&opr)
	if opr.ID != 0 {
		return opr, true
	}
	return models.Operator{}, false
}

// FindOprByAlias 通过别名查找最匹配的干员
func FindOprByAlias(alias string) (opr models.Operator, ok bool) {
	var a models.Alias
	// 完全匹配
	database.DB.Where("lower(alias) = ?", strings.ToLower(alias)).
		Where("type", models.OprAliasType).First(&a)
	if a.Name != "" {
		return FindOprByName(a.Name)
	}

	// 前缀匹配
	database.DB.Where("lower(alias) like ?", strings.ToLower(alias)+"%").
		Where("type", models.OprAliasType).First(&a)
	if a.Name != "" {
		return FindOprByName(a.Name)
	}

	// 后缀匹配
	database.DB.Where("lower(alias) like ?", "%"+strings.ToLower(alias)).
		Where("type", models.OprAliasType).First(&a)
	if a.Name != "" {
		return FindOprByName(a.Name)
	}

	// 中间包含
	database.DB.Where("lower(alias) like ?", "%"+strings.ToLower(alias)+"%").
		Where("type", models.OprAliasType).First(&a)
	if a.Name != "" {
		return FindOprByName(a.Name)
	}
	return models.Operator{}, false
}
