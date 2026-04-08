package websocket

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type PricePublisher interface {
	Add(symbol, price string, timestamp time.Time)
}
type WSClient interface {
	WriteJSON(v interface{}) error
	Close() error
}

type RealWSClient struct {
	conn *websocket.Conn
}

func (c *RealWSClient) WriteJSON(v interface{}) error {
	return c.conn.WriteJSON(v)
}

func (c *RealWSClient) Close() error {
	return c.conn.Close()
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

type ClientManager struct {
	mu      sync.RWMutex
	clients map[WSClient]bool
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[WSClient]bool),
	}
}

func SetupWebSocket(r *gin.Engine, mgr *ClientManager) {
	r.GET("/ws", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		client := &RealWSClient{conn: ws}

		mgr.mu.Lock()
		mgr.clients[client] = true
		mgr.mu.Unlock()

		defer func() {
			mgr.mu.Lock()
			delete(mgr.clients, client)
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

func PublishPriceUpdate(pricePublisher PricePublisher, mgr *ClientManager, symbol, price string, timestamp time.Time) {
	pricePublisher.Add(symbol, price, timestamp)
	msg := gin.H{"event": "price_updated", "symbol": symbol, "price": price}

	mgr.mu.RLock()
	// better to copy clients and don't keep lock during network I/O
	clients := make([]WSClient, 0, len(mgr.clients))
	for c := range mgr.clients {
		clients = append(clients, c)
	}
	mgr.mu.RUnlock()

	for _, client := range clients {
		err := client.WriteJSON(msg)
		if err != nil {
			client.Close()

			mgr.mu.Lock()
			delete(mgr.clients, client)
			mgr.mu.Unlock()
		}
	}
}
