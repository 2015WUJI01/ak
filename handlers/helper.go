package handlers

import (
	"main/scheduler"
)

func sendConfirmMsg(c *scheduler.Context) {
	_, _ = c.Reply("该指令需要您的再次确认方可执行，请输入「确认」执行，「取消」取消")
}
