package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/urfave/cli"
)

var check = cli.Command{
	Name:    "check",
	Usage:   "检查缓存文件",
	Aliases: []string{"c", "chk"},
	Flags:   []cli.Flag{},
	Action: func(c *cli.Context) error {

		if err := initLoadStore(c); err != nil {
			return err
		}

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

		var (
			chCount       = 0
			isDone        = 0
			isExample     = 0
			isDonwExample = 0
		)
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
					isDonwExample++
				}

			}
		}
		if isDone != 0 {
			log.Printf("[读入] 已缓存:%d 样本:%d 完成样本:%d", isDone, isExample, isDonwExample)
		}

		return nil
	},
}
