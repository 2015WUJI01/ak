// Package database 数据库操作
package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"main/models"
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
}

func NewSQLiteDB(dsn string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Silent),
	})
}

func NewMySQL(host, user, pass, dbname string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", user, pass, host, dbname)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Silent),
	})
}
