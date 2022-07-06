package ri

type NameI18n struct {
	En string `json:"en"`
	Ja string `json:"ja"`
	Ko string `json:"ko"`
	Zh string `json:"zh"`
}

type NameI18nExist struct {
	Exist bool `json:"exist"`
}

type NameI18nExistence struct {
	CN NameI18nExist `json:"CN"`
	JP NameI18nExist `json:"JP"`
	KR NameI18nExist `json:"KR"`
	US NameI18nExist `json:"US"`
}

// Material 材料
type Material struct {
	Name string
	Img  string // 图片链接
}

// SkillMaterial 单个技能升级材料，需要知道升级的材料是什么，需要多少个这种材料
type SkillMaterial struct {
	Material Item
	Amount   uint
}
