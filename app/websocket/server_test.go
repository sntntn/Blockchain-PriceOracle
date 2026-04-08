package websocket

import (
	"fmt"
	"testing"
	"time"
)

type MockClient struct {
	written []interface{}
	closed  bool
	fail    bool
}

func (m *MockClient) WriteJSON(v interface{}) error {
	if m.fail {
		return fmt.Errorf("write error")
	}
	m.written = append(m.written, v)
	return nil
}

func (m *MockClient) Close() error {
	m.closed = true
	return nil
}

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
	mgr := NewClientManager()
	mock := &MockPublisher{}

	PublishPriceUpdate(mock, mgr, "BTC", "7200000000000", time.Now())

	if !mock.called {
		t.Errorf("expected Add to be called")
	}
}

func TestPublishPriceUpdate_PassesCorrectData(t *testing.T) {
	mgr := NewClientManager()

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

func TestPublishPriceUpdate_BroadcastsToAllClients(t *testing.T) {
	mgr := NewClientManager()

	c1 := &MockClient{}
	c2 := &MockClient{}

	mgr.clients[c1] = true
	mgr.clients[c2] = true

	mock := &MockPublisher{}

	PublishPriceUpdate(mock, mgr, "BTC", "100", time.Now())

	if len(c1.written) != 1 {
		t.Errorf("expected client1 to receive message")
	}

	if len(c2.written) != 1 {
		t.Errorf("expected client2 to receive message")
	}
}

func TestPublishPriceUpdate_RemovesFailedClient(t *testing.T) {
	mgr := NewClientManager()

	badClient := &MockClient{fail: true}
	goodClient := &MockClient{}

	mgr.clients[badClient] = true
	mgr.clients[goodClient] = true

	mock := &MockPublisher{}

	PublishPriceUpdate(mock, mgr, "BTC", "100", time.Now())

	if !badClient.closed {
		t.Errorf("expected bad client to be closed")
	}

	if len(goodClient.written) != 1 {
		t.Errorf("expected good client to receive message")
	}

	if _, exists := mgr.clients[badClient]; exists {
		t.Errorf("expected bad client to be removed from manager")
	}

	if _, exists := mgr.clients[goodClient]; !exists {
		t.Errorf("good client should still be in manager")
	}
}

func TestPublishPriceUpdate_LockCopySafety(t *testing.T) {
	mgr := NewClientManager()

	// 2 clients
	c1 := &MockClient{}
	c2 := &MockClient{}
	mgr.clients[c1] = true
	mgr.clients[c2] = true

	mock := &MockPublisher{}

	start := make(chan bool)
	done := make(chan bool)
	go func() {
		<-start
		PublishPriceUpdate(mock, mgr, "BTC", "500", time.Now())
		done <- true
	}()

	start <- true
	mgr.mu.Lock()
	delete(mgr.clients, c1)
	mgr.mu.Unlock()

	<-done

	if len(c2.written) != 1 {
		t.Errorf("expected c2 to receive message")
	}
}

func TestPublishPriceUpdate_NoClients_NoPanic(t *testing.T) {
	mgr := NewClientManager()
	mock := &MockPublisher{}

	PublishPriceUpdate(mock, mgr, "DOGE", "0.25", time.Now())
}
