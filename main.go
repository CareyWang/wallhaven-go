package main

import (
	"flag"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	tag := flag.String("tag", "artwork", "分类")
	ps := flag.Int("ps", 1, "开始页")
	pe := flag.Int("pe", 1, "结束页")
	userProxy := flag.Bool("useProxy", false, "是否使用代理")
	flag.Parse()

	fmt.Println(*tag, *ps, *pe, *userProxy)

	var detailPages []string

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36 Edg/83.0.478.45"),
	)
	if *userProxy {
		setProxy(c)
	}

	c.OnHTML(".preview", func(c *colly.HTMLElement) {
		// 详情页 url
		detailUrl := c.Attr("href")
		detailPages = append(detailPages, detailUrl)
	})

	for page := *ps; page <= *pe; page++ {
		toVisit := fmt.Sprintf("https://wallhaven.cc/search?q=%s&page=%d", *tag, page)
		fmt.Println(toVisit)

		_ = c.Visit(toVisit)
	}

	count := 1
	dirCount := 1
	for _, detailPage := range detailPages {
		imageUrls := getImageUrls(c.Clone(), detailPage)

		log.Println(imageUrls)

		if imageUrls != nil {
			for _, imageUrl := range imageUrls {
				// 每100个分个文件夹
				path := fmt.Sprintf("images/%d/", dirCount)
				if count%100 == 0 {
					dirCount++
				}
				count++

				// 下载图片
				go download(path, imageUrl)
				time.Sleep(time.Millisecond * 100)
			}
		}
		time.Sleep(time.Second)
	}

}

// 获取详情页链接
func getDetailPageUrl(c colly.Collector) string {
	var detailUrl string

	c.OnHTML(".preview", func(c *colly.HTMLElement) {
		// 详情页 url
		detailUrl = c.Attr("href")
		log.Print(detailUrl)
	})

	return detailUrl
}

// 获取详情页图片 url
func getImageUrls(c *colly.Collector, pageUrl string) []string {
	var imageUrls []string

	c.OnHTML("#wallpaper", func(c *colly.HTMLElement) {
		// 图片 url
		imageUrl := c.Attr("src")
		imageUrls = append(imageUrls, imageUrl)
	})

	_ = c.Visit(pageUrl)
	return imageUrls
}

// 设置代理
func setProxy(c *colly.Collector) {
	p, err := proxy.RoundRobinProxySwitcher(
		"socks5://127.0.0.1:7891",
		"http://127.0.0.1:7890",
	)

	if err != nil {
		c.SetProxyFunc(p)
	}
}

// 下载文件
func download(path, imageUrl string) {
	// Get the data
	resp, err := http.Get(imageUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	_sl := strings.Split(imageUrl, "/")
	imageName := _sl[len(_sl)-1]

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	// 创建一个文件用于保存
	out, err := os.Create(path + imageName)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
}
