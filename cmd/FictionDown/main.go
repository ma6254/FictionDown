package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ma6254/FictionDown/site"
	"github.com/ma6254/FictionDown/store"

	"github.com/go-yaml/yaml"
	"github.com/urfave/cli"
)

var (
	// Version git or release tag
	Version = ""
	// CommitID latest commit id
	CommitID = ""
	// BuildData build data
	BuildData = ""
)

var (
	filename = ""
	bookurl  = ""
	driver   = ""
	chapter  *store.Store
)

var app = &cli.App{
	Name:  "FictionDown",
	Usage: `https://github.com/ma6254/FictionDown`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "u,url",
			Usage:       "图书链接",
			Destination: &bookurl,
		},
		cli.StringSliceFlag{
			Name:  "tu,turl",
			Usage: "资源网站链接",
		},
		cli.StringFlag{
			Name:        "i,input",
			Usage:       "输入缓存文件",
			Destination: &filename,
		},
		cli.StringFlag{
			Name:  "log",
			Usage: "log file path",
		},
		cli.StringFlag{
			Name:        "driver,d",
			Usage:       "请求方式,support: none,phantomjs,chromedp",
			Destination: &driver,
		},
	},
	Commands: []cli.Command{
		download,
		check,
		edit,
		convert,
		pirate,
		search,
	},
}

