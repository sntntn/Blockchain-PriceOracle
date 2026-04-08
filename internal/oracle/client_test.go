package oracle

import (
	"context"
	"math/big"
	"strings"
	"testing"
	"time"

	"Blockchain-PriceOracle/internal/oracle/mocks"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
)

type MockLimiter struct {
	allow bool
	err   error
}

func (m *MockLimiter) Allow() (bool, time.Duration, error) {
	return m.allow, 0, m.err
}

func (m *MockLimiter) Wait(_ context.Context) error { return nil }

func mustParseABI(t *testing.T) *abi.ABI {
	const jsonABI = `[{"inputs":[{"type":"string"}],"name":"getPrices","outputs":[{"type":"uint256"},{"type":"uint256"}],"stateMutability":"view","type":"function"}]`
	abiObj, err := abi.JSON(strings.NewReader(jsonABI))
	if err != nil {
		t.Fatalf("ABI parse: %v", err)
	}
	return &abiObj
}

func TestGetPrices_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRPC := mocks.NewMockEthereumClient(ctrl)

	abiInst := mustParseABI(t)
	outputData, _ := abiInst.Pack("", big.NewInt(123), big.NewInt(456))

	mockRPC.EXPECT().
		CallContract(gomock.Any(), gomock.Any(), (*big.Int)(nil)).
		Return(outputData, nil)

	limiter := &MockLimiter{allow: true}
	client := &Client{
		rpc:         mockRPC,
		contractABI: abiInst,
		limiter:     limiter,
		addr:        common.HexToAddress("0x1234567890123456789012345678901234567890"),
	}

	onchain, chainlink, err := client.GetPrices("BTC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("onchain price: %s", onchain.String())
	t.Logf("chainlink price: %s", chainlink.String())
	t.Logf("onchain == 123? %v", onchain.Cmp(big.NewInt(123)) == 0)
	t.Logf("chainlink == 456? %v", chainlink.Cmp(big.NewInt(456)) == 0)

	if onchain.Cmp(big.NewInt(123)) != 0 {
		t.Errorf("want 123, got %v", onchain)
	}
	if chainlink.Cmp(big.NewInt(456)) != 0 {
		t.Errorf("want 456, got %v", chainlink)
	}
}
