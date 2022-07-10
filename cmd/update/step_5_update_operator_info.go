package update

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"gorm.io/gorm/clause"
	"main/database"
	"main/logger"
	"main/models"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}
var completed = 0

var operators []models.Operator
var skills []models.Skill
var skillLevels []models.SkillLevel
var skillIcons []models.Skill
var skillLevelMaterials []models.SkillLevelMaterial
var modules []models.Module
var moduleStages []models.ModuleStage
var moduleStageMaterials []models.ModuleStageMaterial

func Step5() {
	var oprs []models.Operator
	database.DB.Find(&oprs)

	fmt.Println("逐条更新干员信息中...")

	c := colly.NewCollector(colly.Async(true))
	c.SetRequestTimeout(5 * time.Second)
	c.OnError(func(r *colly.Response, err error) {
		_ = r.Request.Retry()
	})
	_ = c.Limit(&colly.LimitRule{
		DomainGlob:  "*prts.wiki*",
		Parallelism: 5,
		Delay:       50 * time.Millisecond,
		RandomDelay: 50 * time.Millisecond,
	})

	// 获取并更新干员 wiki 短链接、职业、稀有度、最后编辑时间
	updateOperatorInfo(c)
	for _, opr := range oprs {
		_ = c.Visit(Link("/w/" + opr.Name))
	}

	c.Wait()
	wg.Wait()
	fmt.Println()

	logger.Infof("干员信息补充：%d", len(operators))
	cols := []string{"name", "class", "subclass", "rarity", "wiki_short", "updated_at"}
	database.DB.Select(cols).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns(cols),
	}).Create(&operators)

	logger.Infof("所有技能数量：%d", len(skills))
	cols = []string{"opr_id", "opr_name", "order", "name", "restore", "active"}
	database.DB.Select(cols).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "opr_id"}, {Name: "order"}},
		DoUpdates: clause.AssignmentColumns(cols),
	}).Create(&skills)

	logger.Infof("所有技能图标数量：%d", len(skillIcons))
	cols = []string{"opr_id", "order", "icon"}
	database.DB.Select(cols).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "opr_id"}, {Name: "order"}},
		DoUpdates: clause.AssignmentColumns(cols),
	}).Create(&skillIcons)

	logger.Infof("所有技能各等级数量：%d", len(skillLevels))
	cols = []string{"opr_id", "opr_name", "order", "level", "ori_pt", "cost_pt", "last", "comment"}
	database.DB.Select(cols).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "opr_id"}, {Name: "order"}, {Name: "level"}},
		DoUpdates: clause.AssignmentColumns(cols),
	}).CreateInBatches(&skillLevels, 1000)

	logger.Infof("所有技能各等级升级材料数量：%d", len(skillLevelMaterials))
	cols = []string{"opr_id", "opr_name", "order", "to_level", "item_name", "amount"}
	database.DB.Select(cols).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "opr_name"}, {Name: "order"}, {Name: "to_level"}, {Name: "item_name"}},
		DoUpdates: clause.AssignmentColumns(cols),
	}).CreateInBatches(&skillLevelMaterials, 1000)

	logger.Infof("所有模组数量：%d", len(modules))
	cols = []string{"opr_id", "opr_name", "order", "name", "missions"}
	database.DB.Select(cols).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "opr_name"}, {Name: "order"}},
		DoUpdates: clause.AssignmentColumns(cols),
	}).Create(&modules)

	logger.Infof("所有模组等级：%d", len(moduleStages))
	cols = []string{"opr_name", "opr_id", "module_name", "module_order", "stage", "basic_info", "attribution"}
	database.DB.Select(cols).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "opr_name"}, {Name: "module_name"}, {Name: "stage"}},
		DoUpdates: clause.AssignmentColumns(cols),
	}).Create(&moduleStages)

	logger.Infof("所有模组等级升级材料：%d", len(moduleStageMaterials))
	cols = []string{"opr_id", "opr_name", "module_name", "module_order", "to_stage", "item_name", "amount"}
	database.DB.Select(cols).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "opr_name"}, {Name: "module_name"}, {Name: "to_stage"}, {Name: "item_name"}},
		DoUpdates: clause.AssignmentColumns(cols),
	}).Create(&moduleStageMaterials)

	logger.Info("Step5. 干员基本信息更新完成")
}

