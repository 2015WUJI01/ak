package handlers

import (
	"fmt"
	"main/database"
	"main/pkg/logger"
	"main/ri"
	"main/scheduler"
	"main/wiki"
)

// Operator 目前只发送链接，中期转文字，后期转图片
func Operator(c *scheduler.Context) {
	msg := c.PretreatedMessage

	var oprs []ri.Operator
	database.DB.Where("`uuid` = ?", msg).
		Or("`name` like ?", "%"+msg+"%").
		Or("json_contains(`alias`, ?)", msg).
		Find(&oprs)
	logger.Info(oprs)

	if len(oprs) == 0 {
		_, _ = c.Reply("未找到干员信息，请检查输入的关键字")
	} else if len(oprs) == 1 {
		o := oprs[0]
		database.DB.Table("skill").Where("opr_uuid = ?", o.UUID).Order("`order`").Find(&o.Skill)
		var skillLevels []ri.OprSkillLevel
		database.DB.Table("skill_level as a").
			Select("a.opr_uuid, a.order, a.level, a.comment, a.ori_pt, a.cost_pt, a.last").
			Joins("join (?) b on a.opr_uuid=b.opr_uuid and a.order=b.order and a.level=b.level", database.DB.Table("skill_level").
				Select("opr_uuid, `order`, Max(`level`) `level`").Group("opr_uuid, `order`")).
			Where("a.opr_uuid = ?", o.UUID).Order("`level` desc,`order`").Find(&skillLevels)
		for i := 0; i < len(skillLevels); i++ {
			if len(o.Skill) > i {
				o.Skill[i].Level = append(o.Skill[i].Level, skillLevels[i])
			}
		}
		oprMsg := fmt.Sprintf("%d %s %d★%v-%v", o.UUID, o.Name, o.Rarity, o.Class, o.SubClass)
		oprMsg += fmt.Sprintf("\n\n【技能（max）】")
		if len(o.Skill) == 0 {
			oprMsg += fmt.Sprintf("\n该干员暂无技能信息")
		} else {
			for i := 0; i < len(o.Skill); i++ {
				if len(o.Skill[i].Level) > 0 {
					oprMsg += fmt.Sprintf("\n%v技能「%s」\n%s %s %v/%v/%v%s\n%s", o.Skill[i].Order, o.Skill[i].Name,
						o.Skill[i].Restore, o.Skill[i].Active,
						o.Skill[i].Level[0].OriPt, o.Skill[i].Level[0].CostPt, o.Skill[i].Level[0].Last, "s",
						o.Skill[i].Level[0].Comment)
				} else {
					oprMsg += fmt.Sprintf("\n「%s」\n%s %s -/-/-", o.Skill[i].Name,
						o.Skill[i].Restore, o.Skill[i].Active)
				}
			}
		}
		oprMsg += fmt.Sprintf("【模组】\n")
		oprMsg += wiki.OperatorPage(o.Name)
		_, _ = c.Reply(oprMsg)
	} else {
		oprsMsg := fmt.Sprintf("共找到 %d 个结果：", len(oprs))
		for _, o := range oprs {
			oprsMsg += fmt.Sprintf("\n%d %s", o.UUID, o.Name)
		}
		_, _ = c.Reply(oprsMsg)
	}
}
