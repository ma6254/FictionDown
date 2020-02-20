package output

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	goepub "github.com/bmaupin/go-epub"
	"github.com/ma6254/FictionDown/store"
	"github.com/ma6254/FictionDown/utils"
)

type EPUB struct {
}

func (t *EPUB) Conv(src store.Store, outpath string, opts Option) (err error) {
	e := goepub.NewEpub(src.BookName)
	e.SetLang("中文")
	e.SetAuthor(src.Author)

	if src.CoverURL != "" {

		body, err := utils.GetWebPageBodyReader(src.CoverURL)
		if err != nil {
			return err
		}
		tempfile, err := ioutil.TempFile("", "book_cover_*.jpg")
		if err != nil {
			return err
		}
		coverBuf, _ := ioutil.ReadAll(body)
		ioutil.WriteFile(tempfile.Name(), coverBuf, 0775)

		log.Printf("Save Cover Image: %#v", tempfile.Name())

		e.AddImage(tempfile.Name(), "cover.jpg")
		e.SetCover("cover.jpg", "")
	}

	d := ""
	dlist := strings.Split(src.Description, "\n")
	for _, cc := range dlist {
		d += fmt.Sprintf(`<p style="text-indent:2em">%s</p>`, cc)
	}
	// Description := fmt.Sprintf(`<h1><a href=%#v>%s</a></h1>%s`, src.BiqugeURL, src.BookName, d)
	Description := fmt.Sprintf(`<h1>%s</h1>%s`, src.BookName, d)
	_, err = e.AddSection(Description, "简介", "Cover.xhtml", "")
	if err != nil {
		return err
	}
	for k1, v1 := range src.Volumes {
		for k2, v2 := range v1.Chapters {
			s := ""
			// s += fmt.Sprintf(`<h1><a href=%#v>%s</a></h1>`, v2.URL, v2.Name)
			s += fmt.Sprintf(`<h1>%s</h1>`, v2.Name)
			for _, cc := range v2.Text {
				s += fmt.Sprintf(`<p style="text-indent:2em">%s</p>`, cc)
			}
			_, err = e.AddSection(s, v2.Name, fmt.Sprintf("%d-%d.xhtml", k1, k2), "")
			if err != nil {
				return err
			}
		}
	}
	err = e.Write(outpath)
	return
}
