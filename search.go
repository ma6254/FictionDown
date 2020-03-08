package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ma6254/FictionDown/site"
	"gopkg.in/yaml.v2"

	"github.com/urfave/cli"
)

type SearchBookMeta struct {
	Name   string
	Author string
}

type SearchBookMetaA struct {
	Name   string
	Author string
	Urls   []string
}

type SearchMeteList []SearchBookMetaA

func (a SearchMeteList) Len() int           { return len(a) }
func (a SearchMeteList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SearchMeteList) Less(i, j int) bool { return len(a[i].Urls) > len(a[j].Urls) }

func ConvSearchRequest(r map[SearchBookMeta][]string) []SearchBookMetaA {
	rs := make([]SearchBookMetaA, 0)
	for rk, rv := range r {
		rs = append(rs, SearchBookMetaA{rk.Name, rk.Author, rv})
	}
	if len(rs) < 2 {
		return rs
	}
	sort.Sort(SearchMeteList(rs))
	return rs
}

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
	},
	Action: func(c *cli.Context) error {

		var (
			err        error
			result     = make(map[SearchBookMeta][]string)
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
				case ".yaml":
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
			case "yaml", "yml":
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
			result[SearchBookMeta{nvName, nvAuthor}] = append(result[SearchBookMeta{nvName, nvAuthor}], v.BookURL)
		}
		rr := ConvSearchRequest(result)

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
			rrr := result[SearchBookMeta{chapter.BookName, chapter.Author}]
			chapter.Tmap = rrr
			b, err := yaml.Marshal(chapter)
			if err != nil {
				return err
			}
			return ioutil.WriteFile(filename, b, os.ModePerm)
		}

		fmt.Printf("搜索到%d个内容:\n", len(rr))
		for _, v := range rr {
			fmt.Printf("书名: %s 作者: %s %d个书源\n", v.Name, v.Author, len(v.Urls))
			for _, u := range v.Urls {
				fmt.Printf("\t%s\n", u)
			}
		}
		return nil
	},
}
