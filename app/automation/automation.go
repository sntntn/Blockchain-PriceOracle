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
			// TO DO fallback prices
			<-ticker.C
			continue
		}

		utils.PrintPrices(&cgPrices)

		for symbol, price := range cgPrices {
			if utils.CheckPriceCriteria(symbol, price) {
				log.Printf("%s - PASSED THE CRITERIA - SEND TX NOW!\n", symbol)

				txHash, err := oracle.GetOracleClient().SetPrice(symbol, price)
				if err != nil {
					log.Printf("TX FAILED %s: %v", symbol, err)
				} else {
					log.Printf("TX SENT %s: %s", symbol, txHash.Hex())  //temporary - TO DO - see if TX is reverted
					utils.GetPriceHistory().Add(symbol, price.String()) // TO DO - if not reverted
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
