package oracle

import (
	"log"
	"strings"
)

var contractErrorMap = map[string]string{
	"UnsupportedSymbol": "Unsupported symbol",
	"IncompleteRound":   "Incomplete round",
	"InvalidPrice":      "Invalid price",
	"StalePrice":        "Stale price",
}

func handleContractError(symbol string, err error) {
	for key, msg := range contractErrorMap {
		if strings.Contains(err.Error(), key) {
			log.Printf("%s %s: %v", msg, symbol, err)
			return
		}
	}
}
