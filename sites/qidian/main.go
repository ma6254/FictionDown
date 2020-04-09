package qidian

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"

	fcontext "github.com/ma6254/FictionDown/context"

	"github.com/antchfx/htmlquery"
	"github.com/buger/jsonparser"
	"github.com/ma6254/FictionDown/site"
	"github.com/ma6254/FictionDown/store"
	"github.com/ma6254/FictionDown/utils"
	"gopkg.in/yaml.v2"
)

func SingleSpace(s string) (r string) {
	rex := regexp.MustCompile("[\u0020\u3000]")
	return rex.ReplaceAllString(s, " ")
}

func Site() site.SiteA {
	return site.SiteA{
		Name:     "起点中文网",
		HomePage: "https://www.qidian.com/",
		Tags:     site.AddTag(nil, "正版", "阅文集团"),
		Match: []string{
			`https://book\.qidian\.com/info/\d+/*(#\w+)?`,
			`https://read\.qidian\.com/chapter/[\w_-]+/[\w_-]+/*`,
			`https://vipreader\.qidian\.com/chapter/\d+/\d+/*`,
		},
		BookInfo: func(body io.Reader) (s *store.Store, err error) {

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

			if len(volumes) == 0 {
				bidNode := htmlquery.FindOne(doc, `//*[@id="bookImg"]/@data-bid`)
				if bidNode == nil {
					return nil, fmt.Errorf("not match bid")
				}
				s.Volumes, err = qidianGetChapter(htmlquery.InnerText(bidNode))
				if err != nil {
					return nil, err
				}
			}
			return
		},
		Chapter: func(ctx fcontext.Context) (content []string, err error) {
			doc, err := htmlquery.Parse(ctx.Value(fcontext.KeyBody).(io.Reader))
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
		Search: site.Type1Search("",
			func(s string) *http.Request {
				baseurl, err := url.Parse("https://www.qidian.com/search")
				if err != nil {
					panic(err)
				}
				value := baseurl.Query()
				value.Add("kw", s)
				baseurl.RawQuery = value.Encode()

				req, err := http.NewRequest("GET", baseurl.String(), nil)
				if err != nil {
					panic(err)
				}
				return req
			},
			`//div[(@class="book-mid-info") and (./h4/a/cite/@class="red-kw") ]`,
			`h4/a`,
			`p/a[@class="name"]/text()`),
	}
}

type qidianAPIStatus struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
type qidianAPIChapterData struct {
	Data struct {
		Volumes []qidianAPIVolume `json:"vs"`
	} `json:"data"`
}

type qidianAPIVolume struct {
	Name     string             `json:"vN"`
	IsVIP    bool               `json:"-"`
	Chapters []qidianAPIChapter `json:"cs"`
}

type qidianAPIChapter struct {
	Name  string      `json:"cN"`
	URL   string      `json:"cU"`
	ID    json.Number `json:"id,number"`
	IsVIP bool        `json:"-"`
}

func (q qidianAPIStatus) Error() string {
	return fmt.Sprintf("%d %s", q.Code, q.Msg)
}

func qidianGetChapter(bid string) ([]store.Volume, error) {
	var (
		err                error
		vols               = []store.Volume{}
		baseURLs           = `https://book.qidian.com/ajax/book/category`
		vipChapterBaseURLs = `https://vipreader.qidian.com/chapter/`
		resultState        = qidianAPIStatus{}
		resultData         = qidianAPIChapterData{}

		baseURL           *url.URL
		vipChapterBaseURL *url.URL
	)

	baseURL, err = url.Parse(baseURLs)
	if err != nil {
		panic(err)
	}
	vipChapterBaseURL, err = url.Parse(vipChapterBaseURLs)
	if err != nil {
		panic(err)
	}

	value := baseURL.Query()
	value.Add("bookId", bid)
	baseURL.RawQuery = value.Encode()
	resp, err := utils.RequestGet(baseURL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(body, &resultState); err != nil {
		return nil, err
	}
	if resultState.Code != 0 {
		return nil, resultState
	}
	if err = json.Unmarshal(body, &resultData); err != nil {
		return nil, err
	}

	for vk, vol := range resultData.Data.Volumes {
		isVips, err := jsonparser.GetInt(body, "data", "vs", fmt.Sprintf("[%d]", vk), "vS")
		if err != nil {
			return nil, err
		}
		isVip := true
		if isVips == 0 {
			isVip = false
		}
		if isVip {
			resultData.Data.Volumes[vk].IsVIP = true
		}

		for ck, chapter := range vol.Chapters {

			isFrees, err := jsonparser.GetInt(body, "data", "vs", fmt.Sprintf("[%d]", vk), "cs", fmt.Sprintf("[%d]", ck), "sS")
			if err != nil {
				return nil, err
			}
			if isFrees == 0 {
				resultData.Data.Volumes[vk].Chapters[ck].IsVIP = true
			}

			vipChapterBaseURL, err = url.Parse(vipChapterBaseURLs)
			if err != nil {
				panic(err)
			}
			vipChapterBaseURL.Path = path.Join(vipChapterBaseURL.Path, bid, chapter.ID.String())
			resultData.Data.Volumes[vk].Chapters[ck].URL = vipChapterBaseURL.String()
		}
	}
	volsYaml, err := yaml.Marshal(resultData.Data.Volumes)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(volsYaml, &vols); err != nil {
		return nil, err
	}

	return vols, nil
}
