# FictionDown

FictionDown 是一个命令行界面的小说爬取工具

**用于批量下载盗版网络小说，该软件仅用于数据分析的样本采集，请勿用于其他用途**

**该软件所产生的文档请勿传播，请勿用于数据评估外的其他用途**

[![License](https://img.shields.io/github/license/ma6254/FictionDown.svg)](https://raw.githubusercontent.com/ma6254/FictionDown/master/LICENSE)
[![release_version](https://img.shields.io/github/release/ma6254/FictionDown.svg)](https://github.com/ma6254/FictionDown/releases)
[![last-commit](https://img.shields.io/github/last-commit/ma6254/FictionDown.svg)](https://github.com/ma6254/FictionDown/commits)
[![Download Count](https://img.shields.io/github/downloads/ma6254/FictionDown/total.svg)](https://github.com/ma6254/FictionDown/releases)
[![goproxy.cn](https://goproxy.cn/stats/github.com/ma6254/FictionDown/badges/download-count.svg)](https://goproxy.cn)

[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/ma6254/FictionDown/)
[![QQ 群](https://img.shields.io/badge/qq%E7%BE%A4-934873832-orange.svg)](https://jq.qq.com/?_wv=1027&k=5bN0SVA)

[![Go](https://github.com/ma6254/FictionDown/workflows/Go/badge.svg)](https://github.com/ma6254/FictionDown/actions/runs/39839114)
[![travis-ci](https://www.travis-ci.org/ma6254/FictionDown.svg?branch=master)](https://travis-ci.org/ma6254/FictionDown)
[![Go Report Card](https://goreportcard.com/badge/github.com/ma6254/FictionDown)](https://goreportcard.com/report/github.com/ma6254/FictionDown)

## 文档

文档目前「指南」部分已完成，你可以在[这里](https://ma6254.github.io/FictionDown/)查看。

## 特性

- 以起点为样本，多站点多线程爬取校对
- 支持导出 txt，以兼容大多数阅读器
- 支持导出 epub(还有些问题，某些阅读器无法打开)
- 支持导出 markdown，可以用 pandoc 转换成 epub，附带 epub 的`metadata`，保留书本信息、卷结构、作者信息
- 内置简单的广告过滤（现在还不完善）
- 用 Golang 编写，安装部署方便，可选的外部依赖：Chromedp
- 支持断点续爬，强制结束再爬会在上次结束的地方继续

## 站点支持

- 是否正版：✅ 为正版站点 ❌ 为盗版站点
- 是否分卷：✅ 章节分卷 ❌ 所有章节放在一个卷中不分卷
- 站内搜索：✅ 完全支持 ❌ 不支持 ❔ 站点支持但软件未适配 ⚠️ 站点支持，但不可用或维护中 ⛔ 站点支持搜索，但没有好的适配方案（比如用 Google 做站内搜索）

| 站点名称     | 网址              | 是否正版 | 是否分卷 | 支持站内搜索 | 代码文件                       |
| ------------ | ----------------- | -------- | -------- | ------------ | ------------------------------ |
| 起点中文网   | www.qidian.com    | ✅       | ✅       | ✅           | sites\qidian\main.go           |
| 笔趣阁       | www.biquge5200.cc | ❌       | ❌       | ✅           | sites\biquge5200_cc\main.go    |
| 顶点小说     | www.booktxt.net   | ❌       | ❌       | ✅           | sites\booktxt_net\main.go      |
| 新八一中文网 | www.81new.com     | ❌       | ❌       | ✅           | sites\new81\main.go            |
| 书迷楼       | www.shumil.co     | ❌       | ❌       | ✅           | sites\shumil_co\main.go        |
| 完本神站     | www.wanbentxt.com | ❌       | ❌       | ✅           | site\wanbentxt_com.go          |
| 38 看书      | www.38kanshu.com  | ❌       | ❌       | ⚠️           | sites\com_38kanshu\38kanshu.go |

## 使用注意

- 起点和盗版站的页面可能随时更改，可能会使抓取匹配失效，如果失效请提 issue
- 生成的 EPUB 文件可能过大，市面上大多数阅读器会异常卡顿或者直接崩溃
- 某些过于老的书或者作者频繁修改的书，盗版站都没有收录，也就无法爬取，如能找此书可用的盗版站请提 issue，并写出书名和正版站链接、盗版站链接

## 工作流程

1. 输入起点链接
2. 获取到书本信息，开始爬取每章内容，遇到 vip 章节放入`Example`中作为校对样本
3. 手动设置笔趣阁等盗版小说的对应链接，`tamp`字段
4. 再次启动，开始爬取，只爬取 VIP 部分，并跟`Example`进行校对
5. 手动编辑对应的缓存文件，手动删除广告和某些随机字符(有部分是关键字,可能会导致 pandoc 内存溢出或者样式错误)
6. `conv -f md`生成 markwown
7. 用 pandoc 转换成 epub，`pandoc -o xxxx.epub xxxx.md`

### Example

```bash
> ./FictionDown --url https://book.qidian.com/info/3249362 d # 获取正版信息

# 有时会发生`not match volumes`的错误，请启用Chromedp或者PhantomJS
# Use Chromedp
> ./FictionDown --url https://book.qidian.com/info/3249362 -d chromedp d
# Use PhantomJS
> ./FictionDown --url https://book.qidian.com/info/3249362 -d phantomjs d

> vim 一世之尊.FictionDown # 加入盗版小说链接
> ./FictionDown -i 一世之尊.FictionDown d # 获取盗版内容
# 爬取完毕就可以输出可阅读的文档了
> ./FictionDown -i 一世之尊.FictionDown conv -f txt
# 转换成epub有两种方式
# 1.输出markdown，再用pandoc转换成epub
> ./FictionDown -i 一世之尊.FictionDown conv -f md
> pandoc -o 一世之尊.epub 一世之尊.md
# 某些阅读器需要对章节进行定位,需要加上--epub-chapter-level=2
> pandoc -o 一世之尊.epub --epub-chapter-level=2 一世之尊.md
# 2.直接输出epub（调用Pandoc）
> ./FictionDown -i 一世之尊.FictionDown conv -f epub
```

#### 可直接根据搜索结果直接下载（当存在至少一个正版源时可用）

```bash
> ./FictionDown s -d -k "诡秘之主"
```

#### 站内搜索，然后填入

```bash
> ./FictionDown --url https://book.qidian.com/info/3249362 d # 获取正版信息

# 有时会发生`not match volumes`的错误，请启用Chromedp或者PhantomJS
# Use Chromedp
> ./FictionDown --url https://book.qidian.com/info/3249362 --driver chromedp d
# Use PhantomJS
> ./FictionDown --url https://book.qidian.com/info/3249362 --driver phantomjs d

> ./FictionDown -i 一世之尊.FictionDown s -k 一世之尊 -p # 搜索然后放入
> ./FictionDown -i 一世之尊.FictionDown d # 获取盗版内容
# 爬取完毕就可以输出可阅读的文档了
> ./FictionDown -i 一世之尊.FictionDown conv -f txt
# 转换成epub有两种方式
# 1.输出markdown，再用pandoc转换成epub
> ./FictionDown -i 一世之尊.FictionDown conv -f md
> pandoc -o 一世之尊.epub 一世之尊.md
# 2.直接输出epub（某些阅读器会报错）
> ./FictionDown -i 一世之尊.FictionDown conv -f epub
```

## 未实现

- 爬取正版的时候带上`Cookie`，用于爬取已购买章节
- 支持 晋江文学城
- 支持 纵横中文网
- 支持有毒小说网
- 支持刺猬猫（即“欢乐书客”）
- 整理 main 包中的面条逻辑
- 整理命令行参数风格
- 完善广告过滤
- 简化使用步骤
- 优化 log 输出
- 对于特殊章节，支持手动指定盗版链接或者跳过忽略
- 外部加载匹配规则，让用户可以自己添加正/盗版源
- 支持章节更新
- 章节匹配过程优化

## Usage

```bash
NAME:
   FictionDown - https://github.com/ma6254/FictionDown

USAGE:
    [global options] command [command options] [arguments...]

AUTHOR:
   ma6254 <9a6c5609806a@gmail.com>

COMMANDS:
     download, d, down  下载缓存文件
     check, c, chk      检查缓存文件
     edit, e            对缓存文件进行手动修改
     convert, conv      转换格式输出
     pirate, p          检索盗版站点
     search, s          检索盗版站点
     help, h            Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -u value, --url value     图书链接
   --tu value, --turl value  资源网站链接
   -i value, --input value   输入缓存文件
   --log value               log file path
   --driver value, -d value  请求方式,support: none,phantomjs,chromedp
   --help, -h                show help
   --version, -v             print the version
```

## 安装和编译

程序为单执行文件，命令行 CLI 界面

包管理为 gomod

```bash
go get github.com/ma6254/FictionDown
```

交叉编译这几个平台的可执行文件：`linux/arm` `linux/amd64` `darwin/amd64` `windows/amd64`

```bash
make multiple_build
```
