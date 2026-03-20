package api

import (
	"Blockchain-PriceOracle/internal/coingecko"
	"Blockchain-PriceOracle/internal/oracle"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PriceResponse struct {
	Symbol       string `json:"symbol"`
	OnChainPrice string `json:"onchain_price"`
	Chainlink    string `json:"chainlink_price"` // I added this
	CoinGeckoRaw string `json:"coingecko_raw"`
}

type PriceHistory struct {
	Timestamp    string `json:"timestamp"`
	OnChainPrice string `json:"onchain_price"`
}

type HistoryResponse struct {
	Symbol string         `json:"symbol"`
	Data   []PriceHistory `json:"data"`
}

func GetPricesHandler(c *gin.Context) {
	symbol := c.Param("symbol")

	oracleClient := oracle.GetOracleClient()

	onChainPrice, err := oracleClient.GetOnChainPrice(symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch onchain price: " + err.Error(),
		})
		return
	}

	chainlinkPrice, err := oracleClient.GetChainlinkPrice(symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch Chainlink price: " + err.Error(),
		})
		return
	}

	cgPrices, err := coingecko.FetchPrices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch CoinGecko price: " + err.Error(),
		})
		return
	}

	cgPrice, exists := cgPrices[symbol]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Symbol " + symbol + " not supported by CoinGecko",
		})
		return
	}

	response := PriceResponse{
		Symbol:       symbol,
		OnChainPrice: onChainPrice.String(),
		Chainlink:    chainlinkPrice.String(),
		CoinGeckoRaw: cgPrice.String(),
	}

	c.JSON(http.StatusOK, response)
}

// TO DO - real data
func GetPriceHistoryHandler(c *gin.Context) {
	symbol := c.Param("symbol")

	// Mock data - BTC last 7 days
	mockHistory := []PriceHistory{
		{"2026-03-13 10:00", "7200000000000"},
		{"2026-03-14 10:00", "7250000000000"},
		{"2026-03-15 10:00", "7100000000000"},
		{"2026-03-16 10:00", "7350000000000"},
		{"2026-03-17 10:00", "7280000000000"},
		{"2026-03-18 10:00", "7320000000000"},
		{"2026-03-19 10:00", "7300000000000"},
	}

	response := HistoryResponse{
		Symbol: symbol,
		Data:   mockHistory,
	}

	c.JSON(http.StatusOK, response)
}
