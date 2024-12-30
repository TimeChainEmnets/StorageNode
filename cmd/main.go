package main

import (
	"fmt"
	"config"
	"blockchain"

)

func main() {
	cfg := &config.Config{
		BlockchainConfig: &config.BlockchainConfig{
				NodeURL:         "https://sepolia.infura.io/v3/YOUR-PROJECT-ID",
				ContractAddress: "YOUR-CONTRACT-ADDRESS",
				PrivateKey:      "YOUR-PRIVATE-KEY",
				GasLimit:        3000000,
		},
	}
	
	client := blockchain.NewClient(cfg)
	
	// 注册节点
	err := client.RegisterNode(context.Background(), 
		"192.168.1.1", 
		39.9042, 
		116.4074, 
		big.NewInt(1000000))
}
