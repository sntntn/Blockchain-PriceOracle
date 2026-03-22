package oracle

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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

func (c *Client) SetPrice(symbol string, newPrice *big.Int) (common.Hash, error) {
	config := MustLoadConfig()

	privateKey, err := crypto.HexToECDSA(config.PrivateKey)
	if err != nil {
		return common.Hash{}, fmt.Errorf("invalid private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Hash{}, fmt.Errorf("cannot cast public key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := c.rpc.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return common.Hash{}, fmt.Errorf("nonce error: %w", err)
	}

	gasPrice, err := c.rpc.SuggestGasPrice(context.Background())
	if err != nil {
		return common.Hash{}, fmt.Errorf("gas price error: %w", err)
	}

	// ABI encode
	data, err := c.contractABI.Pack("set", symbol, newPrice)
	if err != nil {
		return common.Hash{}, fmt.Errorf("pack set: %w", err)
	}

	tx := types.NewTransaction(
		nonce,
		c.addr,
		big.NewInt(0), // value (ETH)
		300000,        // gas limit
		gasPrice,
		data,
	)

	chainID, err := c.rpc.NetworkID(context.Background())
	if err != nil {
		return common.Hash{}, fmt.Errorf("chainID error: %w", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return common.Hash{}, fmt.Errorf("sign tx error: %w", err)
	}

	err = c.rpc.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("send tx error: %w", err)
	}

	log.Printf("TX SENT: %s", signedTx.Hash().Hex())
	return signedTx.Hash(), nil
}
