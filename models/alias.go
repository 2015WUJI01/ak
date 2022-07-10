package models

type Alias struct {
	Name  string `json:"name" gorm:"primaryKey;column:name"`
	Alias string `json:"alias" gorm:"primaryKey;column:alias"`
	Type  int    `json:"type" gorm:"primaryKey;column:type"`
}

func (a Alias) TableName() string {
	return "aliases"
}
