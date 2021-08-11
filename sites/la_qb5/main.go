package la_qb5

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/ma6254/FictionDown/site"
	"github.com/ma6254/FictionDown/utils"
)

func Site() site.SiteA {
	return site.SiteA{
		Name:     "全本小说网",
		HomePage: "https://www.qb5.la/",
		Tags:     func() []string { return []string{"盗版", "优质书源"} },
		Match: []string{
			`https://www\.qb5\.la/book_\d+/`,
			`https://www\.qb5\.la/book_\d+/\d+\.html`,
		},
		BookInfo: site.Type1BookInfo(
			`//*[@id="info"]/h1/text()`,
			`//*[@id="picbox"]/div/img`,
			`//*[@id="info"]/h1/small/a/text()`,
			`//div[@class="zjbox"]/dl[@class="zjlist"]/dd/a`),
		Chapter: site.Type1Chapter(`//*[@id="content"]/text()`),
		Search: site.Type1SearchAfter("https://www.qb5.la/modules/article/search.php",
			func(s string) *http.Request {
				baseurl, err := url.Parse("https://www.qb5.la/modules/article/search.php")
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
