package site

import (
	"context"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
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
		html string
		opts []chromedp.Option
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
		opts = []chromedp.Option{
			chromedp.WithLog(nil),
			chromedp.WithErrorf(nil),
		}
	} else {
		clog := log.New(os.Stdout, "", log.LstdFlags)
		opts = []chromedp.Option{
			chromedp.WithLog(clog.Printf),
			chromedp.WithErrorf(clog.Printf),
		}
	}

	opts = append(opts, chromedp.WithRunnerOptions(
		runner.Flag("headless", true),
		runner.Flag("disable-gpu", true),
		runner.Flag("no-sandbox", true),
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

	return Site.BookInfo(strings.NewReader(html))
}
