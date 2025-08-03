#!/bin/bash

# Launchpad Development Setup Script

set -e

echo "🚀 Setting up Launchpad development environment..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    print_error "Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    print_error "Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    print_error "Node.js is not installed. Please install Node.js 18+ first."
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

print_status "All prerequisites found!"

# Create environment files if they don't exist
print_status "Creating environment files..."

# Backend .env
if [ ! -f "backend/.env" ]; then
    cat > backend/.env << EOF
DATABASE_URL=postgres://postgres:postgres@localhost:5432/launchpad?sslmode=disable
REDIS_URL=redis://localhost:6379
RPC_URL=http://localhost:8545
JWT_SECRET=your-secret-key-change-in-production
PORT=8080
EOF
    print_success "Created backend/.env"
else
    print_warning "backend/.env already exists, skipping..."
fi

# Frontend .env.local
if [ ! -f "frontend/.env.local" ]; then
    cat > frontend/.env.local << EOF
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_RPC_URL=http://localhost:8545
EOF
    print_success "Created frontend/.env.local"
else
    print_warning "frontend/.env.local already exists, skipping..."
fi

# Install dependencies
print_status "Installing dependencies..."

# Install contracts dependencies
print_status "Installing contract dependencies..."
cd contracts
npm install
cd ..

# Install frontend dependencies
print_status "Installing frontend dependencies..."
cd frontend
npm install
cd ..

# Install backend dependencies
print_status "Installing backend dependencies..."
cd backend
go mod download
cd ..

print_success "All dependencies installed!"

# Start services with Docker Compose
print_status "Starting database and supporting services..."
docker-compose up -d postgres redis

# Wait for database to be ready
print_status "Waiting for database to be ready..."
sleep 10

print_success "Setup completed!"

echo ""
echo "🎉 Launchpad development environment is ready!"
echo ""
echo "To start the development servers:"
echo ""
echo "1. Start Hardhat node:"
echo "   cd contracts && npm run node"
echo ""
echo "2. Start backend (in new terminal):"
echo "   cd backend && go run main.go"
echo ""
echo "3. Start frontend (in new terminal):"
echo "   cd frontend && npm run dev"
echo ""
echo "4. Open http://localhost:3000 in your browser"
echo ""
echo "Or use Docker Compose to start everything:"
echo "   docker-compose up"
echo ""
print_success "Happy coding! 🚀"