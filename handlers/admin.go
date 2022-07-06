package handlers

import (
	"fmt"
	"log"
	"main/database"
	"main/pkg/help"
	"main/ri"
	"main/scheduler"
	"time"
)

// ConfirmUpdateItemsData 加载 items 数据
func ConfirmUpdateItemsData(c *scheduler.Context) {
	_, _ = c.Reply("准备从企鹅物流获取 items 数据...")
	t := time.Now()
	err := ri.UpdateItemsData()
	if err != nil {
		log.Println(err.Error())
		_, _ = c.Reply("从企鹅物流获取数据失败")
	}
	_, _ = c.Reply(fmt.Sprintf("从企鹅物流更新数据完成，耗时 %vs", help.SpendSeconds(t)))
}

func UpdateAllOperator(c *scheduler.Context) {
	_, _ = c.Reply("开始更新所有干员数据...")
	startAt := time.Now()
	var ops []ri.Operator
	errMsg := "更新过程出现异常："
	res := database.DB.Find(&ops)
	if res.Error != nil {
		_, _ = c.Reply(errMsg + "\n数据库查询出错 " + res.Error.Error())
		return
	}

	// 轮询更新
	var sucCnt, errCnt int
	for _, op := range ops {
		err := op.Fetch()
		if err != nil {
			errMsg += "\n" + err.Error()
			errCnt += 1
		} else {
			sucCnt += 1
		}
	}
	dura := time.Since(startAt).Seconds()
	_, _ = c.Reply(fmt.Sprintf(
		"更新结束。\n共更新 %d 条记录，耗时 %.2fs。其中成功 %d 条，失败 %d 条",
		sucCnt+errCnt, dura, sucCnt, errCnt))
	if errCnt > 0 {
		_, _ = c.Reply(errMsg)
	}
}
