# Blockchain Price Oracle


[![Go](https://img.shields.io/badge/Go-1.22-blue.svg)](https://golang.org)
[![Ethereum](https://img.shields.io/badge/Ethereum-Sepolia-orange.svg)](https://sepolia.etherscan.io)
[![Contract](https://img.shields.io/badge/Contract-Deployed-blueviolet.svg)](https://sepolia.etherscan.io/address/0xFf0fCE651EB5C7C147Af111f89dCB25AB57407e8#code)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)


**Real-time cryptocurrency price oracle**


## **What it does**

1. **Fetches** BTC/ETH prices from CoinGecko every minute  
2. **Updates** smart contract when price differs >2% from on-chain **AND** <20% from Chainlink
3. **Streams** contract price changes (all sources) via WebSocket 
4. **Serves** current prices, time-range history, history of reverts
5. **Zero reverts** - prevents gas waste by pre-TX validation and by waiting previously sent TX for chosen symbol to be completed 
6. **Multi-instance ready** - easy to add distributed lock coordination to keep contract reverts at zero when multiple instances are running at the same time
7. **Extensible** - easy to add new symbols (BTC/ETH/SOL/...) when contract supports it


##  Setup

Create `.env` file in root directory

### ENV
``` env
ANKR_SEPOLIA_API_KEY=your_ankr_key
WSS_URL=wss://eth-sepolia.g.alchemy.com/v2/your_key
CONTRACT_ADDR=0xFf0fCE651EB5C7C147Af111f89dCB25AB57407e8
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
- Ankr Sepolia
- CoinGecko API
- Chainlink Price Feeds
- Storage: In-memory FIFO