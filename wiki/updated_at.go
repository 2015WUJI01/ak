package wiki

import (
	"fmt"
	"strings"
	"time"
)

type UpdatedAtStr string

func (us UpdatedAtStr) Format(layout string) string {
	return us.AsTime().Format(layout)
}

// AsTime 解析源码中的时间字符串为 time.Time 格式
// 原文本为 "此页面最后编辑于2022年5月22日 (星期日) 12:32。"
// 无法直接作为 go 时间解析规则，所以尝试替换 7 次，总有一种能够解析成功
func (us UpdatedAtStr) AsTime() time.Time {
	weeks := []string{"日", "一", "二", "三", "四", "五", "六"}
	timestr := strings.ReplaceAll(strings.TrimSpace(string(us)), "。", " +0800")
	for i := 0; i <= 7; i++ {
		layout := fmt.Sprintf("此页面最后编辑于2006年1月2日 (星期%s) 15:04 -0700", weeks[i])
		t, err := time.Parse(layout, timestr)
		if err == nil {
			return t
		}
	}
	// Warn: 解析时间失败
	return time.Time{}
}
