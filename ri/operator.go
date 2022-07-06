package ri

import "C"
import (
	"fmt"
	"gorm.io/gorm/clause"
	"main/database"
)

type Equipment struct {
	OprUUID     int         `gorm:"index:pk,unique;not null"`
	Order       int         `gorm:"index:pk,unique;not null"`
	Name        string      `gorm:""`
	Missions    []string    `gorm:"serializer:json"`
	Unlock      []ItemGroup `gorm:"serializer:json"`
	Attribution string      `gorm:""`
}

type OprSkill struct {
	OprUUID int             `json:"opr_uuid" gorm:"column:opr_uuid"`
	Order   uint            `json:"order" gorm:"column:order"`
	Name    string          `json:"name" gorm:"column:name"`
	Icon    string          `json:"icon" gorm:"column:icon"`
	Restore string          `json:"restore" gorm:"column:restore"`
	Active  string          `json:"active" gorm:"column:active"`
	Level   []OprSkillLevel `gorm:"-:all"`
}

type OprSkillLevel struct {
	OprUUID     int
	Order       uint
	Level       uint
	UpMaterials []ItemGroup `gorm:"-:all"`
	OriPt       int
	CostPt      int
	Last        int
	Comment     string
}

type Operator struct {
	ModelFetcher
	UUID int `gorm:"primaryKey"`
	// 干员的正式中文名称
	Name string `json:"name" gorm:"column:name;type:string;size:255"`

	// 干员职业名称，例如 "狙击"
	Class string `json:"class" gorm:"column:class;type:string;size:255"`
	// 干员细分职业，例如 "速射狙"
	SubClass string `json:"sub_class" gorm:"column:sub_class;type:string;size:255"`
	Rarity   int    `json:"rarity" gorm:"column:rarity;type:tinyint;size:2"`

	Roguelike bool     `json:"roguelike" gorm:"type:bool"`
	Alias     []string `json:"alias" gorm:"serializer:json"`
	// NameI18n NameI18n `json:"name_i18n"` // 可以做，但没必要
	OperatorEquipment []Equipment `json:"operator_equipment" gorm:"-:all"`

	// 技能相关
	Skill              []OprSkill       `json:"skill" gorm:"-:all"`
	SkillUpMaterials   [6][]ItemGroup   `json:"skill_up_materials" gorm:"-:all"`
	SkillSpecMaterials [][3][]ItemGroup `json:"skill_spec_materials" gorm:"-:all"` // 三个技能，三个等级 e.g.[1][2]表示二技能专三需要的材料

	// 依次对应 1->2 ~ 6->7 级 e.g.[0] 表示 1->2 级需要的材料
	// SkillUpgradeMaterials [6][]ItemGroup `json:"skill_upgrade_materials" gorm:"-:all"`
	// SkillUpgradeMaterials       [6][]SkillMaterial `json:"skill_upgrade_materials" gorm:"-:all"`
	// SkillUpgradeMaterialsString string             `json:"-" gorm:"column:skill_upgrade_materials;type:json"`
	// SkillSpecializeMaterials       [3][3][]SkillMaterial `json:"skill_specialize_materials" gorm:"-:all"`
	// SkillSpecializeMaterialsString string                `json:"-" gorm:"column:skill_specialize_materials;type:json"`

	// Feature string // "攻击造成<span style=\"color string //#00B0FF;\">群体物理伤害</span>" // 特性
	// flex string // "优良" //
	// group string // ""

	// 半身图像
	// Half string // "//wiki.wiki/images/thumb/f/f8/%E5%8D%8A%E8%BA%AB%E5%83%8F_%E8%8F%B2%E4%BA%9A%E6%A2%85%E5%A1%94_1.png/110px-%E5%8D%8A%E8%BA%AB%E5%83%8F_%E8%8F%B2%E4%BA%9A%E6%A2%85%E5%A1%94_1.png"

	// 图标
	// Icon string // "//wiki.wiki/images/b/b8/%E5%A4%B4%E5%83%8F_%E8%8F%B2%E4%BA%9A%E6%A2%85%E5%A1%94.png"

	// ID // 这个 ID 不是唯一的，比如 阿米娅的两个形态，都是阿米娅，但是职位不同，其他的名称、ID 都是一样的
	// Index string // "LT11"

	// 日文名称
	// Jp string // ""

	// 座右铭？
	// Moredes string // "别问她的代号，否则她会抽空把你也干掉。"

	// 种族
	// Nation string // "拉特兰"

	// 初始：
	// OriAtk   int    `gorm:"column:OriAtk"`   // "375" // 攻击
	// OriBlock int    `gorm:"column:OriBlock"` // "1" // 阻挡
	// OriCd    string `gorm:"column:OriCd"`    // "2.8s" // 攻击间隔
	// OriDc    string `gorm:"column:OriDc"`    // "25→27→29" // 部署费用
	// OriDef   int    `gorm:"column:OriDef"`   // "80" // 防御
	// OriDt    int    `gorm:"column:OriDt"`    // "70s" // 再部署 // 获取的是值，存储前需要去掉秒转为数字
	// OriHp    int    `gorm:"column:OriHp"`    // "985" // 生命
	// OriRes   int    `gorm:"column:OriRes"`   // "0" // 法抗
	// Plan      string // "标准" // ?
	// Position  string // "远程位" // 站位
	// Race      string // ['黎博利']
	// Rarity    int    // "5" // 稀有度
	// Sex       string // "女" // 性别
	// Skill     string // "优良" // ?
	// SortId    int    `json:"sort_id"` // "228" // ?
	// Str       string // "标准" // ?
	// Tag       string // ['输出'] // 标签
	// Team      string // ""
	// Tolerance string // "优良"
}

