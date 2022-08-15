package models

import (
	"fmt"
	"github.com/gocolly/colly"
	"main/database"
	"main/wiki"
	"strings"
	"time"
)

func init() {
	_ = database.DB.AutoMigrate(
		&Item{}, &Operator{}, &Alias{},
		&Skill{}, &SkillLevel{}, &SkillLevelMaterial{},
		&Module{}, &ModuleStage{}, &ModuleStageMaterial{},
	)
}

type Item struct {
	Name      string    `json:"name" gorm:"primaryKey;column:name"`
	Group     string    `json:"group" gorm:"column:group"`
	Type      string    `json:"type" gorm:"column:type"`
	Image     string    `json:"image" gorm:"column:image"`
	Wiki      string    `json:"wiki" gorm:"column:wiki"`
	WikiShort string    `json:"wiki_short" gorm:"column:wiki_short"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:false"`
}

func (i Item) TableName() string {
	return "items"
}

func (i Item) Info(alias []string) string {
	msg := fmt.Sprintf("%s", i.Name)
	if len(alias) > 0 {
		msg += fmt.Sprintf("\n别名: %s", strings.Join(alias, ", "))
	}
	msg += fmt.Sprintf("\nimg: %s", i.Image)
	msg += fmt.Sprintf("\nwiki: %s", i.WikiShort)
	return msg
}

func FreshItemWiki(items ...*Item) []*Item {
	if len(items) == 0 {
		for n, w := range wiki.FetchAllWiki() {
			items = append(items, &Item{Name: n, Wiki: w})
		}
		return items
	}
	mwikis := wiki.FetchAllWiki(func() []string {
		var names []string
		for _, item := range items {
			names = append(names, item.Name)
		}
		return names
	}()...)
	for _, item := range items {
		item.Wiki = mwikis[item.Name]
	}
	return items
}

func FreshItemImgSWikiUpdateTime(items ...*Item) {
	res := wiki.FetchItemInfo("理智", "先锋芯片")
	for name, v := range res {
		fmt.Printf("%v:\n\t%v\n\t%v\n\t%v\n",
			name,
			v["image"].(string),
			v["wikishort"].(string),
			v["updatedat"].(time.Time).Format("2006-01-02 15:04:05"),
		)
	}
}

func (i Item) FetchGroup(c *colly.Collector) (g string) {
	// 表格中一行就是一个 group
	c.OnHTML("table.uncollapsed > tbody > tr:has(th:not(.navbox-title))", func(tr *colly.HTMLElement) {
		var group string
		var types []string

		// 获取 Group 名称
		tr.ForEach(`th.navbox-group:not([style="text-align:center; width:5%"])`,
			func(_ int, th *colly.HTMLElement) { group = th.Text })

		// 获取 Type 列表，有可能当前 group 没有 type
		tr.ForEach(`th.navbox-group[style="text-align:center; width:5%"]`,
			func(_ int, th *colly.HTMLElement) { types = append(types, th.Text) })

		if len(types) == 0 {
			// 当 group 没有 type 时，此时的一行即
			tr.ForEach(`td:has(div:not(:empty)[style="padding:0em 0.25em"]) .smw-value`,
				func(_ int, e *colly.HTMLElement) {
					if e.Text == i.Name {
						g = group
					}
				})
		} else {
			// 二级元素
			tr.ForEach("td:has(table) table tr:has(th)", func(ii int, tr *colly.HTMLElement) {
				tr.ForEach("td li .smw-value", func(_ int, e *colly.HTMLElement) {
					if e.Text == i.Name {
						g = group
					}
				})
			})
		}
	})
	_ = c.Visit("https://prts.wiki/w/%E7%90%86%E6%99%BA")
	return g
}
func (i *Item) FreshGroup() {
	i.Group = i.FetchGroup(wiki.NewCollector())
}
func (i *Item) FreshGroupIfEmpty() {
	if i.Group != "" {
		i.FreshGroup()
	}
}
