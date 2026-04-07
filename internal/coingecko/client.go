package coingecko

import (
	"Blockchain-PriceOracle/internal/ratelimit"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"sync"
)

var (
	cgClient *Client
	cgOnce   sync.Once
)

func GetCoinGeckoClient(limiter ratelimit.Limiter) *Client {
	cgOnce.Do(func() {
		url := BuildPriceURL()
		cgClient = NewCgClient(limiter, url, fetchJSON)
	})
	return cgClient
}

type Client struct {
	limiter   ratelimit.Limiter
	url       string
	fetchJSON func(url string) (Response, error)
}

func NewCgClient(limiter ratelimit.Limiter, url string, fetchJSON func(url string) (Response, error)) *Client {
	return &Client{
		limiter:   limiter,
		url:       url,
		fetchJSON: fetchJSON,
	}
}

func (c *Client) FetchPrices() (map[string]*big.Int, error) {
	ok, delay, err := c.limiter.Allow()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("rate limit exceeded, try again in %v", delay)
	}

	cgResp, err := c.fetchJSON(c.url)
	if err != nil {
		return nil, err
	}

	return MapCGtoContract(cgResp), nil
}

func fetchJSON(url string) (Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP: %s", string(body))
	}

	var cgResp Response
	if err := json.Unmarshal(body, &cgResp); err != nil {
		return nil, fmt.Errorf("JSON: %w", err)
	}
	return cgResp, nil
}

func MapCGtoContract(cgResp Response) map[string]*big.Int {
	prices := make(map[string]*big.Int)

	for cgId, contractSymbol := range SupportedSymbols {

		if coinData, ok := cgResp[cgId]; ok {

			if usdPrice, ok := coinData["usd"]; ok {
				contractPrice := float64ToContract(usdPrice)
				prices[contractSymbol] = contractPrice
				log.Printf("%s: $%.2f -> %s (coingecko conversion)\n", contractSymbol, usdPrice, contractPrice.String())
			} else {
				log.Printf("DEBUG: %s USD price not found\n", contractSymbol)
			}

		} else {
			log.Printf("DEBUG: %s not in response\n", cgId)
		}
	}
	return prices
}

func float64ToContract(price float64) *big.Int {
	scaledPrice := big.NewFloat(price)
	scaledPrice.Mul(scaledPrice, big.NewFloat(1e8))
	result := new(big.Int)
	scaledPrice.Int(result)
	return result
}
