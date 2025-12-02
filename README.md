<p align="center">
  <img src="assets/vulcan_logo.png" alt="vulcan Logo" height="100">
</p>

<h1 align="center">Vulcan</h1>

<p align="center">
  <b>Mini Blockchain that Decentralized. Transparent. Secure.</b> 
</p>

<p align="center">
  <a href="https://github.com/OhMyDitzzy/vulcan/stargazers">
    <img src="https://img.shields.io/github/stars/OhMyDitzzy/vulcan?style=social" alt="GitHub Stars">
  </a>
  <a href="https://github.com/OhMyDitzzy/vulcan/network/members">
    <img src="https://img.shields.io/github/forks/OhMyDitzzy/vulcan?style=social" alt="GitHub Forks">
  </a>
</p>

---
> [!WARNING]
> This is a educational purposes project. Do not use in production without significant security hardening.


# Overview

A production-quality mini blockchain implementation with Proof-of-Work consensus, UTXO model, peer-to-peer networking, and a responsive web UI.

## Features

- ✅ Complete blockchain implementation with ECDSA signatures (secp256k1)
- ✅ Proof-of-Work consensus with adjustable difficulty
- ✅ UTXO (Unspent Transaction Output) model with full state management
- ✅ Transaction pool (mempool) with fee prioritization
- ✅ Merkle tree validation for blocks
- ✅ Peer-to-peer networking with gossip protocol
- ✅ Persistent storage using BadgerDB
- ✅ RESTful JSON API server
- ✅ React + TypeScript web UI with Tailwind CSS

## Prerequisites

- Go 1.20 or higher
- Node.js 18 or higher
- Make (optional, for convenience)

## Quick Start

### 1. Clone and Build

```bash
git clone https://github.com/OhMyDitzzy/vulcan.git
cd vulcan
make build
```

### 2. Run a Single Node

```bash
# Start the node with API server
./bin/vulcan --api-port=8080 --port=6000 --db-path=./data/node1

# In another terminal, start the frontend
cd vulcan-web
npm install
npm run dev
```

Visit http://localhost:5173 to access the web UI.

### 3. Using Docker

```bash
# Build and run everything
docker-compose up --build

# Access the UI at http://localhost:5173
# API available at http://localhost:8080
```

## Usage Examples

### Create a Wallet

```bash
curl "http://localhost:8080/wallet/new?consent=true"
```

Response:
```json
{
  "address": "04a1b2c3...",
  "private_key": "e8f7a6b5...",
  "warning": "NEVER share your private key!"
}
```

### Check Balance

```bash
curl "http://localhost:8080/balance/04a1b2c3..."
```

### Create and Sign Transaction

```bash
# Create transaction payload
cat > tx.json <<EOF
{
  "from": "04a1b2c3...",
  "to": "04d4e5f6...",
  "amount": 100,
  "fee": 10
}
EOF

# Sign it
curl -X POST http://localhost:8080/wallet/sign \
  -H "Content-Type: application/json" \
  -d '{
    "private_key": "e8f7a6b5...",
    "transaction": {
      "from": "04a1b2c3...",
      "to": "04d4e5f6...",
      "amount": 100,
      "fee": 10
    }
  }'
```

### Broadcast Transaction

```bash
curl -X POST http://localhost:8080/tx \
  -H "Content-Type: application/json" \
  -d @signed_tx.json
```

### Mine a Block

```bash
curl -X POST http://localhost:8080/mine \
  -H "Content-Type: application/json" \
  -d '{"miner_address": "04a1b2c3..."}'
```

### View Blockchain

```bash
# Get latest blocks
curl "http://localhost:8080/blockchain/blocks?start=0&limit=10"

# Get specific block
curl "http://localhost:8080/blockchain/block/00000abc..."

# Get transaction
curl "http://localhost:8080/blockchain/tx/abc123..."
```

### Manage Peers

```bash
# List peers
curl http://localhost:8080/peers

# Add peer
curl -X POST http://localhost:8080/peers \
  -H "Content-Type: application/json" \
  -d '{"address": "localhost:6001"}'
```

### View Metrics

```bash
curl http://localhost:8080/metrics
```

## Running Multiple Nodes

We've provided a script to easily test peer-to-peer networking:

