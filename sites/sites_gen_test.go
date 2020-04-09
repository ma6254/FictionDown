package sites

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/ma6254/FictionDown/site"
	"github.com/ma6254/FictionDown/store"
	"github.com/ma6254/FictionDown/utils"
)

func GenBookInfoSite(s site.SiteA) func(t *testing.T) {
	type TestCase struct {
		Name       string
		Author     string
		ChapterNum int               // 章节数目
		URL        map[string]string // Key:SiteName Value:BookURL
	}
	dd := []TestCase{
		{
			Name:   "放开那个女巫",
			Author: "二目",
			URL: map[string]string{
				"起点中文网":  "https://book.qidian.com/info/1003306811",
				"笔趣阁1":   "https://www.biquge5200.cc/39_39136/",
				"新八一中文网": "https://www.81new.net/11/11609/",
				"书迷楼":    "https://www.shumil.co/fangkainagenvwu/",
				"顶点小说":   "http://www.booktxt.net/book/goto/id/2600",
				"完本神站":   "https://www.wanbentxt.com/1949/",
				"38看书":   "https://www.38kanshu.com/9839/",
			},
		},
		{
			Name:   "黎明之剑",
			Author: "远瞳",
			URL: map[string]string{
				"起点中文网":  "https://book.qidian.com/info/1010400217",
				"笔趣阁1":   "https://www.biquge5200.cc/95_95192/",
				"顶点小说":   "http://www.booktxt.net/book/goto/id/5414",
				"书迷楼":    "https://www.shumil.co/limingzhijian/",
				"新八一中文网": "https://www.81new.net/44/44290/",
				"完本神站":   "https://www.wanbentxt.com/2817/",
				"38看书":   "https://www.38kanshu.com/1897/",
			},
		},
		{
			Name:   "异常生物见闻录",
			Author: "远瞳",
			URL: map[string]string{
				"起点中文网":  "https://book.qidian.com/info/3242304",
				"顶点小说":   "http://www.booktxt.net/book/goto/id/10",
				"笔趣阁1":   "https://www.biquge5200.cc/0_799/",
				"书迷楼":    "https://www.shumil.co/yichangshengwujianwenlu/",
				"完本神站":   "https://www.wanbentxt.com/643/",
				"新八一中文网": "https://www.81new.net/15/15408/",
				"38看书":   "https://www.38kanshu.com/92489/",
			},
		},
		{
			Name:   "战略级天使",
			Author: "白伯欢",
			URL: map[string]string{
				"完本神站": "https://www.wanbentxt.com/15287/",
			},
		},
	}
	CheckFunc := func(t *testing.T, tc TestCase, st *store.Store) (err error) {
		if tc.Name != st.BookName {
			return fmt.Errorf("BookName check does not match")
		}
		if tc.Author != st.Author {
			return fmt.Errorf("Author check does not match want:%#v but:%#v", tc.Author, st.Author)
		}
		chapterNum := 0
		for _, vol := range st.Volumes {
			chapterNum += len(vol.Chapters)
		}
		if chapterNum == 0 {
			return fmt.Errorf("ChapterNum check does not match")
		}
		t.Logf("%s %s %s %d", tc.Name, tc.Author, st.BookURL, chapterNum)
		return
	}
	return func(t *testing.T) {
		t.Parallel()
		var (
			err       error
			st        *store.Store
			matchSite []*site.SiteA
		)
		for _, d := range dd {
			u, ok := d.URL[s.Name]
			if !ok {
				continue
			}
			matchSite, err = site.MatchSites(site.Sitepool, u)
			if err != nil {
				t.Logf("matchone site fail: %s", err)
				t.Fail()
				continue
			}
			if len(matchSite) != 1 {
				t.Logf("matchone site uniqueness: %d %#v", len(matchSite), matchSite)
				t.Fail()
				continue
			}
			if (matchSite[0].Name != s.Name) || (matchSite[0].HomePage != s.HomePage) {
				t.Logf("match is not equals want:%#v but:%#v", s.Name, matchSite[0].Name)
				t.Fail()
				continue
			}
			if st, err = site.BookInfo(u); err != nil {
				t.Logf("%s", err)
				t.Fail()
				continue
			}
			if err = CheckFunc(t, d, st); err != nil {
				t.Logf("%s %s", d.Name, err)
				t.Fail()
				continue
			}
		}
	}
}

func GenSearchSite(s site.SiteA) func(t *testing.T) {
	dd := []struct {
		Name   string
		Author string
	}{
		// 起点
		{"诡秘之主", "爱潜水的乌贼"},
		{"放开那个女巫", "二目"},
		{"黎明之剑", "远瞳"},
		{"异常生物见闻录", "远瞳"},
		// 有毒
		{"绿龙筑巢记", "归兮北冥"},
		{"底栖魔鱼日记", "辣鸡葱花"},
		// 书客
		{"俺，龙领主", "熊瀚"},
		{"红龙", "接口卡"},
		// 小红花
		{"战略级天使", "白伯欢"},
	}
	return func(t *testing.T) {
		t.Parallel()
		t.Log("============================================")
		t.Logf("Site: %s %s %s", s.Name, s.HomePage, s.File)
		t.Log("============================================")
		if s.Search == nil {
			t.Logf("site search func is empty")
			t.Fail()
			return
		}
		isOK := false
		for _, b := range dd {
			var (
				result []site.ChaperSearchResult
				err    error
			)
			t.Logf(">>>>> %s %s <<<<<", b.Name, b.Author)
			if utils.Retry(3, 1*time.Second, func() error {
				result, err = s.Search(b.Name)
				return err
			}); err != nil {
				t.Fatalf("%s %s %s", s.Name, s.HomePage, err)
			}
			for _, r := range result {
				if (r.Author == b.Author) && (r.BookName == b.Name) {
					isOK = true
					continue
				}
			}
			tmpStr := ""
			for _, v := range result {

				if searchSiteResult[SiteMeta{s.Name, s.HomePage}] == nil {
					searchSiteResult[SiteMeta{s.Name, s.HomePage}] = make(map[SearchSite][]string)
				}
				searchSiteResult[SiteMeta{s.Name, s.HomePage}][SearchSite{v.BookName, v.Author}] =
					append(
						searchSiteResult[SiteMeta{s.Name, s.HomePage}][SearchSite{v.BookName, v.Author}],
						v.BookURL)
				tmpStr += fmt.Sprintf("[%s(%s)](%#v) ", v.BookName, v.Author, v.BookURL)
			}

			t.Logf("%s(%s) %d %s", b.Name, b.Author, len(result), tmpStr)
		}
		if !isOK {
			t.Logf("%s %s %s", s.Name, s.HomePage, "搜索结果无效")
			t.Fail()
			return
		}
		b, err := yaml.Marshal(searchSiteResult)
		if err != nil {
			t.Fatal(err)
		}
		if err := ioutil.WriteFile(fmt.Sprintf("TestSearch.yml"), b, os.ModePerm); err != nil {
			t.Fatal(err)
		}
		t.Run(fmt.Sprintf("BookInfo"), GenBookInfoSite(s))

	}
}
