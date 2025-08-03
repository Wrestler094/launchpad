# Launchpad MVP Platform

A decentralized token launchpad platform that allows users to create ERC-20 tokens and launch presales through a simple web interface.

## 🎯 Features

- **MetaMask Authentication**: Secure login using Ethereum wallet signatures
- **Token Creation**: Deploy custom ERC-20 tokens with configurable parameters
- **Presale Management**: Create and manage token presales with customizable parameters
- **Landing Pages**: Auto-generated landing pages for presale participants
- **Real-time Updates**: Live presale status and participation tracking

## 🏗️ Architecture

The project is organized into the following directories:

```
/launchpad
├── /frontend          # Next.js React application
├── /backend          # Go REST API server  
├── /smart-contracts  # Solidity contracts & Hardhat
├── /deploy          # Docker Compose configs
└── /docs           # Documentation
```

### Components

1. **Smart Contracts** (Solidity + Hardhat)
   - `MyToken.sol`: ERC-20 token with mint/burn functionality
   - `Presale.sol`: Presale contract with caps, deadlines, and refunds
   - `LaunchpadFactory.sol`: Factory for deploying tokens and presales

2. **Backend** (`/backend` - Go + Chi)
   - REST API with JWT authentication
   - PostgreSQL database integration
   - Ethereum blockchain integration
   - MetaMask signature verification

3. **Frontend** (`/frontend` - Next.js + TypeScript)
   - React components with Tailwind CSS
   - wagmi + viem for Web3 integration
   - Responsive design for all devices

4. **Infrastructure**
   - PostgreSQL for metadata storage
   - Local Hardhat network for development
   - Docker Compose for easy deployment

### API Endpoints

```
POST /api/auth/nonce      - Generate nonce for MetaMask
POST /api/auth/login      - Authenticate with signature
POST /api/token/create    - Deploy new token
POST /api/presale/create  - Create presale
GET  /api/presale/{id}    - Get presale details
GET  /api/public/presale/{id} - Public presale info for landing
```

## 🚀 Quick Start

### Prerequisites

- Node.js 18+
- Go 1.21+
- Docker & Docker Compose
- MetaMask browser extension

### Development Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd launchpad
   ```

2. **Start with Docker Compose (Recommended)**
   ```bash
   cd deploy
   docker compose up -d
   ```

   **Or use the root Makefile for convenience**
   ```bash
   make dev     # Start development environment
   make help    # Show all available commands
   ```

3. **Or run components individually:**

   **Smart Contracts (Hardhat Network)**
   ```bash
   cd smart-contracts
   npm install
   npm run node
   ```

   **Backend (Go API)**
   ```bash
   cd backend
   cp .env.example .env
   go mod tidy
   make run
   ```

   **Frontend (Next.js)**
   ```bash
   cd frontend
   npm install
   cp .env.example .env.local
   npm run dev
   ```

4. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Hardhat Network: http://localhost:8545

### MetaMask Setup

1. Add local Hardhat network to MetaMask:
   - Network Name: Hardhat Local
   - RPC URL: http://localhost:8545
   - Chain ID: 31337
   - Currency Symbol: ETH

2. Import test account:
   - Private Key: `0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80`

## 📝 Usage

### Creating a Token

1. Connect MetaMask wallet
2. Sign authentication message
3. Navigate to "Create Token" tab
4. Fill in token details:
   - Name (e.g., "My Token")
   - Symbol (e.g., "MTK")
   - Total Supply (e.g., 1000000)
5. Click "Create Token"

### Creating a Presale

1. Ensure you have created a token first
2. Navigate to "Create Presale" tab
3. Configure presale parameters:
   - Select your token
   - Set rate (tokens per ETH)
   - Set soft cap and hard cap
   - Set deadline
4. Click "Create Presale"
5. Share the generated landing page URL

### Participating in Presales

1. Visit a presale landing page
2. Connect MetaMask wallet
3. Enter ETH amount to contribute
4. Confirm transaction
5. Receive tokens automatically

## 🔧 Configuration

### Environment Variables

**Backend (.env)**
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=launchpad
DB_PASSWORD=password
DB_NAME=launchpad
RPC_URL=http://localhost:8545
PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
PORT=8080
```

**Frontend (.env.local)**
```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

## 🗄️ Database Schema

```sql
-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    address VARCHAR(42) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tokens table
CREATE TABLE tokens (
    id SERIAL PRIMARY KEY,
    address VARCHAR(42) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    symbol VARCHAR(10) NOT NULL,
    total_supply VARCHAR(255) NOT NULL,
    creator_address VARCHAR(42) NOT NULL,
    tx_hash VARCHAR(66) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Presales table
CREATE TABLE presales (
    id SERIAL PRIMARY KEY,
    address VARCHAR(42) UNIQUE NOT NULL,
    token_address VARCHAR(42) NOT NULL,
    creator_address VARCHAR(42) NOT NULL,
    rate VARCHAR(255) NOT NULL,
    soft_cap VARCHAR(255) NOT NULL,
    hard_cap VARCHAR(255) NOT NULL,
    deadline TIMESTAMP NOT NULL,
    active BOOLEAN DEFAULT true,
    finalized BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🛠️ Development

### Building Components

**All Components (from root)**
```bash
make all      # Build everything
make backend  # Build just the backend
make frontend # Build just the frontend
make clean    # Clean all build artifacts
```

**Smart Contracts**
```bash
cd smart-contracts
npm run compile
npm run test
npm run deploy
```

**Backend**
```bash
cd backend
make server    # Build the server
make run       # Build and run the server  
make clean     # Clean build artifacts
```

**Frontend**
```bash
cd frontend
npm run build
npm start
```

### Testing

**Smart Contracts**
```bash
cd smart-contracts
npm run test
```

**Backend**
```bash
go test ./...
```

## 🚧 MVP Limitations

This is an MVP implementation with the following limitations:

1. **Smart Contract Integration**: Currently simulated - contracts are not actually deployed
2. **Transaction Handling**: Mock transaction hashes used instead of real blockchain transactions
3. **Security**: Basic JWT implementation, should be enhanced for production
4. **Scalability**: Single-server setup, needs load balancing for production
5. **Testing**: Limited test coverage, needs comprehensive test suite

## 🔄 Next Steps

1. **Smart Contract Integration**: Connect backend to actual smart contracts
2. **Transaction Processing**: Implement real blockchain transaction handling
3. **Enhanced Security**: Add rate limiting, input validation, CSRF protection
4. **Testing**: Add comprehensive unit and integration tests
5. **Monitoring**: Add logging, metrics, and monitoring
6. **Mobile Support**: Optimize for mobile wallet integration
7. **Multi-chain Support**: Support other EVM-compatible chains

## 📄 License

MIT License - see LICENSE file for details