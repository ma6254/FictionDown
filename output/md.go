package output

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-yaml/yaml"
	"github.com/ma6254/FictionDown/store"
)

type MarkdownEPUBmeta struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Author      string `yaml:"creator"`
	Lang        string `yaml:"lang"`
	Cover       string `yaml:"cover-image"`
}

type Markdown struct {
}

func (t *Markdown) Conv(src store.Store, outpath string, opts Option) (err error) {

	var (
		meta MarkdownEPUBmeta
		temp *template.Template
	)

	f, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer f.Close()

	if !opts.NoEPUBMetadata {
		meta = MarkdownEPUBmeta{
			Title:       src.BookName,
			Description: src.Description,
			Author:      src.Author,
			Lang:        "zh-CN",
		}
		if !opts.IgnoreCover {
			if src.CoverURL != "" {
				client := &http.Client{}
				req, err := http.NewRequest("GET", src.CoverURL, nil)
				if err != nil {
					return err
				}
				req.Header.Add(
					"user-agent",
					"Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Mobile Safari/537.36",
				)
				resp, err := client.Do(req)
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

				meta.Cover = tempfile.Name()
			}
		}
	}

	temp = template.New("markdown_fiction")

	temp = temp.Funcs(template.FuncMap{
		"yaml_marshal": func(in interface{}) (string, error) {
			a, err := yaml.Marshal(in)
			return string(a), err
		},
		"split":    strings.Split,
		"markdown": MarkdownEscape,
	})

	temp, err = temp.Parse(MarkdownTemplate)
	if err != nil {
		return err
	}

	return temp.Execute(
		f, MarkdownTemplateValues{
			Store:    src,
			Opts:     opts,
			EPUBMeta: meta,
		})
}

func MarkdownEscape(s string) string {
	for _, v := range "\\!\"#$%&'()*+,./:;<=>?@[]^_`{|}~-" {
		s = strings.Replace(s, string(v), "\\"+string(v), -1)
	}
	return s
}

// MarkdownTemplate is Markdown format Template
//  Parameters:
//  - `store`: Store.store
//  - `opts`: Options
var MarkdownTemplate = `
{{- if not .Opts.NoEPUBMetadata -}}
---
{{.EPUBMeta | yaml_marshal}}
---
{{end -}}
# 简介

书名: {{.Store.BookName}}
作者: {{.Store.Author}}
简介: 
{{range split .Store.Description "\n" -}}
<p style="text-indent:2em">{{. | markdown}}</p>
{{end -}}
{{range .Store.Volumes }}
# {{.Name | markdown}} {{if .IsVIP}}付费{{else}}免费{{end}}卷
{{range .Chapters}}
## {{.Name | markdown}}

{{range .Text -}}
<p style="text-indent:2em">{{. | markdown}}</p>
{{end}}
{{end}}{{end}}
`

type MarkdownTemplateValues struct {
	Store    store.Store
	Opts     Option
	EPUBMeta MarkdownEPUBmeta
}
