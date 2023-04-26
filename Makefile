.PHONY: build clean run lint test fmt fix pretty vendor mocks

build:
	@go build -race -o downloader

clean:
	go clean -cache -testcache
	go mod tidy

ci: clean lint test

run: build
	./downloader

mocks:
	find . -iname '*_mock.go' -exec rm {} \;
	export PATH=$PATH:$(go env GOPATH)/bin
	go install github.com/golang/mock/mockgen@v1.6.0
	go generate -v ./...

lint:
	golangci-lint --version
	golangci-lint cache clean
	golangci-lint run --timeout=5m --verbose

test:
	go test ./...

fmt:
	@gofmt -w -s .

pretty: fmt fix

vendor:
	@go mod vendor

build-docker: vendor
	@docker build --rm -t downloader .

run-docker:
	docker run -i -t downloader
