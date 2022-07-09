package models

import "time"

type Item struct {
	Name      string    `json:"name" gorm:"primaryKey;column:name"`
	Group     string    `json:"group" gorm:"column:group"`
	Type      string    `json:"type" gorm:"column:type"`
	Image     string    `json:"image" gorm:"column:image"`
	Wiki      string    `json:"wiki" gorm:"column:wiki"`
	WikiShort string    `json:"wiki_short" gorm:"column:wiki_short"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:false"`
}

func (m Item) TableName() string {
	return "items"
}