func main() {

	app.Version = Version
	app.Authors = []cli.Author{
		cli.Author{Name: "ma6254", Email: "9a6c5609806a@gmail.com"},
	}

	if (BuildData != "") && (CommitID != "") {
		app.Description = fmt.Sprintf("BuildData: %s\n   CommitID: %s ", BuildData, CommitID)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

type SyncStore struct {
	lock    *sync.Mutex
	IsTWork bool
	jobs    [][]bool
	Store   *store.Store
}

func (s *SyncStore) Init() {
	s.jobs = make([][]bool, len(s.Store.Volumes))
	s.lock = &sync.Mutex{}
	for k := range s.jobs {
		s.jobs[k] = make([]bool, len(s.Store.Volumes[k].Chapters))
	}
}

func (s *SyncStore) GetJob() (vi, ci int, url string, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for vi, vol := range s.Store.Volumes {
		for ci, ch := range vol.Chapters {
			if !s.jobs[vi][ci] {
				if (len(ch.Text) == 0) && (len(ch.Example) == 0) {
					s.jobs[vi][ci] = true
					// log.Printf("GetJob，%s-%s %#v", s.Store.Volumes[vi].Name, s.Store.Volumes[vi].Chapters[ci].Name, s.Store.Volumes[vi].Chapters[ci].URL)
					return vi, ci, ch.URL, nil
				}
			}
		}
	}
	return 0, 0, "", io.EOF
}

func (s *SyncStore) GetTJob() (vi, ci int, url []string, rawurl string, example []string, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for vi, vol := range s.Store.Volumes {
		for ci, ch := range vol.Chapters {
			if !s.jobs[vi][ci] {
				if (len(ch.Text) == 0) && (len(ch.Example) != 0) {
					s.jobs[vi][ci] = true
					return vi, ci, ch.TURL, ch.URL, ch.Example, nil
				}
			}
		}
	}
	return 0, 0, nil, "", nil, io.EOF
}

func (s *SyncStore) SaveJob(vi, ci int, text []string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.Store.Volumes[vi].IsVIP && !s.IsTWork {
		s.Store.Volumes[vi].Chapters[ci].Example = text
	} else {
		s.Store.Volumes[vi].Chapters[ci].Text = text
	}

	// log.Printf("SaveJob 2")
	bbb, err := yaml.Marshal(*(s.Store))
	if err != nil {
		panic(err)
	}
	// log.Printf("SaveJob 3")
	var filename = fmt.Sprintf("%s.%s", s.Store.BookName, store.FileExt)
	err = ioutil.WriteFile(filename, bbb, 0775)
	if err != nil {
		panic(err)
	}
	// log.Printf("SaveJob，%s-%s", s.Store.Volumes[vi].Name, s.Store.Volumes[vi].Chapters[ci].Name)
	// log.Printf("SaveJob End")
}

func Job(syncStore *SyncStore, jobch chan error) {
	defer func(jobch chan error) {
		jobch <- io.EOF
	}(jobch)
	// defer log.Printf("End Job")

	for {
		vi, ci, BookURL, err := syncStore.GetJob()
		if err != nil {
			if err != io.EOF {
				jobch <- err
			}
			return
		}

	A:
		for {
			content, err := site.Chapter(BookURL)
			if err != nil {
				log.Printf("Error: %s", err)
				time.Sleep(errSleep)
				continue A
			}
			syncStore.SaveJob(vi, ci, content)
			jobch <- nil
			time.Sleep(tSleep)
			break A
		}
	}
}

func TJob(syncStore *SyncStore, jobch chan error) {
	defer func(jobch chan error) {
		// log.Printf("Fuck Exit")
		jobch <- io.EOF
	}(jobch)

	for {
		var (
			deiver      = 0
			errCount    = 0
			MaxErrCount = 5
		)
		vi, ci, BookURL, RawURL, Example, err := syncStore.GetTJob()
		if err != nil {
			if err != io.EOF {
				jobch <- err
			}
			return
		}
		var P = 0
		var RP = make([]bool, len(BookURL))
	A:
		for {
			if len(BookURL) == 0 {
				log.Printf("无源章节 卷: %#v 章节: %#v", syncStore.Store.Volumes[vi].Name, syncStore.Store.Volumes[vi].Chapters[ci].Name)
				break A
			}

			var (
				content []string
				err     error
			)

			switch deiver {
			case 0:
				content, err = site.Chapter(BookURL[P])
			case 1:
				content, err = site.PhChapter(BookURL[P])
			default:
				jobch <- fmt.Errorf("爬取方式错误: %d", deiver)
				break A
			}
			if err != nil {
				errCount++
				log.Printf("Error: %s %s %s",
					syncStore.Store.Volumes[vi].Chapters[ci].Name,
					BookURL[P],
					err,
				)
				if errCount < MaxErrCount {
					time.Sleep(errSleep)
					continue A
				} else {
					deiver++
					log.Printf("错误次数过多，忽略此章节，并尝试更换爬取方式")
					continue A
				}
			}

			//开始对比
			sss := ""
			aaa := ""
			var ok float32
			var fail float32

			for _, v := range content {
				sss += v
			}

			for _, v := range Example {
				aaa += v
			}
			var ee = strings.Split(aaa, "。")
			// ee = SplitXX(ee, "，", "：", "“", "”", "？", "…")
			ee = SplitX(ee, "，")
			ee = SplitX(ee, "：")
			ee = SplitX(ee, "“")
			ee = SplitX(ee, "”")
			ee = SplitX(ee, "？")
			ee = SplitX(ee, "…")
			ee = SplitX(ee, "[")
			ee = SplitX(ee, "]")
			ee = SplitX(ee, "!")
			ee = SplitX(ee, "！")

			for _, v := range ee {
				if strings.Contains(sss, v) {
					ok++
				} else if strings.Contains(v, sss) {
					ok++
				} else {
					fail++
				}
			}

			if (ok / (ok + fail)) < 0.4 {
				RP[P] = true

				isDie := true
			IsDie:
				for _, v := range RP {
					if !v {
						isDie = false
						break IsDie
					}
				}

				if isDie {
					// log.Printf("ok/fail %d/%d", ok, fail)
					log.Printf("全部校验失败 %f Raw: %s", ok/(ok+fail), RawURL)
					log.Printf("BookURL: %#v", BookURL)
					// log.Printf("EEE: %#v", ee)
					// log.Printf("SSS: %s", sss)

					break A
				}

				P++
				log.Printf("校验失败 %f 切换源 %d %s %s %s", ok/(ok+fail), P, RawURL, BookURL[P-1], BookURL[P])
				// log.Printf("EEE: %#v", ee)
				// log.Printf(sss)
				// log.Fatal("Fuck")

				continue A
			}

			syncStore.SaveJob(vi, ci, content)
			jobch <- nil
			time.Sleep(tSleep)
			break A
		}
	}
}

func SplitX(s []string, q string) []string {
	e := []string{}
	for _, v := range s {
		if "" == strings.TrimSpace(v) {
			continue
		}

		eee := []string{}
		ee := strings.Split(v, q)
		for _, vv := range ee {
			if "" == strings.TrimSpace(vv) {
				continue
			}
			eee = append(eee, vv)
		}
		e = append(e, eee...)
	}
	return e
}

func SplitXX(s []string, q ...string) []string {
	e := make([]string, len(s))
	copy(s, e)
	for _, v := range q {
		e = SplitX(e, v)
	}
	return e
}

func initLoadStore(c *cli.Context) error {
	if filename == "" {
		if bookurl == "" {
			return fmt.Errorf("Must Input URL")
		}
		bookURL, err := url.Parse(bookurl)
		if err != nil {
			return err
		}
		log.Printf("URL: %#v", bookURL.String())
		switch c.GlobalString("driver") {
		case "phantomjs":
			log.Printf("Init PhantomJS")
			site.InitPhantomJS()
			defer func() {
				log.Printf("Close PhantomJS")
				site.ClosePhantomJS()
			}()
			for errCount := 0; errCount < 20; errCount++ {
				chapter, err = site.PhBookInfo(bookURL.String())
				if err == nil {
					break
				} else {
					log.Printf("ErrCount: %d Err: %s", errCount, err)
					time.Sleep(1000 * time.Millisecond)
				}
			}
		case "chromedp":
			log.Printf("Chromedp Running...")
			chapter, err = site.ChromedpBookInfo(bookURL.String(), c.String("chromedp-log"))
		default:
			log.Printf("use golang default http")
			for errCount := 0; errCount < 20; errCount++ {
				chapter, err = site.BookInfo(bookURL.String())
				if err == nil {
					break
				} else {
					log.Printf("ErrCount: %d Err: %s", errCount, err)
					time.Sleep(1000 * time.Millisecond)
				}
			}
		}
		if err != nil {
			return err
		}
		chapter.BookURL = bookURL.String()
		filename = fmt.Sprintf("%s.%s", chapter.BookName, store.FileExt)
		filemode, err := os.Stat(filename)
		if err != nil && os.IsNotExist(err) {
			b, err := yaml.Marshal(chapter)
			if err != nil {
				return err
			}
			ioutil.WriteFile(filename, b, 0775)
		} else {
			if filemode.IsDir() {
				return fmt.Errorf("is Dir")
			}
			log.Printf("Loading....")
			b, err := ioutil.ReadFile(filename)
			if err != nil {
				return err
			}
			err = yaml.Unmarshal(b, &chapter)
			if err != nil {
				return err
			}
		}
	} else {
		log.Printf("Loading cache file: %s", filename)
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		err = yaml.Unmarshal(b, &chapter)
		if err != nil {
			return err
		}
	}

	if len(c.GlobalStringSlice("turl")) != 0 {
		chapter.Tmap = []string{}
		for _, v := range c.GlobalStringSlice("turl") {
			chapter.Tmap = append(chapter.Tmap, v)
		}
	}
	return nil
}
