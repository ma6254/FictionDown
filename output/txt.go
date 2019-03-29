package output

import (
	"html/template"
	"os"
	"strings"

	"github.com/ma6254/FictionDown/store"
)

type TXT struct {
}

func (t *TXT) Conv(src store.Store, outpath string, opts Option) (err error) {
	f, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer f.Close()

	temp := template.New("txt_fiction")
	temp = temp.Funcs(template.FuncMap{
		"split": strings.Split,
	})
	temp, err = temp.Parse(TxtTemplate)
	if err != nil {
		return err
	}

	return temp.Execute(
		f, src)
}

var TxtTemplate = `书名：{{.BookName}}
作者：{{.Author}}
链接：{{.BookURL}}
简介：
{{range split .Description "\n"}}	{{.}}
{{end}}
{{- range .Volumes }}
{{if .IsVIP}}付费{{else}}免费{{end}}卷 {{.Name}}
{{range .Chapters}}
{{.Name}}
{{range .Text}}	{{.}}
{{end}}{{end}}{{end}}`
