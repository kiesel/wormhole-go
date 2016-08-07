PACKAGE=github.com/kiesel/wormhole-go/wormhole

BUILD_FLAGS=-ldflags '-X main.VersionString=$(TRAVIS_TAG)'

.PHONY: install test clean build dist

default: test

fetch-dependencies:
	go get gopkg.in/op/go-logging.v1
	go get gopkg.in/yaml.v2
	go get -t ./...

unittest:
	go test -v ./...

test: unittest dist

install:
	go install $(BUILD_FLAGS) ./...

clean:
	rm -rf dist/

build: build-windows build-darwin build-linux
	cp README.md dist/wormhole/
	cp wormhole.yml dist/wormhole/.wormhole.yml

dist: build
	cd dist/ && zip -r wormhole-${TRAVIS_TAG}.zip wormhole/

build-windows:
	GOOS=windows GOARCH=386 go build $(BUILD_FLAGS) -o dist/wormhole/windows_386/wormhole.exe $(PACKAGE)

build-darwin:
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o dist/wormhole/darwin_amd64/wormhole $(PACKAGE)

build-linux:
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o dist/wormhole/linux_amd64/wormhole $(PACKAGE)