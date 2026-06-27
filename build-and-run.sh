#!/bin/bash

set -e

echo "🏗️  Building and running sqlc-demo..."

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 1. Generate sqlc code
echo -e "${YELLOW}🔧 Generating sqlc code...${NC}"
sqlc generate

# 2. Build for ARM64
echo -e "${YELLOW}📦 Building binary...${NC}"
mkdir -p bin
CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$(git describe --tags 2>/dev/null || echo 'dev')" -o bin/app main.go

# 3. Check if PostgreSQL is running
echo -e "${YELLOW}📊 Checking PostgreSQL...${NC}"
if ! pg_ctl -D $PREFIX/var/lib/postgresql status &>/dev/null; then
    echo "Starting PostgreSQL..."
    pg_ctl -D $PREFIX/var/lib/postgresql start
    sleep 2
fi

# 4. Run migrations if needed
echo -e "${YELLOW}🔄 Running migrations...${NC}"
./migrate.sh up 2>/dev/null || echo "Migrations already applied"

# 5. Run the application
echo -e "${GREEN}✅ Build complete!${NC}"
echo -e "${GREEN}🚀 Starting application...${NC}"
echo -e "${YELLOW}📋 Press Ctrl+C to stop${NC}"

./bin/app
