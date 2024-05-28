.PHONY: tidy
tidy:
	go mod tidy

.PHONY: style
style:
	goimports -l -w ./cmd
	goimports -l -w ./internal

.PHONY: clean
clean:
	go clean -testcache

.PHONY: test
test:
	go test -v -race -cover ./...

.PHONY: build
build:
	CGO_ENABLED=0 go build -o /Users/wesleyanderson/repos/github.com/w-h-a/micro/bin/micro ./
