# 🚀 Launchpad - Token & Presale Platform

A complete MVP platform for creating ERC-20 tokens and launching presales with ease. Built with Next.js, Golang, Solidity, and PostgreSQL.

![Launchpad Homepage](https://github.com/user-attachments/assets/d8a9367f-d499-4e34-ac7b-d35437dfec44)

## 🌟 Features

- **Token Creation**: Deploy ERC-20 tokens with custom parameters using OpenZeppelin contracts
- **Presale Management**: Configure and launch token presales with customizable rates and caps
- **Landing Pages**: Automatic landing page generation for presale participants
- **Web3 Integration**: MetaMask authentication and wallet connectivity
- **Local Development**: Complete local development environment with Hardhat

## 🏗️ Architecture

### Components
- **Frontend**: Next.js with TypeScript and Tailwind CSS
- **Backend**: Golang with Gin framework
- **Smart Contracts**: Solidity with Hardhat and OpenZeppelin
- **Database**: PostgreSQL for metadata storage
- **Caching**: Redis (optional)
- **Development**: Docker Compose for local setup

### User Flow
1. **Connect Wallet**: User connects MetaMask wallet
2. **Authenticate**: Web3 signature-based authentication
3. **Create Token**: Deploy ERC-20 token with custom parameters
4. **Configure Presale**: Set rates, caps, and duration
5. **Launch**: Deploy presale contract and get landing page
6. **Share**: Share landing page URL with potential investors

## 🚀 Quick Start

### Prerequisites
- Node.js 18+
- Go 1.21+
- Docker & Docker Compose
- MetaMask browser extension

### Local Development

1. **Clone the repository**
```bash
git clone https://github.com/Wrestler094/launchpad.git
cd launchpad
```

2. **Start all services with Docker Compose**
```bash
docker-compose up -d
```

This will start:
- PostgreSQL database on port 5432
- Redis on port 6379
- Hardhat node on port 8545
- Golang backend on port 8080
- Next.js frontend on port 3000

3. **Or run services individually**

**Start Database**
```bash
docker-compose up postgres redis -d
```

**Start Hardhat Node**
```bash
cd contracts
npm install
npm run node
```

**Start Backend**
```bash
cd backend
go mod tidy
go run main.go
```

**Start Frontend**
```bash
cd frontend
npm install
npm run dev
```

### Environment Variables

Create `.env` files in the respective directories:

**Backend (.env)**
```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/launchpad?sslmode=disable
REDIS_URL=redis://localhost:6379
RPC_URL=http://localhost:8545
JWT_SECRET=your-secret-key
PORT=8080
```

**Frontend (.env.local)**
```
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_RPC_URL=http://localhost:8545
```

## 📁 Project Structure

```
launchpad/
├── contracts/              # Solidity smart contracts
│   ├── contracts/
│   │   ├── MyERC20.sol     # ERC-20 token contract
│   │   └── Presale.sol     # Presale contract
│   ├── ignition/modules/   # Deployment scripts
│   └── test/               # Contract tests
├── backend/                # Golang backend API
│   ├── config/             # Configuration
│   ├── handlers/           # HTTP handlers
│   ├── middleware/         # Auth middleware
│   ├── models/             # Data models
│   ├── services/           # Business logic
│   └── main.go             # Entry point
├── frontend/               # Next.js frontend
│   └── src/app/           # App router pages
├── database/              # Database schemas
│   └── init.sql           # PostgreSQL initialization
└── docker-compose.yml    # Docker services
```

## 🔧 API Endpoints

### Authentication
- `POST /api/auth/login` - Web3 signature authentication

### Token Management
- `POST /api/token/create` - Deploy new ERC-20 token

### Presale Management
- `POST /api/presale/create` - Create new presale
- `GET /api/presale/:id` - Get presale information
- `POST /api/presale/:id/participate` - Participate in presale

## 🧪 Smart Contracts

### MyERC20.sol
- Standard ERC-20 implementation using OpenZeppelin
- Mintable and burnable functionality
- Owner-controlled minting

### Presale.sol
- Secure presale contract with rate-based token distribution
- Hard cap and soft cap enforcement
- Automatic refunds if soft cap not reached
- Time-based presale duration

## 🔒 Security Features

- Web3 signature-based authentication
- JWT token management
- Input validation and sanitization
- Reentrancy protection in smart contracts
- OpenZeppelin security standards

## 🚧 Current Status

This is an MVP implementation with the following completed:
- ✅ Basic project structure and Docker setup
- ✅ Smart contracts (ERC-20 and Presale)
- ✅ Backend API with authentication
- ✅ Frontend UI with responsive design
- ✅ Database schema and migrations
- ⏳ Web3 integration (in progress)
- ⏳ Contract deployment automation
- ⏳ Presale participation flow

## 🛣️ Roadmap

### Phase 1 (Current)
- Complete Web3 wallet integration
- Implement contract deployment via API
- Add presale participation functionality

### Phase 2
- Add more token standards (ERC-721, ERC-1155)
- Multi-chain support
- Advanced presale configurations

### Phase 3
- Governance features
- Staking mechanisms
- Analytics dashboard

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📄 License

MIT License - see LICENSE file for details

## 🆘 Support

For questions or issues:
- Create an issue on GitHub
- Check the documentation
- Review the API endpoints

---

Built with ❤️ for the Web3 community