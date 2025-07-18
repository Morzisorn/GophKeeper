.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/protos/users/users.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/protos/items/items.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/protos/crypto/crypto.proto

.PHONY: testcoverage
testcoverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out   

.PHONY: agent
run-agent:
	go run cmd/agent/main.go

.PHONY: server
run-server:
	go run cmd/server/main.go