func updateOperatorInfo(c *colly.Collector) {
	c.OnHTML("body", func(body *colly.HTMLElement) {

		// 1. 首要任务是在页面中获取干员名称，以便锁定干员
		var opr models.Operator
		body.ForEach("h1#firstHeading", func(_ int, e *colly.HTMLElement) {
			opr.Name = strings.TrimSpace(e.Text)
			database.DB.Model(&models.Operator{}).Where("name", opr.Name).First(&opr)
		})

		// 2. 获取干员职业
		body.ForEach("#charclasstxt a", func(i int, a *colly.HTMLElement) {
			if i == 0 {
				opr.Class = a.Text
			} else if i == 1 {
				opr.Subclass = a.Text
			}
		})

		// 3. 获取 wiki 短链接
		body.ForEach(".copyUrl", func(i int, e *colly.HTMLElement) {
			opr.WikiShort = e.Attr("data-clipboard-text")
		})

		// 4. 获取最后编辑时间
		body.ForEach("#footer-info-lastmod", func(_ int, e *colly.HTMLElement) {
			opr.UpdatedAt = parseTime(e.Text)
		})

		// 5. 获取干员稀有度
		body.ForEach("div#star div.starimg img", func(_ int, img *colly.HTMLElement) {
			icon := img.Attr("src")
			starIdx, _ := strconv.Atoi(icon[len(icon)-5 : len(icon)-4])
			opr.Rarity = starIdx + 1
		})

		// update: 边抓取边更新补充干员信息
		// CreateOrUpdateOperator([]string{"class", "subclass", "rarity", "wiki_short", "updated_at"}, opr)
		operators = append(operators, opr)

		// 6. 遍历获取并更新三个技能
		for i := 1; i <= 3; i++ {
			order := i
			body.ForEach(selector(order), func(i int, table *colly.HTMLElement) {
				// 1. 获取技能名称
				skillName := table.ChildText("tr:nth-of-type(1) td:nth-of-type(2) big")
				// logger.Debugf("获取到干员%d技能名称：%s", order, skillName)

				// 2. 获取技力回复方式 激活方式
				restore := table.ChildText("tr:nth-of-type(1) td:nth-of-type(3) span:nth-of-type(1)")
				active := table.ChildText("tr:nth-of-type(1) td:nth-of-type(3) span:nth-of-type(2)")
				// logger.Debugf("获取到干员%d技能回复方式：%s", order, restore)
				// logger.Debugf("获取到干员%d技能激活方式：%s", order, active)

				// 3. 获取技能 icon 链接
				href := table.ChildAttr("tr:nth-of-type(1) td:nth-child(1) > span > a", "href")
				wg.Add(1)
				go FetchAndUpdateOperatorSkillIcon(href, models.Skill{OprID: opr.ID, Order: order})

				// cols := []string{"opr_id", "opr_name", "order", "name", "restore", "active"}
				// CreateOrUpdateOperatorSkill(cols, models.Skill{
				// 	OprID:   opr.ID,
				// 	OprName: opr.Name,
				// 	Order:   order,
				// 	Name:    skillName,
				// 	Restore: restore,
				// 	Active:  active,
				// })
				skills = append(skills, models.Skill{
					OprID:   opr.ID,
					OprName: opr.Name,
					Order:   order,
					Name:    skillName,
					Restore: restore,
					Active:  active,
				})

				// 获取每一级的技能信息
				table.ForEach("tr", func(i int, tr *colly.HTMLElement) {
					var sk models.SkillLevel
					if i >= 2 && i <= 11 {
						oriPt, _ := strconv.Atoi(tr.ChildText("td:nth-of-type(3)"))
						costPt, _ := strconv.Atoi(tr.ChildText("td:nth-of-type(4)"))
						last, _ := strconv.Atoi(tr.ChildText("td:nth-of-type(5)"))
						sk = models.SkillLevel{
							OprID:   opr.ID,
							OprName: opr.Name,
							Order:   order,
							Level:   i - 1,
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
						// cols := []string{"opr_id", "opr_name", "order", "level", "opr_pt", "cost_pt", "last", "comment"}
						// CreateOrUpdateOperatorSkillLevel(cols, sk)
						skillLevels = append(skillLevels, sk)
						// logger.Debugf("获取到干员%d技能%d级初始技力：%v", order, i-1, sk.OriPt)
						// logger.Debugf("获取到干员%d技能%d级消耗技力：%v", order, i-1, sk.CostPt)
						// logger.Debugf("获取到干员%d技能%d级持续时间：%v", order, i-1, sk.Last)
						// logger.Infof("获取到干员%d技能%d级详细说明：%v", order, i-1, sk.Comment)
					} else if i > 12 {
						// 技能备注信息
					}
				})
			})
		}

		body.ForEach("h2:has(span#技能升级材料) + table tbody tr", func(i int, tr *colly.HTMLElement) {
			// cols := []string{"opr_id", "opr_name", "order", "to_level", "item_name", "amount"}
			if i == 1 || i == 3 { // i == 1 或 3 时表示技能升级材料
				tr.ForEach("td", func(idx int, td *colly.HTMLElement) {
					td.ForEach("div", func(_ int, div *colly.HTMLElement) {
						// CreateOrUpdateOperatorSkillLevelMaterial(cols, models.SkillLevelMaterial{
						// 	OprID:    opr.ID,
						// 	OprName:  opr.Name,
						// 	Order:    0,
						// 	ToLevel:  i*4/3 + idx + 1,
						// 	ItemName: strings.TrimSpace(div.ChildAttr("a", "title")),
						// 	Amount:   transferAmountUnit(div.ChildText("span")),
						// })
						skillLevelMaterials = append(skillLevelMaterials, models.SkillLevelMaterial{
							OprID:    opr.ID,
							OprName:  opr.Name,
							Order:    0,
							ToLevel:  i*4/3 + idx + 1,
							ItemName: strings.TrimSpace(div.ChildAttr("a", "title")),
							Amount:   transferAmountUnit(div.ChildText("span")),
						})
					})
				})
			} else if i == 6 || i == 7 || i == 8 { // i == 1 或 3 时表示技能专精材料
				tr.ForEach("td", func(idx int, td *colly.HTMLElement) {
					td.ForEach("div", func(_ int, div *colly.HTMLElement) {
						// CreateOrUpdateOperatorSkillLevelMaterial(cols, models.SkillLevelMaterial{
						// 	OprID:    opr.ID,
						// 	OprName:  opr.Name,
						// 	Order:    idx + 1,
						// 	ToLevel:  i + 2,
						// 	ItemName: strings.TrimSpace(div.ChildAttr("a", "title")),
						// 	Amount:   transferAmountUnit(div.ChildText("span")),
						// })
						skillLevelMaterials = append(skillLevelMaterials, models.SkillLevelMaterial{
							OprID:    opr.ID,
							OprName:  opr.Name,
							Order:    idx + 1,
							ToLevel:  i + 2,
							ItemName: strings.TrimSpace(div.ChildAttr("a", "title")),
							Amount:   transferAmountUnit(div.ChildText("span")),
						})
					})
				})
			}
		})
		completed++
		fmt.Printf("=")
		if completed%100 == 0 {
			fmt.Println()
		}
	})
	c.OnResponse(func(resp *colly.Response) {
		oprname := regexp.MustCompile("<h1 id=\"firstHeading\".*>(.*)</h1>").
			FindAllStringSubmatch(string(resp.Body), -1)[0][1]
		var opr models.Operator
		database.DB.Model(&models.Operator{}).Where("name", oprname).First(&opr)

		// 查找默认模组，并添加到模组数组中
		mods := searchModulesInRespBody(resp.Body)

		// 默认证章
		// if len(mods) > 0 {
		// 	selector := fmt.Sprintf("h3:has(span#%s) + table tbody tr", modules[0])
		// 	c.OnHTML(selector, func(tr *colly.HTMLElement) {})
		// }

		// 其余证章
		for i := 1; i < len(mods); i++ {
			var om = models.Module{
				OprID:   opr.ID,
				OprName: opr.Name,
				Order:   i,
				Name:    mods[i],
			}
			// cols := []string{"opr_id", "opr_name", "order", "name"}
			// CreateOrUpdateOperatorModule(cols, om)

			selector := fmt.Sprintf("h3:has(span#%s) + p + table tbody tr", om.Name)
			c.OnHTML(selector, func(tr *colly.HTMLElement) {

				if tr.Index == 2 { // stage 1
					// cols := []string{"opr_id", "module_name", "module_order", "stage", "basic_info", "attribution"}
					// CreateOrUpdateOperatorModuleStage(cols, models.ModuleStage{
					// 	OprID:       opr.ID,
					// 	ModuleName:  om.Name,
					// 	ModuleOrder: om.Order,
					// 	Stage:       1,
					// 	BasicInfo:   strings.TrimSpace(tr.ChildText("td:nth-of-type(2)")),
					// 	Attribution: strings.TrimSpace(tr.ChildText("td:nth-of-type(3)")),
					// })
					moduleStages = append(moduleStages, parseModuleStage(tr, om, 1))
				} else if tr.Index == 3 { // stage 2
					// cols := []string{"opr_id", "module_name", "module_order", "stage", "basic_info", "attribution"}
					// CreateOrUpdateOperatorModuleStage(cols, models.ModuleStage{
					// 	OprID:       opr.ID,
					// 	ModuleName:  om.Name,
					// 	ModuleOrder: om.Order,
					// 	Stage:       2,
					// 	BasicInfo:   strings.TrimSpace(tr.ChildText("td:nth-of-type(2)")),
					// 	Attribution: strings.TrimSpace(tr.ChildText("td:nth-of-type(3)")),
					// })
					moduleStages = append(moduleStages, parseModuleStage(tr, om, 2))
				} else if tr.Index == 4 { // stage 3
					// cols := []string{"opr_id", "module_name", "module_order", "stage", "basic_info", "attribution"}
					// CreateOrUpdateOperatorModuleStage(cols, models.ModuleStage{
					// 	OprID:       opr.ID,
					// 	ModuleName:  om.Name,
					// 	ModuleOrder: om.Order,
					// 	Stage:       3,
					// 	BasicInfo:   strings.TrimSpace(tr.ChildText("td:nth-of-type(2)")),
					// 	Attribution: strings.TrimSpace(tr.ChildText("td:nth-of-type(3)")),
					// })
					moduleStages = append(moduleStages, parseModuleStage(tr, om, 3))
				} else if tr.Index == 5 { // mission 1
					mission1 := strings.TrimSpace(tr.ChildText("td"))
					om.Missions = append(om.Missions, mission1)
				} else if tr.Index == 6 { // mission 2
					mission2 := strings.TrimSpace(tr.ChildText("td"))
					om.Missions = append(om.Missions, mission2)
					// CreateOrUpdateOperatorModule([]string{"missions"}, om)
					modules = append(modules, om)
				} else if tr.Index == 7 { // stage 1 unlock
					tr.ForEach("td:nth-of-type(2) div", func(i int, div *colly.HTMLElement) {
						// cols := []string{"opr_id", "opr_name", "module_name", "module_order", "to_stage", "item_name", "amount"}
						// CreateOrUpdateOperatorModuleStageMaterial(cols, models.ModuleStageMaterial{
						// 	OprID:       opr.ID,
						// 	OprName:     opr.Name,
						// 	ModuleName:  om.Name,
						// 	ModuleOrder: om.Order,
						// 	ToStage:     1,
						// 	ItemName:    strings.TrimSpace(div.ChildAttr("a", "title")),
						// 	Amount:      transferAmountUnit(div.ChildText("span")),
						// })
						msm := parseModuleStageMaterials(div, om, 1)
						moduleStageMaterials = append(moduleStageMaterials, msm)
					})
				} else if tr.Index == 8 { // stage 2 & 3 upgrade
					tr.ForEach("td:nth-of-type(1) div", func(i int, div *colly.HTMLElement) {
						// cols := []string{"opr_id", "opr_name", "module_name", "module_order", "to_stage", "item_name", "amount"}
						// CreateOrUpdateOperatorModuleStageMaterial(cols, models.ModuleStageMaterial{
						// 	OprID:       opr.ID,
						// 	OprName:     opr.Name,
						// 	ModuleName:  om.Name,
						// 	ModuleOrder: om.Order,
						// 	ToStage:     2,
						// 	ItemName:    strings.TrimSpace(div.ChildAttr("a", "title")),
						// 	Amount:      transferAmountUnit(div.ChildText("span")),
						// })
						msm := parseModuleStageMaterials(div, om, 2)
						moduleStageMaterials = append(moduleStageMaterials, msm)
					})
					tr.ForEach("td:nth-of-type(2) div", func(i int, div *colly.HTMLElement) {
						// cols := []string{"opr_id", "opr_name", "module_name", "module_order", "to_stage", "item_name", "amount"}
						// CreateOrUpdateOperatorModuleStageMaterial(cols, models.ModuleStageMaterial{
						// 	OprID:       opr.ID,
						// 	OprName:     opr.Name,
						// 	ModuleName:  om.Name,
						// 	ModuleOrder: om.Order,
						// 	ToStage:     3,
						// 	ItemName:    strings.TrimSpace(div.ChildAttr("a", "title")),
						// 	Amount:      transferAmountUnit(div.ChildText("span")),
						// })
						msm := parseModuleStageMaterials(div, om, 3)
						moduleStageMaterials = append(moduleStageMaterials, msm)
					})
				}
			})
		}
	})
}

func transferAmountUnit(s string) int {
	amount, _ := strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(s, "万", "0000"), "千", "000"))
	return amount
}

