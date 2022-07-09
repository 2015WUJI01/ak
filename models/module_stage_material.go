package models

type ModuleStageMaterial struct {

	// 干员
	OprID   int    `json:"opr_id" gorm:"column:opr_id"`
	OprName string `json:"opr_name" gorm:"primaryKey;column:opr_name"`

	// 模组
	ModuleName  string `json:"module_name" gorm:"primaryKey;column:module_name"`
	ModuleOrder int    `json:"module_order" gorm:"column:module_order"`

	// 等级
	ToStage int `json:"to_stage" gorm:"primaryKey;column:to_stage"`

	// 材料
	ItemName string `json:"item_name;" gorm:"primaryKey;column:item_name"`
	Amount   int    `json:"amount" gorm:"column:amount"`
}

func (msm ModuleStageMaterial) TableName() string {
	return "module_stage_materials"
}
