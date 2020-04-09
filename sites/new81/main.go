package new81

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"

	fcontext "github.com/ma6254/FictionDown/context"

	"github.com/antchfx/htmlquery"
	"github.com/ma6254/FictionDown/site"
	"github.com/ma6254/FictionDown/store"
	"github.com/ma6254/FictionDown/utils"
	"golang.org/x/text/transform"
)

func Site() site.SiteA {
	return site.SiteA{
		Name:     "新八一中文网",
		HomePage: "https://www.81new.net/",
		Tags:     func() []string { return []string{"盗版", "优质书源"} },
		Match: []string{
			`https://www\.81new\.net/\d+/\d+/`,
			`https://www\.81new\.net/\d+/\d+/d+\.html`,
		},
		BookInfo: wwww81newcomBookInfo,
		Chapter: func(ctx fcontext.Context) (content []string, err error) {
			doc, err := htmlquery.Parse(ctx.Value(fcontext.KeyBody).(io.Reader))
			if err != nil {
				return nil, err
			}

			M := []string{}
			//list
			// nodeContent := htmlquery.Find(doc, `//div[@id="content"]/text()`)
			nodeContent := htmlquery.Find(doc, `//*[@id="articlecontent"]/text()`)
			if len(nodeContent) == 0 {
				err = fmt.Errorf("No matching content")
				return nil, err
			}
			for _, v := range nodeContent {
				t := htmlquery.InnerText(v)
				t = strings.TrimSpace(t)

				switch t {
				case
					"[八一中文网 请记住",
					"手机版访问 m.81new.com 绿色无弹窗]",
					"":
					continue
				}

				M = append(M, t)
			}

			return M, nil
		},
		Search: func(s string) (result []site.ChaperSearchResult, err error) {
			baseurl, err := url.Parse("https://www.81new.net/modules/article/search.php")
			if err != nil {
				return
			}
			value := baseurl.Query()
			value.Add("searchkey", utils.U8ToGBK(s))
			baseurl.RawQuery = value.Encode()

			// Get WebPage
			resp, err := utils.RequestGet(baseurl.String())
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			var body io.Reader = bytes.NewReader(bodyBytes)
			encode := utils.DetectContentCharset(bytes.NewReader(bodyBytes))
			body = transform.NewReader(body, encode.NewDecoder())

			if regexp.MustCompile(`/modules/article/search\.php`).MatchString(resp.Request.URL.Path) {
				// 多个搜索结果
				doc, err := htmlquery.Parse(body)
				if err != nil {
					return nil, err
				}
				r := htmlquery.Find(doc, `//table[@id="author"]/tbody/tr`)
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
			} else if regexp.MustCompile(`/\d+/\d+/*`).MatchString(resp.Request.URL.Path) {
				// 单个搜索结果
				store, err := wwww81newcomBookInfo(body)
				if err != nil {
					return nil, err
				}
				result = append(result, site.ChaperSearchResult{
					BookName: store.BookName,
					Author:   store.Author,
					BookURL:  resp.Request.URL.String(),
				})
			}

			return
		},
	}
}

func wwww81newcomBookInfo(body io.Reader) (s *store.Store, err error) {
	doc, err := htmlquery.Parse(body)
	if err != nil {
		return
	}

	s = &store.Store{}

	// Book Name
	node_title := htmlquery.Find(doc, `//div[@class="introduce"]/h1`)
	if len(node_title) == 0 {
		err = fmt.Errorf("No matching title")
		return
	}
	s.BookName = htmlquery.InnerText(node_title[0])

	// Description
	node_desc := htmlquery.Find(doc, `//*[@class="jj"]`)
	if len(node_desc) == 0 {
		err = fmt.Errorf("No matching desc")
		return
	}
	s.Description = strings.Replace(
		htmlquery.OutputHTML(node_desc[0], false),
		"<br/>", "\n",
		-1)

	// Author
	var author = htmlquery.Find(doc, `//*[@class="bq"]/span[2]/a`)
	s.Author = htmlquery.OutputHTML(author[0], false)

	// Contents
	node_content := htmlquery.Find(doc, `//*[@class="ml_list"]/ul/li/a`)
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

	s.CoverURL = htmlquery.SelectAttr(htmlquery.FindOne(doc, `//*[@class="pic"]/img`), "src")

	return
}
