GOBUILD=go build
GOGEN=go generate

GENERATED=trace/loglevel_string.go

.PHONY: all
all: test build

$(GENERATED): trace/log.go
	$(GOGEN) trace/log.go

build: out/4gtx out/4gt
out/4gtx: $(GENERATED)
	$(GOBUILD) -o out/4gtx cmd/4gtx/main.go

out/4gt: $(GENERATED)
	$(GOBUILD) -o out/4gt cmd/4gt/main.go

.PHONY: clean
clean:
	rm -rf out

.PHONY: test
test: $(GENERATED)
	go test ./db/... ./cmd/...

.PHONY: fmt
fmt:
	gofmt -s -w ./

install: test build
	mv out/4gtx $(HOME)/bin
	mv out/4gt $(HOME)/bin
