package coingecko

import (
	"math/big"
	"strings"
	"testing"
	"time"
)

type MockLimiter struct {
	allow bool
	delay time.Duration
	err   error
}

func (m *MockLimiter) Allow() (bool, time.Duration, error) {
	return m.allow, m.delay, m.err
}

func MockFetchJSON_1(url string) (Response, error) {
	return Response{
		"bitcoin": {"usd": 68243.00},
	}, nil
}

func MockFetchJSON_3(url string) (Response, error) {
	return Response{
		"bitcoin":   {"usd": 68243.00},
		"ethereum":  {"usd": 2088.85},
		"chainlink": {"usd": 8.71},
	}, nil
}
func TestFetchPrices_RateLimited(t *testing.T) {
	mockLimiter := &MockLimiter{
		allow: false,
		delay: 30 * time.Second,
	}
	mockURL := "http://example.com/mock"

	client := NewCgClient(mockLimiter, mockURL, MockFetchJSON_1)

	_, err := client.FetchPrices()
	if err == nil || !strings.Contains(err.Error(), "rate limit exceeded") {
		t.Fatalf("expected rate limit error, got %v", err)
	}
}

func TestFetchPrices_MapExpectedContractPrice(t *testing.T) {
	mockLimiter := &MockLimiter{
		allow: true,
		delay: 0,
	}

	mockURL := "http://example.com/mock"
	client := NewCgClient(mockLimiter, mockURL, MockFetchJSON_1)

	prices, err := client.FetchPrices()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	btcPrice, exists := prices["BTC"]
	if !exists {
		t.Fatalf("expected BTC price, got nothing")
	}

	expected := int64(6824300000000)
	if btcPrice.Cmp(big.NewInt(expected)) != 0 {
		t.Fatalf("unexpected BTC price: got %v, want %v", btcPrice, expected)
	}
}

func TestFetchPrices_Map3ExpectedContractPrice(t *testing.T) {
	mockLimiter := &MockLimiter{
		allow: true,
		delay: 0,
	}

	mockURL := "http://example.com/mock"
	client := NewCgClient(mockLimiter, mockURL, MockFetchJSON_3)

	prices, err := client.FetchPrices()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedPrices := map[string]int64{
		"BTC":  6824300000000,
		"ETH":  208885000000,
		"LINK": 871000000,
	}

	for symbol, expected := range expectedPrices {
		price, exists := prices[symbol]
		if !exists {
			t.Fatalf("expected %s price, got nothing", symbol)
		}
		if price.Cmp(big.NewInt(expected)) != 0 {
			t.Fatalf("unexpected %s price: got %v, want %v", symbol, price, expected)
		}
	}
}
