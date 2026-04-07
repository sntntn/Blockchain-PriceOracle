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
	SepoliaRPC      string
	ContractAddr    string
	DeploymentBlock uint64
	PrivateKey      string
}

func LoadConfig() Config {

	apiKey := os.Getenv("ANKR_SEPOLIA_API_KEY")
	if apiKey == "" {
		panic("ANKR_SEPOLIA_API_KEY not set")
	}
	rpc := "https://rpc.ankr.com/eth_sepolia/"
	rpc += apiKey

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
		SepoliaRPC:      rpc,
		ContractAddr:    os.Getenv("CONTRACT_ADDR"),
		DeploymentBlock: deploymentBlock,
		PrivateKey:      os.Getenv("PRIVATE_KEY"),
	}
}
