package matching

import (
	"regexp"

	"github.com/ma6254/FictionDown/utils"
)

//TitleAlias 获取标题别称
func TitleAlias(s string) (alias []string) {
	var res = []string{
		`^([第]?([零一二三四五六七八就十百千万0-9]+)[章\,，、]?)[ ]*([^ （）\(\))]+).*$`,
		`([^ （）\(\)]+)`,   // 没有章节号的章节
		`[（\(](.+?)[）\)]`, //括号内的
	}

	for _, v := range res {
		re, err := regexp.Compile(v)
		if err != nil {
			panic(err)
		}
		find := re.FindAllStringSubmatch(s, -1)
		if find == nil {
			continue
		}
		// log.Printf("find: %#v", find)
		for _, v1 := range find {
			for _, v := range v1[1:] {
				if v == "" {
					continue
				}
				if v == s {
					continue
				}
				if utils.StringInSlice(v, alias) {
					continue
				}
				alias = append(alias, v)
			}
		}
	}
	return
}
