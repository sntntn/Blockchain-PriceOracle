package main

import (
	"fmt"
	"log"
	"math/big"
	"time"

	"Blockchain-PriceOracle/app/utils"
	"Blockchain-PriceOracle/internal/coingecko"

	"github.com/joho/godotenv"
)

func printPrices(prices *map[string]*big.Int) {
	for symbol, price := range *prices {
		fmt.Printf("price -> %s: $%s\n", symbol, price.String())
	}
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	utils.InitOracleClient()
	utils.TestOracle()

	// -------------------------------------------
	fmt.Println("----------------------------")
	fmt.Println("CoinGecko Test - PRICE FETCH")

	ticker := time.NewTicker(1 * time.Minute)
	for {
		fmt.Println("\n=== NEW MINUTE ===")

		cgPrices, err := coingecko.FetchPrices()
		if err != nil {
			fmt.Printf("Fetch error: %v\n", err)
			// TO DO fallback prices
			<-ticker.C
			continue
		}

		printPrices(&cgPrices)

		for symbol, price := range cgPrices {
			if utils.CheckPriceCriteria(symbol, price) {
				fmt.Printf("%s - SEND TX NOW!\n", symbol)
			}
		}

		<-ticker.C
	}

}
