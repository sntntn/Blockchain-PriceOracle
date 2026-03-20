package api

import (
	"Blockchain-PriceOracle/internal/oracle"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PriceResponse struct {
	Symbol       string `json:"symbol"`
	OnChainPrice string `json:"onchain_price"`
	Chainlink    string `json:"chainlink_price"` // I added this
	// TO DO - add coingecko response
}

func GetPricesHandler(c *gin.Context) {
	symbol := c.Param("symbol")

	client := oracle.GetOracleClient()

	onChainPrice, err := client.GetOnChainPrice(symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch onchain price: " + err.Error(),
		})
		return
	}

	chainlinkPrice, err := client.GetChainlinkPrice(symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch Chainlink price: " + err.Error(),
		})
		return
	}

	response := PriceResponse{
		Symbol:       symbol,
		OnChainPrice: onChainPrice.String(),
		Chainlink:    chainlinkPrice.String(),
	}

	c.JSON(http.StatusOK, response)
}
