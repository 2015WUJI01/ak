package handlers

import (
	"gorm.io/gorm/clause"
	"main/database"
	"main/models"
	"main/pkg/logger"
	"main/scheduler"
)

func Analysis(c *scheduler.Context) {
	var token string
	var ok bool
	if token, ok = userAkToken(c.GetSenderId()); !ok {
		// 判断第二个参数
		if c.PretreatedMessage == "" {
			// 无附加的 token 参数
			_, _ = c.Reply("未绑定 token，请使用「抽卡 <token>」进行绑定。\ntoken 获取步骤：\n1.登录鹰角官网（https://as.hypergryph.com）并登录；\n2.访问鹰角官方 API（https://as.hypergryph.com/user/info/v1/token_by_cookie）并复制其中的 token 参数，并替换 <token>")
			return
		} else if len(c.PretreatedMessage) == 24 {
			// 传入了第二个参数
			saveUserAkToken(c.GetSenderId(), token)
			_, _ = c.Reply("绑定 token 成功")
		} else {
			logger.Warn(token)
			_, _ = c.Reply("请检查传入的 token 格式是否正确")
			return
		}
	}

}

func userAkToken(qq int64) (string, bool) {
	user := models.User{}
	database.DB.Where("qq = ?", qq).First(&user)
	return user.AkToken, user.AkToken == ""
}

func saveUserAkToken(qq int64, akToken string) {
	database.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&models.User{
		QQ:      qq,
		AkToken: akToken,
	})
}
