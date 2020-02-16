package site

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/ma6254/FictionDown/utils"
)

func init() {
	addSite(SiteA{
		Name:     "书迷楼",
		HomePage: "http://www.shumil.co/",
		Match: []string{
			`http\??://www\.shumil\.co/\w+/`,
			`http\??://www\.shumil\.co/\w+/\d+\.html`,
		},
		BookInfo: Type1BookInfo(
			`//div[@class="content"]/div[@class="list"]/div[@class="tit"]/b`,
			``,
			`//a[starts-with(@href, "/zuozhe/")]/text()`,
			`//div[@class="content"]/div[@class="list"]/ul/li/a`),
		Chapter: Type1Chapter(`//*[@id="content"]/p[1]/text()`),
		Search: Type1SearchAfter("http://www.shumil.co/search.php",
			func(s string) *http.Request {
				baseurl, err := url.Parse("http://www.shumil.co/search.php")
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
			func(r ChaperSearchResult) ChaperSearchResult {
				r.Author = strings.TrimPrefix(r.Author, "/")
				return r
			},
		),
	})
}
