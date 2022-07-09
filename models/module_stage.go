package models

type ModuleStage struct {
	OprID   int    `json:"opr_id" gorm:"column:opr_id"`
	OprName string `json:"opr_name" gorm:"primaryKey;column:opr_name"`

	ModuleName  string `json:"module_name" gorm:"primaryKey;column:module_name"`
	ModuleOrder int    `json:"module_order" gorm:"column:module_order"`

	Stage int `json:"stage" gorm:"primaryKey;column:stage"`

	BasicInfo   string `json:"basic_info" gorm:"column:basic_info"`
	Attribution string `json:"attribution" gorm:"column:attribution"`
}

func (ms ModuleStage) TableName() string {
	return "module_stages"
}
