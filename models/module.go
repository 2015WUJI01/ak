package models

type Module struct {
	OprID   int    `json:"opr_id" gorm:"column:opr_id"`
	OprName string `json:"opr_name" gorm:"primaryKey;column:opr_name"`

	// 第几个模组 默认证章为 0，其余按发布顺序来（以 prts.wiki 为准）
	Order int `json:"order" gorm:"primaryKey;column:order"`

	Name     string   `json:"name" gorm:"column:name"`
	Missions []string `json:"missions" gorm:"serializer:json;column:missions"`
}

func (m Module) TableName() string {
	return "modules"
}
