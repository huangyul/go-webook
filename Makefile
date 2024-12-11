.PHONY: mock
mock:
	@go mod tidy
	@go generate ./...
