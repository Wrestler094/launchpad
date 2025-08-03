# Launchpad Platform Architecture

## System Overview

The Launchpad platform is a decentralized token creation and presale platform built with a modern web3 architecture. It enables users to create ERC-20 tokens and launch presales through a user-friendly interface.

## Architecture Diagram

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│                 │    │                 │    │                 │
│   Frontend      │    │    Backend      │    │  Smart          │
│   (Next.js)     │◄──►│    (Go/Chi)     │◄──►│  Contracts      │
│                 │    │                 │    │  (Solidity)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                        │                        │
         │                        │                        │
         ▼                        ▼                        ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│                 │    │                 │    │                 │
│   MetaMask      │    │   PostgreSQL    │    │   Hardhat       │
│   (Web3)        │    │   (Database)    │    │   (Local EVM)   │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Component Details

### 1. Frontend Layer (Next.js + TypeScript)

**Technology Stack:**
- Next.js 14 with App Router
- TypeScript for type safety
- Tailwind CSS for styling
- wagmi + viem for Web3 integration
- TanStack Query for state management

**Key Components:**
- `pages/index.tsx`: Main dashboard with authentication
- `components/TokenCreator.tsx`: Token creation form
- `components/PresaleCreator.tsx`: Presale configuration form
- `components/Dashboard.tsx`: User's tokens and presales overview
- `pages/presale/[id].tsx`: Public presale landing page

**Features:**
- MetaMask wallet connection
- Signature-based authentication
- Token creation interface
- Presale management dashboard
- Public presale participation pages
- Responsive design

### 2. Backend Layer (Go + Chi Router)

**Technology Stack:**
- Go 1.21+ for performance and reliability
- Chi router for HTTP routing
- PostgreSQL for data persistence
- go-ethereum for blockchain interaction
- JWT for authentication

**Architecture Patterns:**
- Clean Architecture with service layers
- Repository pattern for data access
- Middleware for authentication and CORS
- Error handling with structured responses

**Key Services:**
- `AuthService`: MetaMask signature verification and JWT management
- `TokenService`: Token creation and management
- `PresaleService`: Presale lifecycle management
- `ContractsClient`: Blockchain interaction layer

**API Endpoints:**
```
Authentication:
GET  /api/auth/nonce          - Generate nonce for MetaMask
POST /api/auth/login          - Authenticate with signature
POST /api/auth/verify         - Verify JWT token

Token Management:
POST /api/token/create        - Deploy new token
GET  /api/token/{address}     - Get token details
GET  /api/token/list          - List user's tokens

Presale Management:
POST /api/presale/create      - Create new presale
GET  /api/presale/{id}        - Get presale details
GET  /api/presale/list        - List user's presales
POST /api/presale/{id}/participate - Record participation

Public Endpoints:
GET  /api/public/presale/{id} - Public presale information
```

### 3. Smart Contracts Layer (Solidity)

**Contracts:**

1. **MyToken.sol** (ERC-20 Token)
   - Standard ERC-20 implementation using OpenZeppelin
   - Mintable and burnable functionality
   - Owner-based access control
   - Configurable name, symbol, and supply

2. **Presale.sol** (Presale Contract)
   - ETH-to-token exchange functionality
   - Soft cap and hard cap limits
   - Time-based deadlines
   - Automatic refunds if soft cap not reached
   - Rate-based token pricing

3. **LaunchpadFactory.sol** (Factory Contract)
   - Deploys new tokens and presales
   - Tracks all created contracts
   - User-to-contract mapping
   - Event emission for tracking

**Key Features:**
- Reentrancy protection
- Gas optimization
- Comprehensive error handling
- Event logging for transparency

### 4. Database Layer (PostgreSQL)

**Schema Design:**

```sql
-- Core entities
users (id, address, created_at)
tokens (id, address, name, symbol, total_supply, creator_address, tx_hash, created_at)
presales (id, address, token_address, creator_address, rate, soft_cap, hard_cap, deadline, active, finalized, created_at)
presale_participations (id, presale_id, participant_address, amount_eth, amount_tokens, tx_hash, created_at)

-- Indexes for performance
idx_tokens_creator ON tokens(creator_address)
idx_presales_creator ON presales(creator_address)
idx_presales_token ON presales(token_address)
idx_participations_presale ON presale_participations(presale_id)
```

