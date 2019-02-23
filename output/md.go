package output

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/ma6254/FictionDown/store"
	pb "gopkg.in/cheggaaa/pb.v1"
)

type Markdown struct {
}

func (t *Markdown) Conv(src store.Store, outpath string) (err error) {

	o := ""

	o += fmt.Sprintf("---\n")
	o += fmt.Sprintf("title: %#v\n", src.BookName)
	o += fmt.Sprintf("description: %#v\n", src.Description)
	o += fmt.Sprintf("creator: %#v\n", src.Author)

	if src.CoverURL != "" {
		resp, err := http.Get(src.CoverURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		coverBuf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		tempfile, err := ioutil.TempFile("", "book_cover_*.jpg")
		if err != nil {
			return err
		}

		ioutil.WriteFile(tempfile.Name(), coverBuf, 0775)

		log.Printf("Save Cover Image: %#v", tempfile.Name())

		o += fmt.Sprintf("cover-image: %#v\n", tempfile.Name())
	}
	o += fmt.Sprintf("---\n\n")

	o += fmt.Sprintf("# 简介\n\n")
	dlist := strings.Split(src.Description, "\n")

	for _, cc := range dlist {
		o += fmt.Sprintf("<p style=\"text-indent:2em\">%s</p>\n", cc)
	}
	o += "\n"

	for _, v1 := range src.Volumes {
		o += fmt.Sprintf("# %s\n\n", v1.Name)
		log.Printf("正在转换卷: %s", v1.Name)
		bar := pb.StartNew(len(v1.Chapters))
		for _, v2 := range v1.Chapters {
			// s += fmt.Sprintf(`<h1><a href=%#v>%s</a></h1>`, v2.URL, v2.Name)
			o += fmt.Sprintf("## <a href=%#v>%s</a>\n\n", v2.URL, v2.Name)
			for _, cc := range v2.Text {
				o += fmt.Sprintf("<p style=\"text-indent:2em\">%s</p>\n", cc)
			}
			bar.Increment()
			o += "\n"
		}
		bar.Finish()
	}
	return ioutil.WriteFile(outpath, []byte(o), 0775)
}
