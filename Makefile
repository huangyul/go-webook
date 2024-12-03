.PHONY: mock
mock:
	@mockgen -source=./internal/service/user.go -destination=./internal/service/mock/mock_user.go -package=mocksvc
	@mockgen -source=./internal/service/code.go -destination=./internal/service/mock/mock_code.go -package=mocksvc
