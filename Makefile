GOGEN=go generate

all: generate

generate:
	$(GOGEN) log/log.go


