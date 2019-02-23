package output

import "github.com/ma6254/FictionDown/store"

func init() {
	// RegOutputFormat("txt", &TXT{})
}

type TXT struct {
	src *store.Store
}

func (t *TXT) NewConv(src *store.Store) error {
	t.src = src
	return nil
}

func (t TXT) Output() []byte {
	return nil
}
