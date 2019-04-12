package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/ma6254/FictionDown/matching"
	"github.com/ma6254/FictionDown/site"
	"github.com/ma6254/FictionDown/store"
	"github.com/urfave/cli"
	processbar "gopkg.in/cheggaaa/pb.v1"
)

var (
	tSleep   time.Duration
	errSleep time.Duration
)

var download = cli.Command{
	Name:    "download",
	Usage:   "下载缓存文件",
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

		if err := initLoadStore(c); err != nil {
			return err
		}

		if (CommitID != "") && (BuildData != "") && (Version != "") {
			fmt.Printf("Commit ID: %s\n", CommitID)
			fmt.Printf("Build Data: %s\n", BuildData)
			fmt.Printf("Build Version: %s\n", Version)
		}

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
			err error
		)

		fmt.Printf("书名: %#v\n", chapter.BookName)
		fmt.Printf("作者: %#v\n", chapter.Author)
		fmt.Printf("封面: %s\n", chapter.CoverURL)
		fmt.Printf("简介: \n\t%v\n", strings.Replace(chapter.Description, "\n", "\n\t", -1))
		fmt.Printf("章节数: \n")
		for _, v := range chapter.Volumes {
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
			Store: chapter,
		}
		ssss.Init()

		var chCount = 0
		var isDone = 0
		var isExample = 0
		var isDoneExample = 0
		for _, v := range chapter.Volumes {
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
			for _, v := range chapter.Volumes {
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
		}(chapter)

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

			for k3, vol2 := range chapter.Volumes {
				for k4 := range vol2.Chapters {
					chapter.Volumes[k3].Chapters[k4].TURL = []string{}
				}
			}

			log.Printf("生成别名")
			for k3, vol2 := range chapter.Volumes {
				for k4, chaper := range vol2.Chapters {
					chapter.Volumes[k3].Chapters[k4].Alias = matching.TitleAlias(chaper.Name)
				}
			}

			// 从盗版源下载
			if len(chapter.Tmap) != 0 {

				log.Printf("开始缓存盗版源")

				for k3, vol2 := range chapter.Volumes {
					for k4 := range vol2.Chapters {
						chapter.Volumes[k3].Chapters[k4].TURL = []string{}
					}
				}

				for _, v := range chapter.Tmap {

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

					for k3, vol2 := range chapter.Volumes {
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
										chapter.Volumes[k3].Chapters[k4].TURL = append(chapter.Volumes[k3].Chapters[k4].TURL, ch.URL)
									}
								}
							}
							// if len(Chapter.Volumes[k3].Chapters[k4].TURL) == 0 {
							// 	log.Printf("无源章节: %s %#v", vol2.Name, vol2.Chapters[k4].Name)
							// }
						}
					}

					b, err := yaml.Marshal(chapter)
					if err != nil {
						return err
					}
					ioutil.WriteFile(filename, b, 0775)
				}

				for _, vol2 := range chapter.Volumes {
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
				for _, v := range chapter.Volumes {
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

		for k3, vol2 := range chapter.Volumes {
			for k4, ch2 := range vol2.Chapters {
				newContent := []string{}
				for _, v := range ch2.Text {
					v = strings.TrimSpace(v)
					if v == "" {
						continue
					}
					v = strings.Replace(v, "“”", "", -1)
					if regexp.MustCompile("^[…]+$").MatchString(v) {
						continue
					}
					newContent = append(newContent, v)
				}
				chapter.Volumes[k3].Chapters[k4].Text = newContent
			}
		}
		b, err := yaml.Marshal(chapter)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(filename, b, 0775); err != nil {
			return err
		}

		return nil
	},
}
