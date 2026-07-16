.PHONY: test build mock-verify release-check clean

test:
	go test ./...

build:
	mkdir -p bin
	go build -o bin/figma-cli ./cmd/figma-cli

mock-verify: build
	python3 scripts/mock_verify.py

release-check:
	go run github.com/goreleaser/goreleaser/v2@latest check

clean:
	rm -rf bin dist coverage.out
