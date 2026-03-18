package coingecko

// CoinGecko API response
type Response struct {
	Bitcoin  map[string]float64 `json:"bitcoin"`
	Ethereum map[string]float64 `json:"ethereum"`
}

// Simbol mapping (CoinGecko ID -> Contract simbol)
var SupportedSymbols = map[string]string{
	"bitcoin":  "BTC",
	"ethereum": "ETH",
}
