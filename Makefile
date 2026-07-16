.PHONY: test build clean

test:
	go test ./...

build:
	mkdir -p bin
	go build -o bin/figma-cli ./cmd/figma-cli

clean:
	rm -rf bin dist coverage.out
