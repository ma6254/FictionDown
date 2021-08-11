package sites

import (
	"reflect"
	"runtime"

	"github.com/ma6254/FictionDown/site"
	"github.com/ma6254/FictionDown/sites/cc_b520"
	"github.com/ma6254/FictionDown/sites/co_shumil"
	"github.com/ma6254/FictionDown/sites/com_ddyueshu"
	"github.com/ma6254/FictionDown/sites/com_mijiashe"
	"github.com/ma6254/FictionDown/sites/com_qidian"
	"github.com/ma6254/FictionDown/sites/la_qb5"
	"github.com/ma6254/FictionDown/sites/net_81new"
	"github.com/ma6254/FictionDown/sites/org_wanben"
)

type siteFunc func() site.SiteA

func addSiteFunc(fn siteFunc) {
	s := fn()
	s.File = runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	site.AddSite(s)
}

func InitSites() {
	addSiteFunc(cc_b520.Site)
	addSiteFunc(co_shumil.Site)
	addSiteFunc(com_ddyueshu.Site)
	addSiteFunc(com_mijiashe.Site)
	addSiteFunc(com_qidian.Site)
	addSiteFunc(la_qb5.Site)
	addSiteFunc(net_81new.Site)
	addSiteFunc(org_wanben.Site)
}
