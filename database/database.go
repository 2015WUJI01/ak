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

	_ = DB.AutoMigrate(
		&models.Item{}, &models.Operator{}, &models.Alias{},
		&models.Skill{}, &models.SkillLevel{}, &models.SkillLevelMaterial{},
		&models.Module{}, &models.ModuleStage{}, &models.ModuleStageMaterial{},
	)
}

func NewSQLiteDB(dsn string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: NewGormLogger("logs/sql.log"),
	})
}

func NewMySQL(host, user, pass, dbname string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", user, pass, host, dbname)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: NewGormLogger("logs/sql.log"),
	})
}

func NewGormLogger(filename string) logger.Interface {
	file, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return logger.New(
		log.New(file, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			Colorful:                  false,       // 禁用彩色打印
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
		},
	)
}
