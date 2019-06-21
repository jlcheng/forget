GOBUILD=go build
GOGEN=go generate

GENERATED=trace/loglevel_string.go

# Set to 1 because tests currently share the same GetTempIndexDir(). Parallel test
# execution will introduce unreproducible test failures due to races.
PARALLEL_TESTS=1


.PHONY: all
all: test build

$(GENERATED): trace/log.go
	$(GOGEN) trace/log.go

build: out/4gt
out/4gt: $(GENERATED)
	$(GOBUILD) -o out/4gt cmd/4gt/main.go

.PHONY: clean
clean:
	rm -rf out

.PHONY: test
test: $(GENERATED)
	go test -p $(PARALLEL_TESTS) ./app/... ./cmd/... ./db/...  ./orgmode/... ./watcher/...

.PHONY: fmt
fmt:
	gofmt -s -w ./
	golangci-lint run

install: test build
	mv out/4gt $(HOME)/bin

