# k3d Development Setup for Trainees

This guide explains how to use k3d for GophKeeper

## Quick Start

### 1. Build and Deploy

```bash
make k3d-build
# This builds Docker images and deploys them to k3d
```

### 2. Start the Client

```bash
make k3d-agent
# This opens the password manager client in your terminal
```

### 3. Check Database (optional)

```bash
make k3d-db
# This connects to PostgreSQL to see stored data
```

## What's Running?

When you run `make k3d-build`, these components start in k3d:

- **Server** (`gophkeeper` deployment) - The password manager server
- **Agent** (`agent` pod) - The client application (dormant until you run `make k3d-agent`)
- **Database** (`postgres` deployment) - PostgreSQL database for storing encrypted data

## Typical Development Workflow

1. Make code changes
2. Rebuild and redeploy:
   - `make k3d-build` - Full rebuild (both server and agent)
   - `make k3d-build-server` - Server only (faster for server changes)
   - `make k3d-build-agent` - Agent only (faster for client changes)
3. Run `make k3d-agent` to test your changes
4. Use `make k3d-db` to inspect database changes

## Available Commands

Run `make` or `make help` to see all available commands:

- `make k3d-build` - Build and deploy everything
- `make k3d-build-server` - Rebuild only server (faster for server changes)
- `make k3d-build-agent` - Rebuild only agent (faster for client changes)
- `make k3d-agent` - Start the password manager client
- `make k3d-server` - View server logs (server runs automatically)
- `make k3d-db` - Connect to database

## Prerequisites

- Docker
- k3d
- kubectl

## Files Description

- `deploy.sh` - Main deployment script that creates k3d cluster and deploys all components
- `cleanup.sh` - Cleanup script to remove all resources and cluster
- `verify.sh` - Verification script to check deployment status and connectivity
- `deployment.yaml` - Main application deployment and service
- `postgres.yaml` - PostgreSQL database deployment and service
- `agent-pod.yaml` - Agent (client) pod for testing
- `configmap.yaml` - Application configuration
- `secrets.yaml` - Application secrets
- `init-wrapper.sql` - Database initialization wrapper (imports actual schema files)
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
   # Create database schema configmap from actual schema files
   kubectl create configmap postgres-init \
     --from-file=000_init.sql=k8s/init-wrapper.sql \
     --from-file=001_types.sql=internal/server/repositories/database/schema/001_types.sql \
     --from-file=002_tables.sql=internal/server/repositories/database/schema/002_tables.sql
   
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

- `ADDRESS` - Server bind address and port (0.0.0.0:8080)
- `KEYS_DIR` - Directory where RSA keys are mounted (/etc/keys)
- `DATABASE_URI` - PostgreSQL connection string (from secrets)

### Key Files

The RSA key files (`private_key.pem` and `public_key.pem`) must be present in the project root directory before
deployment. These are automatically mounted into the container at `/etc/keys/`.
