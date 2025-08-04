#!/bin/bash

# Simple verification script for troubleshooting
# For normal development, use: make logs, make db, etc.

echo "🔍 Quick GophKeeper Status Check"
echo "================================"

echo ""
echo "🔧 Pod Status:"
kubectl get pods

echo ""
echo "📋 Recent Server Logs:"
kubectl logs deployment/gophkeeper --tail=5

echo ""
echo "🗄️  Database Status:"
kubectl logs deployment/postgres --tail=3

echo ""
echo "✅ Status check complete!"
echo ""
echo "For development, use these commands instead:"
echo "  make logs     # View server logs"
echo "  make db       # Connect to database"  
echo "  make client   # Start the client"