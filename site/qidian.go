package site

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/ma6254/FictionDown/store"
)

func SingleSpace(s string) (r string) {
	rex := regexp.MustCompile("[\u0020\u3000]")
	return rex.ReplaceAllString(s, " ")
}

var qidian = SiteA{
	Name:     "起点中文网",
	HomePage: "https://www.qidian.com/",
	Match: []string{
		`https://book\.qidian\.com/info/\d+/*(#\w+)?`,
		`https://read\.qidian\.com/chapter/[\w_-]+/[\w_-]+/*`,
		`https://vipreader\.qidian\.com/chapter/\d+/\d+/*`,
	},
	BookInfo: func(body io.Reader) (s *store.Store, err error) {
		// body = transform.NewReader(body, simplifiedchinese.GBK.NewDecoder())
		doc, err := htmlquery.Parse(body)
		if err != nil {
			return
		}

		s = &store.Store{}

		s.BookName = htmlquery.InnerText(htmlquery.FindOne(doc, `//div[@class="book-info "]/h1/em`))

		s.CoverURL = "https:" + strings.TrimSpace(htmlquery.SelectAttr(htmlquery.FindOne(doc, `//div[@class="book-img"]/a/img`), "src"))
		// var author = htmlquery.Find(doc, `//*[@id="info"]/p[1]`)
		s.Author = strings.TrimSpace(htmlquery.InnerText(htmlquery.FindOne(doc, `//a[@class="writer"]`)))

		var desc = ""
		node_desc := htmlquery.Find(doc, `//div[@class="book-intro"]/p/text()`)
		if len(node_desc) == 0 {
			err = fmt.Errorf("No matching desc")
			return
		}
		for _, v := range node_desc {
			desc += strings.TrimSpace(htmlquery.InnerText(v)) + "\n"
		}
		s.Description = desc[:len(desc)-1]

		// ioutil.WriteFile("xxxx.html", []byte(htmlquery.OutputHTML(doc, false)), 0775)

		volumes := htmlquery.Find(doc, `//div[@class="volume"]`)
		for _, volume := range volumes {
			volumeName := htmlquery.InnerText(htmlquery.FindOne(volume, `/h3/text()[2]`))
			volumeName = strings.Replace(volumeName, "\n", "", -1)
			volumeName = strings.TrimSpace(volumeName)

			V := store.Volume{
				Name:     volumeName,
				Chapters: make([]store.Chapter, 0),
			}
			VIPClass := htmlquery.SelectAttr(htmlquery.FindOne(volume, `/h3/span`), "class")
			switch VIPClass {
			case "free":
				V.IsVIP = false
			case "vip":
				V.IsVIP = true
			}

			// fmt.Printf("卷: %#v %v\n", V.Name, V.IsVIP)
			nodeContent := htmlquery.Find(volume, `//div[@class="volume"]/ul/li/a`)
			for _, v := range nodeContent {
				c := store.Chapter{
					Name: strings.TrimSpace(
						SingleSpace(
							htmlquery.InnerText(v),
						),
					),
					URL: htmlquery.SelectAttr(v, "href"),
				}
				// fmt.Printf("%#v %v\n", c.Name, c.URL)
				V.Chapters = append(V.Chapters, c)
			}
			s.Volumes = append(s.Volumes, V)
		}
		return
	},
	Chapter: func(body io.Reader) ([]string, error) {
		// body := transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())
		doc, err := htmlquery.Parse(body)
		if err != nil {
			return nil, err
		}

		M := []string{}
		//list
		nodeContent := htmlquery.Find(doc, `//div[contains(@class,"read-content")]/p`)
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
	},
}
