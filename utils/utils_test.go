package utils

import "testing"

func TestTupleSlice(t *testing.T) {
	var examples = []struct {
		Src []string
		Dst []string
	}{
		{[]string{"a", "b"}, []string{"a", "b"}},
		{[]string{"a", "b", "b"}, []string{"a", "b"}},
		{[]string{"a", "b", "b", "a"}, []string{"a", "b"}},
	}

	for k, v := range examples {
		a := TupleSlice(v.Src)
		b := v.Dst
		if !StringSliceEqual(a, b) {
			t.Fatalf("Fail %d A:%v B:%v", k, a, b)
		}
	}
	return
}
