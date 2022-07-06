package handlers

import (
	"main/ri"
	"main/scheduler"
)

func SkillUpgrade(c *scheduler.Context) {
	var op ri.Operator
	var ok bool
	_, _ = c.Reply(c.PretreatedMessage)
	if op, ok = ri.StaffExisted(c.PretreatedMessage); !ok {
		_, _ = c.Reply("查询不到，请检查干员名称" + c.PretreatedMessage)
		return
	}
	err := op.Fetch()
	if err != nil {
		_, _ = c.Reply("干员数据拉取失败")
		return
	}
	_, _ = c.Reply(op.SkillMaterialGroupsToMsg())
}
