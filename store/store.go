/*Package store is download cache file struct
 */
package store

// FileExt is filename extension (without dot)
const FileExt = "FictionDown"

// Store is store yaml data file format
type Store struct {
	BookURL     string
	BookName    string
	Author      string   // 作者
	CoverURL    string   // 封面链接
	Description string   // 介绍
	Tmap        []string //盗版源
	Volumes     []Volume
}

// Volume 卷
type Volume struct {
	Name     string
	IsVIP    bool
	Chapters []Chapter
}

// Chapter 章节
type Chapter struct {
	Name    string
	URL     string
	TURL    []string
	Text    []string
	Example []string
}
