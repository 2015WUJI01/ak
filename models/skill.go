package models

type Skill struct {
	OprID   int    `json:"opr_id" gorm:"primaryKey;column:opr_id"`
	OprName string `json:"opr_name" gorm:"column:opr_name"`
	Order   int    `json:"order" gorm:"primaryKey;column:order"`
	Name    string `json:"name" gorm:"column:name"`
	Icon    string `json:"icon" gorm:"column:icon"`
	Restore string `json:"restore" gorm:"column:restore"`
	Active  string `json:"active" gorm:"column:active"`
}

func (s Skill) TableName() string {
	return "skills"
}
