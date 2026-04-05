package automation

import (
	"Blockchain-PriceOracle/app/criteria"
	"Blockchain-PriceOracle/app/history"
	"Blockchain-PriceOracle/app/websocket"
	"Blockchain-PriceOracle/internal/coingecko"
	"Blockchain-PriceOracle/internal/oracle"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	history := history.GetPriceHistory()
	oracleClient := oracle.GetOracleClient()

	fromBlock := oracleClient.DeploymentBlock()
	currentLatestBlock := uint64(10588392)
	// currentLatestBlock, err := oracleClient.RPC().BlockNumber(context.Background())
	// if err != nil {
	// 	return fmt.Errorf("latest block: %w", err)
	// }

	if err := history.ReverseSyncFromContract(oracleClient, fromBlock, currentLatestBlock); err != nil {
		log.Printf("Reverse Backfill Failed: %v", err)
	}

	// fromBlock = currentLatestBlock
	fromBlock = uint64(10596423)
	if err := history.ForwardSyncFromContract(oracleClient, fromBlock); err != nil {
		log.Printf("Reverse Backfill Failed: %v", err)
	}

}

func CoinGeckoLoop() {
	time.Sleep(3 * time.Second) //to separate these logs from logs at app start

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	fmt.Println("----------------------------")
	log.Println("CoinGecko loop: every 1min (background)")
	fmt.Println("1min intervals started - CoinGecko")

	for {
		fmt.Println("\n=== NEW MINUTE ===")

		cgPrices, err := coingecko.FetchPrices()
		if err != nil {
			log.Printf("Fetch error: %v\n", err)
			<-ticker.C
			continue
		}

		PrintPrices(&cgPrices)

		for symbol, price := range cgPrices {
			// NOTE: returning Chainlink price captured at validation time;
			// used later for tx result analysis (reverted/confirmed context - task 8)
			ok, clPrice := criteria.CheckPriceCriteria(symbol, price)
			if ok {
				log.Printf("%s - PASSED THE CRITERIA - SEND TX NOW!\n", symbol)

				err := oracle.GetOracleClient().SetPrice(symbol, price, clPrice)
				if err != nil {
					log.Printf("TX FAILED %s: %v", symbol, err)
				}
			}
		}

		<-ticker.C
	}
}

func StartEthereumListener() {
	wssURL := os.Getenv("WSS_URL")
	contractAddr := os.Getenv("CONTRACT_ADDR")

	if contractAddr != "" && wssURL != "" {
		go func() {
			log.Printf("Starting Ethereum listener: %s", wssURL)
			websocket.StartEventListener(contractAddr, wssURL)
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
