package models

type User struct {
	QQ      int64  `json:"qq" gorm:"column:qq;primaryKey;not null"`
	AkToken string `json:"ak_token" gotm:"column:ak_token"`
}

func (u User) TableName() string {
	return "users"
}
