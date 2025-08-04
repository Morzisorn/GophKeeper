
.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/protos/users/users.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/protos/items/items.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/protos/crypto/crypto.proto

.PHONY: testcoverage
testcoverage:
	go test -coverprofile=coverage.out ./...
	@./scripts/filter-coverage.sh coverage.out

.PHONY: agent
run-agent:
	go run cmd/agent/main.go

.PHONY: server
run-server:
	go run cmd/server/main.go

.PHONY: k3d-rebuild
k3d-rebuild:
	docker build --platform linux/arm64 -t gophkeeper:latest --target server .
	docker build --platform linux/arm64 -t gophkeeper-agent:latest -f Dockerfile.agent .
	k3d image import gophkeeper:latest gophkeeper-agent:latest -c gophkeeper
	kubectl rollout restart deployment gophkeeper
	kubectl delete pod agent
	kubectl apply -f k8s/agent-pod.yaml

.PHONY: k3d-agent
k3d-agent:
	kubectl exec -it agent -- env TERM=xterm ./agent

.PHONY: k3d-agent-logs
k3d-agent-logs:
	kubectl logs -f agent

.PHONY: k3d-agent-run
k3d-agent-run:
	kubectl exec -it agent -- ./agent

.PHONY: k3d-server
k3d-server:
	kubectl logs -f deployment/gophkeeper

.PHONY: k3d-db
k3d-db:
	kubectl exec -it deployment/postgres -- psql -U dmitrij -d gophkeeper_db
