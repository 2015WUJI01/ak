package models

import (
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

func (m Item) TableName() string {
	return "items"
}

func (m Item) FetchWiki(c *colly.Collector) (w string) {
	// 获取第一轮的 200 条数据
	c.OnHTML(".mw-category-group ul li a", func(a *colly.HTMLElement) {
		if strings.TrimSpace(a.Text) == m.Name {
			w = a.Attr("href")
		}
	})
	c.OnHTML(`a[title="分类:道具"]:last-of-type`, func(e *colly.HTMLElement) {
		cc := wiki.NewCollector()
		// 获取第二轮的 200 条数据
		cc.OnHTML(".mw-category-group ul li a", func(a *colly.HTMLElement) {
			if strings.TrimSpace(a.Text) == m.Name {
				w = a.Attr("href")
			}
		})
		// 获取第三页链接
		cc.OnHTML(`a[title="分类:道具"]:last-of-type`, func(ee *colly.HTMLElement) {
			ccc := wiki.NewCollector()
			// 获取第三页的 200 条数据
			ccc.OnHTML(".mw-category-group ul li a", func(a *colly.HTMLElement) {
				if strings.TrimSpace(a.Text) == m.Name {
					w = a.Attr("href")
				}
			})
			_ = ccc.Visit(wiki.Link(ee.Attr("href")))
		})
		_ = cc.Visit(wiki.Link(e.Attr("href")))
	})
	_ = c.Visit("https://prts.wiki/w/分类:道具")
	return wiki.Link(w)
}
func (m *Item) FreshWiki() {
	m.Wiki = m.FetchWiki(wiki.NewCollector())
}
func (m *Item) FreshWikiIfEmpty() {
	if m.Wiki != "" {
		m.FreshWiki()
	}
}

func (m Item) FetchGroup(c *colly.Collector) (g string) {
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
					if e.Text == m.Name {
						g = group
					}
				})
		} else {
			// 二级元素
			tr.ForEach("td:has(table) table tr:has(th)", func(i int, tr *colly.HTMLElement) {
				tr.ForEach("td li .smw-value", func(_ int, e *colly.HTMLElement) {
					if e.Text == m.Name {
						g = group
					}
				})
			})
		}
	})
	_ = c.Visit("https://prts.wiki/w/%E7%90%86%E6%99%BA")
	return g
}
func (m *Item) FreshGroup() {
	m.Group = m.FetchGroup(wiki.NewCollector())
}
func (m *Item) FreshGroupIfEmpty() {
	if m.Group != "" {
		m.FreshGroup()
	}
}
