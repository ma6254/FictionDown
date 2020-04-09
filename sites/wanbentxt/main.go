package wanbentxt

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/ma6254/FictionDown/site"
	"github.com/ma6254/FictionDown/utils"
	"golang.org/x/net/html"
)

func Site() site.SiteA {
	return site.SiteA{
		Name:     "完本神站",
		HomePage: "https://www.wanbentxt.com/",
		Tags:     func() []string { return []string{"盗版", "优质书源"} },
		Match: []string{
			`https://www\.wanbentxt\.com/\d+/`,
			`https://www\.wanbentxt\.com/\d+/\d+\.html`,
		},
		BookInfo: site.Type1BookInfo(
			`//div[@class="detailTitle"]/h1/text()`,
			`//div[@class="detailTopLeft"]/img/@src`,
			`//div[@class="detailTopMid"]/div[@class="writer"]/a/text()`,
			`//div[@class="chapter"]/ul/li/a`),
		Chapter: site.Type2Chapter(`//div[@class="readerCon"]/p/text()`, func(preURL *url.URL, doc *html.Node) *html.Node {
			nextNode := htmlquery.FindOne(doc, `//div[@class="readPage"]/a[3]`)
			if nextNode == nil {
				return nil
			}
			nextText := htmlquery.InnerText(nextNode)
			// log.Printf("nextText: %v\n", nextText)
			if strings.Contains(nextText, "下一章") {
				return nil
			} else if strings.Contains(nextText, "下一页") {
				nextURL := htmlquery.SelectAttr(nextNode, "href")
				// log.Printf("nextURL: %v\n", nextURL)
				doc, err := utils.GetWegPageDOM(nextURL)
				if err != nil {
					log.Printf("GetWegPageDOM: %s", err)
					return nil
				}
				return doc
			}
			return nil
		}, func(b []string) []string {
			if strings.HasPrefix(b[0], "一秒记住") {
				b = b[1:]
			}
			return b[:len(b)-1]
		}),
		Search: site.Type1Search(
			"",
			func(s string) *http.Request {
				baseurl, err := url.Parse("https://www.wanbentxt.com/modules/article/search.php")
				if err != nil {
					panic(err)
				}
				value := baseurl.Query()
				value.Add("searchtype", "all")
				value.Add("searchkey", utils.U8ToGBK(s))
				baseurl.RawQuery = value.Encode()

				req, err := http.NewRequest("GET", baseurl.String(), nil)
				if err != nil {
					panic(err)
				}
				return req
			},
			`//div[@class="result"]/div[@class="resultLeft"]/ul/li/div[@class="sortPhr"]`,
			`a`,
			`p[@class="author"]/a/text()`),
	}
}
