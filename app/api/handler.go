package api

import (
	"Blockchain-PriceOracle/app/utils"
	"Blockchain-PriceOracle/internal/coingecko"
	"Blockchain-PriceOracle/internal/oracle"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PriceResponse struct {
	Symbol       string `json:"symbol"`
	OnChainPrice string `json:"onchain_price"`
	Chainlink    string `json:"chainlink_price"` // I added this
	CoinGeckoRaw string `json:"coingecko_raw"`
}

type RangeRequest struct {
	From string `json:"from" binding:"required"` // "2026-03-21T15:00:00Z"
	To   string `json:"to" binding:"required"`   // "2026-03-22T17:00:00Z"
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

func GetPriceRangeHandler(c *gin.Context) {
	symbol := c.Param("symbol")

	var req RangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid from/to format"})
		return
	}

	// RFC3339 = "2026-03-21T15:43:07Z"
	from, err := time.Parse(time.RFC3339, req.From)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "From must be RFC3339"})
		return
	}

	to, err := time.Parse(time.RFC3339, req.To)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "To must be RFC3339"})
		return
	}

	prices := utils.GetPriceHistory().Range(symbol, from, to)

	c.JSON(http.StatusOK, gin.H{
		"symbol": symbol,
		"from":   from.Format(time.RFC3339),
		"to":     to.Format(time.RFC3339),
		"prices": prices,
		"count":  len(prices),
	})
}

func GetLastNHandler(c *gin.Context) {
	symbol := c.Param("symbol")
	n := 10 // Default
	if nStr := c.Query("n"); nStr != "" {
		var err error
		n, err = strconv.Atoi(nStr)
		if err != nil || n <= 0 || n > 1000 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "n must be 1-1000"})
			return
		}
	}

	prices := utils.GetPriceHistory().LastN(symbol, n)
	c.JSON(http.StatusOK, gin.H{
		"symbol": symbol,
		"n":      n,
		"prices": prices,
		"count":  len(prices),
	})
}
