package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	BlockchainConfig BlockchainConfig `json:"blockchain_config"`
	StorageConfig    StorageConfig    `json:"storage_config"`
}

type BlockchainConfig struct {
	NodeURL         string `json:"node_url"`
	ChainID         int64 `json:"chain_id"`
	GasLimit        uint64 `json:"gas_limit"`
	ContractAddress string `json:"contract_address"`
	PrivateKey			string `json:"private_key"`
}

type StorageConfig struct {
	DataDir string `json:"storage_dir"` // 本地数据存储目录
}

func Load(fileName string) (*Config, error) {
	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// 构造配置文件的相对路径
	configPath := filepath.Join(currentDir, fileName)

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
