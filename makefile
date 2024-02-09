VERSION := 0.1.0

build: clean version macos

clean:
	@echo "Cleaning up"
	rm -rf bin

version:
	@echo $(VERSION)
	mkdir -p bin/$(VERSION)

macos:
	@echo "Building for MacOS"
	mkdir -p bin/$(VERSION)/darwin/amd64
	GOOS=darwin GOARCH=amd64 go build -o bin/$(VERSION)/darwin/amd64/macbuilder