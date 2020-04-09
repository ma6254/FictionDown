package com_38kanshu

import (
	"log"
	"net/url"
	"strings"

	"golang.org/x/net/html"

	"github.com/antchfx/htmlquery"
	"github.com/ma6254/FictionDown/site"
	"github.com/ma6254/FictionDown/utils"
)

func Site() site.SiteA {
	return site.SiteA{
		Name:     "38看书",
		HomePage: "https://www.38kanshu.com/",
		Tags:     func() []string { return []string{"盗版", "优质书源"} },
		Match: []string{
			`https://www\.38kanshu\.com/\d+/`,
			`https://www\.38kanshu\.com/\d+/\d+\.html`,
		},
		BookInfo: site.Type1BookInfo(
			`//*[@class="bookPhr"]/h2`,
			`//*[@class="bookImg"]/img/@src`,
			`//meta[@property="og:novel:author"]/@content`,
			`//*[@class="chapterCon"]/ul/li/a`,
		),
		Chapter: site.Type2Chapter(
			`//*[@class="articleCon"]/p/*/text() | //*[@class="articleCon"]/p/text()`,
			func(preURL *url.URL, doc *html.Node) *html.Node {
				nextNode := htmlquery.FindOne(doc, `//*[@class="page"]/a[3]`)
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
					n, err := preURL.Parse(nextURL)
					if err != nil {
						return nil
					}
					doc, err := utils.GetWegPageDOM(n.String())
					if err != nil {
						log.Printf("GetWegPageDOM: %s", err)
						return nil
					}
					return doc
				}
				return nil
			}, nil),
	}
}
