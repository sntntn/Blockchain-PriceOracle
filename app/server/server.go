package server

import (
	"Blockchain-PriceOracle/app/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupServer() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/prices/:symbol", api.GetPricesHandler)
		v1.GET("/prices/:symbol/history", api.GetPriceHistoryHandler)
		v1.GET("/health", healthHandler)
	}

	return r
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":        "OK",
		"price_checker": "running",
	})
}
