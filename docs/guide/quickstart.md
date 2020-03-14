# 快速上手

## 直接爬取盗版源

```bash
NAME:
    download - 下载缓存文件

USAGE:
    download [command options] [arguments...]

OPTIONS:
   -t value              线程数 (default: 10)
   -f value              输出格式
   -o value              输出路径
   --chromedp-log value  Chromedp log file
   --tsleep value        章节爬取间隔 (default: 200ms)
   --errsleep value      章节爬取错误间隔 (default: 500ms)
```

在不指定输出路径时，默认为`{书名}-{作者}-{站点}.FictionDown`

例如如下的命令将输出文件为`奥术神座-爱潜水的乌贼-书迷楼.FictionDown`

```bash
╰─$ FictionDown --url http://www.shumil.co/aoshushenzuo/ d
2020/03/15 00:42:59 URL: "http://www.shumil.co/aoshushenzuo/"
2020/03/15 00:42:59 use golang default http
2020/03/15 00:43:00 Loading....
书名: "奥术神座"
作者: "爱潜水的乌贼"
封面:
简介:

章节数:
        正文卷(免费) 1082章
2020/03/15 00:43:00 线程数: 10,预缓存中...
2020/03/15 00:43:00 [读入] 已缓存:45 样本:0 完成样本:0
 1082 / 1082 [=====================================================================================================================================================================] 100.00% 3m4s
2020/03/15 00:46:05 缓存完成
2020/03/15 00:46:05 生成别名
2020/03/15 00:46:06 [爬取结束] 已缓存:1082 样本:0 完成样本:0
```

## 爬取正版源再补充盗版源

```bash
FictionDown --url https://book.qidian.com/info/1009968948 d
FictionDown -i 绝对交易-隐语者-起点中文网.FictionDown s -k 绝对交易 -p
FictionDown -i 绝对交易-隐语者-起点中文网.FictionDown d
```

## 转换格式输出

Usage 如下：

```bash
NAME:
    convert - 转换格式输出

USAGE:
    convert [command options] [arguments...]

OPTIONS:
   -f value            输出格式
   -o value            输出路径
   --ignore_cover      忽略封面
   --no-EPUB-metadata  禁用EPUB元数据
```

### 输出 TXT

```bash
FictionDown -i 绝对交易-隐语者-起点中文网.FictionDown conv -f txt
```

### 直接输出 EPUB

```bash
FictionDown -i 绝对交易-隐语者-起点中文网.FictionDown conv -f epub
```

## 搜索书籍

```bash
AME:
    search - 检索盗版站点

USAGE:
    search [command options] [arguments...]

OPTIONS:
   -k value, --keyword value  搜索关键词
   --put, -p                  对比并放入缓存文件
   --save value, -s value     searh and save to file 搜索结果存储
   --format value, -f value   save format support json,yaml 存储格式
```

### 搜索结果保存为文件

支持 json 和 yaml 格式，当只指定文件名而不指定格式时，将会识别文件扩展名作为格式

```bash
╰─$ FictionDown s -k "诡秘" -s a.yml
```

```bash
╰─$ FictionDown s -k "诡秘" -s a.json
```

### 输出在终端里

