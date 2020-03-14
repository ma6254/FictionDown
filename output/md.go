package output

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ma6254/FictionDown/store"
	"github.com/ma6254/FictionDown/utils"
	"gopkg.in/yaml.v2"
)

type MarkdownEPUBmeta struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description,omitempty"`
	Author      string `yaml:"creator,omitempty"`
	Lang        string `yaml:"lang,omitempty"`
	Cover       string `yaml:"cover-image,omitempty"`
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

				meta.Cover = filepath.FromSlash(tempfile.Name())
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
{{.EPUBMeta | yaml_marshal -}}
---
{{end -}}
# 简介

	该文档由FictionDown工具生成
	仅供软件测试评估使用
	请勿传播该文档以及用于此软件评估外的任何用途
<https://github.com/ma6254/FictionDown>

书名: [{{.Store.BookName | markdown}}]({.Store.BookURL})

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
