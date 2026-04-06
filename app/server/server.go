package server

import (
	"Blockchain-PriceOracle/app/api"
	"Blockchain-PriceOracle/app/history"
	"Blockchain-PriceOracle/app/websocket"
	"Blockchain-PriceOracle/internal/oracle"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupServer(oracleClient *oracle.Client, priceHistory *history.PriceHistory, revertHistory *oracle.RevertHistory, mgr *websocket.ClientManager) *gin.Engine {
	r := gin.Default()

	h := &api.Handler{
		OracleClient: oracleClient,
		PriceHistory: priceHistory,
		Reverts:      revertHistory,
	}

	v1 := r.Group("/api/v1")
	{
		v1.GET("/prices/:symbol", h.GetPricesHandler)
		v1.POST("/prices/:symbol/range", h.GetPriceRangeHandler)
		v1.GET("/prices/:symbol/last", h.GetLastNHandler)
		v1.GET("/reverts", h.GetRevertsHandler)
		v1.GET("/health", healthHandler)
	}

	websocket.SetupWebSocket(r, mgr)

	return r
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":        "OK",
		"price_checker": "running",
	})
}
