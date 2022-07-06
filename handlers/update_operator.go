package handlers

import (
	"fmt"
	"main/pkg/help"
	"main/ri"
	"main/scheduler"
	"time"
)

// UpdateOperatorsName 增量更新干员名称
func UpdateOperatorsName(c *scheduler.Context) {
	t := time.Now()
	rowsAffected := ri.OC.UpdateOperatorsName()
	_, _ = c.Reply(fmt.Sprintf(
		"增量更新干员名称已完成，本次新增 %d 个干员，耗时 %v", rowsAffected, help.PrintCostTimeSince(t),
	))
}
