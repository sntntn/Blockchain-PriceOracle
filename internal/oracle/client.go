package oracle

import (
	"Blockchain-PriceOracle/internal/ratelimit"
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	oracleClient *Client
	once         sync.Once
	initErr      error
)

func GetOracleClient(reverts *RevertHistory, limiter ratelimit.Limiter) (*Client, error) {
	once.Do(func() {
		config := LoadConfig()

		conn, err := ethclient.Dial(config.SepoliaRpc)
		if err != nil {
			initErr = fmt.Errorf("RPC dial: %w", err)
			return
		}
		connSync, err := ethclient.Dial(config.SepoliaRpcSync)
		if err != nil {
			initErr = fmt.Errorf("RPC SYNC dial: %w", err)
			return
		}

		contractAbi, err := abi.JSON(strings.NewReader(ContractABI))
		if err != nil {
			initErr = fmt.Errorf("ABI parse: %w", err)
			return
		}

		addr := common.HexToAddress(config.ContractAddr)

		oracleClient, err = NewOracleClient(
			reverts,
			limiter,
			addr,
			config.DeploymentBlock,
			&contractAbi,
			conn,
			connSync,
		)
		if err != nil {
			initErr = fmt.Errorf("Oracle client init: %w", err)
			return
		}
	})

	return oracleClient, initErr
}

type Client struct {
	rpc             *ethclient.Client
	rpcSync         *ethclient.Client
	addr            common.Address
	deploymentBlock uint64
	contractABI     *abi.ABI

	reverts RevertHistoryInterface
	limiter ratelimit.Limiter
}

func NewOracleClient(reverts RevertHistoryInterface, limiter ratelimit.Limiter, addr common.Address, deploymentBlock uint64, contractAbi *abi.ABI, conn *ethclient.Client, connSync *ethclient.Client) (*Client, error) {
	return &Client{
		rpc:             conn,
		rpcSync:         connSync,
		addr:            addr,
		deploymentBlock: deploymentBlock,
		contractABI:     contractAbi,
		reverts:         reverts,
		limiter:         limiter,
	}, nil
}

func (c *Client) DeploymentBlock() uint64 {
	return c.deploymentBlock
}

func (c *Client) Address() common.Address {
	return c.addr
}

func (c *Client) RpcSync() *ethclient.Client {
	return c.rpc
}

func (c *Client) GetPrices(symbol string) (onchainPrice, chainlinkPrice *big.Int, err error) {
	if err := c.allow(); err != nil {
		return nil, nil, err
	}

	data, err := c.contractABI.Pack("getPrices", symbol)
	if err != nil {
		return nil, nil, fmt.Errorf("pack: %w", err)
	}

	result, err := c.rpc.CallContract(context.Background(), ethereum.CallMsg{
		To:   &c.addr,
		Data: data,
	}, nil)
	if err != nil {
		if strings.Contains(err.Error(), "UnsupportedSymbol") {
			log.Printf("Unsupported symbol %s: %v", symbol, err)
		} else if strings.Contains(err.Error(), "IncompleteRound") {
			log.Printf("Incomplete round %s: %v", symbol, err)
		} else if strings.Contains(err.Error(), "InvalidPrice") {
			log.Printf("Invalid price %s: %v", symbol, err)
		} else if strings.Contains(err.Error(), "StalePrice") {
			log.Printf("Stale price %s: %v", symbol, err)
		}
		return nil, nil, err
	}

	unpacked, err := c.contractABI.Unpack("getPrices", result)
	if err != nil {
		return nil, nil, fmt.Errorf("unpack: %w", err)
	}

	if len(unpacked) != 2 {
		return nil, nil, fmt.Errorf("expected 2 values")
	}

	onchainPrice = unpacked[0].(*big.Int)
	chainlinkPrice = unpacked[1].(*big.Int)

	return onchainPrice, chainlinkPrice, nil
}

