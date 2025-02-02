# Makefile
BIN_NAME = readTheRoom
VERSION = 1.0.0
BUILD_FLAGS = -trimpath -ldflags="-s -w"

.PHONY: all clean amd64 arm64

all: amd64 arm64

amd64:
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o bin/$(BIN_NAME)-$(VERSION)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o bin/$(BIN_NAME)-$(VERSION)-darwin-amd64 .
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o bin/$(BIN_NAME)-$(VERSION)-windows-amd64.exe .

arm64:
	GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o bin/$(BIN_NAME)-$(VERSION)-linux-arm64 .
	GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o bin/$(BIN_NAME)-$(VERSION)-darwin-arm64 .

clean:
	rm -rf bin/

init:
	mkdir -p bin