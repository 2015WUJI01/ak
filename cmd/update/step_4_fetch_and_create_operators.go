package update

import (
	"github.com/gocolly/colly"
	"gorm.io/gorm/clause"
	"main/database"
	"main/models"
	"main/pkg/logger"
	"strings"
	"time"
)

// Step4 获取干员所有名称，若不存在则插入到数据库中
func Step4() {
	oprs := fetchAllOperators()
	CreateOrUpdateOperators([]string{"name", "id", "roguelike", "wiki"}, oprs)
	logger.Infof("Step4. 共获取到 %d 条干员数据", len(oprs))
	return
}

func fetchAllOperators() []models.Operator {
	var oprs []models.Operator
	c := colly.NewCollector()
	c.SetRequestTimeout(5 * time.Second)
	c.OnError(func(r *colly.Response, err error) {
		_ = r.Request.Retry()
	})
	// 获取所有的干员名称（以及是否是肉鸽限定）
	c.OnHTML("#mw-content-text table tbody tr td:nth-of-type(1) a", func(a *colly.HTMLElement) {
		oprs = append(oprs, models.Operator{
			ID:        a.Index + 1,
			Name:      strings.TrimSpace(a.Text),
			Roguelike: isRoguelike(strings.TrimSpace(a.Text)),
			Wiki:      Link(a.Attr("href")),
		})
	})
	_ = c.Visit(Link("/w/干员上线时间一览"))
	return oprs
}

// CreateOprIfNotExists 若干员不存在则创建该干员
func CreateOprIfNotExists(opr models.Operator) bool {
	var cnt int64
	database.DB.Model(&models.Operator{}).Where("name", opr.Name).Count(&cnt)
	if cnt == 0 {
		database.DB.Select("id", "name", "roguelike", "wiki").Create(&opr)
		return true
	}
	return false
}

// CreateOrUpdateOperator
// 根据主键进行查询，若存在则更新，若不存在则创建
func CreateOrUpdateOperator(cols []string, opr models.Operator) (created bool) {
	var cnt int64
	database.DB.Model(&models.Operator{}).Where("name", opr.Name).Count(&cnt)
	if cnt > 0 {
		database.DB.Select(cols).Updates(&opr)
		return false
	}
	database.DB.Select(cols).Create(&opr)
	return true
}

// CreateOrUpdateOperators
// 根据主键进行查询，若存在则更新，若不存在则创建
func CreateOrUpdateOperators(cols []string, oprs []models.Operator) {
	database.DB.Select(cols).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns(cols),
	}).Create(&oprs)
}

// 以下干员为肉鸽限定，无法获得，共 9 位
var roguelikes = map[string]struct{}{
	"Touch": {}, "Sharp": {}, "Stormeye": {}, "Pith": {}, "暮落(集成战略)": {},
	"预备干员-术师": {}, "预备干员-近战": {}, "预备干员-狙击": {}, "预备干员-后勤": {},
}

func isRoguelike(name string) bool {
	_, ok := roguelikes[name]
	return ok
}
