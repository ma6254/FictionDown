package shumil_co

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/ma6254/FictionDown/site"
	"github.com/ma6254/FictionDown/utils"
)

func Site() site.SiteA {
	return site.SiteA{
		Name:     "书迷楼",
		HomePage: "http://www.shumil.co/",
		Tags:     func() []string { return []string{"盗版", "优质书源"} },
		Match: []string{
			`https://www\.shumil\.co/\w+/`,
			`https://www\.shumil\.co/\w+/\d+\.html`,
		},
		BookInfo: site.Type1BookInfo(
			`//div[@class="content"]/div[@class="list"]/div[@class="tit"]/b`,
			``,
			`//a[starts-with(@href, "/zuozhe/")]/text()`,
			`//div[@class="content"]/div[@class="list"]/ul/li/a`),
		Chapter: site.Type1Chapter(`//*[@id="content"]/p[1]/text()`),
		Search: site.Type1SearchAfter("https://www.shumil.co/search.php",
			func(s string) *http.Request {
				baseurl, err := url.Parse("https://www.shumil.co/search.php")
				if err != nil {
					panic(err)
				}
				value := baseurl.Query()
				value.Add("searchtype", "all")
				value.Add("searchkey", utils.U8ToGBK(s))
				value.Add("sbt", utils.U8ToGBK("搜索"))
				baseurl.RawQuery = value.Encode()

				req, err := http.NewRequest("GET", baseurl.String(), nil)
				if err != nil {
					panic(err)
				}
				return req
			},
			`//div[@class="content"]/div[@class="list"]/ul/li`,
			`a`,
			`text()`,
			func(r site.ChaperSearchResult) site.ChaperSearchResult {
				r.Author = strings.TrimPrefix(r.Author, "/")
				return r
			},
		),
	}
}
