package ri

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"main/database"
	"main/pkg/help"
	"main/pkg/logger"
	"main/wiki"
	"time"
)

var IC ItemCenter

type ItemCenter struct {
	Items []Item
}

// 在 ic 中查询
func (ic ItemCenter) FindItemByName(name string) Item {
	// 在 ic 中则直接返回
	for _, item := range ic.Items {
		if item.Name == name {
			return item
		}
	}
	// 不在 ic 中则使用数据库进行全名的模糊查询
	return FindItemByName(name)
}

// 在数据库中找
func FindItemByName(name string) Item {
	var item Item
	database.DB.Model(&Item{}).Where("name", name).Or("name like ?", "%"+name+"%").First(&item)
	if item.UUID == 0 {
		item.Name = name
		database.DB.Create(&item)
	}
	return item
}

type Item struct {
	UUID   uint   `json:"-" gorm:"primaryKey;"`
	ItemID string `json:"itemId" gorm:"column:item_id"`
	Name   string `json:"name"`
	// NameI18n    NameI18n          `json:"name_i18n"`
	// Existence   NameI18nExistence `json:"existence"`
	ItemType string `json:"itemType"`
	SortId   int    `json:"sortId"`
	GroupID  string `json:"groupID"`

	SpriteCoord  []int `json:"spriteCoord" gorm:"-:all"`
	SpriteCoordX int   `json:"-" gorm:"column:sprite_coord_x;type:int;size:5"`
	SpriteCoordY int   `json:"-" gorm:"column:sprite_coord_y;type:int;size:5"`

	Alias       voice  `json:"alias" gorm:"-:all"` // 别名
	AliasString string `json:"-" gorm:"column:alias;type:text"`
	// Pron        voice             `json:"pron"`  // 发音
}

type voice struct {
	// Ja []string `json:"ja"`
	Zh []string `json:"zh"`
}

// 更新 ic
func (ic ItemCenter) Update() error {
	// 从 API 中获取数据文本
	return database.DB.Find(&IC.Items).Error
}

func (ic ItemCenter) Truncate() {
	database.DB.Exec("TRUNCATE TABLE `items`")
}

func WriteItemsIntoJsonFile(body *[]byte) {
	if err := ioutil.WriteFile("items.json", *body, 0755); err != nil {
		log.Println(err.Error())
		return
	}
}

func WriteItemsIntoDB(body *[]byte) {
	var items []Item
	_ = json.Unmarshal(*body, &items)

	// 预处理
	for i, item := range items {
		// 保存别名
		res, _ := json.Marshal(item.Alias)
		items[i].AliasString = string(res)

		// 保存雪碧图坐标
		if len(item.SpriteCoord) == 2 {
			items[i].SpriteCoordX = item.SpriteCoord[0]
			items[i].SpriteCoordY = item.SpriteCoord[1]
		}
	}

	// 执行 SQL
	IC.Truncate()
	database.DB.Create(&items)

	// 执行补充处理
	database.DB.Exec("insert ignore into items(`item_id`, `name`, `item_type`, `sort_id`, `sprite_coord_x`, `sprite_coord_y`, `alias`) values(?,?,?,?,?,?,?)", "4001", "龙门币", "TEMP", -10000, 4, 7, "{\"zh\":[\"龙门币\",\"钱\"]}")
}

// UpdateItemsData 加载 items 数据
func UpdateItemsData() error {
	// 初始化数据
	logger.Debug("开始从企鹅物流获取 items 数据...")
	t1 := time.Now()
	body, err := wiki.GetRespBodyFromItemsAPI() // 从企鹅物流 API 拿到数据文本
	if err != nil {
		return errors.New("从企鹅物流获取数据失败，" + err.Error())
	}
	logger.Debug(fmt.Sprintf("[%.2fs] 企鹅物流 items 数据源获取完成", help.SpendSeconds(t1)))

	t2 := time.Now()
	WriteItemsIntoJsonFile(&body) // 将数据写到 JSON 中
	logger.Debug(fmt.Sprintf("[%.2fs] 已将 items 数据存为 JSON 文件", help.SpendSeconds(t2)))

	t3 := time.Now()
	WriteItemsIntoDB(&body) // 将数据写到数据库中
	logger.Debug(fmt.Sprintf("[%.2fs] 已将 items 数据存入数据库", help.SpendSeconds(t3)))

	logger.Debug(fmt.Sprintf("初始化数据完成，总耗时 %.2fs", help.SpendSeconds(t1)))

	return nil
}
