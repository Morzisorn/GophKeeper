
# ==============================================================================
# Help
# ==============================================================================

.PHONY: help
help: ## Show this help message
	@echo "GophKeeper Development Commands"
	@echo "=============================="
	@echo ""
	@echo "Quick Start (for trainees):"
	@echo "  make k3d-build  - Build and deploy to k3d"
	@echo "  make k3d-agent  - Start the password manager client"
	@echo "  make k3d-server - View server logs"
	@echo "  make k3d-db     - Connect to database"
	@echo ""
	@echo "All available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

.DEFAULT_GOAL := help

# ==============================================================================
# Code Generation
# ==============================================================================

.PHONY: proto
proto: ## Generate Go code from protobuf files
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/protos/users/users.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/protos/items/items.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/protos/crypto/crypto.proto

# ==============================================================================
# Testing
# ==============================================================================

.PHONY: testcoverage
testcoverage: ## Run tests with coverage report
	go test -coverprofile=coverage.out ./...
	@./scripts/filter-coverage.sh coverage.out

# ==============================================================================
# Local Development (runs on your machine)
# ==============================================================================

.PHONY: run-agent
run-agent: ## Run the client agent locally
	go run cmd/agent/main.go

.PHONY: run-server
run-server: ## Run the server locally
	go run cmd/server/main.go

# ==============================================================================
# k3d Development (runs in Docker containers)
# ==============================================================================

.PHONY: k3d-build
k3d-build: ## Build and deploy Docker images to k3d cluster
	@echo "Updating database schema..."
	kubectl delete configmap postgres-init --ignore-not-found=true
	kubectl create configmap postgres-init \
		--from-file=000_init.sql=k8s/init-wrapper.sql \
		--from-file=001_types.sql=internal/server/repositories/database/schema/001_types.sql \
		--from-file=002_tables.sql=internal/server/repositories/database/schema/002_tables.sql
	@echo "Building server Docker image..."
	docker build --platform linux/arm64 -t gophkeeper:latest --target server .
	@echo "Building agent Docker image..."
	docker build --platform linux/arm64 -t gophkeeper-agent:latest -f Dockerfile.agent .
	@echo "Importing images to k3d cluster..."
	k3d image import gophkeeper:latest gophkeeper-agent:latest -c gophkeeper
	@echo "Restarting deployments..."
	kubectl rollout restart deployment gophkeeper
	kubectl delete pod agent --ignore-not-found=true
	kubectl apply -f k8s/agent-pod.yaml
	@echo "✅ Build complete! Use 'make k3d-agent' to start the client"

.PHONY: k3d-build-server
k3d-build-server: ## Rebuild only the server (faster for server-only changes)
	@echo "Building server Docker image..."
	docker build --platform linux/arm64 -t gophkeeper:latest --target server .
	@echo "Importing server image to k3d cluster..."
	k3d image import gophkeeper:latest -c gophkeeper
	@echo "Restarting server deployment..."
	kubectl rollout restart deployment gophkeeper
	@echo "✅ Server rebuild complete!"

.PHONY: k3d-build-agent
k3d-build-agent: ## Rebuild only the agent (faster for agent-only changes)
	@echo "Building agent Docker image..."
	docker build --platform linux/arm64 -t gophkeeper-agent:latest -f Dockerfile.agent .
	@echo "Importing agent image to k3d cluster..."
	k3d image import gophkeeper-agent:latest -c gophkeeper
	@echo "Recreating agent pod..."
	kubectl delete pod agent --ignore-not-found=true
	kubectl apply -f k8s/agent-pod.yaml
	@echo "✅ Agent rebuild complete! Use 'make k3d-agent' to start the client"

.PHONY: k3d-agent
k3d-agent: ## Run the client agent in k3d (interactive terminal UI)
	kubectl exec -it agent -- env TERM=xterm ./agent

.PHONY: k3d-server
k3d-server: ## Show server logs from k3d
	kubectl logs -f deployment/gophkeeper



.PHONY: k3d-db
k3d-db: ## Connect to PostgreSQL database in k3d
	kubectl exec -it deployment/postgres -- psql -U gophkeeper -d gophkeeper_db

.PHONY: k3d-clean
k3d-clean: ## Clean up k3d cluster and all resources
	./k8s/cleanup.sh

