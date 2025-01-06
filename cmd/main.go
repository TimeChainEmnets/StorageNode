package main

import (
	"fmt"
	"math/big"
	"github.com/TimeChainEmnets/StorageNode/internal/config"
	"github.com/TimeChainEmnets/StorageNode/internal/blockchain"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		panic(err)
	}
	client := blockchain.NewClient(cfg)
	
	// 注册节点
	result, err := client.RegisterNode("192.168.1.1", 
		big.NewInt(39), 
		big.NewInt(116), 
		big.NewInt(1000000))
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
