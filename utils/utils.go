package utils

import (
	"fmt"
	"net/http"
)

func RequestGet(u string) (resp *http.Response, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return
	}
	req.Header.Add(
		"user-agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36",
	)
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("%d %v", resp.StatusCode, resp.Status)
		return
	}
	return resp, nil
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
	b := []string{}
	ia := map[string]int{}
	for _, v := range a {
		if ia[v] == 0 {
			b = append(b, v)
		}
		ia[v]++
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

func StringSliceEqual(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
