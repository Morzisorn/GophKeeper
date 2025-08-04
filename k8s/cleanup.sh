#!/bin/bash
set -e

echo "🧹 Cleaning up GophKeeper deployment..."

echo "🗑️  Removing Kubernetes resources..."
kubectl delete -f k8s/deployment.yaml --ignore-not-found=true
kubectl delete -f k8s/postgres.yaml --ignore-not-found=true
kubectl delete -f k8s/init-db.yaml --ignore-not-found=true
kubectl delete -f k8s/secrets.yaml --ignore-not-found=true
kubectl delete -f k8s/configmap.yaml --ignore-not-found=true
kubectl delete configmap gophkeeper-keys --ignore-not-found=true

echo "🗂️  Removing Docker image from k3d..."
k3d image remove gophkeeper:latest -c gophkeeper || true

echo "🏗️  Stopping and removing k3d cluster..."
k3d cluster delete gophkeeper || echo "Cluster may not exist"

echo "🐳 Removing local Docker image..."
docker rmi gophkeeper:latest || echo "Image may not exist"

echo "✅ Cleanup completed!"