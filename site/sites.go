package site

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	fcontext "github.com/ma6254/FictionDown/context"

	"github.com/antchfx/htmlquery"
	"github.com/ma6254/FictionDown/store"
	"github.com/ma6254/FictionDown/utils"
	"golang.org/x/net/html"
	"golang.org/x/text/transform"
)

// Type1BookInfo 书籍信息页，单页，无翻页，无分卷
func Type1BookInfo(nameExpr, coverExpr, authorExpr, chapterExpr string) func(body io.Reader) (s *store.Store, err error) {
	return func(body io.Reader) (s *store.Store, err error) {
		doc, err := htmlquery.Parse(body)
		if err != nil {
			return
		}
		s = &store.Store{}
		var tmpNode *html.Node

		tmpNode = htmlquery.FindOne(doc, nameExpr)
		if tmpNode == nil {
			err = fmt.Errorf("No matching bookName")
			return
		}
		s.BookName = htmlquery.InnerText(tmpNode)

		if coverExpr == "" {
			// log.Printf("Empty Cover Image Expr")
		} else {
			coverNode := htmlquery.FindOne(doc, coverExpr)
			if coverNode == nil {
				err = fmt.Errorf("No matching author")
				return
			}
			if cu, err := url.Parse(strings.TrimSpace(htmlquery.InnerText(coverNode))); err != nil {
				log.Printf("Cover Image URL Error: %v", err)
			} else {
				s.CoverURL = cu.String()
			}
		}

		// Author
		authorContent := htmlquery.FindOne(doc, authorExpr)
		if authorContent == nil {
			err = fmt.Errorf("No matching author")
			return
		}
		s.Author = strings.TrimSpace(htmlquery.InnerText(authorContent))

		// Contents
		nodeContent := htmlquery.Find(doc, chapterExpr)
		if len(nodeContent) == 0 {
			err = fmt.Errorf("No matching contents")
			return
		}

		var vol = store.Volume{
			Name:     "正文",
			Chapters: make([]store.Chapter, 0),
		}
		for _, v := range nodeContent {
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
	}
}

// Type1Chapter 小说章节段落匹配
func Type1Chapter(expr string) func(ctx fcontext.Context) (content []string, err error) {
	return func(ctx fcontext.Context) (content []string, err error) {
		doc, err := htmlquery.Parse(ctx.Value(fcontext.KeyBody).(io.Reader))
		if err != nil {
			return nil, err
		}

		M := []string{}
		//list
		nodeContent := htmlquery.Find(doc, expr)
		if len(nodeContent) == 0 {
			err = fmt.Errorf("No matching content")
			return nil, err
		}
		for _, v := range nodeContent {
			t := htmlquery.InnerText(v)
			t = strings.TrimSpace(t)

			if t == "" {
				continue
			}

			M = append(M, t)
		}
		return M, nil
	}
}

// Type2Chapter 章节匹配2：单章分多页,
// next函数返回下一个页面的DOM
// block函数用于屏蔽多余的段落
func Type2Chapter(
	expr string,
	next func(preURL *url.URL, doc *html.Node) *html.Node,
	block func([]string) []string,
) func(fcontext.Context) (content []string, err error) {
	return func(ctx fcontext.Context) (content []string, err error) {
		doc, err := htmlquery.Parse(ctx.Value(fcontext.KeyBody).(io.Reader))
		if err != nil {
			return nil, err
		}
		M := []string{}
		if block == nil {
			block = func(a []string) []string { return a }
		}
		for {
			//list
			nodeContent := htmlquery.Find(doc, expr)
			if len(nodeContent) == 0 {
				err = fmt.Errorf("No matching content")
				return nil, err
			}
			MM := []string{}
		loopContent:
			for _, v := range nodeContent {
				t := htmlquery.InnerText(v)
				t = strings.TrimSpace(t)

				if t == "" {
					continue loopContent
				}
				MM = append(MM, t)
			}
			return MM, nil
			M = append(M, block(MM)...)

			if next == nil {
				return M, nil
			}
			doc = next(ctx.Value(fcontext.KeyURL).(*url.URL), doc)
			if doc == nil {
				break
			}
		}
		return M, nil
	}
}

type SearchFunc func(s string) (result []ChaperSearchResult, err error)

// Type1Search 搜索类型1: 搜索后得到302跳转或者列表的
func Type1Search(
	URL string,
	getReq func(s string) *http.Request,
	resultExpr, nameExpr, authorExpr string) func(s string) (result []ChaperSearchResult, err error) {
	return Type1SearchAfter(URL, getReq, resultExpr, nameExpr, authorExpr, nil)
}

// Type1SearchAfter 搜索类型1: 搜索后得到302跳转或者列表的
func Type1SearchAfter(
	URL string,
	getReq func(s string) *http.Request,
	resultExpr, nameExpr, authorExpr string,
	after func(r ChaperSearchResult) ChaperSearchResult) func(s string) (result []ChaperSearchResult, err error) {
	return func(s string) (result []ChaperSearchResult, err error) {
		req := getReq(s)
		var (
			resp *http.Response
		)
		if err = utils.Retry(5, time.Millisecond*500, func() error {
			resp, err = http.DefaultClient.Do(req)
			return err
		}); err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if req.URL.String() != resp.Request.URL.String() {
			// 单个搜索结果
			log.Printf("%s %s", req.URL.String(), resp.Request.URL.String())
			store, e := BookInfo(resp.Request.URL.String())
			if e != nil {
				return nil, e
			}
			r := ChaperSearchResult{
				BookName: store.BookName,
				Author:   store.Author,
				BookURL:  resp.Request.URL.String(),
			}
			if after != nil {
				r = after(r)
			}
			result = append(result, r)
			return result, nil
		}

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		var body io.Reader = bytes.NewReader(bodyBytes)
		encode := utils.DetectContentCharset(bytes.NewReader(bodyBytes))
		body = transform.NewReader(body, encode.NewDecoder())

		doc, err := htmlquery.Parse(body)
		if err != nil {
			return
		}

		r := htmlquery.Find(doc, resultExpr)
		if len(r) == 0 {
			return nil, nil
		}
		for _, v := range r {
			s2 := htmlquery.FindOne(v, nameExpr)
			if s2 == nil {
				return nil, fmt.Errorf("No matching result name")
			}
			s4 := htmlquery.FindOne(v, authorExpr)
			if s4 == nil {
				return nil, fmt.Errorf("No matching result author")
			}

			u1, _ := url.Parse(htmlquery.SelectAttr(s2, "href"))

			r := ChaperSearchResult{
				BookName: htmlquery.InnerText(s2),
				Author:   htmlquery.InnerText(s4),
				BookURL:  resp.Request.URL.ResolveReference(u1).String(),
			}
			if after != nil {
				r = after(r)
			}
			result = append(result, r)
		}
		return
	}
}