```bash
╰─$ FictionDown s -k "诡秘"
2020/03/15 01:09:53 开始搜索站点: 起点中文网 https://www.qidian.com/
2020/03/15 01:09:53 开始搜索站点: 完本神站 https://www.wanbentxt.com/
2020/03/15 01:09:53 开始搜索站点: 书迷楼 http://www.shumil.co/
2020/03/15 01:09:53 开始搜索站点: 新八一中文网 https://www.81new.net/
2020/03/15 01:09:53 开始搜索站点: 顶点小说 https://www.booktxt.net/
2020/03/15 01:09:53 开始搜索站点: 笔趣阁1 https://www.biquge5200.cc/
2020/03/15 01:09:53 搜索站点: 结果: 7 书迷楼 http://www.shumil.co/
2020/03/15 01:09:53 搜索站点: 结果: 10 起点中文网 https://www.qidian.com/
2020/03/15 01:09:54 搜索站点: 结果: 17 笔趣阁1 https://www.biquge5200.cc/
2020/03/15 01:09:54 搜索站点: 结果: 10 顶点小说 https://www.booktxt.net/
2020/03/15 01:09:54 搜索站点: 结果: 5 完本神站 https://www.wanbentxt.com/
2020/03/15 01:09:55 搜索站点: 结果: 24 新八一中文网 https://www.81new.net/
搜索到41个内容:
书名: 诡秘世界之旅 作者: 梦里几度寒秋 6个书源
        "书迷楼" : "http://www.shumil.co/guimishijiezhilv/"
        "起点中文网" : "https://book.qidian.com/info/1015609520"
        "笔趣阁1" : "https://www.biquge5200.cc/123_123176/"
        "顶点小说" : "http://www.booktxt.net/book/goto/id/16070"
        "完本神站" : "https://www.wanbentxt.com/16638/"
        "新八一中文网" : "https://www.81new.net/131/131496/"
书名: 诡秘之主 作者: 爱潜水的乌贼 6个书源
        "书迷楼" : "http://www.shumil.co/guimizhizhu/"
        "起点中文网" : "https://book.qidian.com/info/1010868264"
        "笔趣阁1" : "https://www.biquge5200.cc/94_94525/"
        "顶点小说" : "http://www.booktxt.net/book/goto/id/5552"
        "完本神站" : "https://www.wanbentxt.com/853/"
        "新八一中文网" : "https://www.81new.net/34/34569/"
书名: 诡秘神探 作者: 残剑 5个书源
        "书迷楼" : "http://www.shumil.co/guimishentan/"
        "笔趣阁1" : "https://www.biquge5200.cc/126_126408/"
        "顶点小说" : "http://www.booktxt.net/book/goto/id/18142"
        "完本神站" : "https://www.wanbentxt.com/20149/"
        "新八一中文网" : "https://www.81new.net/139/139778/"
书名: 诡秘无限 作者: 柒月海岸 3个书源
        "起点中文网" : "https://book.qidian.com/info/1016934876"
        "笔趣阁1" : "https://www.biquge5200.cc/130_130271/"
        "顶点小说" : "http://www.booktxt.net/book/goto/id/19660"
书名: 诡秘三千藏 作者: 八黎 3个书源
        "书迷楼" : "http://www.shumil.co/guimisanqiancang/"
        "完本神站" : "https://www.wanbentxt.com/8481/"
        "新八一中文网" : "https://www.81new.net/105/105720/"
书名: 诡秘力量 作者: 潇湘夫子 3个书源
        "起点中文网" : "https://book.qidian.com/info/1016911451"
        "笔趣阁1" : "https://www.biquge5200.cc/127_127572/"
        "新八一中文网" : "https://www.81new.net/143/143310/"
书名: 光掩黑色：诡秘女探 作者: 青丝染霜 3个书源
        "书迷楼" : "http://www.shumil.co/guangyanheiseguiminvtan/"
        "笔趣阁1" : "https://www.biquge5200.cc/103_103830/"
        "新八一中文网" : "https://www.81new.net/99/99938/"
书名: 诡秘邪典 作者: 花生和尚 3个书源
        "起点中文网" : "https://book.qidian.com/info/1017679642"
        "笔趣阁1" : "https://www.biquge5200.cc/131_131743/"
        "新八一中文网" : "https://www.81new.net/147/147554/"
书名: 诡秘妖异之变 作者: 星点烽火 3个书源
        "书迷楼" : "http://www.shumil.co/guimiyaoyizhibian/"
        "笔趣阁1" : "https://www.biquge5200.cc/97_97317/"
        "新八一中文网" : "https://www.81new.net/107/107732/"
书名: 诡秘魔术师 作者: 浮世黄粱 2个书源
        "完本神站" : "https://www.wanbentxt.com/5759/"
        "新八一中文网" : "https://www.81new.net/113/113779/"
书名: 诡秘档案之追凶 作者: 隐蔽者 2个书源
        "笔趣阁1" : "https://www.biquge5200.cc/112_112741/"
        "新八一中文网" : "https://www.81new.net/111/111369/"
书名: 诡秘之王 作者: 妖异君 2个书源
        "起点中文网" : "https://book.qidian.com/info/1015636402"
        "顶点小说" : "http://www.booktxt.net/book/goto/id/14077"
书名: 诡秘邮件 作者: 一口老仙 2个书源
        "起点中文网" : "https://book.qidian.com/info/1016716524"
        "顶点小说" : "http://www.booktxt.net/book/goto/id/19984"
书名: 诡秘之梦 作者: 洒家随风 2个书源
        "笔趣阁1" : "https://www.biquge5200.cc/128_128260/"
        "顶点小说" : "http://www.booktxt.net/book/goto/id/19709"
书名: 一个诡秘作家的自我修养 作者: 月宫陈树 2个书源
        "笔趣阁1" : "https://www.biquge5200.cc/129_129010/"
        "顶点小说" : "http://www.booktxt.net/book/goto/id/20547"
书名: 夜不语诡秘档案 作者: 夜不语 1个书源
        "新八一中文网" : "https://www.81new.net/58/58299/"
书名: 诡秘狂欢 作者: 反派驾到 1个书源
        "笔趣阁1" : "https://www.biquge5200.cc/110_110484/"
书名: 诡秘的禁地 作者: 黎钥 1个书源
        "笔趣阁1" : "https://www.biquge5200.cc/90_90335/"
书名: 灵魂诡秘 作者: 越月约 1个书源
        "新八一中文网" : "https://www.81new.net/127/127170/"
书名: 灵岛诡秘 作者: 落花蝶 1个书源
        "新八一中文网" : "https://www.81new.net/57/57188/"
书名: 诡秘妖主 作者: 槐风眠 1个书源
        "新八一中文网" : "https://www.81new.net/144/144864/"
书名: 诡秘游戏 作者: 九三曰 1个书源
        "笔趣阁1" : "https://www.biquge5200.cc/123_123188/"
书名: 诡秘侦探所 作者: 君公子墨 1个书源
        "起点中文网" : "https://book.qidian.com/info/1019622326"
书名: 诡秘探索 作者: 低等水平 1个书源
        "新八一中文网" : "https://www.81new.net/135/135594/"
书名: 诡秘事件簿 作者: 漠然旅者 1个书源
        "新八一中文网" : "https://www.81new.net/107/107664/"
书名: 诡秘来袭 作者: 六小卿 1个书源
        "起点中文网" : "https://book.qidian.com/info/1017261930"
书名: 道士诡秘手札 作者: 枫叶恋秋落 1个书源
        "新八一中文网" : "https://www.81new.net/59/59174/"
书名: 诡秘怪谈 作者: 神鬼测 1个书源
        "新八一中文网" : "https://www.81new.net/142/142779/"
书名: 仙湖诡秘录 作者: 洛珩书 1个书源
        "新八一中文网" : "https://www.81new.net/72/72591/"
书名: 诡秘迷宫 作者: 磐生莲儿 1个书源
        "新八一中文网" : "https://www.81new.net/109/109610/"
书名: 诡秘的旅程 作者: 今天吃鱼 1个书源
        "起点中文网" : "https://book.qidian.com/info/1017733082"
书名: 创造诡秘世界 作者: 奶茶会挖煤 1个书源
        "笔趣阁1" : "https://www.biquge5200.cc/130_130587/"
书名: 诡秘18月 作者: 磐生莲儿 1个书源
        "新八一中文网" : "https://www.81new.net/99/99820/"
书名: 万佛诡秘 作者: 梦无命 1个书源
        "新八一中文网" : "https://www.81new.net/77/77826/"
书名: 诡秘盒子 作者: 蛋煎平底锅 1个书源
        "顶点小说" : "http://www.booktxt.net/book/goto/id/16639"
书名: 诡秘昆仑 作者: 夏二公子 1个书源
        "新八一中文网" : "https://www.81new.net/142/142010/"
书名: 捉诡秘记 作者: 天涯无月 1个书源
        "新八一中文网" : "https://www.81new.net/28/28810/"
书名: 诡秘管理员 作者: 陈知道 1个书源
        "顶点小说" : "http://www.booktxt.net/book/goto/id/20919"
书名: 诡秘之主笔趣 作者: 爱潜水的乌贼 1个书源
        "书迷楼" : "http://www.shumil.co/guimizhizhubiqu/"
书名: 能量进化 作者: 诡秘OL 1个书源
        "笔趣阁1" : "https://www.biquge5200.cc/64_64476/"
书名: 时空穿梭的诡秘者 作者: 秀屿熙 1个书源
        "笔趣阁1" : "https://www.biquge5200.cc/49_49766/"
```
