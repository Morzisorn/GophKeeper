# GophKeeper Kubernetes Deployment

This directory contains Kubernetes manifests and scripts for deploying GophKeeper using k3d.

## Prerequisites

- Docker
- k3d
- kubectl

## Quick Start

1. **Deploy the application:**
   ```bash
   ./k8s/deploy.sh
   ```

2. **Verify deployment:**
   ```bash
   ./k8s/verify.sh
   ```

3. **Access the application:**
    - gRPC Server: localhost:8080 (NodePort)
    - Or use port-forwarding: `kubectl port-forward svc/gophkeeper-service 8080:80`
    - PostgreSQL: localhost:5432 (if needed for debugging)
    - Connect CLI agent to: localhost:8080 (gRPC endpoint)

4. **Clean up:**
   ```bash
   ./k8s/cleanup.sh
   ```

## Files Description

- `deploy.sh` - Main deployment script that creates k3d cluster and deploys all components
- `cleanup.sh` - Cleanup script to remove all resources and cluster
- `verify.sh` - Verification script to check deployment status and connectivity
- `deployment.yaml` - Main application deployment and service
- `postgres.yaml` - PostgreSQL database deployment and service
- `configmap.yaml` - Application configuration
- `secrets.yaml` - Application secrets (DB connection, JWT key)
- `init-db.yaml` - Database initialization scripts
- `README.md` - This file

## Manual Deployment Steps

If you prefer to deploy manually:

1. **Create k3d cluster:**
   ```bash
   k3d cluster create gophkeeper --port "8080:30080@agent:0" --port "5432:30432@agent:0"
   ```

2. **Build and import Docker image:**
   ```bash
   docker build --platform linux/arm64 -t gophkeeper:latest .
   k3d image import gophkeeper:latest -c gophkeeper
   ```

3. **Deploy configurations:**
   ```bash
   kubectl apply -f k8s/configmap.yaml
   kubectl apply -f k8s/secrets.yaml
   ```

4. **Create keys configmap:**
   ```bash
   kubectl create configmap gophkeeper-keys \
     --from-file=private.pem=private_key.pem \
     --from-file=public.pem=public_key.pem \
     --dry-run=client -o yaml | kubectl apply -f -
   ```

5. **Deploy PostgreSQL:**
   ```bash
   kubectl apply -f k8s/init-db.yaml
   kubectl apply -f k8s/postgres.yaml
   kubectl rollout status deployment/postgres --timeout=300s
   ```

6. **Deploy application:**
   ```bash
   kubectl apply -f k8s/deployment.yaml
   kubectl rollout status deployment/gophkeeper --timeout=300s
   ```

## Troubleshooting

### Check logs

```bash
# Application logs
kubectl logs -f deployment/gophkeeper

# Database logs
kubectl logs -f deployment/postgres
```

### Check pod status

```bash
kubectl get pods
kubectl describe pod <pod-name>
```

### Port forwarding (alternative access)

```bash
# Application
kubectl port-forward svc/gophkeeper-service 8080:80

# Database
kubectl port-forward svc/postgres 5432:5432
```

### Database connection

```bash
# Connect to database for debugging
kubectl exec -it deployment/postgres -- psql -U dmitrij -d gophkeeper_db
```

## Configuration

### Environment Variables

The application uses the following environment variables configured in `configmap.yaml` and `secrets.yaml`:

- `SERVER_HOST` - gRPC server bind address (default: 0.0.0.0)
- `SERVER_PORT` - gRPC server port (default: 8080)
- `LOG_LEVEL` - Logging level (default: debug)
- `AUTH_TOKEN_EXPIRY` - JWT token expiry time (default: 24h)
- `DB_MAX_CONNECTIONS` - Maximum database connections (default: 25)
- `DB_TIMEOUT` - Database operation timeout (default: 30s)
- `AUTH_PRIVATE_KEY_PATH` - Path to RSA private key (mounted from configmap)
- `AUTH_PUBLIC_KEY_PATH` - Path to RSA public key (mounted from configmap)
- `AUTH_SECRET_KEY` - JWT signing secret
- `DATABASE_URI` - PostgreSQL connection string

### Key Files

The RSA key files (`private_key.pem` and `public_key.pem`) must be present in the project root directory before
deployment. These are automatically mounted into the container at `/etc/keys/`.

## Security Notes

⚠️ **Warning**: The current configuration uses hardcoded database credentials for development purposes. In production:

1. Use Kubernetes secrets with base64 encoded values
2. Use strong, randomly generated passwords
3. Enable SSL/TLS for database connections
4. Use RBAC for proper access control
5. Store sensitive data in a proper secret management system
