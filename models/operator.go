package models

import "time"

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
