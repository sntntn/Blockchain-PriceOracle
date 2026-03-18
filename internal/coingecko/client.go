package coingecko

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchPrices() (map[string]float64, error) {
	url := "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin,ethereum&vs_currencies=usd"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}

	var cgResp Response
	if err := json.Unmarshal(body, &cgResp); err != nil {
		return nil, fmt.Errorf("JSON error: %w", err)
	}

	prices := make(map[string]float64)

	if price, ok := cgResp.Ethereum["usd"]; ok {
		prices["ETH"] = price
		//fmt.Printf("ETHEREUM: $%.2f\n", price)
	} else {
		fmt.Println("ETH price not found")
	}

	if price, ok := cgResp.Bitcoin["usd"]; ok {
		prices["BTC"] = price
		//fmt.Printf("BITCOIN: $%.2f\n", price)
	} else {
		fmt.Println("BTC price not found")
	}

	return prices, nil
}