func (c *Client) SetPrice(symbol string, newPrice *big.Int, clPrice *big.Int) error {
	lock := GetTxLock()

	if lock.IsLocked(symbol) {
		log.Printf("SKIP sending TX for %s: waiting for previous %s tx to finish", symbol, symbol)
		return nil
	}
	lock.Lock(symbol)

	config := LoadConfig()

	privateKey, err := crypto.HexToECDSA(config.PrivateKey)
	if err != nil {
		lock.Unlock(symbol)
		return fmt.Errorf("invalid private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		lock.Unlock(symbol)
		return fmt.Errorf("cannot cast public key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	if err := c.allow(); err != nil {
		lock.Unlock(symbol)
		return err
	}
	nonce, err := c.rpc.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		lock.Unlock(symbol)
		return fmt.Errorf("nonce error: %w", err)
	}

	if err := c.allow(); err != nil {
		lock.Unlock(symbol)
		return err
	}
	gasPrice, err := c.rpc.SuggestGasPrice(context.Background())
	if err != nil {
		lock.Unlock(symbol)
		return fmt.Errorf("gas price error: %w", err)
	}

	// ABI encode
	data, err := c.contractABI.Pack("set", symbol, newPrice)
	if err != nil {
		lock.Unlock(symbol)
		return fmt.Errorf("pack set: %w", err)
	}

	tx := types.NewTransaction(
		nonce,
		c.addr,
		big.NewInt(0), // value (ETH)
		300000,        // gas limit
		gasPrice,
		data,
	)

	if err := c.allow(); err != nil {
		lock.Unlock(symbol)
		return err
	}
	chainID, err := c.rpc.NetworkID(context.Background())
	if err != nil {
		lock.Unlock(symbol)
		return fmt.Errorf("chainID error: %w", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		lock.Unlock(symbol)
		return fmt.Errorf("sign tx error: %w", err)
	}

	if err := c.allow(); err != nil {
		lock.Unlock(symbol)
		return err
	}
	err = c.rpc.SendTransaction(context.Background(), signedTx)
	if err != nil {
		lock.Unlock(symbol)
		return fmt.Errorf("send tx error: %w", err)
	}

	log.Printf("TX SENT %s: %s", symbol, signedTx.Hash().Hex())
	go c.waitForTxResult(signedTx, symbol, newPrice, clPrice)

	return nil
}

func (c *Client) waitForTxResult(tx *types.Transaction, symbol string, price *big.Int, clPrice *big.Int) {
	defer GetTxLock().Unlock(symbol)

	log.Printf("WAITING FOR STATUS of TX %s: %s", symbol, tx.Hash().Hex())

	if err := c.limiter.Wait(context.Background()); err != nil {
		log.Printf("ERROR: rate limiter wait failed (WaitMined): %v", err)
		log.Printf("DEBUG: revert is now possible for %s (track: %s)", symbol, tx.Hash().Hex())
		return
	}
	receipt, err := bind.WaitMined(context.Background(), c.rpc, tx.Hash())
	if err != nil {
		log.Printf("WAIT MINED TX ERROR: %s | %v", tx.Hash().Hex(), err)
		log.Printf("DEBUG: revert is now possible for %s (track: %s)", symbol, tx.Hash().Hex())
		return
	}

	if receipt.Status == 0 {
		entry := fmt.Sprintf(
			"tx=%s | symbol=%s | sent price=%s | chainlink price=%s",
			tx.Hash().Hex(),
			symbol,
			price.String(),
			clPrice.String(),
		)

		c.reverts.Add(entry)

		log.Printf("REVERTED: %s", entry)
	} else {
		log.Printf("CONFIRMED: %s | %s -> %s",
			tx.Hash().Hex(),
			symbol,
			price.String(),
		)
	}
}

func (c *Client) allow() error {
	ok, delay, err := c.limiter.Allow()
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("rpc rate limit exceeded, try again in %v", delay)
	}
	return nil
}
