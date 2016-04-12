PACKAGE=github.com/kiesel/wormhole-go/wormhole

install:
	go install $(PACKAGE)

dest:
	mkdir -p dest/windows dest/linux-amd64 dest/darwin

clean:
	rm -rf dest/

build: dest
	go build -o dest/wormhole $(PACKAGE)

build-all: build-windows build-darwin build-linux 

build-windows: dest
	GOOS=windows GOARCH=386 go build -v -o dest/windows/wormhole.exe $(PACKAGE)

build-darwin: dest
	GOOS=darwin GOARCH=amd64 go build -v -o dest/darwin/wormhole $(PACKAGE)

build-linux: dest
	GOOS=linux GOARCH=amd64 go build -v -o dest/linux-amd64/wormhole $(PACKAGE)