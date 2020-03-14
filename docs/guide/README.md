# 概览

FictionDown 是一个命令行界面的小说爬取工具

::: warning
用于批量下载盗版网络小说，该软件仅用于数据分析的样本采集，请勿用于其他用途
:::

::: warning
该软件所产生的文档请勿传播，请勿用于数据评估外的其他用途
:::

[![License](https://img.shields.io/github/license/ma6254/FictionDown.svg)](https://raw.githubusercontent.com/ma6254/FictionDown/master/LICENSE)
[![release_version](https://img.shields.io/github/release/ma6254/FictionDown.svg)](https://github.com/ma6254/FictionDown/releases)
[![last-commit](https://img.shields.io/github/last-commit/ma6254/FictionDown.svg)](https://github.com/ma6254/FictionDown/commits)
[![Download Count](https://img.shields.io/github/downloads/ma6254/FictionDown/total.svg)](https://github.com/ma6254/FictionDown/releases)

[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/ma6254/FictionDown/)
[![QQ 群](https://img.shields.io/badge/qq%E7%BE%A4-934873832-orange.svg)](https://jq.qq.com/?_wv=1027&k=5bN0SVA)

[![Go](https://github.com/ma6254/FictionDown/workflows/Go/badge.svg)](https://github.com/ma6254/FictionDown/actions/runs/39839114)
[![travis-ci](https://www.travis-ci.org/ma6254/FictionDown.svg?branch=master)](https://travis-ci.org/ma6254/FictionDown)
[![Go Report Card](https://goreportcard.com/badge/github.com/ma6254/FictionDown)](https://goreportcard.com/report/github.com/ma6254/FictionDown)

## 特性

- 以起点为样本，多站点多线程爬取校对
- 支持导出 txt，以兼容大多数阅读器
- 支持导出 epub(还有些问题，某些阅读器无法打开)
- 支持导出 markdown，可以用 pandoc 转换成 epub，附带 epub 的`metadata`，保留书本信息、卷结构、作者信息
- 内置简单的广告过滤（现在还不完善）
- 用 Golang 编写，安装部署方便，可选的外部依赖：PhantomJS、Chromedp
- 支持断点续爬，强制结束再爬会在上次结束的地方继续

## 站点支持

- 是否正版：✅ 为正版站点 ❌ 为盗版站点
- 是否分卷：✅ 章节分卷 ❌ 所有章节放在一个卷中不分卷
- 站内搜索：✅ 完全支持 ❌ 不支持 ❔ 站点支持但软件未适配 ⚠️ 站点支持，但不可用或维护中 ⛔ 站点支持搜索，但没有好的适配方案（比如用 Google 做站内搜索）

| 站点名称     | 网址              | 是否正版 | 是否分卷 | 支持站内搜索 | 代码文件              |
| ------------ | ----------------- | -------- | -------- | ------------ | --------------------- |
| 起点中文网   | www.qidian.com    | ✅       | ✅       | ✅           | site\qidian.go        |
| 笔趣阁       | www.biquge5200.cc | ❌       | ❌       | ✅           | site\biquge.go        |
| 笔趣阁 5200  | www.bqg5200.com   | ❌       | ❌       | ❔           | site\biquge2.go       |
| 笔趣阁       | www.biqiuge.com   | ❌       | ❌       | ⚠️           | site\biquge3.go       |
| 顶点小说     | www.booktxt.net   | ❌       | ❌       | ✅           | site\dingdian1.go     |
| 新八一中文网 | www.81new.com     | ❌       | ❌       | ✅           | site\81new.go         |
| 书迷楼       | www.shumil.co     | ❌       | ❌       | ✅           | site\shumil_co.go     |
| 完本神站     | www.wanbentxt.com | ❌       | ❌       | ✅           | site\wanbentxt_com.go |

## 使用注意

- 起点和盗版站的页面可能随时更改，可能会使抓取匹配失效，如果失效请提 issue
- 生成的 EPUB 文件可能过大，市面上大多数阅读器会异常卡顿或者直接崩溃
- 某些过于老的书或者作者频繁修改的书，盗版站都没有收录，也就无法爬取，如能找此书可用的盗版站请提 issue，并写出书名和正版站链接、盗版站链接
