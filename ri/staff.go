package ri

import (
	"fmt"
	"log"
	"strings"
)

type Staff struct {
	Name                     string
	SkillUpgradeMaterials    [6][]SkillMaterial    // 依次对应 1->2 ~ 6->7 级 e.g.[0] 表示 1->2 级需要的材料
	SkillSpecializeMaterials [3][3][]SkillMaterial // 三个技能，三个等级 e.g.[1][2]表示二技能专三需要的材料
	// Scraper
}

func (s *Staff) SkillSpecializeMaterialsToMsg() string {
	log.Println(s.SkillSpecializeMaterials)
	msg := fmt.Sprintf("【%s技能专精材料】", s.Name)
	for skillIdx, skillMaterials := range s.SkillSpecializeMaterials {
		if len(skillMaterials[skillIdx]) > 0 {
			if skillIdx > 0 {
				msg += "\n"
			}
			msg += fmt.Sprintf("\n「%d技能」：", skillIdx+1)
			for lv, materials := range skillMaterials {
				msg += fmt.Sprintf("\n%d→%d ", lv+7, lv+8)
				var strs []string
				for _, material := range materials {
					strs = append(strs, fmt.Sprintf("%s×%d", material.Material.Name, material.Amount))
				}
				msg = msg + strings.Join(strs, " ")
			}
		}
	}
	return msg
}
