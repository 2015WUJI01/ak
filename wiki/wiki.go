package wiki

import (
	"io/ioutil"
	"main/pkg/logger"
	"net/http"
	"net/url"
)

type Wiki struct {
}

func Domain() string {
	return "https://prts.wiki"
}

func Page(uri string) string {
	return Domain() + uri
}

func OperatorPage(name string) string {
	return Page("/w/" + url.QueryEscape(name))
}

// GetRespBodyFromItemsAPI 从企鹅物流 API 拿数据
func GetRespBodyFromItemsAPI() ([]byte, error) {
	// 从企鹅物流 API 拿到数据
	resp, err := http.Get("https://penguin-stats.io/PenguinStats/api/v2/items")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读出 body 内容，直接返回
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return body, nil
}
