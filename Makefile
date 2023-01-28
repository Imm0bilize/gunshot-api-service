cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

protogen:
	protoc -I . api/proto/v1/api_service.proto --go_out=pkg/ --go-grpc_out=pkg/ --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative

mockgen:
	mockgen -source=./internal/infrastructure/repository/repository.go -destination=internal/infrastructure/repository/mocks/mockrepository.go

.PHONY: cover protogen