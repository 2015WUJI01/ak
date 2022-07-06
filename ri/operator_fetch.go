package ri

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"main/pkg/logger"
	"main/wiki"
	"regexp"
	"strconv"
	"strings"
)

// Fetch 获取干员数据
func (o *Operator) Fetch() error {
	logger.Debugf("开始在线获取%s干员数据", o.Name)
	return o.FetchClassName().FetchSkill().FetchEquipment().FetchStar().FetchLimited().Visit()
}

// Visit 开始访问干员页面
func (o Operator) Visit() error { return o.Fetcher().Visit(wiki.OperatorPage(o.Name)) }

// FetchBasicValue 基础数值
func (o *Operator) FetchBasicValue(f *Fetcher) *Operator {

	// f.OnHTML("table.char-extra-attr-table tr td", func(e *colly.HTMLElement) {
	// 	switch e.Index {
	// 	case 0:
	// 		o.OriDt, _ = strconv.Atoi(strings.TrimSuffix(strings.TrimSpace(e.Text), "s"))
	// 	case 1:
	// 		o.OriDc = strings.TrimSpace(e.Text)
	// 	case 2:
	// 		o.OriBlock, _ = strconv.Atoi(strings.TrimSpace(e.Text))
	// 	case 3:
	// 		o.OriCd = strings.TrimSpace(e.Text)
	// 	}
	// })
	// c.OnHTML("table.char-base-attr-table tbody tr", func(tr *colly.HTMLElement) {
	// 	if tr.Index > 0 {
	// 		val := 0
	// 		add := 0
	// 		tr.ForEach("td", func(i int, td *colly.HTMLElement) {
	// 			if i == 3 {
	// 				// 精二满级数值
	// 				val, _ = strconv.Atoi(strings.TrimSpace(td.Text))
	// 			} else if i == 4 {
	// 				// 信赖加成
	// 				add, _ = strconv.Atoi(strings.TrimSpace(td.Text))
	// 			}
	// 		})
	// 		switch tr.Index {
	// 		case 1:
	// 			o.OriHp = val
	// 		case 2:
	// 			o.OriAtk = val
	// 		case 3:
	// 			o.OriDef = val
	// 		case 4:
	// 			o.OriRes = val
	// 		}
	// 	}
	// })
	return o
}

// FetchClassName 获取干员职业
func (o *Operator) FetchClassName() *Operator {
	o.Fetcher().OnHTML("#charclasstxt a", func(a *colly.HTMLElement) {
		switch a.Index {
		case 0:
			o.Class = a.Text
			logger.Debugf("获取到干员主职业：%s", o.Class)
		case 1:
			o.SubClass = a.Text
			logger.Debugf("获取到干员副职业：%s", o.SubClass)
		}
	})
	return o
}

// FetchStar 获取干员稀有度
func (o *Operator) FetchStar() *Operator {
	o.Fetcher().OnHTML("div#star div.starimg img", func(img *colly.HTMLElement) {
		icon := img.Attr("src")
		starIdx, _ := strconv.Atoi(icon[len(icon)-5 : len(icon)-4])
		o.Rarity = starIdx + 1
	})
	return o
}

// FetchLimited 获取是否是限定
func (o *Operator) FetchLimited() *Operator {
	o.Fetcher().OnResponse(func(r *colly.Response) {
		// 判定是否为限定需要正则在 js 代码中匹配，html 代码无法区分
	})
	return o
}

