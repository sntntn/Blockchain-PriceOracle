package history

import (
	"container/list"
	"fmt"
	"testing"
	"time"
)

func testHistory() *PriceHistory {
	return &PriceHistory{
		data: make(map[string]*list.List),
	}
}

func TestAddAndLastN(t *testing.T) {
	h := testHistory()

	now := time.Now()

	h.Add("BTC", "100", now)
	h.Add("BTC", "200", now.Add(time.Second))
	h.Add("BTC", "300", now.Add(2*time.Second))

	res := h.LastN("BTC", 2)

	if len(res) != 2 {
		t.Fatalf("expected 2, got %d", len(res))
	}

	if res[0].Price != "200" || res[1].Price != "300" {
		t.Errorf("wrong order or values")
	}
}

func TestRangeInSorted(t *testing.T) {
	h := testHistory()

	base := time.Now()

	h.Add("BTC", "100", base)
	h.Add("BTC", "200", base.Add(1*time.Second)) // IN
	h.Add("BTC", "300", base.Add(2*time.Second)) // IN
	h.Add("BTC", "400", base.Add(3*time.Second))

	from := base
	to := base.Add(3 * time.Second)

	res := h.Range("BTC", from, to)

	if len(res) != 2 {
		t.Fatalf("expected 2 results, got %d", len(res))
	}

	if res[0].Price != "200" || res[1].Price != "300" {
		t.Errorf("unexpected values: %+v", res)
	}
}

func TestRangeInUnsorted(t *testing.T) {
	h := testHistory()

	base := time.Now()

	h.Add("BTC", "100", base)
	h.Add("BTC", "200", base.Add(1*time.Second)) // IN
	h.Add("BTC", "300", base.Add(2*time.Second)) // IN
	h.Add("BTC", "400", base.Add(4*time.Second))
	h.Add("BTC", "500", base.Add(3*time.Second)) // IN

	from := base
	to := base.Add(4 * time.Second)

	res := h.Range("BTC", from, to)

	if len(res) != 3 {
		t.Fatalf("expected 3 results, got %d", len(res))
	}

	if res[0].Price != "200" || res[1].Price != "300" || res[2].Price != "500" {
		t.Errorf("unexpected values: %+v", res)
	}
}

func TestAddFront(t *testing.T) {
	h := testHistory()

	base := time.Now()

	h.Add("BTC", "300", base.Add(2*time.Second))
	h.Add("BTC", "400", base.Add(3*time.Second))

	err := h.AddFront("BTC", "200", base.Add(1*time.Second))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	res := h.LastN("BTC", 3)

	if len(res) != 3 {
		t.Fatalf("expected 3, got %d", len(res))
	}

	if res[0].Price != "200" || res[1].Price != "300" || res[2].Price != "400" {
		t.Errorf("wrong order: %+v", res)
	}
}

func TestAddFront_Full(t *testing.T) {
	h := testHistory()

	base := time.Now()

	for i := 0; i < MAX_HISTORY_SIZE; i++ {
		h.Add("BTC", "100", base.Add(time.Duration(i)*time.Second))
	}

	err := h.AddFront("BTC", "999", base)

	if err == nil {
		t.Errorf("expected error when history is full")
	}
}

func TestBackfill_Add_and_AddFront(t *testing.T) {
	h := testHistory()

	base := time.Now()

	h.AddFront("BTC", "200", base.Add(1*time.Second))
	h.AddFront("BTC", "100", base)

	h.Add("BTC", "300", base.Add(2*time.Second))
	h.Add("BTC", "400", base.Add(3*time.Second))

	res := h.Range("BTC", base.Add(-time.Second), base.Add(4*time.Second))

	if len(res) != 4 {
		t.Fatalf("expected 4, got %d", len(res))
	}

	if res[0].Price != "100" ||
		res[1].Price != "200" ||
		res[2].Price != "300" ||
		res[3].Price != "400" {
		t.Errorf("wrong order: %+v", res)
	}
}

func TestAddFrontFill_ThenFIFO(t *testing.T) {
	h := testHistory()

	base := time.Now()

	var err error
	for i := 2000; i > 0; i-- {
		ts := base.Add(-time.Duration(i) * time.Second)
		price := fmt.Sprintf("%d", i)
		err = h.AddFront("BTC", price, ts)
		if err != nil {
			break
		}
	}

	if err == nil {
		t.Fatalf("expected error when exceeding max size")
	}

	res := h.LastN("BTC", MAX_HISTORY_SIZE)
	if len(res) != MAX_HISTORY_SIZE {
		t.Fatalf("expected %d elements, got %d", MAX_HISTORY_SIZE, len(res))
	}

	for i := 2001; i <= 2005; i++ {
		price := fmt.Sprintf("%d", i)
		ts := base.Add(time.Duration(i) * time.Second)

		h.Add("BTC", price, ts)
	}

	res = h.LastN("BTC", MAX_HISTORY_SIZE)

	if len(res) != MAX_HISTORY_SIZE {
		t.Fatalf("expected %d elements after FIFO, got %d", MAX_HISTORY_SIZE, len(res))
	}

	last := res[len(res)-6:]

	expectedLast := []string{"2000", "2001", "2002", "2003", "2004", "2005"}

	for i := 0; i < 5; i++ {
		if last[i].Price != expectedLast[i] {
			t.Errorf("expected last[%d]=%s, got %s", i, expectedLast[i], last[i].Price)
		}
	}

	first := res[:5]

	expectedFirst := []string{"1006", "1007", "1008", "1009", "1010"}

	for i := 0; i < 5; i++ {
		if first[i].Price != expectedFirst[i] {
			t.Errorf("expected first[%d]=%s, got %s", i, expectedFirst[i], first[i].Price)
		}
	}
}
