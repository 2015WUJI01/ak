package models

const (
	ItemAliasType = 1 + iota
	OprAliasType
)

type Alias struct {
	Name  string `json:"name" gorm:"primaryKey;column:name"`
	Alias string `json:"alias" gorm:"primaryKey;column:alias"`

	// Type 别名类型 1=item 2=opr
	Type int `json:"type" gorm:"primaryKey;column:type"`
}

func (a Alias) TableName() string {
	return "aliases"
}
