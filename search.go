package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ma6254/FictionDown/site"
	"gopkg.in/yaml.v2"

	"github.com/urfave/cli"
)

var search = cli.Command{
	Name:    "search",
	Aliases: []string{"s"},
	Usage:   "检索盗版站点",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "k,keyword",
			Usage: "搜索关键词",
		},
		cli.BoolFlag{
			Name:  "put,p",
			Usage: "对比并放入缓存文件",
		},
		cli.StringFlag{
			Name:  "save,s",
			Usage: "searh and save to file 搜索结果存储",
		},
		cli.StringFlag{
			Name:  "format,f",
			Usage: "save format support json,yaml 存储格式",
		},
		cli.BoolFlag{
			Name:  "download,d",
			Usage: "直接下载正版",
		},
	},
	Action: func(c *cli.Context) error {

		var (
			err        error
			result     = make(map[site.SearchBookMeta][]string)
			keyword    = c.String("keyword")
			savepath   = c.String("save")
			saveformat = c.String("format")
			dumpFunc   func(v interface{}) ([]byte, error)
		)

		if savepath != "" {
			fmt.Printf("Save File: %s\n", savepath)
			if saveformat == "" {
				fileext := filepath.Ext(savepath)
				switch fileext {
				case "", ".json":
					saveformat = "json"
				case ".yaml", ".yml":
					saveformat = "yaml"
				default:
					return fmt.Errorf("Unsupported file extension %s", fileext)
				}
			}
			fmt.Printf("Save Format: %s\n", saveformat)
			switch saveformat {
			case "json":
				dumpFunc = func(v interface{}) ([]byte, error) {
					return json.MarshalIndent(v, "", "\t")
				}
			case "yaml":
				dumpFunc = yaml.Marshal
			default:
				return fmt.Errorf("unsupport marshal format: %s", saveformat)
			}
		}

		r, err := site.Search(keyword)
		if err != nil {
			return err
		}

		for _, v := range r {
			nvName := strings.TrimSuffix(v.BookName, "【完结】")
			// nvAuthor := strings.TrimLeft(v.Author, `\/`)
			nvAuthor := v.Author
			result[site.SearchBookMeta{Name: nvName, Author: nvAuthor}] = append(result[site.SearchBookMeta{Name: nvName, Author: nvAuthor}], v.BookURL)
		}
		rr := site.ConvSearchRequest(result)

		if savepath != "" {
			var body []byte
			if body, err = dumpFunc(rr); err != nil {
				return err
			}
			return ioutil.WriteFile(savepath, body, os.ModePerm)
		}

		if c.Bool("put") {
			if err = initLoadStore(c); err != nil {
				return err
			}
			rrr := result[site.SearchBookMeta{Name: chapter.BookName, Author: chapter.Author}]
			chapter.Tmap = rrr
			b, err := yaml.Marshal(chapter)
			if err != nil {
				return err
			}
			return ioutil.WriteFile(filename, b, os.ModePerm)
		}

		if c.Bool("download") {
			fmt.Printf("开始缓存正版源\n")
			vipBookURL, tBookURL :=
				SearchMatchSites(rr[0].Urls)
			if len(vipBookURL) > 1 {
				return fmt.Errorf("多个正版源")
			}
			fmt.Printf("正版：\n")
			for _, u := range vipBookURL {
				matchSite, err := site.MatchOne(site.Sitepool, u)
				n := matchSite.Name
				if err != nil {
					n = "无匹配站点"
				}
				fmt.Printf("\t%#v : %#v\n", n, u)
			}
			fmt.Printf("盗版：\n")
			for _, u := range tBookURL {
				matchSite, err := site.MatchOne(site.Sitepool, u)
				n := matchSite.Name
				if err != nil {
					n = "无匹配站点"
				}
				fmt.Printf("\t%#v : %#v\n", n, u)
			}
			if len(vipBookURL) == 0 {
				return fmt.Errorf("无正版站点")
			}
			bookurl = vipBookURL[0]
			if err = download.Run(c); err != nil {
				return err
			}
			chapter.Tmap = tBookURL
			b, err := yaml.Marshal(chapter)
			if err != nil {
				return err
			}
			if err = ioutil.WriteFile(filename, b, os.ModePerm); err != nil {
				return err
			}
			if err = download.Run(c); err != nil {
				return err
			}
			return nil
		}

		fmt.Printf("搜索到%d个内容:\n", len(rr))
		for _, v := range rr {
			fmt.Printf("书名: %s 作者: %s %d个书源\n", v.Name, v.Author, len(v.Urls))
			for _, u := range v.Urls {
				matchSite, err := site.MatchOne(site.Sitepool, u)
				n := matchSite.Name
				if err != nil {
					n = "无匹配站点"
				}
				fmt.Printf("\t%#v : %#v\n", n, u)
			}
		}

		return nil
	},
}

func SearchMatchSites(urls []string) (vipURL, tURL []string) {
	for _, u := range urls {
		matchSite, err := site.MatchOne(site.Sitepool, u)
		if err != nil {
			continue
		}
		isvip := false
		for _, v := range matchSite.Tags() {
			if v == "正版" {
				isvip = true
				break
			}
		}
		if isvip {
			vipURL = append(vipURL, u)
			continue
		}
		tURL = append(tURL, u)
	}
	return
}
