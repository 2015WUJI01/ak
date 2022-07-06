package wiki

import "C"
import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"regexp"
	"strings"
	"time"
)

type Scraper struct {
	C            *colly.Collector
	Debug        bool
	RetryCounter int
	URL          string
}

func DefaultScraper(url string) *Scraper {
	c := colly.NewCollector(
		colly.AllowedDomains("prts.wiki"),
	)
	return &Scraper{
		C:   c,
		URL: url,
	}
}

func (s *Scraper) SetDebug(debug bool) *Scraper {
	s.Debug = debug
	return s
}

func (s *Scraper) Visit() error {
	return s.C.Visit(s.URL)
}
func (s *Scraper) PrintBody() *Scraper {
	s.C.OnResponse(func(r *colly.Response) {
		log.Println(string(r.Body))
	})
	return s
}

func (s *Scraper) AutoRetry(t time.Duration) *Scraper {
	s.C.OnError(func(r *colly.Response, err error) {
		time.Sleep(t)
		s.RetryCounter++
		if s.Debug {
			log.Println(fmt.Sprintf("访问 %s 失败，正在第 %d 次重试", r.Request.URL, s.RetryCounter))
		}
		_ = r.Request.Retry()
	})
	return s
}

// GetOprEquip 查询模组
func (s *Scraper) GetOprEquip(hasEquip *bool, equipNames *[]string, equipNum *int) *Scraper {
	// 要用正则来提出模组了
	s.C.OnResponse(func(response *colly.Response) {
		var r *regexp.Regexp
		// 查找是否开启模组，并添加到模组数组中
		r = regexp.MustCompile(`var equipItems = .*\n.*{default:(.*)},\n.*label: '(.*)'`)
		for _, items := range r.FindAllStringSubmatch(string(response.Body), -1) {
			if items[1] == "true" {
				*hasEquip = true
			}
			defaultEquip := strings.TrimSpace(items[2])
			*equipNames = append(*equipNames, defaultEquip)
		}
		// 查找其他模组名称，添加到模组数组中
		r = regexp.MustCompile(`equipNames = \["(.*)",\n.*"(.*)",\n.*"(.*)",\n.*"(.*)",\n.*"(.*)"]`)
		for _, items := range r.FindAllStringSubmatch(string(response.Body), -1) {
			for i, str := range items {
				if i > 0 && strings.TrimSpace(str) != "" {
					*equipNames = append(*equipNames, strings.TrimSpace(str))
				}
			}
		}
		if *hasEquip {
			*equipNum = len(*equipNames)
		}
	})
	return s
}

// 技能升级材料
// func (s *Scraper) FetchOprSkillUpgradeMaterials(skillMaterials, subClass *string) *Scraper {
// 	s.C.OnHTML("h2:has(span#技能升级材料) + table tbody tr", func(tr *colly.HTMLElement) {
// 		if tr.Index == 1 || tr.Index == 3 {
// 			tr.ForEach("td", func(idx int, td *colly.HTMLElement) {
// 				level := tr.Index*4/3 + idx
// 				var skillMaterials []SkillMaterial
// 				td.ForEach("div", func(idx int, div *colly.HTMLElement) {
// 					title := div.ChildAttr("a", "title")
// 					amount, _ := strconv.Atoi(div.ChildText("span"))
// 					skillMaterial := SkillMaterial{
// 						Material: *ItemNameList[title],
// 						Amount:   uint(amount),
// 					}
// 					skillMaterials = append(skillMaterials, skillMaterial)
// 				})
// 				o.SkillUpgradeMaterials[level-1] = skillMaterials
// 			})
// 		}
// 	})
// }
