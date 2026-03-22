package oracle

import "sync"

const MAX_REVERT_HISTORY = 1000

type RevertHistory struct {
	mu   sync.RWMutex
	data []string
}

var (
	revertHistory *RevertHistory
	revertOnce    sync.Once
)

func GetRevertHistory() *RevertHistory {
	revertOnce.Do(func() {
		revertHistory = &RevertHistory{
			data: make([]string, 0),
		}
	})
	return revertHistory
}

func (h *RevertHistory) Add(entry string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// FIFO
	if len(h.data) >= MAX_REVERT_HISTORY {
		h.data = h.data[1:]
	}

	h.data = append(h.data, entry)
}

func (h *RevertHistory) All() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make([]string, len(h.data))
	copy(result, h.data)
	return result
}
