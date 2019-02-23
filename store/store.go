/*Package store is download cache file struct
 */
package store

import (
	"time"
)

// FileExt is filename extension (without dot)
const FileExt = "FictionDown"

// Store is store yaml data file format
type Store struct {
	BiqugeURL   string
	BookName    string
	CoverURL    string    // 封面链接
	Description string    // 介绍
	Author      string    // 作者
	LastUpdate  time.Time // 最后更新时间
	Volumes     []Volume
}

// Volume 卷
type Volume struct {
	Name     string
	Chapters []Chapter
}

// Chapter 章节
type Chapter struct {
	Name string
	URL  string
	Text []string
}
