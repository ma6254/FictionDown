package site

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/ma6254/FictionDown/store"
)

var regMap = map[string]Site{
	"www.biqiuge.com":      &Biquge3{},
	"www.booktxt.net":      &DingDian1{},
	"www.biquge5200.cc":    &Biquge_1{},
	"www.bqg5200.com":      &Biquge_2{},
	"www.81new.com":        &Www81newCom{},
	"book.qidian.com":      &QiDian{},
	"read.qidian.com":      &QiDian{},
	"vipreader.qidian.com": &QiDian{},
}

func MatchSite(m string) (*Site, bool) {
	for reg, site := range regMap {
		ok, err := path.Match(reg, m)
		if err != nil {
			return nil, false
		}
		return &site, ok
	}
	return nil, false
}

type ErrUnsupportSite struct {
	Site string
}

func (e ErrUnsupportSite) Error() string {
	return fmt.Sprintf("UnSupport Site: %#v", e.Site)
}

// Site 小说站点
type Site interface {
	BookInfo(body io.Reader) (s *store.Store, err error)
	Chapter(body io.Reader) (content []string, err error)
}

// BookInfo 获取小说信息
func BookInfo(BookURL string) (s *store.Store, err error) {
	u, err := url.Parse(BookURL)
	if err != nil {
		return nil, err
	}
	site, ok := regMap[u.Host]
	if !ok {
		return nil, ErrUnsupportSite{u.Host}
	}

	// Get WebPage
	client := &http.Client{}
	req, err := http.NewRequest("GET", BookURL, nil)
	if err != nil {
		return
	}
	req.Header.Add(
		"user-agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36",
	)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%d %v", resp.StatusCode, resp.Status)
	}

	chapter, err := site.BookInfo(resp.Body)
	chapter.BookURL = BookURL

	for v1, k1 := range chapter.Volumes {
		for v2, k2 := range k1.Chapters {
			u1, _ := url.Parse(k2.URL)
			chapter.Volumes[v1].Chapters[v2].URL = u.ResolveReference(u1).String()
			// if !u.IsAbs() {
			// 	u1.Scheme = resp.Request.URL.Scheme
			// 	u1.Host = resp.Request.URL.Host
			// 	chapter.Volumes[v1].Chapters[v2].URL = u.String()
			// }
		}
	}

	if len(chapter.Volumes) == 0 {
		return nil, fmt.Errorf("not match volumes")
	}

	return chapter, err
}

// Chapter 获取小说章节内容
func Chapter(BookURL string) (content []string, err error) {
	u, err := url.Parse(BookURL)
	if err != nil {
		return nil, err
	}
	site, ok := regMap[u.Host]
	if !ok {
		return nil, ErrUnsupportSite{u.Host}
	}

	// Get WebPage
	client := &http.Client{}
	req, err := http.NewRequest("GET", BookURL, nil)
	if err != nil {
		return
	}
	req.Header.Add(
		"user-agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36",
	)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", resp.Status)
	}
	return site.Chapter(resp.Body)
}
