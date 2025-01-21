package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "github.com/TimeChainEmnets/StorageNode/internel/ipfs"
)

func main() {
    // 命令行参数
    repoPath := flag.String("repo", "", "IPFS仓库路径")
    initRepo := flag.Bool("init", false, "初始化新仓库")
    flag.Parse()

    // 如果未指定仓库路径，使用默认路径
    if *repoPath == "" {
        homeDir, err := os.UserHomeDir()
        if err != nil {
            log.Fatal(err)
        }
        *repoPath = filepath.Join(homeDir, ".myipfs")
    }

    // 创建仓库管理器
    repoManager := ipfs.NewRepoManager(*repoPath)

    // 如果指定了初始化标志或仓库不存在
    if *initRepo {
        exists, _ := repoManager.CheckRepo()
        if exists {
            log.Fatal("仓库已存在，如需重新初始化请先删除现有仓库")
        }
        if err := repoManager.InitRepo(); err != nil {
            log.Fatal("初始化仓库失败:", err)
        }
        fmt.Println("仓库初始化成功")
        return
    }

    // 初始化IPFS客户端
    client, err := ipfs.InitIPFS(*repoPath)
    if err != nil {
        log.Fatal("初始化IPFS客户端失败:", err)
    }
    defer client.Close()

    fmt.Println("IPFS客户端启动成功")
    // 这里可以添加更多的命令行功能...
}