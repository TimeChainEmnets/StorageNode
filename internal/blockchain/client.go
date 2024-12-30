package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"strings"
	//"storage-node/internal/config"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	config     *config.Config
	eth        *ethclient.Client
	contract   *bind.BoundContract
	privateKey *ecdsa.PrivateKey
	address    common.Address
	gasLimit   uint64
}

func NewClient(cfg *config.Config) *Client {
	// 连接到以太坊网络
	client, err := ethclient.Dial(cfg.BlockchainConfig.NodeURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// 加载私钥
	privateKey, err := crypto.HexToECDSA(cfg.BlockchainConfig.PrivateKey)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	// 从私钥获取公钥和地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Failed to get public key")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 创建合约实例
	parsed, err := abi.JSON(strings.NewReader(ContractABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	contractAddress := common.HexToAddress(cfg.BlockchainConfig.ContractAddress)
	contract := bind.NewBoundContract(contractAddress, parsed, client, client, client)

	return &Client{
		config:     cfg,
		eth:        client,
		contract:   contract,
		privateKey: privateKey,
		address:    address,
		gasLimit:   cfg.BlockchainConfig.GasLimit,
	}
}