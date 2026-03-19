package coingecko

// CoinGecko API response
//
//	 {
//	   "bitcoin":{"usd":65234},
//	   "ethereum":{"usd":2330},
//	    ...
//	}
type Response map[string]map[string]float64

// Simbol mapping (CoinGecko ID -> Contract simbol)
var SupportedSymbols = map[string]string{
	"bitcoin":  "BTC",
	"ethereum": "ETH",
	// "solana":   "SOL",
}
