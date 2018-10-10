GOBUILD=go build
GOGEN=go generate

all: generate build

generate:
	$(GOGEN) debug/log.go

build: 4gt


4gt:
	$(GOBUILD) -o out/4gt main.go

clean:
	rm -rf out
