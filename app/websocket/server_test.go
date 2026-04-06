package websocket

import (
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

type MockPublisher struct {
	called    bool
	symbol    string
	price     string
	timestamp time.Time
}

func (m *MockPublisher) Add(symbol, price string, timestamp time.Time) {
	m.called = true
	m.symbol = symbol
	m.price = price
	m.timestamp = timestamp
}

func TestPublishPriceUpdate_CallsAdd(t *testing.T) {
	mgr := &ClientManager{
		clients: make(map[*websocket.Conn]bool),
	}
	mock := &MockPublisher{}

	PublishPriceUpdate(mock, mgr, "BTC", "7200000000000", time.Now())

	if !mock.called {
		t.Errorf("expected Add to be called")
	}
}

func TestPublishPriceUpdate_PassesCorrectData(t *testing.T) {
	mgr := &ClientManager{
		clients: make(map[*websocket.Conn]bool),
	}

	mock := &MockPublisher{}
	now := time.Now()

	PublishPriceUpdate(mock, mgr, "ETH", "200000000000", now)

	if mock.symbol != "ETH" {
		t.Errorf("expected symbol ETH, got %s", mock.symbol)
	}

	if mock.price != "200000000000" {
		t.Errorf("expected price 200, got %s", mock.price)
	}

	if !mock.timestamp.Equal(now) {
		t.Errorf("timestamp mismatch")
	}
}