// Update 增量更新干员数据
func (o *Operator) Update() {
	// 保存干员基本信息
	database.DB.Exec("UPDATE `operators` SET `class`=?, `sub_class`=?, `rarity`=? WHERE `uuid`=?", o.Class, o.SubClass, o.Rarity, o.UUID)

	// 保存技能信息
	for _, skill := range o.Skill {
		database.DB.Exec(
			"INSERT INTO `skill`(`opr_uuid`, `order`, `name`, `icon`, `restore`, `active`) VALUES(?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `name` = ?, `icon` = ?, `restore` = ?, `active` = ?",
			skill.OprUUID, skill.Order, skill.Name, skill.Icon, skill.Restore, skill.Active,
			skill.Name, skill.Icon, skill.Restore, skill.Active,
		)
		// 插入每一级的信息
		for _, lv := range skill.Level {
			database.DB.Exec(
				"INSERT INTO `skill_level`(`opr_uuid`, `order`, `level`, `ori_pt`, `cost_pt`, `last`, `comment`) VALUES(?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `ori_pt` = ?, `cost_pt` = ?, `last` = ?, `comment` = ?",
				lv.OprUUID, lv.Order, lv.Level, lv.OriPt, lv.CostPt, lv.Last, lv.Comment,
				lv.OriPt, lv.CostPt, lv.Last, lv.Comment,
			)
		}
	}

	// 保存模组信息
	for _, equip := range o.OperatorEquipment {

		// database.DB.Exec(
		// 	"INSERT INTO `equipments`(`opr_uuid`, `order`, `name`, `missions`, `unlock`) VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE `name` = ?, `missions` = ?, `unlock` = ?",
		// 	o.UUID, i, equip.Name, equip.Missions, equip.Unlock,
		// 	equip.Name, equip.Missions, equip.Unlock,
		// )
		database.DB.Table("equipments").Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&equip)
	}
}

func (o *Operator) SkillMaterialGroupsToMsg() string {
	msg := fmt.Sprintf("【%s技能升级材料】", o.Name)
	// for i, skillMaterials := range o.SkillUpgradeMaterials {
	// 	msg += fmt.Sprintf("\n%d→%d ", i+1, i+2)
	// 	var strs []string
	// 	for _, skillMaterial := range skillMaterials {
	// 		strs = append(strs, fmt.Sprintf("%s×%d", skillMaterial.Material.Name, skillMaterial.Amount))
	// 	}
	// 	msg = msg + strings.Join(strs, " ")
	// }
	return msg
}

func (o *Operator) SkillSpecializeMaterialsToMsg() string {
	// log.Println(o.SkillSpecializeMaterials)
	// msg := fmt.Sprintf("【%s技能专精材料】", o.Name)
	// for skillIdx, skillMaterials := range o.SkillSpecializeMaterials {
	// 	if len(skillMaterials[skillIdx]) > 0 {
	// 		if skillIdx > 0 {
	// 			msg += "\n"
	// 		}
	// 		msg += fmt.Sprintf("\n「%d技能」：", skillIdx+1)
	// 		for lv, materials := range skillMaterials {
	// 			msg += fmt.Sprintf("\n%d→%d ", lv+7, lv+8)
	// 			var strs []string
	// 			for _, material := range materials {
	// 				strs = append(strs, fmt.Sprintf("%o×%d", material.Material.Name, material.Amount))
	// 			}
	// 			msg = msg + strings.Join(strs, " ")
	// 		}
	// 	}
	// }
	// return msg
	return ""
}
