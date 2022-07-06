package handlers

import (
	"main/scheduler"
)

// Alias 添加、删除干员别名
// alias 241
func Alias(c *scheduler.Context) {
	// 解析参数
	// c.PretreatedMessage
	argNotAllowed := true
	// 参数不符合规则，发送提示
	if argNotAllowed {
		sendAliasHelperMsg(c)
		return
	}
	// 参数合法，开始处理增加或删除别名

	// 正则提出干员uuid或名称，找到 opr 对象

	// 正则提出新增的干员名称

	// 正则提出删除的干员名称

	// 开始执行 db 操作

	// 操作完成，提示成功

	// 将最近的别名导出为新的 json 文件（opr_alias_202202022222.json），并保持指定个文件
}

func sendAliasHelperMsg(c *scheduler.Context) {
	msg := "【指令「别名」的说明】" +
		"\n\n功能 1：查询干员的别名" +
		"\n「/别名 <干员>」" +
		"\n例：/别名 令" +
		"\n\n功能 2：新增或删除别名" +
		"\n「/别名 <干员> +/-<别名>」" +
		"\n例：/别名 浊心斯卡蒂 +红蒂" +
		"\n例：/别名 阿米娅 -驴" +
		"\n\nP.S." +
		"\n1. <干员>可以用干员编号代替，例如「/别名 201 +红蒂」" +
		"\n2. <干员>也可以用已有的别名代替，例如「/别名 红蒂 +浊蒂」" +
		"\n3. +表示新增别名，-表示删除别名" +
		"\n4. 可以同时新增或删除多个别名，但每个别名前都需要+/-进行标识，例如「/别名 浊心斯卡蒂 +红蒂 +浊蒂 +蒂蒂 -弟弟」" +
		"\n5. 可以使用英文alias或a代替关键字「别名」，例如「/a 令」"
	_, _ = c.Reply(msg)
}
