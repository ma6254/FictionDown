package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
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
	err = Retry(10, time.Millisecond*500, func() error {
		resp, err = client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			return fmt.Errorf("%d %v", resp.StatusCode, resp.Status)
		}
		return err
	})
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

func U8ToGBK(a string) string {
	b, _, _ := transform.String(simplifiedchinese.GBK.NewEncoder(), a)
	return b
}

func detectContentCharset(body io.Reader) encoding.Encoding {
	data, err := bufio.NewReader(body).Peek(1024)
	if err != nil {
		panic(err)
	}
	e, _, _ := charset.DetermineEncoding(data, "")
	return e
}

func GetWebPageBodyReader(u string) (r io.Reader, err error) {

	var resp *http.Response

	if err = Retry(5, time.Millisecond*500, func() (err error) {
		resp, err = RequestGet(u)
		return err
	}); err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var body io.Reader = bytes.NewReader(bodyBytes)

	if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		encode := detectContentCharset(bytes.NewReader(bodyBytes))
		body = transform.NewReader(body, encode.NewDecoder())
	}
	return body, nil
}

func GetWegPageDOM(u string) (node *html.Node, err error) {
	body, err := GetWebPageBodyReader(u)
	if err != nil {
		return
	}
	return htmlquery.Parse(body)
}
