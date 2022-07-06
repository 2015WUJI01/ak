package conf

import (
	"gopkg.in/ini.v1"
	"gorm.io/gorm/logger"
)

var Cfg Config

type Config struct {
	App AppConfig
	DB  DBConfig
}

type AppConfig struct {
	Debug bool
	Port  int
}

type DBConfig struct {
	// User     string
	// Pass     string
	// Addr     string
	// Port     string
	// Name     string
	LogLevel logger.LogLevel
}

func InitConfig(path string) error {
	file, err := ini.Load(path)
	if err != nil {
		return err
	}
	Cfg.App.LoadConfig(file)
	Cfg.DB.LoadConfig(file)
	return nil
}

func (c *AppConfig) LoadConfig(file *ini.File) {
	c.Debug = file.Section("app").Key("debug").MustBool()
	c.Port = file.Section("app").Key("port").MustInt()
}

func (c *DBConfig) LoadConfig(file *ini.File) {
	// c.User = file.Section("database").Key("user").String()
	// c.Pass = file.Section("database").Key("pass").String()
	// c.Addr = file.Section("database").Key("addr").String()
	// c.Port = file.Section("database").Key("port").String()
	// c.Name = file.Section("database").Key("name").String()
	c.LogLevel = logger.LogLevel(file.Section("database").Key("log_level").MustInt())
}
