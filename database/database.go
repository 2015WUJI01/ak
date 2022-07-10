// Package database 数据库操作
package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"main/models"
	"os"
	"time"
)

// DB 对象
var DB *gorm.DB

func init() {
	var err error
	DB, err = NewSQLiteDB("arknights.db")
	// DB, err = NewMySQL("127.0.0.1", "root", "", "test")
	if err != nil {
		fmt.Println("数据库初始化异常", err.Error())
	}

	_ = DB.AutoMigrate(&models.Item{})
	_ = DB.AutoMigrate(&models.Operator{})
	_ = DB.AutoMigrate(&models.Skill{})
	_ = DB.AutoMigrate(&models.SkillLevel{})
	_ = DB.AutoMigrate(&models.SkillLevelMaterial{})
	_ = DB.AutoMigrate(&models.Module{})
	_ = DB.AutoMigrate(&models.ModuleStage{})
	_ = DB.AutoMigrate(&models.ModuleStageMaterial{})
	_ = DB.AutoMigrate(&models.Alias{})
}

func NewSQLiteDB(dsn string) (*gorm.DB, error) {
	file, _ := os.OpenFile("sql.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(file, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
			logger.Config{
				SlowThreshold: time.Second, // 慢 SQL 阈值
				LogLevel:      logger.Info, // 日志级别
				Colorful:      false,       // 禁用彩色打印
				// 忽略ErrRecordNotFound（记录未找到）错误
				IgnoreRecordNotFoundError: true,
			},
		),
	})
}

func NewMySQL(host, user, pass, dbname string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", user, pass, host, dbname)
	file, _ := os.OpenFile("sql.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(file, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
			logger.Config{
				SlowThreshold: time.Second, // 慢 SQL 阈值
				LogLevel:      logger.Info, // 日志级别
				Colorful:      false,       // 禁用彩色打印
				// 忽略ErrRecordNotFound（记录未找到）错误
				IgnoreRecordNotFoundError: true,
			},
		),
	})
}
