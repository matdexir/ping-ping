GO ?= go
EXECUTABLE := ping-ping
GOFILES := $(shell find . -type f -name "*.go")

build: $(EXECUTABLE)

$(EXECUTABLE): $(GOFILES)
	$(GO) build -v -o bin/ ./...

test: 
	@$(GO) test -v ./... && echo "\n==>\033[32m Ok\033[m\n" || exit 1

