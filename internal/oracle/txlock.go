package oracle

import "sync"

type TxLock struct {
	mu   sync.RWMutex
	data map[string]bool // symbol -> locked ( [ETH]=true -> there is pending tx for ETH)
}

var (
	txLock     *TxLock
	txLockOnce sync.Once
)

func GetTxLock() *TxLock {
	txLockOnce.Do(func() {
		txLock = &TxLock{
			data: make(map[string]bool),
		}
	})
	return txLock
}

// true = locked
func (l *TxLock) IsLocked(symbol string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.data[symbol]
}

func (l *TxLock) Lock(symbol string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.data[symbol] = true
}

func (l *TxLock) Unlock(symbol string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.data[symbol] = false
}
