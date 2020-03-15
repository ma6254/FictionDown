# 添加自定义书源

目前不支持从网络或者从配置文件加载书源，可通过修改源代码来添加书源

## 通过修改源代码添加书源

### 添加包

示例：<https://github.com/ma6254/FictionDown/blob/master/sites/shumil_co/main.go>

### 加入到导入列表中

<https://github.com/ma6254/FictionDown/blob/master/sites/imports.go>

<<< @/sites/imports.go{highlightLines}

## 通过配置文件添加书源

未支持，可在 Issue：<https://github.com/ma6254/FictionDown/issues/9> 中讨论相关方案
