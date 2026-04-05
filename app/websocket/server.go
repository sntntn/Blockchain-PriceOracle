package websocket

import (
	"Blockchain-PriceOracle/app/history"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

type ClientManager struct {
	mu      sync.RWMutex
	clients map[*websocket.Conn]bool
}

var (
	clientMgr *ClientManager
	once      sync.Once
)

func GetClientManager() *ClientManager {
	once.Do(func() {
		clientMgr = &ClientManager{
			clients: make(map[*websocket.Conn]bool),
		}
	})
	return clientMgr
}

func SetupWebSocket(r *gin.Engine) {
	r.GET("/ws", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		mgr := GetClientManager()
		mgr.mu.Lock()
		mgr.clients[ws] = true
		mgr.mu.Unlock()

		defer func() {
			mgr.mu.Lock()
			delete(mgr.clients, ws)
			mgr.mu.Unlock()
			ws.Close()
		}()

		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				break
			}
		}
	})
}

func PublishPriceUpdate(symbol, price string, timestamp time.Time) {
	history.GetPriceHistory().Add(symbol, price, timestamp)
	msg := gin.H{"event": "price_updated", "symbol": symbol, "price": price}

	mgr := GetClientManager()
	mgr.mu.RLock()
	defer mgr.mu.RUnlock()

	for client := range mgr.clients {
		err := client.WriteJSON(msg)
		if err != nil {
			client.Close()
		}
	}
}