func parseModuleStage(e *colly.HTMLElement, om models.Module, stage int) models.ModuleStage {
	return models.ModuleStage{
		OprID:       om.OprID,
		OprName:     om.OprName,
		ModuleName:  om.Name,
		ModuleOrder: om.Order,
		Stage:       stage,
		BasicInfo:   strings.TrimSpace(e.ChildText("td:nth-of-type(2)")),
		Attribution: strings.TrimSpace(e.ChildText("td:nth-of-type(3)")),
	}
}

func parseModuleStageMaterials(e *colly.HTMLElement, om models.Module, toStage int) models.ModuleStageMaterial {
	return models.ModuleStageMaterial{
		OprID:       om.OprID,
		OprName:     om.OprName,
		ModuleName:  om.Name,
		ModuleOrder: om.Order,
		ToStage:     toStage,
		ItemName:    strings.TrimSpace(e.ChildAttr("a", "title")),
		Amount:      transferAmountUnit(e.ChildText("span")),
	}
}

// 技能序号选择器函数
func selector(o int) string {
	return "h2:has(span#技能) + p + table" + strings.Repeat(" + table + p + table", o-1)
}

func CreateOrUpdateOperatorSkill(cols []string, sk models.Skill) (created bool) {
	var cnt int64
	database.DB.Model(&models.Skill{}).Where("opr_id", sk.OprID).Where("order", sk.Order).Count(&cnt)
	if cnt > 0 {
		database.DB.Model(&models.Skill{}).Select(cols).Where("opr_id", sk.OprID).Where("order", sk.Order).Updates(&sk)
		return false
	}
	database.DB.Model(&models.Skill{}).Select(cols).Create(&sk)
	return true
}

