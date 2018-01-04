.PHONY: test

tgdocs2go:
	go build

test:
	@echo "Running the tests with gofmt, go vet and golint..."
	@test -z $(shell gofmt -s -l tgdocs2go.go)
	@go vet ./...
	@golint -set_exit_status $(shell go list ./...)
