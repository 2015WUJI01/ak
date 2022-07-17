package update

import (
	"fmt"
	"main/database"
	"main/logger"
	"main/models"
	repo "main/repositories"
	"main/wiki"
	"strings"
	"time"
)

// Step2 从数据库中获取所有道具，对每个道具进行数据更新
func Step2() {
	fmt.Println("逐条采集道具信息中...")

	var items []models.Item
	_ = database.DB.Find(&items)
	minfo := wiki.FetchItemInfo(func() []string {
		var names []string
		for _, item := range items {
			names = append(names, item.Name)
		}
		return names
	}()...)
	for i, item := range items {
		items[i].Image = minfo[item.Name]["name"].(string)
		items[i].WikiShort = minfo[item.Name]["wikishort"].(string)
		items[i].UpdatedAt = minfo[item.Name]["updatedat"].(time.Time)
	}
	repo.CreateOrUpdateItems([]string{"name", "image", "wiki_short", "updated_at"}, items)
	logger.Infof("Step2. 道具基本信息更新完成")
}

// parseTime 解析源码中的时间字符串
// 原文本为 "此页面最后编辑于2022年5月22日 (星期日) 12:32。"
// 无法直接作为 go 时间解析规则，所以尝试替换 7 次，总有一种能够解析成功
func parseTime(timeStr string) time.Time {
	weeks := []string{"日", "一", "二", "三", "四", "五", "六"}
	var err error
	var t time.Time
	for i := 0; i <= 7; i++ {
		layout := fmt.Sprintf("此页面最后编辑于2006年1月2日 (星期%s) 15:04 -0700", weeks[i])
		t, err = time.Parse(layout, strings.ReplaceAll(strings.TrimSpace(timeStr), "。", " +0800"))
		if err == nil {
			break
		}
	}
	return t
}
