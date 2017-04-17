COMMIT = $$(git describe --always)
PKG = github.com/yuuki/hworq
PKGS = $$(go list ./... | grep -v vendor)

all: build

dep:
	go get github.com/jteeuwen/go-bindata/...

.PHONY: build
build: assets
	go build -ldflags "-X main.GitCommit=\"$(COMMIT)\"" $(PKG)/cmd/...

.PHONY: test
test:
	go test -v $(PKGS)

.PHONY: test-race
test-race:
	go test -v -race $(PKGS)

.PHONY: test-all
test-all: vet test-race

.PHONY: dep assets
assets: data/schema/
	go-bindata -o pkg/data/bindata.go -pkg data data/schema/

.PHONY: cover
cover:
	go test -cover $(PKGS)

.PHONY: fmt
fmt:
	gofmt -s -w $$(git ls | grep -e '\.go$$' | grep -v -e vendor)

.PHONY: imports
imports:
	goimports -w $$(git ls | grep -e '\.go$$' | grep -v -e vendor)

.PHONY: lint
lint:
	golint $(PKGS)

.PHONY: vet
vet:
	go tool vet -all -printfuncs=Wrap,Wrapf,Errorf $$(find . -maxdepth 1 -mindepth 1 -type d | grep -v -e "^\.\/\." -e vendor)
