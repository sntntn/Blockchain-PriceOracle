package coingecko

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
)

func FetchPrices() (map[string]*big.Int, error) {
	url := BuildPriceURL()

	cgResp, err := fetchJSON(url)
	if err != nil {
		return nil, err
	}

	return MapCGtoContract(cgResp), nil
}

func fetchJSON(url string) (Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP: %s", string(body))
	}

	var cgResp Response
	if err := json.Unmarshal(body, &cgResp); err != nil {
		return nil, fmt.Errorf("JSON: %w", err)
	}
	return cgResp, nil
}

func MapCGtoContract(cgResp Response) map[string]*big.Int {
	prices := make(map[string]*big.Int)

	for cgId, contractSymbol := range SupportedSymbols {

		if coinData, ok := cgResp[cgId]; ok {

			if usdPrice, ok := coinData["usd"]; ok {
				contractPrice := float64ToContract(usdPrice)
				prices[contractSymbol] = contractPrice
				fmt.Printf("%s: $%.2f -> %s (conversion)\n", contractSymbol, usdPrice, contractPrice.String())
			} else {
				fmt.Printf("ERROR: %s USD price not found\n", contractSymbol)
			}

		} else {
			fmt.Printf("ERROR: %s not in response\n", cgId)
		}
	}
	return prices
}

func float64ToContract(price float64) *big.Int {
	scaledPrice := big.NewFloat(price)
	scaledPrice.Mul(scaledPrice, big.NewFloat(1e8))
	result := new(big.Int)
	scaledPrice.Int(result)
	return result
}
