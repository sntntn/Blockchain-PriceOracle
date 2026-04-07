# Blockchain Price Oracle


[![Go](https://img.shields.io/badge/Go-1.22-blue.svg)](https://golang.org)
[![Ethereum](https://img.shields.io/badge/Ethereum-Sepolia-orange.svg)](https://sepolia.etherscan.io)
[![Contract](https://img.shields.io/badge/Contract-Deployed-blueviolet.svg)](https://sepolia.etherscan.io/address/0x0a7cf8518eca70cfe64ff9b28cb0d8162f31a41d#code)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)


**Real-time cryptocurrency price oracle**


## **What it does**

1. **Historical sync** - backfills price history from contract deployment block up to the current state
2. **Fetches** BTC/ETH/LINK prices from CoinGecko every minute  
3. **Updates** smart contract when price differs >2% from on-chain **AND** <20% from Chainlink
4. **Streams** contract price changes (all sources) via WebSocket 
5. **Serves** current prices, time-range history, history of reverts
6. **Zero reverts** - prevents gas waste by pre-TX validation and by waiting previously sent TX for chosen symbol to be completed 
7. **Multi-instance ready** - easy to add distributed lock coordination to keep contract reverts at zero when multiple instances are running at the same time
8. **Extensible** - easy to add new symbols (BTC/ETH/LINK/...) when contract supports it
9. **Rate limited RPC & API usage** - protects external calls (CoinGecko, Ankr RPC..) from hitting provider limits
10. **Isolated rate limits per usage** - separate limiters for app production mode and sync process (supports different API keys)


##  Setup

Create `.env` file in root directory

### ENV
``` env
ANKR_SEPOLIA_API_KEY=your_ankr_key
ANKR_SEPOLIA_SYNC_API_KEY=your_2_ankr_key
COINGECKO_API_KEY=your_coingecko_key
WSS_URL=wss://eth-sepolia.g.alchemy.com/v2/your_key
CONTRACT_ADDR=0x0A7cF8518Eca70cfe64Ff9B28Cb0D8162F31A41D
DEPLOYMENT_BLOCK=10588291
PRIVATE_KEY=0xabcdef...  # 64 hex chars (no 0x prefix)
```


## Run
``` bash
go mod tidy
go run app/main.go
```


## Provided API Endpoints

| Endpoint                          | Method | Description                    |
|-----------------------------------|--------|--------------------------------|
| `/api/v1/prices/:symbol`          | `GET`  | Current prices (on-chain+CL+CG)|
| `/api/v1/prices/:symbol/last`     | `GET`  | Last N prices (`?n=10`)        |
| `/api/v1/prices/:symbol/range`    | `POST` | Time range query for prices    |
| `/api/v1/reverts`                 | `GET`  | Revert history (`?n=100`)      |
| `/api/v1/health`                  | `GET`  | Server status                  |


## Tech Stack
- Go 1.25.0
- Gin  
- Ethereum Sepolia
- Solidity ^0.8.20
- go-ethereum
- Alchemy WSS
- WebSockets streaming
- Ankr Sepolia
- CoinGecko API
- Chainlink Price Feeds
- Storage: In-memory FIFO
- Rate Limiting
- Concurrency
- Testing (mocks / dependency injection)