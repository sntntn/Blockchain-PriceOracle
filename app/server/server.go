package server

import (
	"Blockchain-PriceOracle/app/api"
	"Blockchain-PriceOracle/app/websocket"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupServer() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/prices/:symbol", api.GetPricesHandler)
		v1.POST("/prices/:symbol/range", api.GetPriceRangeHandler)
		v1.GET("/health", healthHandler)
	}

	websocket.SetupWebSocket(r)

	return r
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":        "OK",
		"price_checker": "running",
	})
}
