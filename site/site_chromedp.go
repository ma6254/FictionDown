package site

import (
	"context"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/ma6254/FictionDown/store"
)

func ChromedpBookInfo(BookURL string, logfile string) (s *store.Store, err error) {

	u, err := url.Parse(BookURL)
	if err != nil {
		return nil, err
	}
	Site, ok := regMap[u.Host]
	if !ok {
		return nil, ErrUnsupportSite{u.Host}
	}

	var (
		// BookName string
		// Author   string
		html   string
		logOpt []chromedp.Option
	)

	tasks := chromedp.Tasks{
		chromedp.Navigate(BookURL),
		// chromedp.Text(`html`, &html, chromedp.ByQuery),
		chromedp.OuterHTML(`html`, &html, chromedp.ByQuery),
		// chromedp.WaitVisible(`html`, chromedp.ByQuery),
		// chromedp.ActionFunc(func(ctxt context.Context, c cdp.Executor) error {
		// 	html, err := dom.GetOuterHTML().WithNodeID(cdp.NodeID(0)).Do(ctxt, c)
		// 	return nil
		// }),
	}

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	if logfile == "" {
		logOpt = []chromedp.Option{}
	} else {
		clog := log.New(os.Stdout, "", log.LstdFlags)
		logOpt = []chromedp.Option{
			chromedp.WithLog(clog.Printf),
		}
	}

	// create chrome instance
	c, err := chromedp.New(ctxt, logOpt...)
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

	return Site.BookInfo(strings.NewReader(html))
}
