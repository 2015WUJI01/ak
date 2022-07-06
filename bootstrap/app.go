package bootstrap

import (
	"main/conf"
	"main/pkg/logger"
)

func Initialize() {
	var err error

	// 初始化配置
	if err = conf.InitConfig("./conf/config.ini"); err != nil {
		logger.Warnf("配置文件不存在", err.Error())
	} else {
		logger.Debug("配置初始化完成")
	}

	// 初始化日志
	if conf.Cfg.App.Debug == false {
		logger.SetLogger(logger.LogConfig{
			Level: logger.InfoLevel,
		})
	}

	// 初始化数据库
	if err = LoadDatabase(); err != nil {
		logger.Warnf("数据库初始化异常", err.Error())
	} else {
		logger.Debug("数据库初始化完成")
	}

}
