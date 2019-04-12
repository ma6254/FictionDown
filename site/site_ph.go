package site

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/benbjohnson/phantomjs"
	"github.com/ma6254/FictionDown/store"
)

func InitPhantomJS() {
	if err := phantomjs.DefaultProcess.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ClosePhantomJS() {
	phantomjs.DefaultProcess.Close()
}

func phGetPageBody(u string) (body string, err error) {
	page, err := phantomjs.CreateWebPage()
	if err != nil {
		return "", err
	}
	defer page.Close()

	Header := http.Header{}
	Header.Add("user-agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36",
	)
	if err := page.SetCustomHeaders(Header); err != nil {
		return "", err
	}

	if err := page.SetViewportSize(1024, 800); err != nil {
		return "", err
	}

	// Open a URL.
	if err := page.Open(u); err != nil {
		return "", err
	}

	content, err := page.Content()
	if err != nil {
		return
	}

	// f, err := ioutil.TempFile("./tmp", "*.png")
	// if err != nil {
	// 	return
	// }
	// log.Printf("Render: %s", f.Name())
	// if err = page.Render(f.Name(), "png", 100); err != nil {
	// 	return "", err
	// }

	return content, err
}

func PhBookInfo(BookURL string) (s *store.Store, err error) {

	content, err := phGetPageBody(BookURL)
	if err != nil {
		return
	}

	bu, err := url.Parse(BookURL)
	if err != nil {
		return
	}

	ms, err := MatchOne(Sitepool, BookURL)
	if err != nil {
		return nil, err
	}

	chapter, err := ms.BookInfo(strings.NewReader(content))
	chapter.BookURL = BookURL

	for v1, k1 := range chapter.Volumes {
		for v2, k2 := range k1.Chapters {
			u, _ := url.Parse(k2.URL)
			if !u.IsAbs() {
				u.Scheme = bu.Scheme
				u.Host = bu.Host
				chapter.Volumes[v1].Chapters[v2].URL = u.String()
			}
		}
	}

	if len(chapter.Volumes) == 0 {
		// fmt.Printf(content)
		return nil, fmt.Errorf("not match volumes")
	}

	return chapter, err
}

func PhChapter(BookURL string) (content []string, err error) {

	c, err := phGetPageBody(BookURL)
	if err != nil {
		return
	}

	ms, err := MatchOne(Sitepool, BookURL)
	if err != nil {
		return nil, err
	}

	return ms.Chapter(strings.NewReader(c))
}
