// Package database 数据库操作
package database

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
)

// DB 对象
var DB *gorm.DB

// SQLDB 通用数据库接口
var SQLDB *sql.DB

// Connect 连接数据库
func Connect(dbConfig gorm.Dialector, gormConfig gorm.Config) error {
	// 使用 gorm.Open 连接数据库
	var err error
	if DB, err = gorm.Open(dbConfig, &gormConfig); err != nil {
		return err
	}
	if SQLDB, err = DB.DB(); err != nil {
		return err
	}
	return nil
}

func DNS(user, pass, addr, port, db, char, loc string) string {
	return fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?parseTime=true&multiStatements=true&charset=%v&loc=%v",
		user, pass, addr, port, db, char, loc,
	)
}