// FetchSkill 获取技能相关信息
func (o *Operator) FetchSkill() *Operator {
	// 技能序号选择器函数
	selector := func(o int) string {
		return "h2:has(span#技能) + p + table" + strings.Repeat(" + table + p + table", o-1)
	}

	// 遍历三个技能
	for i := 1; i <= 3; i++ {
		order := i
		o.Fetcher().OnHTML(selector(order), func(table *colly.HTMLElement) {
			// 获取技能名称
			skillName := table.ChildText("tr:nth-of-type(1) td:nth-of-type(2) big")
			logger.Debugf("获取到干员%d技能名称：%s", order, skillName)

			// 获取技力回复方式 激活方式
			restore := table.ChildText("tr:nth-of-type(1) td:nth-of-type(3) span:nth-of-type(1)")
			active := table.ChildText("tr:nth-of-type(1) td:nth-of-type(3) span:nth-of-type(2)")
			logger.Debugf("获取到干员%d技能回复方式：%s", order, restore)
			logger.Debugf("获取到干员%d技能激活方式：%s", order, active)

			// 获取技能 icon 链接
			href := table.ChildAttr("tr:nth-of-type(1) td:nth-child(1) > span > a", "href")
			icon := getSkillIcon(wiki.Page(href))
			logger.Debugf("获取到干员%d技能图标：%s", order, icon)

			skill := OprSkill{
				OprUUID: o.UUID,
				Order:   uint(order),
				Name:    skillName,
				Icon:    icon,
				Restore: restore,
				Active:  active,
			}
			// 获取每一级的技能信息
			table.ForEach("tr", func(i int, tr *colly.HTMLElement) {
				if i >= 2 && i <= 11 {
					oriPt, _ := strconv.Atoi(tr.ChildText("td:nth-of-type(3)"))
					costPt, _ := strconv.Atoi(tr.ChildText("td:nth-of-type(4)"))
					last, _ := strconv.Atoi(tr.ChildText("td:nth-of-type(5)"))
					sk := OprSkillLevel{
						OprUUID: o.UUID,
						Order:   uint(order),
						Level:   uint(i - 1),
						OriPt:   oriPt,
						CostPt:  costPt,
						Last:    last,
						Comment: func() (cmt string) {
							// 技能描述中，有隐藏的文本，但是直接使用 Text 会获取到。所以暂时只能通过手动排除的方式去掉这个元素
							tr.ForEach("td:nth-of-type(2)", func(i int, e *colly.HTMLElement) {
								e.DOM.Each(func(i int, s *goquery.Selection) {
									s.Children().Each(func(ii int, ss *goquery.Selection) {
										if sty, ok := ss.Attr("style"); ok && sty == "display:none;" {
											ss.Remove()
										}
									})
								})
								cmt = strings.TrimSpace(e.DOM.Text())
							})
							return
						}(),
					}
					skill.Level = append(skill.Level, sk)
					logger.Debugf("获取到干员%d技能%d级初始技力：%v", order, i-1, sk.OriPt)
					logger.Debugf("获取到干员%d技能%d级消耗技力：%v", order, i-1, sk.CostPt)
					logger.Debugf("获取到干员%d技能%d级持续时间：%v", order, i-1, sk.Last)
					logger.Infof("获取到干员%d技能%d级详细说明：%v", order, i-1, sk.Comment)
				} else if i > 12 {
					// 技能备注信息
				}
			})
			o.Skill = append(o.Skill, skill)
		})
	}
	// 获取技能的升级和专精材料
	o.fetchSkillUpgradeMaterials().fetchSkillSpecMaterials()
	return o
}

// 获取高清技能图标
func getSkillIcon(link string) (icon string) {
	f := NewFetcher().AutoRetry(0)
	f.OnHTML("#file.fullImageLink a[href]", func(a *colly.HTMLElement) {
		icon = wiki.Page(a.Attr("href"))
	})
	_ = f.Visit(link)
	return
}

// 获取技能专精材料
func (o *Operator) fetchSkillSpecMaterials() *Operator {
	o.Fetcher().OnHTML("h2:has(span#技能升级材料) + table tbody tr", func(tr *colly.HTMLElement) {
		if tr.Index == 6 || tr.Index == 7 || tr.Index == 8 {
			lv := tr.Index - 6
			skillCnt := 0
			tr.ForEach("td", func(i int, td *colly.HTMLElement) {
				skillCnt++
			})
			o.SkillSpecMaterials = make([][3][]ItemGroup, skillCnt)
			tr.ForEach("td", func(idx int, td *colly.HTMLElement) {
				skillIdx := idx
				var skillMaterials []ItemGroup
				td.ForEach("div", func(idx int, div *colly.HTMLElement) {
					item := IC.FindItemByName(strings.TrimSpace(div.ChildAttr("a", "title")))
					amount := transferAmountUnit(div.ChildText("span"))
					skillMaterials = append(skillMaterials, ItemGroup{
						ItemID:   fmt.Sprintf(item.ItemID),
						ItemName: item.Name,
						Amount:   uint(amount),
					})
				})
				o.SkillSpecMaterials[skillIdx][lv] = skillMaterials
				logger.Debugf("已获取干员%d技能%d级专精材料：%v", skillIdx+1, lv+1, skillMaterials)
			})
		}
	})
	return o
}

// 获取技能升级材料
func (o *Operator) fetchSkillUpgradeMaterials() *Operator {
	o.Fetcher().OnHTML("h2:has(span#技能升级材料) + table tbody tr", func(tr *colly.HTMLElement) {
		// tr.Index == 1或3时表示技能升级材料
		if tr.Index == 1 || tr.Index == 3 {
			tr.ForEach("td", func(idx int, td *colly.HTMLElement) {
				var upMaterials []ItemGroup // 升一级所需的材料组
				td.ForEach("div", func(idx int, div *colly.HTMLElement) {
					item := IC.FindItemByName(strings.TrimSpace(div.ChildAttr("a", "title")))
					amount := transferAmountUnit(div.ChildText("span"))
					upMaterials = append(upMaterials, ItemGroup{
						ItemID:   fmt.Sprintf(item.ItemID),
						ItemName: item.Name,
						Amount:   uint(amount),
					})
				})
				// 下标 0-5：表示从 idx+1 级升到 idx+2 级所需要的材料组
				o.SkillUpMaterials[tr.Index*4/3+idx-1] = upMaterials
				logger.Debugf("已获取干员技能%d级升级材料：%v", tr.Index*4/3+idx, upMaterials)
			})
		}
	})
	return o
}

