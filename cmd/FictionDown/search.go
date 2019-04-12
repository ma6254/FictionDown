package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ma6254/FictionDown/utils"

	"github.com/go-yaml/yaml"
	"github.com/ma6254/FictionDown/site"

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
	},
	Action: func(c *cli.Context) error {
		keyword := c.String("keyword")
		r, err := site.Search(keyword)
		if err != nil {
			return err
		}
		if !c.Bool("put") {
			fmt.Printf("搜索到%d个内容:\n", len(r))
			for _, v := range r {
				fmt.Printf("%s %s %s\n", v.BookURL, v.BookName, v.Author)
			}
		} else {
			err := initLoadStore(c)
			if err != nil {
				return err
			}
			rrr := []site.ChaperSearchResult{}
			for _, v := range r {
				if (v.Author == chapter.Author) && (v.BookName == chapter.BookName) {
					log.Printf("%s %s %s", v.BookURL, v.BookName, v.Author)
					rrr = append(rrr, v)
				}
			}
			for _, v := range rrr {
				chapter.Tmap = append(chapter.Tmap, v.BookURL)
			}
			chapter.Tmap = utils.TupleSlice(chapter.Tmap)
			b, err := yaml.Marshal(chapter)
			if err != nil {
				return err
			}
			ioutil.WriteFile(filename, b, 0775)
		}
		return nil
	},
}
