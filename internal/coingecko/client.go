package coingecko

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func FetchPrices() (map[string]float64, error) {
	ids := make([]string, 0, len(SupportedSymbols))
	for cgId := range SupportedSymbols {
		ids = append(ids, cgId)
	}

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd",
		strings.Join(ids, ","))

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

	for cgId, contractSymbol := range SupportedSymbols {

		if coinData, ok := cgResp[cgId]; ok {

			if usdPrice, ok := coinData["usd"]; ok {
				prices[contractSymbol] = usdPrice
				fmt.Printf("%s: $%.2f\n", contractSymbol, usdPrice)
			} else {
				fmt.Printf("ERROR: %s USD price not found\n", contractSymbol)
			}

		} else {
			fmt.Printf("ERROR: %s not in response\n", cgId)
		}
	}

	return prices, nil
}
