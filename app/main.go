package main

import (
	"fmt"
	"time"

	"Blockchain-PriceOracle/internal/coingecko"
)

func printPrices(prices *map[string]float64) {
	for symbol, price := range *prices {
		fmt.Printf("%s: $%.2f\n", symbol, price)
	}
}

func main() {

	fmt.Println("CoinGecko Test - PRICE FETCH")

	ticker := time.NewTicker(1 * time.Minute)
	for {
		fmt.Println("\n=== NEW MINUTE ===")
		if prices, err := coingecko.FetchPrices(); err != nil {
			fmt.Printf("Fetch error: %v\n", err)
		} else {
			printPrices(&prices)
		}
		<-ticker.C
	}

}
