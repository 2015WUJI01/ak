package common

import (
	"fmt"
	"strings"
	"time"
)

type LastMod string

func (l LastMod) Time() time.Time {
	return parseTime(string(l))
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
