package coingecko

import "strings"

const (
	BaseURL    = "https://api.coingecko.com/api/v3"
	PricePath  = "/simple/price"
	VsCurrency = "usd"
)

// Symbol mapping (CoinGecko ID -> Contract symbol)
var SupportedSymbols = map[string]string{
	"bitcoin":  "BTC",
	"ethereum": "ETH",
	// "solana":   "SOL",
}

func GetCoinGeckoIDs() []string {
	ids := make([]string, 0, len(SupportedSymbols))
	for cgId := range SupportedSymbols {
		ids = append(ids, cgId)
	}
	return ids
}

func BuildPriceURL() string {
	ids := GetCoinGeckoIDs()

	return BaseURL + PricePath +
		"?ids=" + strings.Join(ids, ",") +
		"&vs_currencies=" + VsCurrency
}
