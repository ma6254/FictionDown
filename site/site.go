package site

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ma6254/FictionDown/store"
)

var regMap = map[string]Site{
	"www.biquge11.com":  &Biquge{},
	"www.biquge5200.cc": &Biquge{},
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
	Chapter(body io.Reader) (chaper *store.Chapter, err error)
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
	req, err := http.NewRequest("GET", u.Host, nil)
	if err != nil {
		return
	}
	req.Header.Add(
		"user-agent",
		"Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Mobile Safari/537.36",
	)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%d %v", resp.StatusCode, resp.Status)
	}

	return site.BookInfo(resp.Body)
}

// Chapter 获取小说章节内容
func Chapter(BookURL string) (chaper *store.Chapter, err error) {
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
	req, err := http.NewRequest("GET", u.Host, nil)
	if err != nil {
		return
	}
	req.Header.Add(
		"user-agent",
		"Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Mobile Safari/537.36",
	)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%d %v", resp.StatusCode, resp.Status)
	}
	return site.Chapter(resp.Body)
}
