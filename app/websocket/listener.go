package websocket

import (
	"Blockchain-PriceOracle/internal/oracle"
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
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

	contractAbi, err := abi.JSON(strings.NewReader(oracle.ContractABI))
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

			eventData := struct {
				OldPrice *big.Int `json:"oldPrice"`
				NewPrice *big.Int `json:"newPrice"`
			}{}

			err := contractAbi.UnpackIntoInterface(&eventData, "PriceUpdated", vlog.Data)

			if err != nil {
				log.Printf("Parse error: %v", err)
				continue
			}

			symbolBytes := vlog.Topics[1].Bytes()
			symbol := strings.TrimLeft(string(symbolBytes), "\x00")
			symbol = strings.TrimRight(symbol, "\x00")

			log.Printf("✅ REAL: %s → %s → %s",
				symbol, eventData.OldPrice.String(), eventData.NewPrice.String())

			PublishPriceUpdate(symbol, eventData.NewPrice.String())
		}
	}
}
