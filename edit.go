package main

import (
	"github.com/urfave/cli"
)

var edit = cli.Command{
	Name:    "edit",
	Usage:   "对缓存文件进行手动修改",
	Flags:   []cli.Flag{},
	Aliases: []string{"e"},
	Action: func(c *cli.Context) error {
		if err := initLoadStore(c); err != nil {
			return err
		}
		return nil
	},
}
