package handlers

import "main/scheduler"

var InsMap = make(map[int64]Confirmation)

type Confirmation struct {
	DoFunc func()
	Cancel func()
}

func Confirm(c *scheduler.Context) {
	if C, ok := InsMap[c.GetSenderId()]; ok {
		C.DoFunc()
		delete(InsMap, c.GetSenderId())
	} else {
		_, _ = c.Reply("暂未找到需要执行的指令")
	}
}

func Cancel(c *scheduler.Context) {
	if C, ok := InsMap[c.GetSenderId()]; ok {
		C.Cancel()
		delete(InsMap, c.GetSenderId())
	} else {
		_, _ = c.Reply("暂未找到需要执行的指令")
	}
}
