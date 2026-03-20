package automation

import (
	"Blockchain-PriceOracle/app/utils"
	"Blockchain-PriceOracle/internal/coingecko"
	"fmt"
	"log"
	"math/big"
	"time"
)

func CoinGeckoLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	log.Println("CoinGecko loop: every 1min (background)")
	fmt.Println("----------------------------")
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

		printPrices(&cgPrices)

		for symbol, price := range cgPrices {
			if utils.CheckPriceCriteria(symbol, price) {
				log.Printf("%s - SEND TX NOW!\n", symbol)
			}
		}

		<-ticker.C
	}
}

func printPrices(prices *map[string]*big.Int) {
	for symbol, price := range *prices {
		fmt.Printf("price -> %s: $%s\n", symbol, price.String())
	}
}
