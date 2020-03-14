package site

// TagFunc return tag strrings
type TagFunc func() []string

//AddTag a add a tag end
func AddTag(fn TagFunc, s ...string) TagFunc {
	return func() []string {
		if fn == nil {
			return s
		}
		return append(fn(), s...)
	}
}
