package com_ddyueshu

import (
	"net/http"
	"net/url"

	"github.com/ma6254/FictionDown/site"
)

func Site() site.SiteA {
	return site.SiteA{
		Name:     "顶点小说",
		HomePage: "https://www.ddyueshu.com/",
		Tags: func() []string {
			return []string{
				"盗版",
				"一般书源",
				"PTCMS",
				"顶点小说",
			}
		},
		Match: []string{
			`https://www\.ddyueshu\.com/\d+_\d+/*`,
			`https://www\.ddyueshu\.com/\d+_\d+/\d+\.html/*`,
			`http://www\.ddyueshu\.com/book/goto/id/\d+`,
		},
		BookInfo: site.Type1BookInfo(
			`//div[@id="info"]/h1`,
			`//*[@id="fmimg"]/img`,
			`//meta[@property="og:novel:author"]/@content`,
			`//*[@id="list"]/dl/dd/a`,
		),
		Chapter: site.Type1Chapter(`//*[@id="content"]/text()`),
		Search: site.Type1Search("",
			func(s string) *http.Request {
				baseurl, err := url.Parse("https://so.biqusoso.com/s1.php")
				if err != nil {
					panic(err)
				}
				value := baseurl.Query()
				value.Add("ie", "utf-8")
				value.Add("siteid", "ddyueshu.com")
				value.Add("q", s)
				baseurl.RawQuery = value.Encode()

				req, err := http.NewRequest("GET", baseurl.String(), nil)
				if err != nil {
					panic(err)
				}
				return req
			},
			`//div[@class="search-list"]/ul/li[position()>1]`,
			`*[@class="s2"]/a`,
			`*[@class="s4"]`,
		),
	}
}
