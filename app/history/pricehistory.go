package history

import (
	"container/list"
	"sync"
	"time"
)

const MAX_HISTORY_SIZE = 1000 //per symbol

type PricePoint struct {
	Timestamp time.Time `json:"timestamp"`
	Price     string    `json:"price"` // string - simplify
}

type PriceHistory struct {
	mu   sync.RWMutex
	data map[string]*list.List // BTC -> [p1, p2, p3...] sorted by time
}

var (
	priceHistory *PriceHistory
	once         sync.Once
)

func GetPriceHistory() *PriceHistory {
	once.Do(func() {
		priceHistory = &PriceHistory{
			data: make(map[string]*list.List),
		}
	})
	return priceHistory
}

func (h *PriceHistory) Add(symbol string, price string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	l, exists := h.data[symbol]
	if !exists {
		l = list.New()
		h.data[symbol] = l
	}

	// FIFO
	if l.Len() >= MAX_HISTORY_SIZE {
		l.Remove(l.Front())
	}

	l.PushBack(PricePoint{
		Timestamp: time.Now(),
		Price:     price,
	})
}

func (h *PriceHistory) Range(symbol string, from, to time.Time) []PricePoint {
	h.mu.RLock()
	defer h.mu.RUnlock()

	l, exists := h.data[symbol]
	if !exists {
		return nil
	}

	var result []PricePoint
	for e := l.Front(); e != nil; e = e.Next() {
		p := e.Value.(PricePoint)
		if p.Timestamp.After(from) && p.Timestamp.Before(to) {
			result = append(result, p)
		}
	}
	return result
}

func (h *PriceHistory) LastN(symbol string, n int) []PricePoint {
	h.mu.RLock()
	defer h.mu.RUnlock()

	l, exists := h.data[symbol]
	if !exists {
		return nil
	}

	var result []PricePoint
	count := 0
	for e := l.Back(); e != nil && count < n; e = e.Prev() {
		result = append([]PricePoint{e.Value.(PricePoint)}, result...)
		count++
	}
	return result
}
