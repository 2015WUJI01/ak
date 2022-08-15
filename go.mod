module main

go 1.16

replace github.com/2015WUJI01/looog v0.0.0-20220807104046-3ec4caaa1860 => ./pkg/looog

require (
	github.com/PuerkitoBio/goquery v1.8.0
	github.com/gocolly/colly v1.2.0
	github.com/spf13/cobra v1.5.0
	go.uber.org/zap v1.21.0
	gorm.io/driver/mysql v1.3.2
	gorm.io/driver/sqlite v1.3.6
	gorm.io/gorm v1.23.4
)

require (
	github.com/2015WUJI01/looog v0.0.0-20220807104046-3ec4caaa1860
	github.com/antchfx/htmlquery v1.2.4 // indirect
	github.com/antchfx/xmlquery v1.3.10 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/kennygrant/sanitize v1.2.4 // indirect
	github.com/rivo/uniseg v0.3.1 // indirect
	github.com/saintfish/chardet v0.0.0-20120816061221-3af4cd4741ca // indirect
	github.com/schollz/progressbar/v3 v3.8.7
	github.com/stretchr/testify v1.7.0
	github.com/temoto/robotstxt v1.1.2 // indirect
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa // indirect
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f // indirect
	golang.org/x/sys v0.0.0-20220730100132-1609e554cd39 // indirect
	golang.org/x/term v0.0.0-20220722155259-a9ba230a4035 // indirect
	google.golang.org/appengine v1.6.7 // indirect
)
