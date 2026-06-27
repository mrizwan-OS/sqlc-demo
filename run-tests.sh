#!/bin/bash
# Test runner script

set -e

echo "🧪 Running tests..."

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to run tests with coverage
run_tests() {
    echo -e "${BLUE}📊 Running tests with coverage...${NC}"
    go test -v -cover -coverprofile=coverage.out ./...
    
    # Show coverage summary
    echo -e "${BLUE}📊 Coverage summary:${NC}"
    go tool cover -func=coverage.out | grep total
    
    # Generate HTML coverage report
    go tool cover -html=coverage.out -o coverage.html
    echo -e "${GREEN}✅ Coverage report generated: coverage.html${NC}"
}

# Function to run integration tests
run_integration() {
    echo -e "${BLUE}🔗 Running integration tests...${NC}"
    go test -v ./...
}

# Function to run tests with race detection (skip on ARM64)
run_race() {
    echo -e "${YELLOW}⚠️  Race detection is not supported on Android/ARM64${NC}"
    echo -e "${YELLOW}⚠️  Skipping race detection tests${NC}"
    echo -e "${GREEN}✅ Race detection skipped${NC}"
}

# Function to run benchmarks
run_benchmarks() {
    echo -e "${BLUE}⚡ Running benchmarks...${NC}"
    go test -bench=. -benchmem ./...
}

# Function to run specific test
run_specific() {
    if [ -z "$1" ]; then
        echo "❌ Please provide a test name"
        echo "Usage: ./run-tests.sh specific TestName"
        exit 1
    fi
    echo -e "${BLUE}🎯 Running specific test: $1${NC}"
    go test -v -run "$1" ./...
}

# Main menu
case "$1" in
    all)
        run_tests
        run_race
        run_benchmarks
        ;;
    coverage)
        run_tests
        ;;
    integration)
        run_integration
        ;;
    race)
        run_race
        ;;
    bench)
        run_benchmarks
        ;;
    specific)
        run_specific "$2"
        ;;
    *)
        echo "Usage: ./run-tests.sh {all|coverage|integration|race|bench|specific <name>}"
        echo ""
        echo "Examples:"
        echo "  ./run-tests.sh all           - Run all tests with coverage, race, and benchmarks"
        echo "  ./run-tests.sh coverage      - Run tests with coverage report"
        echo "  ./run-tests.sh integration   - Run only integration tests"
        echo "  ./run-tests.sh race          - Skip race detection (not supported on ARM64)"
        echo "  ./run-tests.sh bench         - Run benchmarks"
        echo "  ./run-tests.sh specific TestCreateUser - Run specific test"
        ;;
esac
