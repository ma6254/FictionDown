package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/ma6254/FictionDown/output"
	"github.com/ma6254/FictionDown/store"

	"github.com/go-yaml/yaml"
	"github.com/ma6254/FictionDown/biquge"
	"github.com/urfave/cli"
	processbar "gopkg.in/cheggaaa/pb.v1"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "url",
			Usage: "笔趣阁链接",
		},
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name: "check",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "f",
				},
			},
			Action: func(c *cli.Context) error {
				var (
					Chapter *store.Store
				)
				var filename = c.String("f")
				b, err := ioutil.ReadFile(filename)
				if err != nil {
					return err
				}
				err = yaml.Unmarshal(b, &Chapter)
				if err != nil {
					return err
				}
				for _, vol := range Chapter.Volumes {
					for ci, ch := range vol.Chapters {
						var l string
						if len(ch.Text) == 0 {
							l = "未缓存"
						} else {
							l = fmt.Sprintf("%d", len(ch.Text))
						}
						fmt.Printf("%s %d %#v %s\n", vol.Name, ci, ch.Name, l)
					}
				}
				return nil
			},
		},
		cli.Command{
			Name:    "download",
			Aliases: []string{"d", "down"},
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "t",
					Usage: "线程数",
					Value: 10,
				},
				cli.StringFlag{
					Name:  "f",
					Usage: "输出格式",
				},
				cli.StringFlag{
					Name:  "o",
					Usage: "输出路径",
				},
			},
			Action: func(c *cli.Context) error {

				var (
					err       error
					URLString string
					bookURL   *url.URL
					Chapter   *store.Store
				)

				URLString = c.GlobalString("url")
				if c.GlobalString("url") == "" {
					log.Printf("Must Input URL")
					fmt.Scanln(&URLString)
				}
				bookURL, err = url.Parse(URLString)
				if err != nil {
					return err
				}
				log.Printf("URL: %#v", bookURL.String())
				Chapter, err = biquge.BookInfo(bookURL.String())
				if err != nil {
					return err
				}
				var filename = fmt.Sprintf("%s.%s", Chapter.BookName, store.FileExt)
				filemode, err := os.Stat(filename)
				if err != nil && os.IsNotExist(err) {
					b, err := yaml.Marshal(Chapter)
					if err != nil {
						return err
					}
					ioutil.WriteFile(filename, b, 0775)
				} else {
					if filemode.IsDir() {
						return fmt.Errorf("is Dir")
					} else {
						log.Printf("Loading....")
						b, err := ioutil.ReadFile(filename)
						if err != nil {
							return err
						}
						err = yaml.Unmarshal(b, &Chapter)
						if err != nil {
							return err
						}
					}
				}
				fmt.Printf("书名: %#v\n", Chapter.BookName)
				fmt.Printf("作者: %#v\n", Chapter.Author)
				fmt.Printf("封面: %s\n", Chapter.CoverURL)
				fmt.Printf("简介: \n\t%v\n", strings.Replace(Chapter.Description, "\n", "\n\t", -1))
				fmt.Printf("章节数:\n")
				for _, v := range Chapter.Volumes {
					fmt.Printf("\t%s %d章\n", v.Name, len(v.Chapters))
				}
				log.Printf("Working...\n")
				log.Printf("routine: %d", c.Int("t"))
				ssss := &SyncStore{
					Store: Chapter,
				}
				ssss.Init()

				var chCount = 0
				var isDone = 0
				for _, v := range Chapter.Volumes {
					chCount += len(v.Chapters)
					for _, v2 := range v.Chapters {
						if len(v2.Text) != 0 {
							isDone++
						}
					}
				}
				if isDone != 0 {
					log.Printf("读入缓存: %d", isDone)
				}
				bar := processbar.StartNew(chCount)
				bar.Set(isDone)

				Jobch := make(chan error)
				for i := 0; i < c.Int("t"); i++ {
					go Job(ssss, Jobch)
				}
				cc := make(chan os.Signal)
				signal.Notify(cc, os.Interrupt)

				var ii = 0
			AA:
				for {
					select {
					case ccc := <-cc:
						log.Printf("进程信号: %v", ccc)
						return nil
					case err := <-Jobch:
						if err != nil {
							if err == io.EOF {
								ii++
								if ii >= c.Int("t") {
									bar.Finish()
									log.Printf("缓存完成")
									break AA
								}
							} else {
								log.Printf("Job Error: %s", err)
							}
						} else {
							bar.Increment()
						}
					}
				}

				if c.String("f") == "" {
					return nil
				}

				var ConversionFileName string
				if c.String("o") == "" {
					ConversionFileName = fmt.Sprintf("%s.%s", Chapter.BookName, c.String("f"))
				}
				log.Printf("Start Conversion: Format:%#v OutPath:%#v", c.String("f"), ConversionFileName)
				return output.Output(*Chapter, c.String("f"), ConversionFileName)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

type SyncStore struct {
	lock  *sync.Mutex
	jobs  [][]bool
	Store *store.Store
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
				if len(ch.Text) == 0 {
					s.jobs[vi][ci] = true
					// log.Printf("GetJob，%s-%s %#v", s.Store.Volumes[vi].Name, s.Store.Volumes[vi].Chapters[ci].Name, s.Store.Volumes[vi].Chapters[ci].URL)
					return vi, ci, ch.URL, nil
				}
			}
		}
	}
	return 0, 0, "", io.EOF
}

func (s *SyncStore) SaveJob(vi, ci int, text []string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.Store.Volumes[vi].Chapters[ci].Text = text

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
	// defer log.Printf("End Job")

	for {
		vi, ci, url, err := syncStore.GetJob()
		if err != nil {
			jobch <- err
		}

	A:
		for {
			_, texts, err := biquge.Chapter(url)
			if err != nil {
				// log.Printf("Error: %s", err)
				time.Sleep(500 * time.Millisecond)
				continue A
			}
			syncStore.SaveJob(vi, ci, texts)
			jobch <- nil
			time.Sleep(100 * time.Millisecond)
			break A
		}
	}
}
