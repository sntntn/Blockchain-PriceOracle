package automation

import (
	"Blockchain-PriceOracle/app/criteria"
	"Blockchain-PriceOracle/app/history"
	"Blockchain-PriceOracle/app/websocket"
	"Blockchain-PriceOracle/internal/coingecko"
	"Blockchain-PriceOracle/internal/oracle"
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"
)

func Sync(oracleClient *oracle.Client, priceHistory *history.PriceHistory) {

	fromBlock := oracleClient.DeploymentBlock()
	currentLatestBlock, err := oracleClient.RPC().BlockNumber(context.Background())
	if err != nil {
		log.Printf("error on latest block fetch: %v", err)
		return
	}

	fromBlock = currentLatestBlock - 60 // TO DO - REMOVE THIS LINE IN PRODUCTION - rewritten just for development
	if err := priceHistory.ReverseSyncFromContract(oracleClient, fromBlock, currentLatestBlock); err != nil {
		log.Printf("Reverse Backfill Failed: %v", err)
	}

	fromBlock = currentLatestBlock
	if err := priceHistory.ForwardSyncFromContract(oracleClient, fromBlock); err != nil {
		log.Printf("Forward Sync Failed: %v", err)
	}

}

func CoinGeckoLoop(oracleClient *oracle.Client, cgClient *coingecko.Client) {
	time.Sleep(3 * time.Second) //to separate these logs from logs at app start

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	fmt.Println("----------------------------")
	log.Println("CoinGecko loop: every 1min (background)")
	fmt.Println("1min intervals started - CoinGecko")

	for {
		fmt.Println("\n=== NEW MINUTE ===")

		cgPrices, err := cgClient.FetchPrices()
		if err != nil {
			log.Printf("Fetch error: %v\n", err)
			<-ticker.C
			continue
		}

		PrintPrices(&cgPrices)

		for symbol, price := range cgPrices {
			// NOTE: returning Chainlink price captured at validation time;
			// used later for tx result analysis (reverted/confirmed context - task 8)
			ok, clPrice := criteria.CheckPriceCriteria(oracleClient, symbol, price)
			if ok {
				log.Printf("%s - PASSED THE CRITERIA - SEND TX NOW!\n", symbol)

				err := oracleClient.SetPrice(symbol, price, clPrice)
				if err != nil {
					log.Printf("TX FAILED %s: %v", symbol, err)
				}
			}
		}

		<-ticker.C
	}
}

func StartEthereumListener(pricePublisher websocket.PricePublisher, mgr *websocket.ClientManager) {
	wssURL := os.Getenv("WSS_URL")
	contractAddr := os.Getenv("CONTRACT_ADDR")

	if contractAddr != "" && wssURL != "" {
		go func() {
			log.Printf("Starting Ethereum listener: %s", wssURL)
			websocket.StartEventListener(pricePublisher, mgr, contractAddr, wssURL)
		}()
	} else {
		log.Println("Skipping Ethereum listener - missing ORACLE_CONTRACT_ADDR or ethereum WSS_URL")
	}
}

func PrintPrices(prices *map[string]*big.Int) {
	for symbol, price := range *prices {
		fmt.Printf("price -> %s: $%s\n", symbol, price.String())
	}
}
