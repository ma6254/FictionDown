package biquge

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/ma6254/FictionDown/store"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func BookInfo(bookURL string) (s *store.Store, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", bookURL, nil)
	if err != nil {
		return
	}
	req.Header.Add("user-agent",
		"Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Mobile Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("%v", resp.Status)
		return
	}

	var body = transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())

	doc, err := htmlquery.Parse(body)
	if err != nil {
		return
	}

	s = &store.Store{}

	node_title := htmlquery.Find(doc, `//*[@id="info"]/h1`)
	if len(node_title) == 0 {
		err = fmt.Errorf("No matching title")
		return
	}
	s.BookName = htmlquery.InnerText(node_title[0])

	node_desc := htmlquery.Find(doc, `//*[@id="intro"]/p`)
	if len(node_desc) == 0 {
		err = fmt.Errorf("No matching desc")
		return
	}
	s.Description = strings.Replace(
		htmlquery.OutputHTML(node_desc[0], false),
		"<br/>", "\n",
		-1)

	var author = htmlquery.Find(doc, `//*[@id="info"]/p[1]`)
	s.Author = strings.TrimLeft(htmlquery.OutputHTML(author[0], false), "作\u00a0\u00a0\u00a0\u00a0者：")

	s.BiqugeURL = bookURL

	node_content := htmlquery.Find(doc, `//*[@id="list"]/dl/dd/a`)
	if len(node_desc) == 0 {
		err = fmt.Errorf("No matching contents")
		return
	}

	var vol = store.Volume{
		Name:     "正文",
		Chapters: make([]store.Chapter, 0),
	}

	for _, v := range node_content[9:] {
		//fmt.Printf("href: %v\n", chapter_u)
		vol.Chapters = append(vol.Chapters, store.Chapter{
			Name: htmlquery.InnerText(v),
			URL:  htmlquery.SelectAttr(v, "href"),
		})
	}
	s.Volumes = append(s.Volumes, vol)

	s.CoverURL = htmlquery.SelectAttr(htmlquery.FindOne(doc, `//*[@id="fmimg"]/img`), "src")

	return
}

func Chapter(ChapterURL string) (title string, text []string, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", ChapterURL, nil)
	if err != nil {
		return
	}
	req.Header.Add("user-agent",
		"Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Mobile Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("%v", resp.Status)
		return
	}

	var body = transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())

	doc, err := htmlquery.Parse(body)
	if err != nil {
		return
	}
	M := []string{}
	//list
	node_content := htmlquery.Find(doc, `//*[@id="content"]/p`)
	if len(node_content) == 0 {
		err = fmt.Errorf("No matching content")
		return
	}
	for _, v := range node_content {
		t := htmlquery.InnerText(v)
		t = strings.TrimSpace(t)
		M = append(M, t)
	}

	return "", M, err
}
