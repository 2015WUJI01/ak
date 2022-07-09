package models

import "fmt"

type SkillLevelMaterial struct {
	OprID   int    `json:"opr_id" gorm:"column:opr_id"`
	OprName string `json:"opr_name" gorm:"primaryKey;column:opr_name"`

	Order int `json:"order" gorm:"primaryKey;column:order"`

	ToLevel int `json:"to_level" gorm:"primaryKey;column:to_level"`

	ItemName string `json:"item_name;" gorm:"primaryKey;column:item_name"`
	Amount   int    `json:"amount" gorm:"column:amount"`
}

func (slm SkillLevelMaterial) TableName() string {
	return "skill_level_materials"
}

func (slm SkillLevelMaterial) Echo() {
	order := []string{"", "1", "2", "3"}[slm.Order]
	// level := []string{"", "", "2", "3", "4", "5", "6", "7", "专精一", "专精二", "专精三"}[sm.ToLevel]
	fmt.Printf("%s 的 %s 技能升到 %v 级需要 %d 个 %s\n", slm.OprName, order, slm.ToLevel, slm.Amount, slm.ItemName)
}
