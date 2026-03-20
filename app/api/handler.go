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
