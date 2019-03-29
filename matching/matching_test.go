package matching

import (
	"log"
	"testing"
)

func TestTitleAlias(t *testing.T) {
	type Example struct {
		Src string
		Dst []string
	}
	var (
		examples = []Example{
			{"第一章 神奇的时代",
				[]string{
					"第一章",
					"神奇的时代",
				}},
			{"第1章 神奇的时代",
				[]string{
					"第1章",
					"神奇的时代",
				}},
			{"第二百二十四章 那个男人（6万月票加更）",
				[]string{
					"第二百二十四章",
					"那个男人",
				}},
			{"第二百二十四章 那个男人(6万月票加更)",
				[]string{
					"第二百二十四章", "那个男人",
				}},
			{"第二百二十四章 那个男人 (6万月票加更)",
				[]string{
					"第二百二十四章", "那个男人",
				}},
			{"462 飞升", []string{
				"462", "飞升",
			}},
			{"465 宇宙如卵，大千如池（END）",
				[]string{
					"465",
					"宇宙如卵，大千如池",
					"END",
				}},
			{"又是周一，拜求点击推荐（章节已更）",
				[]string{
					"又是周一，拜求点击推荐",
					"章节已更",
				}},
			{"咳，一个报（e）告（hao）",
				[]string{}},
			{"1，雾都孤儿",
				[]string{
					"1，",
					"雾都孤儿",
				}},
			{"第67章 第四层 上",
				[]string{}},
		}
	)
	for _, v := range examples {
		ret := TitleAlias(v.Src)
		log.Printf("alias: %#v", ret)
		for _, v := range v.Dst {
			if !StringInSlice(v, ret) {
				t.Fatalf("want %#v in %#v", v, ret)
			}
		}
	}
}
