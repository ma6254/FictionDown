package site

import (
	"fmt"
	"io"
	"testing"

	"github.com/ma6254/FictionDown/store"
	"github.com/ma6254/FictionDown/utils"
)

func TestSitesBookEmpty(t *testing.T) {
	for _, s := range Sitepool {
		if s.Name == "" {
			t.Fatalf("Site Name cannot be empty")
		}
		if s.HomePage == "" {
			t.Fatalf("Site HomePage cannot be empty")
		}
		if s.BookInfo == nil {
			t.Fatalf("%s(%s) BookInfo cannot be empty", s.Name, s.HomePage)
		}
		if s.Chapter == nil {
			t.Fatalf("%s(%s) Chapter cannot be empty", s.Name, s.HomePage)
		}
	}
}

func TestAlreadyExistName(t *testing.T) {
	dd := map[string][]SiteA{}
	for _, s := range Sitepool {
		dd[s.Name] = append(dd[s.Name], s)
	}
	for name, d := range dd {
		if len(d) > 1 {
			s := ""
			for _, v := range d {
				s += v.File + " "
			}
			t.Logf("already exist Name: %#v %d %s", name, len(d), s)
			t.Fail()
		}
	}
}
func TestAlreadyExistHomePage(t *testing.T) {
	dd := map[string][]SiteA{}
	for _, s := range Sitepool {
		dd[s.HomePage] = append(dd[s.HomePage], s)
	}
	for name, d := range dd {
		if len(d) > 1 {
			s := ""
			for _, v := range d {
				s += v.File + " "
			}
			t.Logf("already exist HomePage: %#v %d %s", name, len(d), s)
			t.Fail()
		}
	}
}

func GenBookInfoSite(s SiteA) func(t *testing.T) {
	type TestCase struct {
		Name       string
		Author     string
		ChapterNum int               // 章节数目
		URL        map[string]string // Key:SiteName Value:BookURL
	}
	CheckFunc := func(tc TestCase, st *store.Store) (err error) {
		if tc.Name != st.BookName {
			return fmt.Errorf("BookName check does not match")
		}
		if tc.Author != st.Author {
			return fmt.Errorf("Author check does not match")
		}
		chapterNum := 0
		for _, vol := range st.Volumes {
			chapterNum += len(vol.Chapters)
		}
		if chapterNum == 0 {
			return fmt.Errorf("ChapterNum check does not match")
		}
		return
	}
	dd := []TestCase{
		{
			Name:   "放开那个女巫",
			Author: "二目",
			URL: map[string]string{
				"书迷楼":    "http://www.shumil.co/fangkainagenvwu/",
				"完本神站":   "https://www.wanbentxt.com/1949/",
				"新八一中文网": "https://www.81new.net/11/11609/",
				// "起点中文网": "https://book.qidian.com/info/1003306811",
			},
		},
	}

	return func(t *testing.T) {
		t.Parallel()
		var (
			err error
			r   io.Reader
			st  *store.Store
		)
		for _, d := range dd {
			u, ok := d.URL[s.Name]
			if !ok {
				continue
			}
			if r, err = utils.GetWebPageBodyReader(u); err != nil {
				t.Logf("%s", err)
				t.Fail()
				continue
			}
			if st, err = s.BookInfo(r); err != nil {
				t.Logf("%s", err)
				t.Fail()
				continue
			}
			if err = CheckFunc(d, st); err != nil {
				t.Logf("%s", err)
				t.Fail()
				continue
			}
		}
	}
}

func GenSearchSite(s SiteA) func(t *testing.T) {
	dd := []struct {
		Name   string
		Author string
	}{
		{"诡秘之主", "爱潜水的乌贼"},
		{"黎明之剑", "远瞳"},
		{"绿龙筑巢记", "归兮北冥"},
		{"底栖魔鱼日记", "辣鸡葱花"},
		{"俺，龙领主", "熊瀚"},
		{"红龙", "接口卡"},
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
			result, err := s.Search(b.Name)
			if err != nil {
				t.Fatalf("%s %s %s", s.Name, s.HomePage, err)
			}
			for _, r := range result {
				if (r.Author == b.Author) && (r.BookName == b.Name) {
					isOK = true
					continue
				}
			}
			s := ""
			for _, v := range result {
				s += fmt.Sprintf("[%s(%s)](%#v) ", v.BookName, v.Author, v.BookURL)
			}
			t.Logf("%s(%s) %d %s", b.Name, b.Author, len(result), s)
		}
		if !isOK {
			t.Logf("%s %s %s", s.Name, s.HomePage, "搜索结果无效")
			t.Fail()
			return
		}
		t.Run(fmt.Sprintf("BookInfo"), GenBookInfoSite(s))
	}
}

func TestSearch(t *testing.T) {
	for _, s := range Sitepool {
		if s.Search == nil {
			continue
		}
		t.Run(fmt.Sprintf("%s", s.File), GenSearchSite(s))
	}
}
