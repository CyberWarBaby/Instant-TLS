# Railway Deployment Guide

This guide explains how to deploy the InstantTLS application to Railway. InstantTLS is a monorepo with two separate services that need to be deployed independently.

## Architecture

The application consists of:
- **API Service** - Go backend API (`apps/api`)
- **Web Service** - Next.js frontend dashboard (`apps/web`)
- **Database** - PostgreSQL (Railway Plugin)

## Quick Deploy

### Prerequisites

1. A Railway account ([railway.app](https://railway.app))
2. Railway CLI installed (optional but recommended): `npm i -g @railway/cli`

### Option 1: Deploy via Railway Dashboard (Recommended)

#### Step 1: Create a New Project

1. Go to [railway.app](https://railway.app)
2. Click "New Project"
3. Select "Deploy from GitHub repo"
4. Connect and select the `CyberWarBaby/Instant-TLS` repository

#### Step 2: Add PostgreSQL Database

1. In your Railway project, click "New"
2. Select "Database" → "Add PostgreSQL"
3. Railway will automatically provision a PostgreSQL database

#### Step 3: Deploy the API Service

1. Click "New" → "GitHub Repo"
2. Select your InstantTLS repository
3. In the service settings:
   - **Name**: `instanttls-api`
   - **Root Directory**: `apps/api`
   - **Start Command**: `./api` (this will be auto-detected)
   - **Build Command**: `go build -o api .` (this will be auto-detected)

4. Add environment variables:
   ```
   DATABASE_URL=${{Postgres.DATABASE_URL}}
   PORT=8081
   JWT_SECRET=<generate-a-secure-random-string>
   ENV=production
   CORS_ORIGINS=https://<your-web-service>.railway.app
   ```
   
   Replace `<your-web-service>` with the domain Railway assigns to your Web service.
   
   **Note**: You can add multiple origins separated by commas:
   ```
   CORS_ORIGINS=https://your-web.railway.app,https://your-custom-domain.com
   ```

5. Click "Deploy"

#### Step 4: Deploy the Web Service

1. Click "New" → "GitHub Repo"
2. Select your InstantTLS repository again
3. In the service settings:
   - **Name**: `instanttls-web`
   - **Root Directory**: `apps/web`
   - **Start Command**: `npm start` (auto-detected)
   - **Build Command**: `npm install && npm run build` (auto-detected)

4. Add environment variables:
   ```
   NEXT_PUBLIC_API_URL=https://<your-api-service>.railway.app
   API_URL=https://<your-api-service>.railway.app
   ```
   
   Replace `<your-api-service>` with the domain Railway assigns to your API service.

5. Click "Deploy"

### Option 2: Deploy via Railway CLI

```bash
# Login to Railway
railway login

# Create a new project
railway init

# Link to the project
railway link

# Deploy the API
cd apps/api
railway up

# In Railway dashboard:
# 1. Add PostgreSQL plugin
# 2. Set environment variables for API service
# 3. Go back to root directory

# Deploy the Web
cd ../web
railway up

# Set environment variables for Web service in Railway dashboard
```

## Environment Variables

### API Service

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgresql://user:pass@host:5432/db` |
| `PORT` | Server port | `8081` |
| `JWT_SECRET` | Secret key for JWT tokens | `your-super-secret-jwt-key` |
| `ENV` | Environment mode | `production` |
| `CORS_ORIGINS` | Allowed CORS origins (comma-separated) | `https://web.railway.app` |

**Note**: Railway will automatically provide `DATABASE_URL` when you add the PostgreSQL plugin. Reference it as `${{Postgres.DATABASE_URL}}`.

### Web Service

| Variable | Description | Example |
|----------|-------------|---------|
| `NEXT_PUBLIC_API_URL` | Public API URL (client-side) | `https://api.railway.app` |
| `API_URL` | Server-side API URL (optional) | `https://api.railway.app` |

## Post-Deployment

### 1. Run Database Migrations

The API service will automatically run migrations on startup. If you need to run them manually:

```bash
# Via Railway CLI
railway run ./api migrate up
```

### 2. Verify Deployment

1. Check API health:
   ```bash
   curl https://<your-api-domain>.railway.app/health
   ```

2. Visit your web dashboard:
   ```
   https://<your-web-domain>.railway.app
   ```

### 3. Create Your First Account

1. Navigate to your web dashboard
2. Click "Register" and create an account
3. Login and create a Personal Access Token
4. Use the token with the CLI

## Troubleshooting

### Build Failures

**Problem**: Go build fails in API service
- **Solution**: Ensure the root directory is set to `apps/api` in service settings

**Problem**: Next.js build fails in Web service
- **Solution**: Ensure the root directory is set to `apps/web` in service settings

### Database Connection Issues

**Problem**: API can't connect to database
- **Solution**: Make sure you've added the PostgreSQL plugin and set `DATABASE_URL=${{Postgres.DATABASE_URL}}`

### Environment Variable Issues

**Problem**: Web can't connect to API
- **Solution**: Double-check that `NEXT_PUBLIC_API_URL` is set to your API service's public URL

**Problem**: CORS errors when Web tries to access API
- **Solution**: Make sure you've set `CORS_ORIGINS` on the API service to include your Web service URL:
  ```
  CORS_ORIGINS=https://your-web-service.railway.app
  ```

### Port Issues

**Problem**: Service won't start
- **Solution**: Railway automatically sets the `PORT` environment variable. Make sure your app uses it:
  - API: Uses `PORT` from env (already implemented)
  - Web: Next.js automatically uses Railway's `PORT`

## Monorepo Configuration

This repository uses the following Railway configurations:

- `apps/api/railway.json` - API service configuration
- `apps/api/nixpacks.toml` - API build configuration
- `apps/web/railway.json` - Web service configuration
- `apps/web/nixpacks.toml` - Web build configuration

These files configure Nixpacks (Railway's build system) to properly build and deploy each service.

## Custom Domains

To add custom domains:

1. Go to your service settings in Railway
2. Click "Settings" → "Networking"
3. Add your custom domain
4. Update DNS records as instructed
5. Update environment variables to use the new domain

## Scaling

Railway automatically scales your services. For production workloads:

1. Consider upgrading to Railway Pro for better resources
2. Monitor your database connection pool
3. Use Railway's built-in metrics to track performance

## Support

- Railway Docs: [docs.railway.app](https://docs.railway.app)
- Railway Discord: [discord.gg/railway](https://discord.gg/railway)
- Project Issues: [GitHub Issues](https://github.com/CyberWarBaby/Instant-TLS/issues)
