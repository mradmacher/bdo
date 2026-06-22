.PHONY: all windows linux macos clean

all: windows linux macos

assets:
	npm run build

windows: assets
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 \
	CC=x86_64-w64-mingw32-gcc \
	go build -ldflags="-linkmode external -extldflags '-static'" -o build/bdo.exe ./

linux: assets
	GOOS=linux GOARCH=amd64 go build -o build/bdo-linux-amd64 ./

macos: assets
	GOOS=darwin GOARCH=amd64 go build -o build/bdo-darwin-amd64 ./
	GOOS=darwin GOARCH=arm64 go build -o build/bdo-darwin-arm64 ./

clean:
	rm  build/*
