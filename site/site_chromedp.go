package site

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/chromedp"

	fcontext "github.com/ma6254/FictionDown/context"
	"github.com/ma6254/FictionDown/store"
	"github.com/ma6254/FictionDown/utils"
)

var (
	ChromedpCtx context.Context
)

func ChromedpInit() {
	opts := append(
		chromedp.DefaultExecAllocatorOptions,
		// []chromedp.ExecAllocatorOption{},
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
	)

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ChromedpCtx, _ = chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

}

func ChromedpBookInfo(BookURL string, logfile string) (s *store.Store, err error) {

	ms, err := MatchOne(Sitepool, BookURL)
	if err != nil {
		return nil, err
	}

	var (
		// BookName string
		// Author   string
		html string
		u    *url.URL
	)

	u, _ = url.Parse(BookURL)

	tasks := chromedp.Tasks{
		chromedp.Navigate(BookURL),
		chromedp.Sleep(2 * time.Second),
		// chromedp.Text(`html`, &html, chromedp.ByQuery),
		chromedp.OuterHTML(`html`, &html, chromedp.ByQuery),
		// chromedp.WaitVisible(`html`, chromedp.ByQuery),
	}

	if err = utils.Retry(5, time.Millisecond*500, func() error {
		// create chrome instance
		return chromedp.Run(ChromedpCtx, tasks...)
	}); err != nil {
		return nil, err
	}

	chapter, err := ms.BookInfo(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(chapter.BookName) == "" {
		err = fmt.Errorf("BookInfo Name is empty")
		return
	}

	for v1, k1 := range chapter.Volumes {
		for v2, k2 := range k1.Chapters {
			u1, _ := url.Parse(k2.URL)
			chapter.Volumes[v1].Chapters[v2].URL = u.ResolveReference(u1).String()
		}
	}

	if len(chapter.Volumes) == 0 {
		// fmt.Printf(content)
		return nil, fmt.Errorf("not match volumes")
	}

	return chapter, nil
}

func ChromedpChapter(BookURL string) (content []string, err error) {
	var (
		html string
	)
	ms, err := MatchOne(Sitepool, BookURL)
	if err != nil {
		return nil, err
	}

	// Get WebPage
	tasks := chromedp.Tasks{
		chromedp.Navigate(BookURL),
		chromedp.Sleep(2 * time.Second),
		// chromedp.WaitVisible(`html`, chromedp.ByQuery),
		chromedp.OuterHTML(`html`, &html, chromedp.ByQuery),
	}

	if err = utils.Retry(5, time.Millisecond*500, func() error {
		// create chrome instance
		return chromedp.Run(ChromedpCtx, tasks...)
	}); err != nil {
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
	ctx = context.WithValue(ctx, fcontext.KeyBody, strings.NewReader(html))

	return ms.Chapter(ctx)
}
