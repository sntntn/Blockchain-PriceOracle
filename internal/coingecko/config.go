package coingecko

import (
	"os"
	"strings"
)

const (
	BaseURL                     = "https://api.coingecko.com/api/v3"
	PricePath                   = "/simple/price"
	VsCurrency                  = "usd"
	CoinGeckoRateLimitPerMinute = 30
	CoinGeckoRateLimitBurst     = 30
)

// Symbol mapping (CoinGecko ID -> Contract symbol)
var SupportedSymbols = map[string]string{
	"bitcoin":   "BTC",
	"ethereum":  "ETH",
	"chainlink": "LINK",
	// add new symbol here
}

func GetCoinGeckoIDs() []string {
	ids := make([]string, 0, len(SupportedSymbols))
	for cgId := range SupportedSymbols {
		ids = append(ids, cgId)
	}
	return ids
}

func BuildPriceURL() (string, error) {
	ids := GetCoinGeckoIDs()

	apiKey := os.Getenv("COINGECKO_API_KEY")
	if apiKey == "" {
		panic("COINGECKO_API_KEY not set")
	}

	url := BaseURL + PricePath +
		"?ids=" + strings.Join(ids, ",") +
		"&vs_currencies=" + VsCurrency +
		"&x_cg_demo_api_key=" + apiKey

	return url, nil
}
