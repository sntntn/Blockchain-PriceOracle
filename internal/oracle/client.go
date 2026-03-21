package oracle

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	oracleClient *Client
	once         sync.Once
)

func GetOracleClient() *Client {
	once.Do(func() {
		var err error
		oracleClient, err = NewClient()
		if err != nil {
			log.Fatalf("Oracle client init: %v", err)
		}
	})
	return oracleClient
}

type Client struct {
	rpc         *ethclient.Client
	addr        common.Address
	contractABI *abi.ABI
}

func NewClient() (*Client, error) {

	config := MustLoadConfig()

	conn, err := ethclient.Dial(config.SepoliaRPC)
	if err != nil {
		return nil, fmt.Errorf("RPC dial: %w", err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(ContractABI))
	if err != nil {
		return nil, fmt.Errorf("ABI parse: %w", err)
	}

	return &Client{
		rpc:         conn,
		addr:        common.HexToAddress(config.ContractAddr),
		contractABI: &contractAbi,
	}, nil
}

func (c *Client) GetOnChainPrice(symbol string) (*big.Int, error) {
	data, err := c.contractABI.Pack("getPrice", symbol)
	if err != nil {
		return nil, fmt.Errorf("pack getPrice: %w", err)
	}

	result, err := c.rpc.CallContract(context.Background(), ethereum.CallMsg{
		To:   &c.addr,
		Data: data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("call getPrice: %w", err)
	}

	var price big.Int
	price.SetBytes(result)
	return &price, nil
}

func (c *Client) GetChainlinkPrice(symbol string) (*big.Int, error) {
	data, err := c.contractABI.Pack("getChainlinkPrice", symbol)
	if err != nil {
		return nil, fmt.Errorf("pack getChainlinkPrice: %w", err)
	}

	result, err := c.rpc.CallContract(context.Background(), ethereum.CallMsg{
		To:   &c.addr,
		Data: data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("call getChainlinkPrice: %w", err)
	}

	var price big.Int
	price.SetBytes(result)
	return &price, nil
}
