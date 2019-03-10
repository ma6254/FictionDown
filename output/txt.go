package output

import (
	"fmt"
	"os"
	"strings"

	"github.com/ma6254/FictionDown/store"
)

type TXT struct {
}

func (t *TXT) Conv(src store.Store, outpath string) (err error) {
	f, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "%s\n", src.BookName)
	fmt.Fprintf(f, "作者: %s\n", src.Author)
	fmt.Fprintf(f, "书本链接: %s\n", src.BookURL)
	fmt.Fprintf(f, "简介:\n")
	for _, d := range strings.Split(src.Description, "\n") {
		fmt.Fprintf(f, "\t%s\n", d)
	}
	fmt.Fprintf(f, "\n\n")
	for _, v1 := range src.Volumes {
		var VIP string
		if v1.IsVIP {
			VIP = "收费"
		} else {
			VIP = "免费"
		}
		fmt.Fprintf(f, "卷: %#v %s\n", v1.Name, VIP)
		fmt.Fprintf(f, "\n\n")
		for _, v2 := range v1.Chapters {
			fmt.Fprintf(f, "%s\n", v2.Name)
			for _, cc := range v2.Text {
				fmt.Fprintf(f, "\t%s\n",
					strings.Replace(cc, "*", "□", -1),
				)
			}
			fmt.Fprintf(f, "\n\n")
		}
	}
	return nil
}
