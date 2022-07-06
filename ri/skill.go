package ri

type ItemGroup struct {
	ItemID   string `json:"item_id" gorm:""`
	ItemName string `json:"item_name" gorm:""`
	Amount   uint   `json:"amount" gorm:""`
}
