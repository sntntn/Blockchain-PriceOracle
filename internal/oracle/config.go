package oracle

import (
	"os"
	"strconv"
)

const (
	AnkrRateLimitPerMinute = 1800
	AnkrRateLimitBurst     = 1800
)

type Config struct {
	SepoliaRpc      string
	SepoliaRpcSync  string
	ContractAddr    string
	DeploymentBlock uint64
	PrivateKey      string
}

func buildSepoliaRPC(apiKey string) string {
	if apiKey == "" {
		panic("ANKR API key not provided")
	}

	return "https://rpc.ankr.com/eth_sepolia/" + apiKey
}

func LoadConfig() Config {

	rpc := buildSepoliaRPC(os.Getenv("ANKR_SEPOLIA_API_KEY"))
	rpcSync := buildSepoliaRPC(os.Getenv("ANKR_SEPOLIA_SYNC_API_KEY"))

	deploymentBlockStr := os.Getenv("DEPLOYMENT_BLOCK")
	var deploymentBlock uint64
	if deploymentBlockStr != "" {
		parsed, err := strconv.ParseUint(deploymentBlockStr, 10, 64)
		if err != nil {
			panic("invalid DEPLOYMENT_BLOCK")
		}
		deploymentBlock = parsed
	}

	return Config{
		SepoliaRpc:      rpc,
		SepoliaRpcSync:  rpcSync,
		ContractAddr:    os.Getenv("CONTRACT_ADDR"),
		DeploymentBlock: deploymentBlock,
		PrivateKey:      os.Getenv("PRIVATE_KEY"),
	}
}
