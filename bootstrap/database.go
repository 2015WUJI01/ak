package bootstrap

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"main/database"
	"os"
	"time"
)

// LoadDatabase 初始化加载 Database
func LoadDatabase() error {
	// 连接数据库
	var err error
	err = database.Connect(dbConfig(), gormConfig())

	// 数据库配置
	database.SQLDB.SetMaxIdleConns(20)                  // 连接池
	database.SQLDB.SetMaxOpenConns(100)                 // 最大连接数
	database.SQLDB.SetConnMaxLifetime(30 * time.Minute) // 连接可复用的最大时间

	// 自动迁移
	autoMigrate()
	return err
}

// 配置 dbConfig
func dbConfig() gorm.Dialector {
	return sqlite.Open("arknights.db")
}

// 配置 gormConfig
func gormConfig() gorm.Config {
	return gorm.Config{
		PrepareStmt: true,
		// 自定义日志
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				Colorful:                  true,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: false,
			},
		),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: false,
			NameReplacer:  nil,
			NoLowerCase:   false,
		},
	}
}

// 配置自动迁移
func autoMigrate() {
	// database.DB.AutoMigrate(&models.User{})
}
