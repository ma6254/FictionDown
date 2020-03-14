# 编译安装

以下为通过编译的方式安装 FictionDown 的流程，与常规`Golang`语言项目编译流程一致

## 直接 Go Get

```bash
go get -v github.com/ma6254/FictionDown@latest
```

## 可能会遇到各种问题

以下是可能的解决方法，如仍旧无法解决请提 issue

### 提示无某函数或某函数参数问题

可能是没有启用 gomod

```bash
go env -w GO111MODULE=on
```

### 网络错误

设置 goproxy 即可

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

## Clone 源码编译

```bash
git clone https://github.com/ma6254/FictionDown.git
cd FictionDown
```

然后就可以编译了

### 编译并安装到 GOPATH 中

```bash
go install -v .
```

### 在当前目录下生成可执行文件

```bash
go build -v .
```

### 生成多平台的可执行文件

先安装`goreleaser` <https://goreleaser.com/install/> 再按照如下执行即可

```bash
goreleaser release --skip-publish --skip-validate --rm-dist
```

```bash
╰─$ goreleaser release --skip-publish --skip-validate --rm-dist

   • releasing using goreleaser dev...
   • loading config file       file=.goreleaser.yml
   • RUNNING BEFORE HOOKS
      • running go mod tidy
      • running go generate ./...
   • LOADING ENVIRONMENT VARIABLES
      • pipe skipped              error=publishing is disabled
   • GETTING AND VALIDATING GIT STATE
      • releasing v0.1.3, commit 7b08832377af64b1e7b54c2ed751b6274a94870d
      • pipe skipped              error=validation is disabled
   • PARSING TAG
   • SETTING DEFAULTS
      • LOADING ENVIRONMENT VARIABLES
      • SNAPSHOTING
      • GITHUB/GITLAB/GITEA RELEASES
      • PROJECT NAME
      • BUILDING BINARIES
      • ARCHIVES
      • LINUX PACKAGES WITH NFPM
      • SNAPCRAFT PACKAGES
      • CALCULATING CHECKSUMS
      • SIGNING ARTIFACTS
      • DOCKER IMAGES
      • ARTIFACTORY
      • BLOB
      • HOMEBREW TAP FORMULA
      • SCOOP MANIFEST
   • SNAPSHOTING
      • pipe skipped              error=not a snapshot
   • CHECKING ./DIST
   • WRITING EFFECTIVE CONFIG FILE
      • writing                   config=dist/config.yaml
   • GENERATING CHANGELOG
      • writing                   changelog=dist/CHANGELOG.md
   • BUILDING BINARIES
      • building                  binary=/mnt/c/Users/mjc/git/FictionDown/dist/FictionDown_windows_386/FictionDown.exe
      • building                  binary=/mnt/c/Users/mjc/git/FictionDown/dist/FictionDown_linux_arm_7/FictionDown
      • building                  binary=/mnt/c/Users/mjc/git/FictionDown/dist/FictionDown_linux_amd64/FictionDown
      • building                  binary=/mnt/c/Users/mjc/git/FictionDown/dist/FictionDown_darwin_amd64/FictionDown
      • building                  binary=/mnt/c/Users/mjc/git/FictionDown/dist/FictionDown_windows_amd64/FictionDown.exe
      • building                  binary=/mnt/c/Users/mjc/git/FictionDown/dist/FictionDown_linux_arm64/FictionDown
      • building                  binary=/mnt/c/Users/mjc/git/FictionDown/dist/FictionDown_linux_386/FictionDown
      • building                  binary=/mnt/c/Users/mjc/git/FictionDown/dist/FictionDown_linux_arm_6/FictionDown
   • ARCHIVES
      • creating                  archive=dist/FictionDown_0.1.3_Windows_i386.zip
      • creating                  archive=dist/FictionDown_0.1.3_Linux_armv6.tar.gz
      • creating                  archive=dist/FictionDown_0.1.3_Windows_x86_64.zip
      • creating                  archive=dist/FictionDown_0.1.3_Linux_arm64.tar.gz
      • creating                  archive=dist/FictionDown_0.1.3_Linux_x86_64.tar.gz
      • creating                  archive=dist/FictionDown_0.1.3_Linux_i386.tar.gz
      • creating                  archive=dist/FictionDown_0.1.3_Linux_armv7.tar.gz
      • creating                  archive=dist/FictionDown_0.1.3_Darwin_x86_64.tar.gz
   • LINUX PACKAGES WITH NFPM
   • SNAPCRAFT PACKAGES
   • CALCULATING CHECKSUMS
      • checksumming              file=FictionDown_0.1.3_Darwin_x86_64.tar.gz
      • checksumming              file=FictionDown_0.1.3_Windows_x86_64.zip
      • checksumming              file=FictionDown_0.1.3_Linux_x86_64.tar.gz
      • checksumming              file=FictionDown_0.1.3_Linux_arm64.tar.gz
      • checksumming              file=FictionDown_0.1.3_Linux_i386.tar.gz
      • checksumming              file=FictionDown_0.1.3_Linux_armv6.tar.gz
      • checksumming              file=FictionDown_0.1.3_Windows_i386.zip
      • checksumming              file=FictionDown_0.1.3_Linux_armv7.tar.gz
   • SIGNING ARTIFACTS
      • pipe skipped              error=artifact signing is disabled
   • DOCKER IMAGES
      • pipe skipped              error=docker section is not configured
   • PUBLISHING
      • pipe skipped              error=publishing is disabled
   • release succeeded after 156.01s
```
