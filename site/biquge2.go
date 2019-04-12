package site

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/ma6254/FictionDown/store"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var biquge2 = SiteA{
	Name:     "笔趣阁5200",
	HomePage: "http://www.bqg5200.com/",
	Match: []string{
		`https://www\.bqg5200\.com/xiaoshuo/\d+/\d+/*`,
		`https://www\.bqg5200\.com/xiaoshuo/\d+/\d+/\d+\.html`,
	},
	BookInfo: func(body io.Reader) (s *store.Store, err error) {
		body = transform.NewReader(body, simplifiedchinese.GBK.NewDecoder())
		doc, err := htmlquery.Parse(body)
		if err != nil {
			return
		}

		s = &store.Store{}

		s.BookName = htmlquery.InnerText(htmlquery.FindOne(doc, `//*[@id="smallcons"]/h1`))

		// var author = htmlquery.Find(doc, `//*[@id="info"]/p[1]`)
		s.Author = strings.TrimSpace(htmlquery.InnerText(htmlquery.FindOne(doc, `//*[@id="smallcons"]/span[1]`)))

		node_content := htmlquery.Find(doc, `//*[@id="readerlist"]/ul/li/a`)
		if len(node_content) == 0 {
			err = fmt.Errorf("No matching contents")
			return
		}

		var vol = store.Volume{
			Name:     "正文",
			Chapters: make([]store.Chapter, 0),
		}

		for _, v := range node_content {
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

		return
	},
	Chapter: func(body io.Reader) ([]string, error) {
		body = transform.NewReader(body, simplifiedchinese.GBK.NewDecoder())
		doc, err := htmlquery.Parse(body)
		if err != nil {
			return nil, err
		}

		M := []string{}
		//list
		nodeContent := htmlquery.Find(doc, `//*[@id="content"]/text()`)
		if len(nodeContent) == 0 {
			err = fmt.Errorf("No matching content")
			return nil, err
		}
		for _, v := range nodeContent {
			t := htmlquery.InnerText(v)
			t = strings.TrimSpace(t)

			if strings.HasPrefix(t, "…") {
				continue
			}

			t = strings.Replace(t, "…", "", -1)
			t = strings.Replace(t, ".......", "", -1)
			t = strings.Replace(t, "...”", "”", -1)

			if t == "" {
				continue
			}

			M = append(M, t)
		}
		return M, nil
	},
}
