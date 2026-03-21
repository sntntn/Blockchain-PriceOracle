package websocket

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func StartEventListener(oracleAddr, wsURL string) {
	client, err := ethclient.Dial(wsURL)
	if err != nil {
		log.Fatal(err)
	}

	addr := common.HexToAddress(oracleAddr)
	query := ethereum.FilterQuery{Addresses: []common.Address{addr}}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening -> %s (Alchemy WSS)", oracleAddr)

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vlog := <-logs:
			fmt.Printf("Event: tx=%s, topics=%d\n", vlog.TxHash.Hex(), len(vlog.Topics))

			// TO DO - Parse event (topic0 = event signature, topic1 = indexed symbol)
			if len(vlog.Topics) > 1 {
				symbolBytes := vlog.Topics[1].Bytes()                 // indexed symbol
				symbol := string(bytes.TrimLeft(symbolBytes, "\x00")) // BTC, ETH...

				// Mock newPrice iz Data (za sada)
				newPrice := big.NewInt(7500000000000).String()

				PublishPriceUpdate(symbol, newPrice)
				log.Printf("Published: %s → %s", symbol, newPrice)
			}
		}
	}
}
