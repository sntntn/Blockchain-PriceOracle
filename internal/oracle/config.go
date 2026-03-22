package oracle

import (
	"os"
)

type Config struct {
	SepoliaRPC   string
	ContractAddr string
	PrivateKey   string
}

func LoadConfig() Config {
	apiKey := os.Getenv("ANKR_SEPOLIA_API_KEY")
	rpc := "https://rpc.ankr.com/eth_sepolia/"

	if apiKey != "" {
		rpc += apiKey
	}

	return Config{
		SepoliaRPC:   rpc,
		ContractAddr: os.Getenv("CONTRACT_ADDR"),
		PrivateKey:   os.Getenv("PRIVATE_KEY"),
	}
}

func MustLoadConfig() Config {
	config := LoadConfig()
	if config.SepoliaRPC == "https://rpc.ankr.com/eth_sepolia/" {
		panic("ANKR_SEPOLIA_API_KEY not set")
	}
	return config
}
