package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/ma6254/FictionDown/output"
	"github.com/ma6254/FictionDown/site"
	"github.com/ma6254/FictionDown/store"

	"github.com/go-yaml/yaml"
	"github.com/urfave/cli"
	processbar "gopkg.in/cheggaaa/pb.v1"
)

var (

	// Version git or release tag
	Version = ""
	// CommitID latest commit id
	CommitID = ""
	// BuildData build data
	BuildData = ""

	tSleep   time.Duration
	errSleep time.Duration
)

func main() {

	app := cli.NewApp()

	app.Usage = `https://github.com/ma6254/FictionDown`
	app.Version = Version

	app.Description = fmt.Sprintf("BuildData: %s\n   CommitID: %s ", BuildData, CommitID)

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "url",
			Usage: "笔趣阁链接",
		},
		cli.StringSliceFlag{
			Name:  "turl",
			Usage: "盗版网站链接",
		},
		cli.StringFlag{
			Name:  "i",
			Usage: "输入缓存文件",
		},
		cli.StringFlag{
			Name:  "log",
			Usage: "log file path",
		},
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:    "check",
			Aliases: []string{"c"},
			Flags:   []cli.Flag{},
			Action: func(c *cli.Context) error {
				var (
					Chapter  *store.Store
					filename = c.GlobalString("i")
				)

				log.Printf("Loading cache file: %s", filename)
				b, err := ioutil.ReadFile(filename)
				if err != nil {
					return err
				}
				err = yaml.Unmarshal(b, &Chapter)

				fmt.Printf("书名: %#v\n", Chapter.BookName)
				fmt.Printf("作者: %#v\n", Chapter.Author)
				fmt.Printf("封面: %s\n", Chapter.CoverURL)
				fmt.Printf("简介: \n\t%v\n", strings.Replace(Chapter.Description, "\n", "\n\t", -1))
				fmt.Printf("章节数: \n")
				for _, v := range Chapter.Volumes {
					var VIP string
					if v.IsVIP {
						VIP = "VIP"
					} else {
						VIP = "免费"
					}
					fmt.Printf("\t%s卷(%s) %d章\n", v.Name, VIP, len(v.Chapters))
				}

				var (
					chCount       = 0
					isDone        = 0
					isExample     = 0
					isDonwExample = 0
				)
				for _, v := range Chapter.Volumes {
					chCount += len(v.Chapters)
					for _, v2 := range v.Chapters {
						if len(v2.Text) != 0 {
							isDone++
						}

						if len(v2.Example) != 0 {
							isExample++
						}
						if (len(v2.Example) != 0) && (len(v2.Text) != 0) {
							isDonwExample++
						}

					}
				}
				if isDone != 0 {
					log.Printf("[读入] 已缓存:%d 样本:%d 完成样本:%d", isDone, isExample, isDonwExample)
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
					Name:  "driver",
					Usage: "请求方式,support: none,phantomjs,chromedp",
				},
				cli.StringFlag{
					Name:  "f",
					Usage: "输出格式",
				},
				cli.StringFlag{
					Name:  "o",
					Usage: "输出路径",
				},
				cli.StringFlag{
					Name:  "chromedp-log",
					Usage: "Chromedp log file",
				},
				cli.DurationFlag{
					Name:        "tsleep",
					Usage:       "章节爬取间隔",
					Value:       200 * time.Millisecond,
					Destination: &tSleep,
				},
				cli.DurationFlag{
					Name:        "errsleep",
					Usage:       "章节爬取错误间隔",
					Value:       500 * time.Millisecond,
					Destination: &errSleep,
				},
			},
			Action: func(c *cli.Context) error {

				fmt.Printf("Commit ID: %s\n", CommitID)
				fmt.Printf("Build Data: %s\n", BuildData)
				fmt.Printf("Build Version: %s\n", Version)

				if logfile := c.GlobalString("log"); logfile != "" {
					fmt.Printf("Set log file: %s\n", logfile)
					f, err := os.Create(logfile)
					if err != nil {
						return err
					}
					defer f.Close()
					log.SetOutput(f)
				}

				var (
					err       error
					URLString string
					bookURL   *url.URL
					Chapter   *store.Store
					filename  = c.GlobalString("i")
				)

				if filename == "" {
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
					switch c.String("driver") {
					case "phantomjs":
						log.Printf("Init PhantomJS")
						site.InitPhantomJS()
						defer func() {
							log.Printf("Close PhantomJS")
							site.ClosePhantomJS()
						}()
						for errCount := 0; errCount < 20; errCount++ {
							Chapter, err = site.PhBookInfo(bookURL.String())
							if err == nil {
								break
							} else {
								log.Printf("ErrCount: %d Err: %s", errCount, err)
								time.Sleep(1000 * time.Millisecond)
							}
						}
					case "chromedp":
						log.Printf("Chromedp Running...")
						Chapter, err = site.ChromedpBookInfo(bookURL.String(), c.String("chromedp-log"))
					default:
						log.Printf("use golang default http")
						for errCount := 0; errCount < 20; errCount++ {
							Chapter, err = site.BookInfo(bookURL.String())
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
					filename = fmt.Sprintf("%s.%s", Chapter.BookName, store.FileExt)
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
						}
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
				} else {
					log.Printf("Loading cache file: %s", filename)
					b, err := ioutil.ReadFile(filename)
					if err != nil {
						return err
					}
					err = yaml.Unmarshal(b, &Chapter)
					if err != nil {
						return err
					}
				}

				if len(c.GlobalStringSlice("turl")) != 0 {
					Chapter.Tmap = []string{}
					for _, v := range c.GlobalStringSlice("turl") {
						Chapter.Tmap = append(Chapter.Tmap, v)
					}
				}

				fmt.Printf("书名: %#v\n", Chapter.BookName)
				fmt.Printf("作者: %#v\n", Chapter.Author)
				fmt.Printf("封面: %s\n", Chapter.CoverURL)
				fmt.Printf("简介: \n\t%v\n", strings.Replace(Chapter.Description, "\n", "\n\t", -1))
				fmt.Printf("章节数: \n")
				for _, v := range Chapter.Volumes {
					var VIP string
					if v.IsVIP {
						VIP = "VIP"
					} else {
						VIP = "免费"
					}
					fmt.Printf("\t%s卷(%s) %d章\n", v.Name, VIP, len(v.Chapters))
				}

				log.Printf("线程数: %d,预缓存中...\n", c.Int("t"))
				ssss := &SyncStore{
					Store: Chapter,
				}
				ssss.Init()

				var chCount = 0
				var isDone = 0
				var isExample = 0
				var isDoneExample = 0
				for _, v := range Chapter.Volumes {
					chCount += len(v.Chapters)
					for _, v2 := range v.Chapters {
						if len(v2.Text) != 0 {
							isDone++
						}
						if len(v2.Example) != 0 {
							isExample++
						}
						if (len(v2.Example) != 0) && (len(v2.Text) != 0) {
							isDoneExample++
						}
					}
				}
				if isDone != 0 {
					log.Printf("[读入] 已缓存:%d 样本:%d 完成样本:%d", isDone, isExample, isDoneExample)
				}

				// End Print
				defer func(s *store.Store) {
					var chCount = 0
					var isDone = 0
					var isExample = 0
					var isDoneExample = 0
					for _, v := range Chapter.Volumes {
						chCount += len(v.Chapters)
						for _, v2 := range v.Chapters {
							if len(v2.Text) != 0 {
								isDone++
							}
							if len(v2.Example) != 0 {
								isExample++
							}
							if (len(v2.Example) != 0) && (len(v2.Text) != 0) {
								isDoneExample++
							}
						}
					}
					if isDone != 0 {
						log.Printf("[爬取结束] 已缓存:%d 样本:%d 完成样本:%d", isDone, isExample, isDoneExample)
					}
				}(Chapter)

				if isDone < chCount {

					bar := processbar.StartNew(chCount)
					bar.Set(isDone + isExample - isDoneExample)

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

					close(Jobch)

					for k3, vol2 := range Chapter.Volumes {
						for k4 := range vol2.Chapters {
							Chapter.Volumes[k3].Chapters[k4].TURL = []string{}
						}
					}

					// 从盗版源下载
					if len(Chapter.Tmap) != 0 {
						log.Printf("开始缓存盗版源")

						for k3, vol2 := range Chapter.Volumes {
							for k4 := range vol2.Chapters {
								Chapter.Volumes[k3].Chapters[k4].TURL = []string{}
							}
						}

						for _, v := range Chapter.Tmap {

							var (
								ts *store.Store
							)
							ts, err = site.BookInfo(v)
							if err != nil {
								return err
							}
							log.Printf("请求盗版源信息: %s 书名:%#v 作者:%#v\n", ts.BookURL, ts.BookName, ts.Author)

							rr, err := regexp.Compile(`（[\S ]*）`)
							if err != nil {
								return err
							}

							cc, err := regexp.Compile(`[•、 ，,!！。\.]+`)
							if err != nil {
								return err
							}

							for k3, vol2 := range Chapter.Volumes {
								for k4, ch2 := range vol2.Chapters {
									for _, vol := range ts.Volumes {
										for _, ch := range vol.Chapters {

											var (
												Name1 string
												Name2 string
											)

											Name1 = ch.Name
											Name2 = ch2.Name

											// sa := "第一百零七章"
											// if strings.Contains(ch.Name, sa) && strings.Contains(ch2.Name, sa) {
											// 	log.Printf("Fuuuuck 1. %#v 2. %#v", Name1, Name2)
											// }

											Name1 = rr.ReplaceAllString(Name1, "")
											Name1 = cc.ReplaceAllString(Name1, "")

											Name2 = rr.ReplaceAllString(Name2, "")
											Name2 = cc.ReplaceAllString(Name2, "")

											var ok = false
											if Name1 == Name2 {
												ok = true
											} else if strings.Contains(Name1, Name2) {
												ok = true
											} else if strings.Contains(Name2, Name1) {
												ok = true
											}

											if ok {
												Chapter.Volumes[k3].Chapters[k4].TURL = append(Chapter.Volumes[k3].Chapters[k4].TURL, ch.URL)
											}
										}
									}
									// if len(Chapter.Volumes[k3].Chapters[k4].TURL) == 0 {
									// 	log.Printf("无源章节: %s %#v", vol2.Name, vol2.Chapters[k4].Name)
									// }
								}
							}

							b, err := yaml.Marshal(Chapter)
							if err != nil {
								return err
							}
							ioutil.WriteFile(filename, b, 0775)
						}

						for _, vol2 := range Chapter.Volumes {
							for _, ch2 := range vol2.Chapters {
								if !vol2.IsVIP {
									continue
								}
								if len(ch2.TURL) == 0 {
									log.Printf("无源章节: %s %#v", vol2.Name, ch2.Name)
								}
							}
						}

						log.Printf("盗版源信息获取完成")

						bar := processbar.StartNew(chCount)
						bar.Set(isDone)

						ssss.IsTWork = true

						Jobch := make(chan error)
						for i := 0; i < c.Int("t"); i++ {
							go TJob(ssss, Jobch)
						}
						cc := make(chan os.Signal)
						signal.Notify(cc, os.Interrupt)

						var ii = 0
					BB:
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
											break BB
										}
									} else {
										log.Printf("Job Error: %s", err)
									}
								} else {
									bar.Increment()
								}
							}
						}
						isDone = 0
						isExample = 0
						chCount = 0
						isDoneExample = 0
						for _, v := range Chapter.Volumes {
							chCount += len(v.Chapters)
							for _, v2 := range v.Chapters {
								if len(v2.Text) != 0 {
									isDone++
								}

								if len(v2.Example) != 0 {
									isExample++
								}
								if (len(v2.Example) != 0) && (len(v2.Text) != 0) {
									isDoneExample++
								}

							}
						}
						log.Printf("[完成] 已缓存:%d 样本:%d 完成样本:%d", isDone, isExample, isDoneExample)
					}
				}

				for k3, vol2 := range Chapter.Volumes {
					for k4, ch2 := range vol2.Chapters {
						newContent := []string{}
						for _, v := range ch2.Text {
							v = strings.TrimSpace(v)
							if v == "" {
								continue
							}
							if strings.Contains(v, "biqiuge") {
								continue
							}
							v = strings.Replace(v, "“”", "", -1)
							newContent = append(newContent, v)
						}
						Chapter.Volumes[k3].Chapters[k4].Text = newContent
					}
				}
				b, err := yaml.Marshal(Chapter)
				if err != nil {
					return err
				}
				if err := ioutil.WriteFile(filename, b, 0775); err != nil {
					return err
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
