package sites

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ma6254/FictionDown/site"
)

type SearchSite struct {
	Name   string
	Author string
}
type SiteMeta struct {
	Name     string
	HomePage string
}

var (
	searchSiteResult = make(map[SiteMeta]map[SearchSite][]string)
)

func init() {
	fmt.Println("Site init before testing")
	InitSites()
}

func TestSourceYaml(t *testing.T) {
	b, err := json.MarshalIndent(site.Sitepool, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", b)
}

func TestSitesBookEmpty(t *testing.T) {
	for _, s := range site.Sitepool {
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
	dd := map[string][]site.SiteA{}
	for _, s := range site.Sitepool {
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
	dd := map[string][]site.SiteA{}
	for _, s := range site.Sitepool {
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

func TestSearch(t *testing.T) {
	for _, s := range site.Sitepool {
		if s.Search == nil {
			continue
		}
		t.Run(fmt.Sprintf("%s", s.File), GenSearchSite(s))
	}
}
