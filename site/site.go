package site

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	fcontext "github.com/ma6254/FictionDown/context"

	"github.com/ma6254/FictionDown/store"
	"github.com/ma6254/FictionDown/utils"
	"golang.org/x/text/transform"
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
	// biquge2,
}

func AddSite(site SiteA) {
	if site.File == "" {
		_, filename, _, _ := runtime.Caller(1)
		site.File = filename
	}
	Sitepool = append(Sitepool, site)
}

type SiteA struct {
	Name     string // 站点名称
	HomePage string // 站点首页

	File string

	// match url, look that https://godoc.org/path#Match
	Match []string

	// search book on site
	Search func(s string) (result []ChaperSearchResult, err error) `json:"-"`

	// parse fiction info by page body
	BookInfo func(body io.Reader) (s *store.Store, err error) `json:"-"`

	// parse fiction chaper content by page body
	Chapter func(fcontext.Context) (content []string, err error) `json:"-"`

	// get site tags
	Tags func() []string `json:"-"`
}

// MatchOne match one site, is use `MatchSites` first result
func MatchOne(pool []SiteA, u string) (*SiteA, error) {
	a, err := MatchSites(pool, u)
	if err != nil {
		return nil, err
	}
	if len(a) < 1 {
		return nil, ErrUnsupportSite{u}
	}
	return a[0], nil
}

// MatchSites match all site
func MatchSites(pool []SiteA, u string) ([]*SiteA, error) {
	var result = []*SiteA{}
	for k := range pool {
		ok, err := pool[k].match(u)
		if err != nil {
			return nil, err
		}
		if ok {
			result = append(result, &pool[k])
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
	resp, err := utils.RequestGet(BookURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if ms.BookInfo == nil {
		return nil, ErrMethodMissing{ms}
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var body io.Reader = bytes.NewReader(bodyBytes)

	if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		encode := utils.DetectContentCharset(bytes.NewReader(bodyBytes))
		body = transform.NewReader(body, encode.NewDecoder())
	}

	chapter, err := ms.BookInfo(body)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(chapter.BookName) == "" {
		err = fmt.Errorf("BookInfo Name is empty")
		return
	}

	chapter.Author = strings.Replace(chapter.Author, "\u00a0", "", -1)

	chapter.BookURL = BookURL

	for v1, k1 := range chapter.Volumes {
		for v2, k2 := range k1.Chapters {
			u1, err := resp.Request.URL.Parse(k2.URL)
			if err != nil {
				return nil, err
			}
			chapter.Volumes[v1].Chapters[v2].URL = u1.String()
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
	body, err := utils.GetWebPageBodyReader(BookURL)
	if err != nil {
		return nil, err
	}

	if ms.Chapter == nil {
		return nil, fmt.Errorf("Site %s Chapter Func is empty", ms.Name)
	}

	bu, err := url.Parse(BookURL)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, fcontext.KeyURL, bu)
	ctx = context.WithValue(ctx, fcontext.KeyBody, body)

	return ms.Chapter(ctx)
}

func Search(s string) (result []ChaperSearchResult, err error) {
	var (
		wg         sync.WaitGroup
		resultLock = sync.Mutex{}
	)

	searchFunc := func(fn func([]ChaperSearchResult), is SiteA) {
		defer wg.Done()
		var (
			err error
			rr  []ChaperSearchResult
		)

		if err := utils.Retry(3, 1*time.Second, func() error {
			rr, err = is.Search(s)
			if err != nil {
				log.Printf("Error: 搜索站点: %s %s %s", is.Name, is.HomePage, err)
				return err
			}
			return nil
		}); err != nil {
			return
		}
		log.Printf("搜索站点: 结果: %d %s %s", len(rr), is.Name, is.HomePage)
		resultLock.Lock()
		defer resultLock.Unlock()
		fn(rr)
	}

	for _, v := range Sitepool {
		if v.Search == nil {
			continue
		}
		log.Printf("开始搜索站点: %s %s", v.Name, v.HomePage)
		wg.Add(1)
		go searchFunc(func(r []ChaperSearchResult) {
			result = append(result, r...)
		}, v)
	}
	wg.Wait()
	return
}
