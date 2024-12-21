.PHONY: mock
mock:
	@go mod tidy
	@go generate ./...

.PHONE: grpc
grpc:
	@buf generate api/proto