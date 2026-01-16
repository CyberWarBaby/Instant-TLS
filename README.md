# InstantTLS

**Trusted HTTPS locally with zero browser warnings.**

InstantTLS is a developer tool that generates trusted local certificates, installs them in your OS trust store, and manages your development SSL/TLS workflow.

## 🚀 Quick Start (3 Commands!)

```bash
# 1. Install the CLI
go install github.com/CyberWarBaby/Instant-TLS/cli/cmd/instanttls@latest

# 2. One-time setup (login + create CA + trust in all browsers)
sudo instanttls setup

# 3. Generate a certificate
instanttls cert myapp.local
```

That's it! Your browser will now show a **green lock** 🔒 with zero warnings.

## Features

- 🔐 **Local CA Generation** - Create a trusted Certificate Authority on your machine
- 🌐 **Wildcard Certificates** - Generate certs like `*.local.test` for all your local domains
- 💻 **Cross-Platform** - Works on macOS, Linux, and Windows
- 🔄 **Auto Trust** - Automatically installs CA in Chrome, Firefox, and system store
- ☁️ **Cloud Licensing** - Free/Pro/Team plans with account management
- 🎨 **Beautiful CLI** - Polished terminal experience with colors and spinners
- 🖥️ **Web Dashboard** - Modern Next.js dashboard for token and account management

## Using the Certificate

After running `instanttls cert myapp.local`, use the certificates in your project:

**Node.js:**
```javascript
const https = require('https');
const fs = require('fs');
const path = require('path');

const options = {
  key: fs.readFileSync(path.join(process.env.HOME, '.instanttls/certs/myapp.local/key.pem')),
  cert: fs.readFileSync(path.join(process.env.HOME, '.instanttls/certs/myapp.local/cert.pem'))
};

https.createServer(options, (req, res) => {
  res.end('Hello HTTPS!');
}).listen(443);
```

**Don't forget to add your domain to /etc/hosts:**
```bash
echo "127.0.0.1 myapp.local" | sudo tee -a /etc/hosts
```

## Development Setup

### Prerequisites

- Go 1.21+
- Node.js 18+
- Docker & Docker Compose
- PostgreSQL (via Docker)

### 1. Clone and Setup

```bash
cd instanttls
cp .env.example .env
```

### 2. Start Development Environment

```bash
make dev
```

This starts:
- PostgreSQL on port 5433
- Go API on port 8081
- Next.js web on port 3000

### 3. Create an Account

1. Visit http://localhost:3000
2. Register a new account
3. Go to Tokens page and create a Personal Access Token
4. Copy the token (shown only once!)

### 4. Use the CLI

```bash
# Build the CLI
make build-cli

# One-time setup
sudo ./bin/instanttls setup

# Generate a wildcard certificate
./bin/instanttls cert "*.local.test"

# Check everything is working
./bin/instanttls doctor
```

## Project Structure

```
instanttls/
├── apps/
│   ├── api/           # Go Gin API server
│   └── web/           # Next.js dashboard
├── cli/               # Go CLI (Cobra)
├── infra/
│   └── docker-compose.yml
├── .env.example
├── Makefile
└── README.md
```

## CLI Commands

| Command | Description |
|---------|-------------|
| `instanttls login` | Authenticate with your Personal Access Token |
| `instanttls whoami` | Display current user and plan |
| `instanttls init` | Generate and install local CA |
| `instanttls cert <domain>` | Generate certificate for domain |
| `instanttls trust` | Re-install CA in OS trust store |
| `instanttls renew` | Renew expiring certificates |
| `instanttls doctor` | Diagnose setup issues |

## Plans

| Feature | Free | Pro | Team |
|---------|------|-----|------|
| Wildcard Certs | 1 | Unlimited | Unlimited |
| Local CA | ✅ | ✅ | ✅ |
| Auto-renew | ✅ | ✅ | ✅ |
| Priority Support | ❌ | ✅ | ✅ |
| Team Management | ❌ | ❌ | ✅ |

## Development

### API Only

```bash
make run-api
```

### Web Only

```bash
make run-web
```

### Database Migrations

```bash
make migrate-up
make migrate-down
```

### Build CLI

```bash
make build-cli
```

## API Endpoints

### Auth
- `POST /v1/auth/register` - Register new user
- `POST /v1/auth/login` - Login user

### User (requires auth)
- `GET /v1/me` - Get current user (PAT auth)

### Tokens (requires web auth)
- `GET /v1/tokens` - List tokens
- `POST /v1/tokens` - Create token
- `DELETE /v1/tokens/:id` - Revoke token

### License (requires PAT)
- `GET /v1/license` - Get plan and limits

### Machines (requires PAT)
- `POST /v1/machines/ping` - Register/update machine

## End-to-End Demo Workflow

Here's a complete workflow to test everything:

```bash
# 1. Start the environment
./scripts/setup.sh

# 2. In Terminal 1 - Start API
cd apps/api && go run .

# 3. In Terminal 2 - Start Web
cd apps/web && npm run dev

# 4. In Terminal 3 - Build and use CLI
cd cli && go build -o ../bin/instanttls .

# 5. Open browser and create a token
#    Go to http://localhost:3000
#    Login with: demo@instanttls.dev / demo1234
#    Navigate to Tokens page
#    Create a new token and copy it

# 6. Login with CLI
../bin/instanttls login
# Enter: http://localhost:8081 (or press enter for default)
# Paste your token

# 7. Initialize local CA
../bin/instanttls init

# 8. Generate a wildcard certificate
../bin/instanttls cert "*.local.test"

# 9. Verify everything
../bin/instanttls doctor
../bin/instanttls whoami
```

## Troubleshooting

### PostgreSQL connection issues
```bash
# Check if PostgreSQL is running
docker ps

# Restart PostgreSQL
docker-compose -f infra/docker-compose.yml down
docker-compose -f infra/docker-compose.yml up -d
```

### Trust store issues on Linux
```bash
# For Debian/Ubuntu
sudo cp ~/.instanttls/ca/ca.crt /usr/local/share/ca-certificates/instanttls.crt
sudo update-ca-certificates

# For Fedora/RHEL
sudo cp ~/.instanttls/ca/ca.crt /etc/pki/ca-trust/source/anchors/
sudo update-ca-trust
```

### Firefox not trusting certificates
1. Open `about:config` in Firefox
2. Set `security.enterprise_roots.enabled` to `true`
3. Restart Firefox

## License

MIT
