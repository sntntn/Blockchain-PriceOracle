package coingecko

// CoinGecko API response
//
//	 {
//	   "bitcoin":{"usd":65234},
//	   "ethereum":{"usd":2330},
//	    ...
//	}
type Response map[string]map[string]float64
