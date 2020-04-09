package biquge5200_cc

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	fcontext "github.com/ma6254/FictionDown/context"

	"github.com/antchfx/htmlquery"
	"github.com/ma6254/FictionDown/site"
	"github.com/ma6254/FictionDown/store"
	"github.com/ma6254/FictionDown/utils"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func Site() site.SiteA {
	return site.SiteA{
		Name:     "笔趣阁1",
		HomePage: "https://www.biquge5200.cc/",
		Tags: func() []string {
			return []string{
				"盗版",
				"一般书源",
				"PTCMS",
				"笔趣阁",
			}
		},
		Match: []string{
			`https://www\.biquge5200\.cc/\d+_\d+/*`,
			`https://www\.biquge5200\.cc/\d+_\d+/\d+\.html/*`,
		},
		BookInfo: func(body io.Reader) (s *store.Store, err error) {
			doc, err := htmlquery.Parse(body)
			if err != nil {
				return
			}

			s = &store.Store{}

			nodeTitle := htmlquery.Find(doc, `//*[@id="info"]/h1`)
			if len(nodeTitle) == 0 {
				err = fmt.Errorf("No matching title")
				return
			}
			s.BookName = htmlquery.InnerText(nodeTitle[0])

			nodeDesc := htmlquery.Find(doc, `//*[@id="intro"]/p`)
			if len(nodeDesc) == 0 {
				err = fmt.Errorf("No matching desc")
				return
			}
			s.Description = strings.Replace(
				htmlquery.OutputHTML(nodeDesc[0], false),
				"<br/>", "\n",
				-1)

			var author = htmlquery.Find(doc, `//*[@id="info"]/p[1]`)
			s.Author = strings.TrimLeft(htmlquery.OutputHTML(author[0], false), "作\u00a0\u00a0\u00a0\u00a0者：")

			nodeContent := htmlquery.Find(doc, `//*[@id="list"]/dl/dd/a`)
			if len(nodeDesc) == 0 {
				err = fmt.Errorf("No matching contents")
				return
			}

			var vol = store.Volume{
				Name:     "正文",
				Chapters: make([]store.Chapter, 0),
			}

			for _, v := range nodeContent[9:] {
				//fmt.Printf("href: %v\n", chapter_u)
				chapterURL, err := url.Parse(htmlquery.SelectAttr(v, "href"))
				if err != nil {
					return nil, err
				}

				vol.Chapters = append(vol.Chapters, store.Chapter{
					Name: strings.TrimSpace(htmlquery.InnerText(v)),
					URL:  chapterURL.String(),
				})
			}
			s.Volumes = append(s.Volumes, vol)

			s.CoverURL = htmlquery.SelectAttr(htmlquery.FindOne(doc, `//*[@id="fmimg"]/img`), "src")

			return
		},
		Chapter: func(ctx fcontext.Context) (content []string, err error) {
			doc, err := htmlquery.Parse(ctx.Value(fcontext.KeyBody).(io.Reader))
			if err != nil {
				return nil, err
			}

			M := []string{}
			//list
			// nodeContent := htmlquery.Find(doc, `//div[@id="content"]/text()`)
			nodeContent := htmlquery.Find(doc, `//div[@id="content"]/p`)
			if len(nodeContent) == 0 {
				err = fmt.Errorf("No matching content")
				return nil, err
			}
			for _, v := range nodeContent {
				t := htmlquery.InnerText(v)
				t = strings.TrimSpace(t)
				M = append(M, t)
			}

			return M, nil
		},
		Search: func(s string) (result []site.ChaperSearchResult, err error) {
			baseurl, err := url.Parse("https://www.biquge5200.cc/modules/article/search.php")
			if err != nil {
				return
			}
			value := baseurl.Query()
			value.Add("searchkey", s)
			baseurl.RawQuery = value.Encode()

			resp, err := utils.RequestGet(baseurl.String())
			if err != nil {
				return
			}
			defer resp.Body.Close()
			body := transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())
			doc, err := htmlquery.Parse(body)
			if err != nil {
				return
			}

			// fmt.Printf("%s", htmlquery.OutputHTML(doc, true))
			r := htmlquery.Find(doc, `//*[@id="hotcontent"]/table/tbody/tr`)
			if len(r) == 0 {
				return nil, nil
			}

			for _, v := range r[1:] {
				a := htmlquery.FindOne(v, `/*[1]/a`)
				r := site.ChaperSearchResult{
					BookName: htmlquery.InnerText(a),
					Author:   htmlquery.InnerText(htmlquery.FindOne(v, `/*[3]`)),
					BookURL:  htmlquery.SelectAttr(a, "href"),
				}
				result = append(result, r)
			}
			return
		},
	}
}