func CreateOrUpdateOperatorSkillLevel(cols []string, sl models.SkillLevel) (created bool) {
	var cnt int64
	database.DB.Model(&models.SkillLevel{}).Where("opr_id", sl.OprID).Where("order", sl.Order).Where("level", sl.Level).Count(&cnt)
	if cnt > 0 {
		database.DB.Model(&models.SkillLevel{}).
			Select(cols).Where("opr_id", sl.OprID).Where("order", sl.Order).Where("level", sl.Level).
			Updates(&sl)
		return false
	}
	database.DB.Model(&models.SkillLevel{}).Select(cols).Create(&sl)
	return true
}

func FetchAndUpdateOperatorSkillIcon(url string, sk models.Skill) {
	c := colly.NewCollector()
	c.SetRequestTimeout(5 * time.Second)
	c.OnError(func(r *colly.Response, err error) {
		_ = r.Request.Retry()
	})
	c.OnHTML("#file.fullImageLink a[href]", func(a *colly.HTMLElement) {
		sk.Icon = Link(a.Attr("href"))
		// CreateOrUpdateOperatorSkill([]string{"opr_id", "order", "icon"}, sk)
		skillIcons = append(skillIcons, sk)
	})
	_ = c.Visit(Link(url))
	wg.Done()
}

