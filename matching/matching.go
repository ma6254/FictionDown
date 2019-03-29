package matching

import (
	"regexp"
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
				if StringInSlice(v, alias) {
					continue
				}
				alias = append(alias, v)
			}
		}
	}
	return
}

// StringInSlice string in []stirng like python "if a in b" keyword
func StringInSlice(s string, ss []string) bool {
	for _, v := range ss {
		if s == v {
			return true
		}
	}
	return false
}

//TupleSlice 去除重复字符串
func TupleSlice(a []string) []string {
	b := make([]string, len(a))
	ia := make([]int, len(a))
	for k, v := range a {
		if ia[k] == 0 {
			b = append(b, v)
		}
		ia[k]++
	}
	return b
}

//SimilarSlice 对比两个字符串组，得到其中相等字符串的数量，"i < len(a)" and "i < len(b)"
func SimilarSlice(a, b []string) (i int) {
	a = TupleSlice(a)
	b = TupleSlice(b)
	for _, va := range a {
	B:
		for _, vb := range b {
			if va == vb {
				i++
				break B
			}
		}
	}
	return
}