// FetchEquipment 获取模组相关信息
func (o *Operator) FetchEquipment() *Operator {
	o.Fetcher().OnResponse(func(response *colly.Response) {
		var r *regexp.Regexp
		// 查找默认模组，并添加到模组数组中
		r = regexp.MustCompile(`var equipItems = .*\n.*{default:(.*)},\n.*label: '(.*)'`)
		for _, items := range r.FindAllStringSubmatch(string(response.Body), -1) {
			o.OperatorEquipment = append(o.OperatorEquipment, Equipment{
				OprUUID:  o.UUID,
				Order:    len(o.OperatorEquipment),
				Name:     strings.TrimSpace(items[2]),
				Missions: nil,
				Unlock:   nil,
			})
		}
		// 查找其他模组名称，添加到模组数组中
		r = regexp.MustCompile(`equipNames = \["(.*)",\n.*"(.*)",\n.*"(.*)",\n.*"(.*)",\n.*"(.*)"]`)
		for _, items := range r.FindAllStringSubmatch(string(response.Body), -1) {
			for i, str := range items {
				if i > 0 && strings.TrimSpace(str) != "" {
					o.OperatorEquipment = append(o.OperatorEquipment, Equipment{
						OprUUID:  o.UUID,
						Order:    len(o.OperatorEquipment),
						Name:     strings.TrimSpace(str),
						Missions: nil,
						Unlock:   nil,
					})
				}
			}
		}

		if len(o.OperatorEquipment) == 0 {
			logger.Debug("未获取到任何干员模组")
		} else {
			var names []string
			for _, oe := range o.OperatorEquipment {
				names = append(names, oe.Name)
			}
			logger.Debugf("获取到%d个干员模组，分别为：%s", len(o.OperatorEquipment), strings.Join(names, "、"))
		}

		// 默认证章
		if len(o.OperatorEquipment) > 0 {
			selector := fmt.Sprintf("h3:has(span#%s) + table tbody tr", o.OperatorEquipment[0].Name)
			o.Fetcher().OnHTML(selector, func(tr *colly.HTMLElement) {})
		}

		// 其余证章
		for order := 1; order < len(o.OperatorEquipment); order++ {
			i := order
			selector := fmt.Sprintf("h3:has(span#%s) + p +section + p + table tbody tr", o.OperatorEquipment[i].Name)
			o.Fetcher().OnHTML(selector, func(tr *colly.HTMLElement) {
				if tr.Index == 1 {
					// 2022/04/21 22:12:14 诗短梦长 1 基础数值变化：生命上限 +100攻击 +30
				} else if tr.Index == 2 {
					// 分支特性追加：召唤物持有上限+3，召唤物部署费用减少※“清平”和“逍遥”的费用-3；“弦惊”的费用-5
					o.OperatorEquipment[i].Attribution = strings.TrimSpace(strings.ReplaceAll(tr.Text, "※", "\n※"))
					logger.Debugf("获取到模组 %s 分支特性：%v", o.OperatorEquipment[i].Name, o.OperatorEquipment[i].Attribution)
				} else if tr.Index == 3 {
					// 2022/04/21 22:12:14 诗短梦长 3 模组任务
				} else if tr.Index == 4 {
					// 2022/04/21 22:12:14 诗短梦长 4 完成5次战斗；必须编入非助战令并上场，且使用令与令的召唤物歼灭至少7名敌人
					o.OperatorEquipment[i].Missions = append(o.OperatorEquipment[i].Missions, strings.TrimSpace(tr.Text))
				} else if tr.Index == 5 {
					// 2022/04/21 22:12:14 诗短梦长 5 3星通关主题曲3-4；必须编入非助战令并上场，且整场战斗仅部署过令与至多4位其他干员
					o.OperatorEquipment[i].Missions = append(o.OperatorEquipment[i].Missions, strings.TrimSpace(tr.Text))
					logger.Debugf("获取到模组 %s 任务：%v", o.OperatorEquipment[i].Name, o.OperatorEquipment[i].Missions)
				} else if tr.Index == 6 {
					// 2022/04/21 22:12:14 诗短梦长 6 解锁需求和材料消耗
				} else if tr.Index == 7 {
					// 2022/04/21 22:12:14 诗短梦长 7 完成该模组所有模组任务达到精英阶段2 60级信赖值达到100%
					tr.ForEach("td:nth-of-type(2) div", func(_ int, div *colly.HTMLElement) {
						item := IC.FindItemByName(strings.TrimSpace(div.ChildAttr("a", "title")))
						amount := transferAmountUnit(div.ChildText("span"))
						o.OperatorEquipment[i].Unlock = append(o.OperatorEquipment[i].Unlock, ItemGroup{
							ItemID:   fmt.Sprintf(item.ItemID),
							ItemName: item.Name,
							Amount:   uint(amount),
						})
					})
					logger.Debugf("获取到模组 %s 解锁耗材：%v", o.OperatorEquipment[i].Name, o.OperatorEquipment[i].Unlock)
				}
			})
		}
	})
	return o
}

func transferAmountUnit(s string) int {
	amount, _ := strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(s, "万", "0000"), "千", "000"))
	return amount
}
