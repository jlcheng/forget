GOBUILD=go build
GOGEN=go generate

all: generate build

generate: trace/loglevel_string.go
trace/loglevel_string.go:
	$(GOGEN) trace/log.go

build: out/4gt
out/4gt:
	$(GOBUILD) -o out/4gt main.go

.PHONY: clean
clean:
	rm -rf out

install: all
	mv out/4gt $(HOME)/bin
