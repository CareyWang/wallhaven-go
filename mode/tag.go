package mode

import (
	"github.com/gocolly/colly"
	"strings"
)

const tagListPage = "https://wallhaven.cc/tags"
const tagSearchCommonUrl = "https://wallhaven.cc/search?q=id:"


// tag 搜索页
func GetAllTagSearchUrls() []string {
	homePages := getAllTagHomePages()

	if homePages == nil {
		return nil
	}

	return makeTagSearchUrl(homePages)
}

func GetSearchUrlsByTagId(tagId string) string {
	return tagSearchCommonUrl + tagId
}

// tag 主页
func getAllTagHomePages() []string {

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36 Edg/83.0.478.45"),
	)

	var homePages []string
	c.OnHTML("#taglist", func(c *colly.HTMLElement) {
		tagUrl := c.ChildAttr("div>span>a", "href")

		homePages = append(homePages, tagUrl)
	})

	_ = c.Visit(tagListPage)

	return homePages
}

// 构造 tag 搜索页链接
func makeTagSearchUrl(tagUrls []string) []string {
	var tagSearchUrls []string

	for _, tagUrl := range tagUrls {
		_sl := strings.Split(tagUrl, "/")
		tagId := _sl[len(_sl)-1]

		tagSearchUrls = append(tagSearchUrls, tagSearchCommonUrl+tagId)
	}

	return tagSearchUrls
}