**Features:**
- Automated migrations
- Connection pooling
- Transaction support
- Optimized queries with proper indexing

### 5. Development Infrastructure

**Local Development:**
- Hardhat local blockchain network
- Hot reload for frontend and backend
- Database migrations
- Environment configuration

**Docker Compose Setup:**
- PostgreSQL database service
- Redis cache service (optional)
- Hardhat blockchain service
- Backend API service
- Frontend web service

## Data Flow

### 1. User Authentication Flow

```
1. User connects MetaMask
2. Frontend requests nonce from backend
3. User signs authentication message
4. Frontend sends signature to backend
5. Backend verifies signature and issues JWT
6. JWT used for subsequent authenticated requests
```

### 2. Token Creation Flow

```
1. User submits token creation form
2. Frontend sends request to backend API
3. Backend validates request and user authentication
4. Backend calls smart contract factory (simulated in MVP)
5. Token deployment transaction broadcasted
6. Backend stores token metadata in database
7. Frontend displays success and token address
```

### 3. Presale Creation Flow

```
1. User selects token and configures presale parameters
2. Frontend validates form and sends to backend
3. Backend creates presale contract (simulated in MVP)
4. Presale metadata stored in database
5. Landing page URL generated and returned
6. User can share landing page for participation
```

### 4. Presale Participation Flow

```
1. Participant visits presale landing page
2. Page loads presale data from public API
3. User connects wallet and enters contribution amount
4. Transaction sent to presale contract (simulated)
5. Backend records participation in database
6. User receives confirmation of token allocation
```

## Security Considerations

### 1. Authentication Security
- Nonce-based signature verification prevents replay attacks
- JWT tokens with expiration times
- Signature validation using recover functionality
- Address-based user identification

### 2. Input Validation
- Comprehensive form validation on frontend and backend
- SQL injection prevention with parameterized queries
- XSS protection with proper output encoding
- CORS configuration for cross-origin requests

### 3. Smart Contract Security
- OpenZeppelin battle-tested contracts
- Reentrancy guards on critical functions
- Access control with owner permissions
- Emergency pause functionality

### 4. Infrastructure Security
- Environment variable configuration
- Database connection encryption
- Rate limiting on API endpoints
- Error handling without information leakage

## Scalability Considerations

### 1. Database Optimization
- Proper indexing strategy
- Connection pooling
- Query optimization
- Database replication for read scalability

### 2. API Performance
- Efficient HTTP routing
- Response caching strategies
- Connection reuse for blockchain calls
- Pagination for large datasets

### 3. Frontend Optimization
- Code splitting and lazy loading
- Static generation for landing pages
- CDN deployment for assets
- Progressive Web App features

## Monitoring and Observability

### 1. Application Metrics
- API response times and error rates
- Database query performance
- Smart contract gas usage
- User authentication success rates

### 2. Business Metrics
- Token creation volume
- Presale participation rates
- Platform usage analytics
- Revenue tracking

### 3. Infrastructure Monitoring
- Server resource utilization
- Database performance metrics
- Network latency and throughput
- Container health status

## Deployment Strategy

### 1. Local Development
- Docker Compose for full stack
- Hot reload for rapid development
- Local blockchain with test accounts
- Database migrations and seeding

### 2. Staging Environment
- Testnet deployment for smart contracts
- Production-like infrastructure
- End-to-end testing
- Performance validation

### 3. Production Deployment
- Mainnet smart contract deployment
- Load balancer configuration
- Database replication and backup
- CDN and caching layer
- Monitoring and alerting setup

## Future Enhancements

### 1. Technical Improvements
- Multi-chain support (Polygon, BSC, Arbitrum)
- Advanced smart contract features
- Mobile app development
- API rate limiting and quotas

### 2. Feature Additions
- Vesting schedules for presales
- Whitelist functionality
- Token staking mechanisms
- Governance token integration

### 3. Business Features
- Platform fees and revenue sharing
- Advanced analytics dashboard
- Marketing tools for creators
- Community features and social proof