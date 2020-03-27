# 添加自定义书源

::: tip
目前不支持从网络或者从配置文件加载书源，可通过修改源代码来添加书源
:::

一个书源由以下几部分构成：

- `BookInfo`书籍信息匹配：在书籍信息页面中的获取书名、作者名、章节目录以及对应章节的页面链接,有需要的话还有封面图片链接和简介
- `Chapter`小说章节匹配：在章节页面得到每个段落的内容
- `Search`搜索结果匹配：获取该站站内搜索结果的每条结果，以及结果内的书名、书籍信息链接、作者名
- `Tag`书源标签：可能的话，还需要填写书源的`Tag`，例如：是否正版、是否带分卷信息、是否是优质书源

## 通过修改源代码添加书源

### 添加包

定义一个 Package，示例：<https://github.com/ma6254/FictionDown/blob/master/sites/shumil_co/main.go>

### 加入到导入列表中

<https://github.com/ma6254/FictionDown/blob/master/sites/imports.go>

<<< @/../sites/imports.go

## 通过配置文件添加书源

::: tip
未支持，可在 Issue：<https://github.com/ma6254/FictionDown/issues/9> 中讨论相关方案
:::
