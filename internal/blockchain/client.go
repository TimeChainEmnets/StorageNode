package blockchain

import (
	"fmt"
	"crypto/ecdsa"
	"context"
	"log"
	"strings"
	"math/big"

	"github.com/TimeChainEmnets/StorageNode/internal/config"

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

// RegisterNode 方法用于注册存储节点
func (c *Client) RegisterNode(ipAddress string, latitude, longitude *big.Int, capacity *big.Int) (*types.Transaction, error) {
	// 获取当前的 nonce
	nonce, err := c.eth.PendingNonceAt(context.Background(), c.address)
	if err != nil {
			return nil, fmt.Errorf("failed to get nonce: %v", err)
	}

	// 获取当前的 gas price
	gasPrice, err := c.eth.SuggestGasPrice(context.Background())
	if err != nil {
			return nil, fmt.Errorf("failed to suggest gas price: %v", err)
	}

	// 获取链ID
	chainID, err := c.eth.ChainID(context.Background())
	if err != nil {
			return nil, fmt.Errorf("failed to get chain ID: %v", err)
	}

	// 创建交易选项
	auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, chainID)
	if err != nil {
			return nil, fmt.Errorf("failed to create authorized transactor: %v", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // 不发送 ETH
	auth.GasLimit = c.gasLimit      // 使用配置的 gas limit
	auth.GasPrice = gasPrice

	// 创建和发送交易
	tx, err := c.contract.Transact(auth, "registerNode", ipAddress, latitude, longitude, capacity)
	if err != nil {
			return nil, fmt.Errorf("failed to send transaction: %v", err)
	}

	return tx, nil
}