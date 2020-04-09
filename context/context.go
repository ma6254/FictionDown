package context

import (
	goctx "context"
)

type Context goctx.Context

type CtxKey string

var (
	KeyStore = CtxKey("store")
	KeyURL   = CtxKey("url")
	KeyBody  = CtxKey("body")
)