func CreateOrUpdateOperatorSkillLevelMaterial(cols []string, sm models.SkillLevelMaterial) (created bool) {
	var cnt int64
	database.DB.Model(&models.SkillLevelMaterial{}).
		Where("opr_id", sm.OprID).Where("order", sm.Order).
		Where("to_level", sm.ToLevel).Where("item_name", sm.ItemName).Count(&cnt)
	if cnt > 0 {
		database.DB.Model(&models.SkillLevelMaterial{}).Select(cols).
			Where("opr_id", sm.OprID).Where("order", sm.Order).
			Where("to_level", sm.ToLevel).Where("item_name", sm.ItemName).Updates(sm)
		return false
	}
	database.DB.Model(&models.SkillLevelMaterial{}).Select(cols).Create(&sm)
	return true
}

// 在返回的 response body 中查找模组
// 因为有的符合这种匹配规则，有的符合另一种匹配规则，所以单独拎个函数出来处理
func searchModulesInRespBody(body []byte) []string {

	var names []string
	modules := make(map[string]struct{})

	defaultModule := ""

	var r *regexp.Regexp
	r = regexp.MustCompile(`var equipItems = \[new OO.ui.MenuOptionWidget\({\n +data: {default:(.*)},\n +label: '(.*)'`)
	for _, items := range r.FindAllStringSubmatch(string(body), -1) {
		if strings.TrimSpace(items[2]) != "" {
			defaultModule = strings.TrimSpace(items[2])
		}
	}

	if defaultModule == "" {
		r = regexp.MustCompile(`var equipItems =\[new OO.ui.MenuOptionWidget\({data:{default:(.*)},label:'(.*)'`)
		for _, items := range r.FindAllStringSubmatch(string(body), -1) {
			if strings.TrimSpace(items[2]) != "" {
				defaultModule = strings.TrimSpace(items[2])
			}
		}
	}

	r = regexp.MustCompile(`equipNames = \["(.*)",\n.*"(.*)",\n.*"(.*)",\n.*"(.*)",\n.*"(.*)"]`)
	for _, items := range r.FindAllStringSubmatch(string(body), -1) {
		for i, str := range items {
			if i > 0 && strings.TrimSpace(str) != "" {
				modules[strings.TrimSpace(str)] = struct{}{}
			}
		}
	}

	r = regexp.MustCompile(`var equipNames =\["(.*?)","(.*?)","(.*?)","(.*?)","(.*?)"]`)
	for _, items := range r.FindAllStringSubmatch(string(body), -1) {
		for i, str := range items {
			if i > 0 && strings.TrimSpace(str) != "" {
				modules[strings.TrimSpace(str)] = struct{}{}
			}
		}
	}

	if defaultModule == "" {
		return []string{}
	}

	names = append(names, defaultModule)
	for k, _ := range modules {
		names = append(names, k)
	}
	return names
}

