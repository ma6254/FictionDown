package site

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"github.com/ma6254/FictionDown/store"
)

type ErrUnsupportSite struct {
	Site string
}

func (e ErrUnsupportSite) Error() string {
	return fmt.Sprintf("UnSupport Site: %#v", e.Site)
}

type ErrMethodMissing struct {
	Site *SiteA
}

func (e ErrMethodMissing) Error() string {
	return fmt.Sprintf("Method Missing: %s %#v", e.Site.Name, e.Site.HomePage)
}

var Sitepool = []SiteA{
	qidian,
	wwww81newcom,
	dingdian,
	biquge1,
	biquge2,
	biquge3,
}

type SiteA struct {
	Name     string // 站点名称
	HomePage string // 站点首页

	// match url, look that https://godoc.org/path#Match
	Match []string

	// search book on site
	Search func(s string) (result []ChaperSearchResult, err error)

	// parse fiction info by page body
	BookInfo func(body io.Reader) (s *store.Store, err error)

	// parse fiction chaper content by page body
	Chapter func(body io.Reader) (content []string, err error)
}

// MatchOne match one site, is use `MatchSites` first result
func MatchOne(pool []SiteA, u string) (*SiteA, error) {
	a, err := MatchSites(pool, u)
	if err != nil {
		return nil, err
	}
	if len(a) == 0 {
		return nil, ErrUnsupportSite{u}
	}
	return &a[0], nil
}

// MatchSites match all site
func MatchSites(pool []SiteA, u string) ([]SiteA, error) {
	var result = []SiteA{}
	for _, v := range pool {
		ok, err := v.match(u)
		if err != nil {
			return nil, err
		}
		if ok {
			result = append(result, v)
		}
	}
	return result, nil
}

func (s SiteA) match(u string) (bool, error) {
	for _, v := range s.Match {
		re, err := regexp.Compile(v)
		if err != nil {
			return false, err
		}
		if re.MatchString(u) {
			return true, nil
		}
	}
	return false, nil
}

type ChaperSearchResult struct {
	BookName string
	Author   string
	BookURL  string
}

// BookInfo 获取小说信息
func BookInfo(BookURL string) (s *store.Store, err error) {
	ms, err := MatchOne(Sitepool, BookURL)
	if err != nil {
		return nil, err
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

	if ms.BookInfo == nil {
		return nil, ErrMethodMissing{ms}
	}

	chapter, err := ms.BookInfo(resp.Body)
	chapter.BookURL = BookURL

	for v1, k1 := range chapter.Volumes {
		for v2, k2 := range k1.Chapters {
			u1, _ := url.Parse(k2.URL)
			chapter.Volumes[v1].Chapters[v2].URL = resp.Request.URL.ResolveReference(u1).String()
			// if !u.IsAbs() {
			// 	u1.Scheme = resp.Request.URL.Scheme
			// 	u1.Host = resp.Request.URL.Host
			// 	chapter.Volumes[v1].Chapters[v2].URL = u.String()
			// }
		}
	}

	CoverURL, err := url.Parse(chapter.CoverURL)
	if err != nil {
		return nil, err
	}

	if chapter.CoverURL != "" {
		chapter.CoverURL = resp.Request.URL.ResolveReference(CoverURL).String()
	}

	if len(chapter.Volumes) == 0 {
		return nil, fmt.Errorf("not match volumes")
	}
	return chapter, err
}

// Chapter 获取小说章节内容
func Chapter(BookURL string) (content []string, err error) {
	ms, err := MatchOne(Sitepool, BookURL)
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("%#v %s", BookURL, resp.Status)
	}
	return ms.Chapter(resp.Body)
}

func Search(s string) (result []ChaperSearchResult, err error) {
	for _, v := range Sitepool {
		if v.Search == nil {
			continue
		}
		r, err := v.Search(s)
		if err != nil {
			return nil, err
		}
		result = append(result, r...)
	}
	return
}
