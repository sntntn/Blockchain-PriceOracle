package oracle

import (
	"os"
)

type Config struct {
	SepoliaRPC   string
	ContractAddr string
}

func LoadConfig() Config {
	apiKey := os.Getenv("ANKR_SEPOLIA_API_KEY")
	rpc := "https://rpc.ankr.com/eth_sepolia/"

	if apiKey != "" {
		rpc += apiKey
	}

	return Config{
		SepoliaRPC:   rpc,
		ContractAddr: "0x19D1c199cFEC4022A76045aEc281ccb63F18387B",
	}
}

func MustLoadConfig() Config {
	config := LoadConfig()
	if config.SepoliaRPC == "https://rpc.ankr.com/eth_sepolia/" {
		panic("ANKR_SEPOLIA_API_KEY not set")
	}
	return config
}
