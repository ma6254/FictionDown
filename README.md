# FictionDown

**用于批量下载盗版网络小说，该软件仅用于数据分析的样本采集，请勿用于其他用途**

**该软件所产生的文档请勿传播，请勿用于数据评估外的其他用途**

[![License](https://img.shields.io/github/license/ma6254/FictionDown.svg)](https://raw.githubusercontent.com/ma6254/FictionDown/master/LICENSE)[![release_version](https://img.shields.io/github/release/ma6254/FictionDown.svg)](https://github.com/ma6254/FictionDown/releases)[![last-commit](https://img.shields.io/github/last-commit/ma6254/FictionDown.svg)](https://github.com/ma6254/FictionDown/commits)[![Download Count](https://img.shields.io/github/downloads/ma6254/FictionDown/total.svg)](https://github.com/ma6254/FictionDown/releases)

[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/ma6254/FictionDown/)[![QQ 群](https://img.shields.io/badge/qq%E7%BE%A4-934873832-orange.svg)](https://jq.qq.com/?_wv=1027&k=5bN0SVA)

[![travis-ci](https://www.travis-ci.org/ma6254/FictionDown.svg?branch=master)](https://travis-ci.org/ma6254/FictionDown)[![Go Report Card](https://goreportcard.com/badge/github.com/ma6254/FictionDown)](https://goreportcard.com/report/github.com/ma6254/FictionDown)

## 特性

- 以起点为样本，多站点多线程爬取校对
- 支持导出txt，以兼容大多数阅读器
- 支持导出epub(还有些问题，某些阅读器无法打开)
- 支持导出markdown，可以用pandoc转换成epub，附带epub的`metadata`，保留书本信息、卷结构、作者信息
- 内置简单的广告过滤（现在还不完善）
- 用Golang编写，安装部署方便，可选的外部依赖：PhantomJS、Chromedp
- 支持断点续爬，强制结束再爬会在上次结束的地方继续

## 使用注意

- 起点和盗版站的页面可能随时更改，可能会使抓取匹配失效，如果失效请提issue
- 生成的EPUB文件可能过大，市面上大多数阅读器会异常卡顿或者直接崩溃
- 某些过于老的书或者作者频繁修改的书，盗版站都没有收录，也就无法爬取，如能找此书可用的盗版站请提issue，并写出书名和正版站链接、盗版站链接

## 使用流程

1. 输入起点链接
2. 获取到书本信息，开始爬取每章内容，遇到vip章节放入`Example`中作为校对样本
3. 手动设置笔趣阁等盗版小说的对应链接，`tamp`字段
4. 再次启动，开始爬取，只爬取VIP部分，并跟`Example`进行校对
5. 手动编辑对应的缓存文件，手动删除广告和某些随机字符(有部分是关键字,可能会导致pandoc内存溢出或者样式错误)
6. `conv -f md`生成markwown
7. 用pandoc转换成epub，`pandoc -o xxxx.epub xxxx.md`

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

#### 现在支持小说站内搜索，可以不用手动填入了

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

- 爬取起点的时候带上`Cookie`，用于爬取已购买章节
- 支持 晋江文学城
- 支持 纵横中文网
- 支持刺猬猫（即“欢乐书客”）
- ~~支持小说站内搜索~~
- 整理main包中的面条逻辑
- 整理命令行参数风格
- 完善广告过滤
- 简化使用步骤
- 优化log输出
- 对于特殊章节，支持手动指定盗版链接或者跳过忽略
- 外部加载匹配规则，让用户可以自己添加正/盗版源
- 支持章节更新
- 章节匹配过程优化

## 编译

包管理采用godep

1. `dep ensure -v`
2. `make` or `make build` 当前目录下就会产生可执行文件

### 交叉编译

需要安装gox

`make multiple_build`

## 某些匹配问题

小说《一世之尊》

卷: "第三卷 满堂花醉三千客" 章节: "第十四章 姚家小鬼"

在盗版站的章节名均为`"第14章 姚家小鬼"`



## 支持的盗版站点

随机挑选了几个

- www.biqiuge.com
- www.biquge5200.cc
- www.bqg5200.com
- www.booktxt.net
- www.81new.com
- www.shumil.co
- www.wanbentxt.com