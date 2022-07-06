package cmd

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strings"
	"time"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "更新数据",
	RunE: func(cmd *cobra.Command, args []string) error {
		initDB()
		// updateItems()
		res := DB.Find(&itemsCollection)
		fmt.Printf("共查询到 %d 条 items 记录\n", res.RowsAffected)

		c := colly.NewCollector(colly.Async(true))
		c.OnHTML("body", func(body *colly.HTMLElement) {
			var name, image, wikishort string
			var updatedAt time.Time
			body.ForEach("#firstHeading", func(i int, e *colly.HTMLElement) {
				name = e.Text
				name = strings.ReplaceAll(name, "α", "Α")
				name = strings.ReplaceAll(name, "β", "Β")
				name = strings.ReplaceAll(name, "γ", "Γ")
				// 	fmt.Printf("%s ", e.Text)
			})
			body.ForEach(".nomobile > a.image > img.lazyload", func(i int, e *colly.HTMLElement) {
				image = link(e.Attr("data-src"))
				// fmt.Printf("图片链接：%s\n", link(e.Attr("data-src")))
			})
			body.ForEach(".copyUrl", func(i int, e *colly.HTMLElement) {
				wikishort = e.Attr("data-clipboard-text")
			})
			body.ForEach("#footer-info-lastmod", func(i int, e *colly.HTMLElement) {
				// 此页面最后编辑于2022年5月22日 (星期日) 12:32。
				weeks := []string{"日", "一", "二", "三", "四", "五", "六"}
				var err error
				var t time.Time
				for i := 0; i <= 7; i++ {
					layout := fmt.Sprintf("此页面最后编辑于2006年1月2日 (星期%s) 15:04 -0700", weeks[i])
					t, err = time.Parse(layout, strings.ReplaceAll(strings.TrimSpace(e.Text), "。", " +0800"))
					if err == nil {
						break
					}
				}
				if t.IsZero() {
					fmt.Printf("error time is: %s\n", e.Text)
					return
				}
				updatedAt = t
			})

			var updated bool
			if updatedAt.IsZero() {
				updated = UpdateItemInColumns(
					&Item{Name: name, Image: image, WikiShort: wikishort, UpdatedAt: updatedAt},
					[]string{"name", "image", "wiki_short"})
			} else {
				updated = UpdateItemInColumns(
					&Item{Name: name, Image: image, WikiShort: wikishort, UpdatedAt: updatedAt},
					[]string{"name", "image", "wiki_short", "updated_at"})
			}
			// fmt.Printf("%s \n\t短链接: %s\n\t图片链接: %s\n", name, wikishort, image)
			if !updated {
				fmt.Printf("Error: %s 更新失败\n", name)
			} else {
				fmt.Printf("%s 更新完成\n", name)
			}
		})
		_ = c.Limit(&colly.LimitRule{
			DomainGlob:  "*prts.wiki*",
			Parallelism: 15,
			Delay:       500 * time.Millisecond,
		})
		for _, item := range itemsCollection {
			// c.OnRequest(func(r *colly.Request) {
			// r.Ctx.Put("name", item.Name)
			// r.Ctx.Put("name", itemsCollection[0].Name)
			// })
			_ = c.Visit(item.Wiki)
		}
		// c.OnHTML(".copyUrl", func(e *colly.HTMLElement) {
		// 	name := e.Response.Ctx.Get("name")
		// 	fmt.Printf("%s 短链接：%s\n", name, e.Attr("data-clipboard-text"))
		// })
		_ = c.Visit(itemsCollection[0].Wiki)
		c.Wait()

		return nil
	},
}

var DB *gorm.DB

func initDB() {
	var err error
	if DB, err = gorm.Open(sqlite.Open("arknights.db"), &gorm.Config{}); err != nil {
		fmt.Println("数据库初始化异常", err.Error())
	}
	_ = DB.AutoMigrate(&Item{})
}

type Item struct {
	Name      string    `json:"name" gorm:"column:name"`
	Image     string    `json:"image" gorm:"column:image"`
	Wiki      string    `json:"wiki" gorm:"column:wiki"`
	WikiShort string    `json:"wiki_short" gorm:"column:wiki_short"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:false"`
}

func (m Item) TableName() string {
	return "items"
}

func CreateItemIfNotExists(item *Item) {
	var cnt int64
	DB.Where("name = ?", item.Name).Count(&cnt)
	if cnt == 0 {
		DB.Create(&item)
	}
}

func UpdateItemInColumns(item *Item, cols []string) bool {
	res := DB.Model(&Item{}).Select(cols).Where("name", item.Name).Updates(&item)
	return res.RowsAffected > 0
}

func updateMaterials() {
	c := colly.NewCollector()
	c.OnHTML(".mw-category-group ul li a", func(a *colly.HTMLElement) {
		fmt.Printf("%s\n", a.Text)
	})
	_ = c.Visit("https://prts.wiki/w/分类:材料")
}

var itemsCollection []Item

func updateItems() {

	// 1. 到分类:道具页面找到所有的道具名称和链接
	c := colly.NewCollector()
	count := 0
	c.OnHTML(".mw-category-group ul li a", func(a *colly.HTMLElement) {
		count++
		item := Item{
			Name: strings.TrimSpace(a.Text),
			Wiki: link(a.Attr("href")),
		}
		// itemsCollection = append(itemsCollection, item)
		CreateItemIfNotExists(&item)
		fmt.Printf("No.%d %s %s\n", count, a.Text, link(a.Attr("href")))
	})
	c.OnHTML(`a[title="分类:道具"]:last-of-type`, func(element *colly.HTMLElement) {
		cc := colly.NewCollector()
		cc.OnHTML(".mw-category-group ul li a", func(a *colly.HTMLElement) {
			count++
			item := Item{
				Name: strings.TrimSpace(a.Text),
				Wiki: link(a.Attr("href")),
			}
			// itemsCollection = append(itemsCollection, item)
			CreateItemIfNotExists(&item)
			fmt.Printf("No.%d %s %s\n", count, a.Text, link(a.Attr("href")))
		})
		cc.OnHTML(`a[title="分类:道具"]:last-of-type`, func(element *colly.HTMLElement) {
			ccc := colly.NewCollector()
			ccc.OnHTML(".mw-category-group ul li a", func(a *colly.HTMLElement) {
				count++
				item := Item{
					Name: strings.TrimSpace(a.Text),
					Wiki: link(a.Attr("href")),
				}
				// itemsCollection = append(itemsCollection, item)
				CreateItemIfNotExists(&item)
				fmt.Printf("No.%d %s %s\n", count, a.Text, link(a.Attr("href")))
			})
			_ = ccc.Visit(link(element.Attr("href")))
		})
		_ = cc.Visit(link(element.Attr("href")))
	})
	_ = c.Visit("https://prts.wiki/w/分类:道具")
	fmt.Printf("共找到 %d 条数据\n\n", count)
}

func link(uri string) string {
	return "https://prts.wiki" + uri
}
