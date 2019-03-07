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

// Www81newCom 新八一中文网
type Www81newCom struct {
}

func (b *Www81newCom) BookInfo(body io.Reader) (s *store.Store, err error) {
	body = transform.NewReader(body, simplifiedchinese.GBK.NewDecoder())
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
	if len(node_desc) == 0 {
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

func (b *Www81newCom) Chapter(body io.Reader) ([]string, error) {
	body = transform.NewReader(body, simplifiedchinese.GBK.NewDecoder())
	doc, err := htmlquery.Parse(body)
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
}
