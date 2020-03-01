BUILD_TIME := $(shell date "+%F %T")
COMMIT_ID := $(shell git rev-parse HEAD)
BUILD_VERSION := $(shell git describe --tags $(COMMIT_ID))

Package := "github.com/ma6254/FictionDown"

LDFLAG := "\
-s -w \
-X 'main.version=${BUILD_VERSION}' \
-X 'main.commit=${COMMIT_ID}' \
-X 'main.date=${BUILD_TIME}' \
-X 'main.builBy=${git config user.name}' \
"

build_tool:
	go get -v -u github.com/mitchellh/gox
	go mod vendor

build:
	go build -v --ldflags $(LDFLAG) $(Package)

multiple_build:
	gox -osarch="linux/arm" -osarch="linux/amd64" --osarch="darwin/amd64" -osarch="windows/amd64" -ldflags $(LDFLAG) -output "{{.Dir}}_$(BUILD_VERSION)_{{.OS}}_{{.Arch}}" $(Package)

install:
	go install --ldflags $(LDFLAG) $(Package)