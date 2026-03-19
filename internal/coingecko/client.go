package coingecko

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchPrices() (map[string]float64, error) {
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

	var cgResp Response
	if err := json.Unmarshal(body, &cgResp); err != nil {
		return nil, fmt.Errorf("JSON: %w", err)
	}
	return cgResp, nil
}

func MapCGtoContract(cgResp Response) map[string]float64 {
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
	return prices
}
