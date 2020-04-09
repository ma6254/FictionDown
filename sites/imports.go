package sites

import (
	"reflect"
	"runtime"

	"github.com/ma6254/FictionDown/site"
	"github.com/ma6254/FictionDown/sites/biquge5200_cc"
	"github.com/ma6254/FictionDown/sites/booktxt_net"
	"github.com/ma6254/FictionDown/sites/com_38kanshu"
	"github.com/ma6254/FictionDown/sites/new81"
	"github.com/ma6254/FictionDown/sites/qidian"
	"github.com/ma6254/FictionDown/sites/shumil_co"
	"github.com/ma6254/FictionDown/sites/wanbentxt"
)

type siteFunc func() site.SiteA

func addSiteFunc(fn siteFunc) {
	s := fn()
	s.File = runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	site.AddSite(s)
}

func InitSites() {
	addSiteFunc(qidian.Site)
	addSiteFunc(wanbentxt.Site)
	addSiteFunc(shumil_co.Site)
	addSiteFunc(new81.Site)
	addSiteFunc(booktxt_net.Site)
	addSiteFunc(biquge5200_cc.Site)
	addSiteFunc(com_38kanshu.Site)
}
