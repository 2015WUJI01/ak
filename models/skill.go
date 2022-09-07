package models

import "strings"

type Skill struct {
	OprID   int    `json:"opr_id" gorm:"primaryKey;column:opr_id"`
	OprName string `json:"opr_name" gorm:"column:opr_name"`
	Order   int    `json:"order" gorm:"primaryKey;column:order"`
	Name    string `json:"name" gorm:"column:name"`
	Icon    string `json:"icon" gorm:"column:icon"`
	Restore string `json:"restore" gorm:"column:restore"`
	Active  string `json:"active" gorm:"column:active"`
}

func (s *Skill) TableName() string {
	return "skills"
}

func (s *Skill) Trigger() string {
	var restore, active string
	switch s.Restore {
	case "攻击回复":
		restore = "攻回"
	case "自动回复":
		restore = "自回"
	case "被动":
		restore = "被动"
	default:
		restore = s.Restore
	}
	switch s.Active {
	case "自动触发":
		active = "自动"
	case "手动触发":
		active = "手动"
	case "":
		active = ""
	default:
		active = s.Active
	}
	return strings.TrimSpace(strings.Join([]string{active, restore}, " "))
}
