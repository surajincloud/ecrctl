SOURCES := $(shell find . -name '*.go')
BINARY := ecrctl

build: ecrctl

clean:
	@rm -rf $(BINARY)

$(BINARY): $(SOURCES)
	CGO_ENABLED=0 go build -o $(BINARY) -ldflags="-s -w" main.go
