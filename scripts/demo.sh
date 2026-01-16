#!/bin/bash
set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔══════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║              InstantTLS Demo Script                      ║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════════════════════════╝${NC}"
echo ""

# Check if we're in the right directory
if [ ! -f "Makefile" ]; then
    echo "Error: Please run this script from the instanttls root directory"
    exit 1
fi

echo -e "${YELLOW}Step 1: Starting infrastructure...${NC}"
echo "Starting PostgreSQL with Docker Compose..."
docker-compose -f infra/docker-compose.yml up -d
sleep 3

echo ""
echo -e "${YELLOW}Step 2: Running database migrations...${NC}"
cd apps/api
go run . migrate up
cd ../..

echo ""
echo -e "${YELLOW}Step 3: Seeding demo user...${NC}"
cd apps/api
go run . seed
cd ../..
echo -e "${GREEN}✓ Demo user created: demo@instanttls.dev / demo1234${NC}"

echo ""
echo -e "${YELLOW}Step 4: Building CLI...${NC}"
cd cli
go build -o ../bin/instanttls .
cd ..
echo -e "${GREEN}✓ CLI built to ./bin/instanttls${NC}"

echo ""
echo -e "${YELLOW}Step 5: Starting API server (background)...${NC}"
cd apps/api
go run . &
API_PID=$!
cd ../..
sleep 3
echo -e "${GREEN}✓ API running on http://localhost:8081${NC}"

echo ""
echo -e "${BLUE}╔══════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                    Demo Ready!                           ║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════════════════════════╝${NC}"
echo ""
echo "Demo credentials:"
echo "  Email:    demo@instanttls.dev"
echo "  Password: demo1234"
echo ""
echo "Next steps:"
echo "  1. Open http://localhost:3000 in your browser"
echo "  2. Login with demo credentials"
echo "  3. Create a Personal Access Token"
echo "  4. Use the CLI:"
echo ""
echo "     ./bin/instanttls login"
echo "     ./bin/instanttls init"
echo "     ./bin/instanttls cert \"*.local.test\""
echo "     ./bin/instanttls doctor"
echo ""
echo -e "${YELLOW}To start the web dashboard, run in a new terminal:${NC}"
echo "  cd apps/web && npm install && npm run dev"
echo ""
echo "Press Ctrl+C to stop the API server..."

# Wait for the API server
wait $API_PID
