package automation

import (
	"Blockchain-PriceOracle/app/utils"
	"Blockchain-PriceOracle/app/websocket"
	"Blockchain-PriceOracle/internal/coingecko"
	"Blockchain-PriceOracle/internal/oracle"
	"fmt"
	"log"
	"os"
	"time"
)

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

		utils.PrintPrices(&cgPrices)

		for symbol, price := range cgPrices {
			// NOTE: returning Chainlink price captured at validation time;
			// used later for tx result analysis (reverted/confirmed context - task 8)
			ok, clPrice := utils.CheckPriceCriteria(symbol, price)
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
