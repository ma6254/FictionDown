package site

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/ma6254/FictionDown/store"
)

func ChromedpBookInfo(BookURL string, logfile string) (s *store.Store, err error) {

	ms, err := MatchOne(Sitepool, BookURL)
	if err != nil {
		return nil, err
	}

	var (
		// BookName string
		// Author   string
		html string
		opts []chromedp.Option
		u    *url.URL
	)

	u, _ = url.Parse(BookURL)

	tasks := chromedp.Tasks{
		chromedp.Navigate(BookURL),
		// chromedp.Text(`html`, &html, chromedp.ByQuery),
		chromedp.OuterHTML(`html`, &html, chromedp.ByQuery),
		// chromedp.WaitVisible(`html`, chromedp.ByQuery),
	}

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	if logfile == "" {
		// opts = []chromedp.Option{
		// 	chromedp.WithLog(nil),
		// 	chromedp.WithErrorf(nil),
		// }
	} else {
		clog := log.New(os.Stdout, "", log.LstdFlags)
		opts = []chromedp.Option{
			chromedp.WithLog(clog.Printf),
			chromedp.WithErrorf(clog.Printf),
		}
	}

	opts = append(opts, chromedp.WithRunnerOptions(
	// runner.Flag("headless", true),
	// runner.Flag("disable-gpu", true),
	// runner.Flag("no-sandbox", true),
	))

	// create chrome instance
	c, err := chromedp.New(ctxt, opts...)
	if err != nil {
		log.Fatal(err)
	}

	err = c.Run(ctxt, tasks)
	if err != nil {
		log.Fatal(err)
	}

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}

	chapter, err := ms.BookInfo(strings.NewReader(html))
	if err != nil {
		return nil, err
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