func CreateOrUpdateOperatorModule(cols []string, om models.Module) bool {
	var cnt int64
	database.DB.Model(&models.Module{}).
		Where("opr_id", om.OprID).
		Where("order", om.Order).
		Count(&cnt)
	if cnt > 0 {
		database.DB.Model(&models.Module{}).Select(cols).
			Where("opr_id", om.OprID).
			Where("order", om.Order).
			Updates(&om)
		return false
	}
	database.DB.Model(&models.Module{}).Select(cols).Create(&om)
	return true
}

func CreateOrUpdateOperatorModuleStage(cols []string, oms models.ModuleStage) bool {
	var cnt int64
	database.DB.Model(&models.ModuleStage{}).
		Where("opr_id", oms.OprID).
		Where("module_order", oms.ModuleOrder).
		Where("stage", oms.Stage).
		Count(&cnt)
	if cnt > 0 {
		database.DB.Model(&models.ModuleStage{}).Select(cols).
			Where("opr_id", oms.OprID).
			Where("module_order", oms.ModuleOrder).
			Where("stage", oms.Stage).
			Updates(&oms)
		return false
	}
	database.DB.Model(&models.ModuleStage{}).Select(cols).Create(&oms)
	return true
}

func CreateOrUpdateOperatorModuleStageMaterial(cols []string, omsm models.ModuleStageMaterial) bool {
	var cnt int64
	database.DB.Model(&models.ModuleStageMaterial{}).
		Where("opr_id", omsm.OprID).
		Where("module_order", omsm.ModuleOrder).
		Where("to_stage", omsm.ToStage).
		Where("item_name", omsm.ItemName).
		Count(&cnt)
	if cnt > 0 {
		database.DB.Model(&models.ModuleStageMaterial{}).Select(cols).
			Where("opr_id", omsm.OprID).
			Where("module_order", omsm.ModuleOrder).
			Where("to_stage", omsm.ToStage).
			Where("item_name", omsm.ItemName).
			Updates(&omsm)
		return false
	}
	database.DB.Model(&models.ModuleStageMaterial{}).Select(cols).Create(&omsm)
	return true
}
