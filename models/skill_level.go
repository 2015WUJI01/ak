package models

type SkillLevel struct {
	OprID   int    `json:"opr_id" gorm:"column:opr_id"`
	OprName string `json:"opr_name" gorm:"primaryKey;column:opr_name"`
	Order   int    `json:"order" gorm:"primaryKey;column:order"`
	Level   int    `json:"level" gorm:"primaryKey;column:level"`
	OriPt   int    `json:"ori_pt" gorm:"column:ori_pt"`
	CostPt  int    `json:"cost_pt" gorm:"column:cost_pt"`
	Last    int    `json:"last" gorm:"column:last"`
	Comment string `json:"comment" gorm:"column:comment"`
}

func (sl SkillLevel) TableName() string {
	return "skill_levels"
}
