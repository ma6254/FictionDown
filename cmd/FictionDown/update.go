package main

import (
	"fmt"
	"strings"

	"github.com/ma6254/FictionDown/site"
	"github.com/ma6254/FictionDown/store"
	"github.com/urfave/cli"
)

var update = cli.Command{
	Name:    "update",
	Aliases: []string{"u"},
	Usage:   "检查更新信息",
	Action: func(c *cli.Context) error {

		if err := initLoadStore(c); err != nil {
			return err
		}

		updateStore, err := site.BookInfo(chapter.BookURL)
		if err != nil {
			return nil
		}

		chapterLen := 0
		for _, v := range updateStore.Volumes {
			chapterLen += len(v.Chapters)
		}

		oldChapterLen := 0
		for _, v := range chapter.Volumes {
			oldChapterLen += len(v.Chapters)
		}

		fmt.Printf("书名: %#v\n", chapter.BookName)
		fmt.Printf("作者: %#v\n", chapter.Author)
		fmt.Printf("封面: %s\n", chapter.CoverURL)
		fmt.Printf("简介: \n\t%v\n", strings.Replace(chapter.Description, "\n", "\n\t", -1))
		fmt.Printf("章节数: \n")
		for k, v := range chapter.Volumes {
			var VIP string
			if v.IsVIP {
				VIP = "收费"
			} else {
				VIP = "免费"
			}
			fmt.Printf("\t%s(%s) %d章 => ", v.Name, VIP, len(v.Chapters))

			for kk, vv := range updateStore.Volumes {
				if (vv.Name == v.Name) && (kk == k) {
					fmt.Printf("%d章 ", len(vv.Chapters))
				}
			}
			fmt.Println()
		}

		diffVol := []store.Volume{}
		copy(diffVol, updateStore.Volumes)

		for k, v := range chapter.Volumes {
			for kk, vv := range updateStore.Volumes {
				if (vv.Name == v.Name) && (kk == k) {

				}
			}
		}

		return nil
	},
}
