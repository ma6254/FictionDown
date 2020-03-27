package site

import (
	"sort"
)

type SearchBookMeta struct {
	Name   string
	Author string
}

type SearchBookMetaA struct {
	Name   string
	Author string
	Urls   []string
}

type SortSearchMeteList []SearchBookMetaA

func (a SortSearchMeteList) Len() int           { return len(a) }
func (a SortSearchMeteList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortSearchMeteList) Less(i, j int) bool { return len(a[i].Urls) > len(a[j].Urls) }

func ConvSearchRequest(r map[SearchBookMeta][]string) []SearchBookMetaA {
	rs := make([]SearchBookMetaA, 0)
	for rk, rv := range r {
		rs = append(rs, SearchBookMetaA{rk.Name, rk.Author, rv})
	}
	if len(rs) < 2 {
		return rs
	}
	sort.Sort(SortSearchMeteList(rs))
	return rs
}
