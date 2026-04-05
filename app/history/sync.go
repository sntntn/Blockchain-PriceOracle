package history

import (
	"Blockchain-PriceOracle/internal/coingecko"
	"Blockchain-PriceOracle/internal/oracle"
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const (
	REVERSE_BATCH_BLOCKS = 10
	FORWARD_BATCH_BLOCKS = 10
	BACKOFF_MS           = 100
)

func (h *PriceHistory) ReverseSyncFromContract(oracleClient *oracle.Client, stopBlock uint64, latestBlock uint64) error {

	log.Printf("REVERSE BACKFILL: %d -> %d", latestBlock, stopBlock)

	contractAbi, err := abi.JSON(strings.NewReader(oracle.ContractABI))
	if err != nil {
		return fmt.Errorf("ABI parse: %w", err)
	}

	eventSig := contractAbi.Events["PriceUpdated"].ID

	currentBlock := latestBlock
	totalProcessed := 0
	batchCount := 0
	fullSymbols := make(map[string]bool)

	for currentBlock >= stopBlock {
		fromBlock := currentBlock - REVERSE_BATCH_BLOCKS + 1
		if fromBlock < stopBlock {
			fromBlock = stopBlock
		}

		log.Printf("Reverse Batch #%d: %d -> %d (%d blocks)",
			batchCount, fromBlock, currentBlock, currentBlock-fromBlock+1)

		query := ethereum.FilterQuery{
			FromBlock: new(big.Int).SetUint64(fromBlock),
			ToBlock:   new(big.Int).SetUint64(currentBlock),
			Addresses: []common.Address{oracleClient.Address()},
			Topics:    [][]common.Hash{{eventSig}},
		}

		logs, err := oracleClient.RPC().FilterLogs(context.Background(), query)
		if err != nil {
			log.Printf("Reverse Batch #%d failed: %v", batchCount, err)
			time.Sleep(time.Duration(BACKOFF_MS) * time.Millisecond)
			currentBlock -= REVERSE_BATCH_BLOCKS
			batchCount++
			continue
		}

		batchProcessed := 0
		for i := len(logs) - 1; i >= 0; i-- {
			vlog := logs[i]

			if vlog.Topics[0] != eventSig {
				continue
			}

			eventData := struct {
				Symbol   string   `json:"symbol"`
				OldPrice *big.Int `json:"oldPrice"`
				NewPrice *big.Int `json:"newPrice"`
			}{}

			if err := contractAbi.UnpackIntoInterface(&eventData, "PriceUpdated", vlog.Data); err != nil {
				continue
			}

			eventTime := time.Unix(int64(vlog.BlockTimestamp), 0)

			if err := h.AddFront(eventData.Symbol, eventData.NewPrice.String(), eventTime); err != nil {
				if !fullSymbols[eventData.Symbol] {
					fullSymbols[eventData.Symbol] = true
				}
				continue
			}

			batchProcessed++
		}

		totalProcessed += batchProcessed
		batchCount++

		log.Printf("Reverse Batch #%d: %d events -> TOTAL: %d",
			batchCount-1, batchProcessed, totalProcessed)

		currentBlock -= REVERSE_BATCH_BLOCKS

		time.Sleep(time.Duration(BACKOFF_MS) * time.Millisecond)

		if len(fullSymbols) >= len(coingecko.SupportedSymbols) { // BTC, ETH, LINK ...
			log.Println("ALL symbols full - early exit!")
			break
		}
	}

	log.Printf("REVERSE BACKFILL COMPLETE: %d batches, %d events!", batchCount, totalProcessed)
	return nil
}

func (h *PriceHistory) ForwardSyncFromContract(oracleClient *oracle.Client, startBlock uint64) error {
	log.Printf("(LATEST) FORWARD SYNC: starting from block %d", startBlock)

	contractAbi, err := abi.JSON(strings.NewReader(oracle.ContractABI))
	if err != nil {
		return fmt.Errorf("ABI parse: %w", err)
	}

	eventSig := contractAbi.Events["PriceUpdated"].ID

	currentBlock := startBlock
	totalProcessed := 0
	batchCount := 0

	for {
		latestBlock, err := oracleClient.RPC().BlockNumber(context.Background())
		if err != nil {
			return fmt.Errorf("latest block: %w", err)
		}

		if currentBlock > latestBlock {
			break
		}

		toBlock := currentBlock + FORWARD_BATCH_BLOCKS
		if toBlock > latestBlock {
			toBlock = latestBlock
		}

		log.Printf("Forward Batch #%d: %d -> %d (%d blocks)",
			batchCount, currentBlock, toBlock, toBlock-currentBlock+1)

		query := ethereum.FilterQuery{
			FromBlock: new(big.Int).SetUint64(currentBlock),
			ToBlock:   new(big.Int).SetUint64(toBlock),
			Addresses: []common.Address{oracleClient.Address()},
			Topics:    [][]common.Hash{{eventSig}},
		}

		logs, err := oracleClient.RPC().FilterLogs(context.Background(), query)
		if err != nil {
			log.Printf("Forward Batch #%d failed: %v", batchCount, err)
			time.Sleep(time.Duration(BACKOFF_MS) * time.Millisecond)
			currentBlock = toBlock + 1
			batchCount++
			continue
		}

		batchProcessed := 0
		for _, vlog := range logs {
			if vlog.Topics[0] != eventSig {
				continue
			}

			eventData := struct {
				Symbol   string   `json:"symbol"`
				OldPrice *big.Int `json:"oldPrice"`
				NewPrice *big.Int `json:"newPrice"`
			}{}

			if err := contractAbi.UnpackIntoInterface(&eventData, "PriceUpdated", vlog.Data); err != nil {
				continue
			}

			eventTime := time.Unix(int64(vlog.BlockTimestamp), 0)

			h.Add(eventData.Symbol, eventData.NewPrice.String(), eventTime)
			batchProcessed++
		}

		totalProcessed += batchProcessed
		batchCount++

		log.Printf("Forward Batch #%d: %d events -> TOTAL: %d",
			batchCount-1, batchProcessed, totalProcessed)

		currentBlock = toBlock + 1
		time.Sleep(time.Duration(BACKOFF_MS) * time.Millisecond)
	}

	log.Printf("FORWARD SYNC COMPLETE: %d batches, %d events (caught up to %d)!",
		batchCount, totalProcessed, currentBlock-1)
	return nil
}
