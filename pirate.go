package main

import (
	"log"

	"github.com/ma6254/FictionDown/site"
	"github.com/urfave/cli"
)

var pirate = cli.Command{
	Name:    "pirate",
	Aliases: []string{"p"},
	Usage:   "检索盗版站点",
	Flags:   []cli.Flag{},
	Action: func(c *cli.Context) error {
		a := "https://www.biqiuge.com/book/4772/480965712.html/"

		s, err := site.MatchOne(site.Sitepool, a)
		if err != nil {
			return err
		}
		log.Printf("匹配站点: %s %#v", s.Name, s.HomePage)
		return nil
	},
}
