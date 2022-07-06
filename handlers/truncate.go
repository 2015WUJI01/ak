package handlers

import (
	"errors"
	"fmt"
	"main/conf"
	"main/database"
	"main/scheduler"
	"strings"
)

type truncateController struct {
	c      *scheduler.Context
	tables []string
}

func (tc truncateController) DoTruncate() {
	// 判断是否需要清空指定表
	var tbs []string
	for _, t := range strings.Split(tc.c.PretreatedMessage, " ") {
		name := strings.TrimSpace(t)
		if name != "" {
			tbs = append(tbs, name)
		}
	}

	if len(tbs) > 0 {
		for _, t := range tbs {
			if err := tc.truncateTableByTableName(t); err != nil {
				_, _ = tc.c.Reply("清空表 " + t + " 过程中发生错误：" + err.Error())
			} else {
				_, _ = tc.c.Reply("已清空表 " + t)
			}
		}
	} else {
		tc.truncateAllTables()
		_, _ = tc.c.Reply("已清空所有表数据")
	}

}

func (tc truncateController) Cancle() {
	if _, ok := InsMap[tc.c.GetSenderId()]; ok {
		delete(InsMap, tc.c.GetSenderId())
	}
	_, _ = tc.c.Reply("已取消指令「" + tc.c.GetRawMessage() + "」")
}

func (tc truncateController) truncateAllTables() {
	tbs := tables()
	for _, t := range tbs {
		database.DB.Exec(fmt.Sprintf("TRUNCATE TABLE `%s`;", t))
	}
}

func (tc truncateController) truncateTableByTableName(name string) error {
	tbs := tables()
	for _, t := range tbs {
		if name == t {
			sql := fmt.Sprintf("TRUNCATE TABLE %s", name)
			return database.DB.Exec(sql).Error
		}
	}
	return errors.New("该表不存在或暂不允许被清空")
}

func tables() []string {
	var tbs []string
	database.DB.Raw(fmt.Sprintf("show tables from %s;", conf.Cfg.DB.Name)).Scan(&tbs)
	return tbs
}

// ConfirmTruncate 确认清空表
func ConfirmTruncate(c *scheduler.Context) {
	tc := truncateController{c: c}
	sendConfirmMsg(c)
	InsMap[c.GetSenderId()] = Confirmation{
		DoFunc: tc.DoTruncate,
		Cancel: tc.Cancle,
	}
}
