package websocket

import (
	"Blockchain-PriceOracle/internal/oracle"
	"context"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func StartEventListener(pricePublisher PricePublisher, mgr *ClientManager, oracleAddr, wsURL string) {
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

	eventSig := contractAbi.Events["PriceUpdated"].ID

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vlog := <-logs:
			//fmt.Printf("Event: tx=%s, topics=%d\n", vlog.TxHash.Hex(), len(vlog.Topics))
			if vlog.Topics[0] != eventSig {
				continue
			}

			eventData := struct {
				Symbol   string   `json:"symbol"`
				OldPrice *big.Int `json:"oldPrice"`
				NewPrice *big.Int `json:"newPrice"`
			}{}

			err := contractAbi.UnpackIntoInterface(&eventData, "PriceUpdated", vlog.Data)

			if err != nil {
				log.Printf("Parse error: %v", err)
				continue
			}

			log.Printf("SOCKET: OnChain price is CHANGED: %s → %s → %s",
				eventData.Symbol,
				eventData.OldPrice.String(),
				eventData.NewPrice.String(),
			)

			eventTime := time.Unix(int64(vlog.BlockTimestamp), 0)
			PublishPriceUpdate(pricePublisher, mgr, eventData.Symbol, eventData.NewPrice.String(), eventTime)
		}
	}
}
