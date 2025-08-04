#!/bin/bash
set -e

# Check if k3d cluster exists, create if not
if ! k3d cluster list | grep -q gophkeeper; then
    echo "🏗️  Creating k3d cluster..."
    k3d cluster create gophkeeper --port "8080:30080@agent:0" --port "5432:30432@agent:0"
else
    echo "✅ k3d cluster 'gophkeeper' already exists"
fi

# Ensure cluster is started
echo "🔄 Starting k3d cluster..."
k3d cluster start gophkeeper

echo "🏗️  Building Docker images..."
docker build --platform linux/arm64 -t gophkeeper:latest --target server .
docker build --platform linux/arm64 -t gophkeeper-agent:latest -f Dockerfile.agent .

echo "📦 Importing images to k3d..."
k3d image import gophkeeper:latest gophkeeper-agent:latest -c gophkeeper

echo "🔧 Applying configurations..."
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secrets.yaml

echo "🔑 Creating keys configmap..."
if [ -f "private_key.pem" ] && [ -f "public_key.pem" ]; then
    kubectl create configmap gophkeeper-keys \
      --from-file=private.pem=private_key.pem \
      --from-file=public.pem=public_key.pem \
      --dry-run=client -o yaml | kubectl apply -f -
else
    echo "⚠️  Key files not found. Please ensure private_key.pem and public_key.pem exist in the project root."
    exit 1
fi

echo "🗄️  Setting up database schema..."
kubectl delete configmap postgres-init --ignore-not-found=true
kubectl create configmap postgres-init \
  --from-file=000_init.sql=k8s/init-wrapper.sql \
  --from-file=001_types.sql=internal/server/repositories/database/schema/001_types.sql \
  --from-file=002_tables.sql=internal/server/repositories/database/schema/002_tables.sql

echo "🗄️  Deploying PostgreSQL..."
kubectl apply -f k8s/postgres.yaml
kubectl rollout status deployment/postgres --timeout=300s

echo "🔧 Waiting for database to be ready..."
sleep 15  # Wait for postgres to be ready
kubectl wait --for=condition=available --timeout=300s deployment/postgres

echo "🚀 Deploying application..."
kubectl apply -f k8s/deployment.yaml
kubectl rollout status deployment/gophkeeper --timeout=300s

echo "🤖 Deploying agent pod..."
kubectl apply -f k8s/agent-pod.yaml

echo "📋 Getting status..."
kubectl get pods,svc

echo "✅ Deployment completed!"
echo "🤖 To use the password manager client:"
echo "   make k3d-agent"
echo ""
echo "📝 Useful commands:"
echo "   Server logs: make k3d-server"
echo "   Database:    make k3d-db"
echo "   Rebuild:     make k3d-build"
