package main

import (
	"fmt"
	"log"

	"github.com/ma6254/FictionDown/output"

	"github.com/urfave/cli"
)

var convert = cli.Command{
	Name:    "Convert",
	Usage:   "转换格式输出",
	Aliases: []string{"conv"},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "f",
			Usage: "输出格式",
		},
		cli.StringFlag{
			Name:  "o",
			Usage: "输出路径",
		},
		cli.BoolFlag{
			Name:        "ignore_cover",
			Usage:       "忽略封面",
			Destination: &outputOpt.IgnoreCover,
		},
		cli.BoolFlag{
			Name:        "no-EPUB-metadata",
			Usage:       "禁用EPUB元数据",
			Destination: &outputOpt.NoEPUBMetadata,
		},
	},
	Action: func(c *cli.Context) error {

		if err := initLoadStore(c); err != nil {
			return err
		}

		var (
			format     = c.String("f")
			outputpath = c.String("o")
		)
		if format == "" {
			return nil
		}

		var ConversionFileName string
		if outputpath == "" {
			ConversionFileName = fmt.Sprintf("%s.%s", chapter.BookName, format)
		}
		log.Printf("Start Conversion: Format:%#v OutPath:%#v", c.String("f"), ConversionFileName)
		return output.Output(*chapter, format, ConversionFileName, outputOpt)
	},
}

var outputOpt output.Option
