package models

import (
	"fmt"
	"strings"
	"time"
)

type Operator struct {
	ID int `json:"id" gorm:"column:id"`

	Name string `json:"name" gorm:"primaryKey;column:name"`

	// 干员职业名称，例如 "狙击"
	Class string `json:"class" gorm:"column:class;type:string;size:255"`
	// 干员细分职业，例如 "速射狙"
	Subclass string `json:"subclass" gorm:"column:subclass;type:string;size:255"`
	Rarity   int    `json:"rarity" gorm:"column:rarity;type:tinyint;size:2"`

	Roguelike bool `json:"roguelike" gorm:"column:roguelike;type:bool"`

	Wiki      string    `json:"wiki" gorm:"column:wiki"`
	WikiShort string    `json:"wiki_short" gorm:"column:wiki_short"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:false"`
}

func (o Operator) TableName() string {
	return "operators"
}

// Info 标准的 opr 的信息显示
// No.001 6★ 史尔特尔
// 近卫 术战者 / 肉鸽限定干员
// 别名：XXX, XXX
func (o Operator) Info(alias []string) string {
	msg := fmt.Sprintf("No.%03d %d★ %s", o.ID, o.Rarity, o.Name)
	msg += fmt.Sprintf("\n%s %s", o.Class, o.Subclass)
	if o.Roguelike {
		msg += " / 肉鸽限定干员"
	}
	if len(alias) > 0 {
		msg += fmt.Sprintf("\n别名: %s", strings.Join(alias, ", "))
	}
	msg += fmt.Sprintf("\nwiki: %s", o.WikiShort)
	msg += fmt.Sprintf("\nupdated: %s", o.UpdatedAt.Format("2006-01-02 15:04"))
	return msg
}
