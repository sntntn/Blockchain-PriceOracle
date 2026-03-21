package utils

import (
	"Blockchain-PriceOracle/internal/oracle"
	"fmt"
	"log"
	"math/big"
)

// func TestOracle() {
// 	fmt.Println("Oracle Test - Smart Contract PRICE")

// 	btcPrice, err := oracleClient.GetPrice("BTC")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("BTC Onchain: %s\n", btcPrice.String())

// 	btcCL, err := oracleClient.GetChainlinkPrice("BTC")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("BTC Chainlink: %s\n", btcCL.String())
// }

func CheckPriceCriteria(symbol string, newPrice *big.Int) bool {
	oracleClient := oracle.GetOracleClient()
	onChainPrice, err := oracleClient.GetOnChainPrice(symbol)
	if err != nil {
		log.Printf("GetPrice %s: %v", symbol, err)
		return false
	}

	clPrice, err := oracleClient.GetChainlinkPrice(symbol)
	if err != nil {
		log.Printf("GetChainlink %s: %v", symbol, err)
		return false
	}

	// 2% on chain check
	absDiffOnChain := absDiff(onChainPrice, newPrice)
	percentChain := percentBigInt(absDiffOnChain, onChainPrice) // just for print

	minChange := new(big.Int).Div(
		new(big.Int).Mul(onChainPrice, big.NewInt(2)),
		big.NewInt(100),
	)
	if absDiffOnChain.Cmp(minChange) <= 0 {
		fmt.Printf("%s: %.2f%% < 2%% (skip)\n", symbol, percentChain)
		return false
	}

	// 20% CL check
	absDiffCL := absDiff(clPrice, newPrice)
	percentCL := percentBigInt(absDiffCL, clPrice) // just for print

	maxChange := new(big.Int).Div(
		new(big.Int).Mul(clPrice, big.NewInt(20)),
		big.NewInt(100),
	)
	if absDiffCL.Cmp(maxChange) > 0 {
		fmt.Printf("%s: %.2f%% > 20%% CL (skip)\n", symbol, percentCL)
		return false
	}

	fmt.Printf("%s: APPROVED - difference (%.2f%% on Chain, %.2f%% CL)\n", symbol, percentChain, percentCL)
	return true
}

// just for print
func percentBigInt(numerator, denominator *big.Int) float64 {
	num := new(big.Float).SetInt(numerator)
	den := new(big.Float).SetInt(denominator)
	ratio := new(big.Float).Quo(num, den)
	percent := new(big.Float).Mul(ratio, big.NewFloat(100))
	result, _ := percent.Float64()
	return result
}

func absDiff(a, b *big.Int) *big.Int {
	if a.Cmp(b) >= 0 {
		return new(big.Int).Sub(a, b)
	}
	return new(big.Int).Sub(b, a)
}

func PrintPrices(prices *map[string]*big.Int) {
	for symbol, price := range *prices {
		fmt.Printf("price -> %s: $%s\n", symbol, price.String())
	}
}
