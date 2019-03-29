package output

import (
	"fmt"
	"reflect"

	"github.com/ma6254/FictionDown/store"
)

func init() {
	RegOutputFormat("epub", &EPUB{})
	RegOutputFormat("md", &Markdown{})
	RegOutputFormat("txt", &TXT{})
}

// Option is Convert output options
type Option struct {
	IgnoreCover    bool // 忽略封面
	NoEPUBMetadata bool // 不添加EPUB元数据
}

type Conversion interface {
	Conv(src store.Store, outpath string, opts Option) error
}

var formatMap = map[string]Conversion{}

var (
	ErrUnsupportFormat = fmt.Errorf("Unsupport Conversion Format")
)

func RegOutputFormat(s string, conv Conversion) {
	formatMap[s] = conv
}

func Output(src store.Store, format string, outpath string, opts Option) (err error) {
	var c Conversion
	conver, ok := formatMap[format]
	if !ok {
		err = ErrUnsupportFormat
	}
	c = reflect.New(reflect.TypeOf(conver).Elem()).Interface().(Conversion)
	return c.Conv(src, outpath, opts)
}