```bash
# Start three interconnected nodes
./scripts/start_nodes.sh

# Node 1: API on 8080, P2P on 6000
# Node 2: API on 8081, P2P on 6001
# Node 3: API on 8082, P2P on 6002
```

Now broadcast a transaction to one node and watch it propagate:

```bash
# Create and broadcast to node 1
curl -X POST http://localhost:8080/tx -d @signed_tx.json

# Check mempool on node 2
curl http://localhost:8081/mempool

# Mine on node 3
curl -X POST http://localhost:8082/mine -d '{"miner_address": "04abc..."}'

# Verify block propagated to node 1
curl http://localhost:8080/blockchain/blocks
```

## Development

### Run Tests

```bash
# Backend tests
make test

# Frontend tests
cd vulcan-web
npm test

# E2E tests
cd vulcan-web
npm run test:e2e
```

### Run with Auto-Reload

```bash
# Backend (using air or similar)
air

# Frontend
cd vulcan-web
npm run dev
```

### Linting

```bash
# Go
golangci-lint run

# TypeScript
cd vulcan-web
npm run lint
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/blockchain/blocks` | List blocks (paginated) |
| GET | `/blockchain/block/:hash` | Get block by hash |
| GET | `/blockchain/tx/:txid` | Get transaction by ID |
| GET | `/wallet/new` | Create new wallet (requires `?consent=true`) |
| POST | `/wallet/sign` | Sign transaction with private key |
| POST | `/tx` | Broadcast signed transaction |
| GET | `/mempool` | List pending transactions |
| POST | `/mine` | Trigger mining |
| GET | `/balance/:address` | Get address balance and UTXOs |
| GET | `/peers` | List connected peers |
| POST | `/peers` | Add new peer |
| GET | `/metrics` | Prometheus metrics |

## Configuration

Configuration can be provided via CLI flags or environment variables:

| Flag | Environment Variable | Default | Description |
|------|---------------------|---------|-------------|
| `--api-port` | `API_PORT` | `8080` | API server port |
| `--port` | `P2P_PORT` | `6000` | P2P network port |
| `--db-path` | `DB_PATH` | `./data` | Database directory |
| `--peers` | `BOOTSTRAP_PEERS` | `` | Comma-separated peer addresses |
| `--mining` | `ENABLE_MINING` | `false` | Enable automatic mining |
| `--miner-address` | `MINER_ADDRESS` | `` | Address for mining rewards |
| `--difficulty` | `DIFFICULTY` | `4` | PoW difficulty (leading zeros) |

## Architecture

### Data Flow

1. **Transaction Creation**: User creates and signs transaction using wallet
2. **Broadcast**: Transaction submitted to API, validated, added to mempool
3. **Propagation**: Transaction gossiped to all connected peers
4. **Mining**: Miner selects transactions from mempool, creates block, solves PoW
5. **Validation**: Block validated by all nodes (PoW, transactions, UTXO state)
6. **Consensus**: Nodes accept valid blocks, update UTXO state
7. **Persistence**: Block and state changes written to BadgerDB

## Testing

### Unit Tests

```bash
# Run all Go tests
go test ./... -v -cover

# Run specific package
go test ./core -v

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Integration Tests

```bash
# API integration tests
go test ./api -tags=integration -v
```

### E2E Tests

```bash
cd vulcan-web
npm run test:e2e
```

## Performance

Typical performance on modern hardware:

- **Block validation**: ~1-5ms per block
- **Transaction validation**: ~0.5-1ms per transaction
- **PoW mining**: Varies by difficulty (4 leading zeros: ~1-10 seconds)
- **Merkle tree computation**: ~0.1ms for 1000 transactions
- **P2P broadcast**: ~50-200ms to propagate to all peers

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see LICENSE file for details

## Acknowledgments

- Inspired by Bitcoin and Ethereum architectures
- Uses secp256k1 elliptic curve cryptography
- Built with Go and React for Web Interfaces

## Roadmap

Future enhancements we're considering:
- [ ] Smart contract support
- [ ] Proof-of-Stake consensus option
- [ ] Light client implementation
- [ ] Mobile wallet app
- [ ] Enhanced privacy features (zero-knowledge proofs)
- [ ] Cross-chain bridges
- [ ] Sharding for scalability