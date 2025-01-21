package ipfs

import (
    "context"
    "fmt"
    "os"
    "io"

    "github.com/ipfs/go-ipfs-config"
    files "github.com/ipfs/go-ipfs-files"
    "github.com/ipfs/boxo/path"
    gcid "github.com/ipfs/go-cid"
    "github.com/ipfs/kubo/core"
    "github.com/ipfs/kubo/core/coreapi"
    "github.com/ipfs/kubo/core/node/libp2p"
    "github.com/ipfs/kubo/plugin/loader"
    "github.com/ipfs/kubo/repo/fsrepo"
)

// IPFS客户端结构体
type IPFSClient struct {
    ctx    context.Context
    cancel context.CancelFunc
    api    coreapi.CoreAPI
    node   *core.IpfsNode
}

// InitIPFS 初始化IPFS客户端
// repoPath: IPFS仓库路径
// 返回: IPFS客户端实例和错误信息

func InitIPFS(repoPath string) (*IPFSClient, error) {
		ctx, cancel := context.WithCancel(context.Background())

		// 创建客户端实例
		client := &IPFSClient{
				ctx:    ctx,
				cancel: cancel,
		}

		// 初始化插件系统
		plugins, err := loader.NewPluginLoader("")
		if (err != nil) {
				return nil, fmt.Errorf("初始化插件失败: %s", err)
		}
		if err := plugins.Initialize(); err != nil {
				return nil, fmt.Errorf("加载插件失败: %s", err)
		}

		// 检查仓库是否已经存在
		if !fsrepo.IsInitialized(repoPath) {
				// 仓库不存在时才初始化
				if err := setupRepo(repoPath); err != nil {
						return nil, fmt.Errorf("设置仓库失败: %s", err)
				}
		}

		// 打开已存在的仓库
		repo, err := fsrepo.Open(repoPath)
		if err != nil {
				return nil, fmt.Errorf("打开仓库失败: %s", err)
		}

		// 创建节点
		nodeOptions := &core.BuildCfg{
				Online:  true,
				Routing: libp2p.DHTOption,
				Repo:    repo,
		}

		node, err := core.NewNode(ctx, nodeOptions)
		if err != nil {
				return nil, fmt.Errorf("创建IPFS节点失败: %s", err)
		}
		client.node = node

        // 创建API实例
        api, err := coreapi.NewCoreAPI(node)
        if err != nil {
                return nil, fmt.Errorf("创建API失败: %s", err)
        }
        client.api = api

		return client, nil
}

// SaveFile 保存文件到IPFS
// filePath: 要保存的文件路径
// 返回: 文件的CID(内容标识符)和错误信息
func (c *IPFSClient) SaveFile(filePath string) (string, error) {
    // 打开文件
    file, err := os.Open(filePath)
    if err != nil {
        return "", fmt.Errorf("打开文件失败: %s", err)
    }
    defer file.Close()

    // 创建IPFS文件节点
    fileNode := files.NewReaderFile(file)
    path, err := c.api.Unixfs().Add(c.ctx, fileNode)
    if err != nil {
        return "", fmt.Errorf("添加文件到IPFS失败: %s", err)
    }

    return path.String(), nil
}

// GetFile 从IPFS获取文件
// cid: 文件的内容标识符
// outputPath: 文件保存路径
// 返回: 错误信息
func (c *IPFSClient) GetFile(cid string, outputPath string) error {
    // 使用IpfsPath创建正确的IPFS路径
    decodedCid, err := gcid.Decode(cid)
    if err != nil {
        return fmt.Errorf("无效的CID: %v", err)
    }

    ipfsPath := path.FromCid(decodedCid)

    // 获取文件
    node, err := c.api.Unixfs().Get(c.ctx, ipfsPath)
    if err != nil {
        return fmt.Errorf("从IPFS获取文件失败: %s", err)
    }

    // 创建输出文件
    file, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("创建输出文件失败: %s", err)
    }
    defer file.Close()

    // 写入文件
    if _, err := io.Copy(out, file); err != nil {
        return fmt.Errorf("写入文件失败: %v", err)
    }

    return nil
}

// Close 关闭IPFS客户端
func (c *IPFSClient) Close() error {
    c.cancel()
    if err := c.node.Close(); err != nil {
        return fmt.Errorf("关闭节点失败: %s", err)
    }
    return nil
}

// 内部函数：设置IPFS仓库
func setupRepo(repoPath string) error {
    if _, err := os.Stat(repoPath); !os.IsNotExist(err) {
        return nil
    }

    // 创建仓库目录
    if err := os.MkdirAll(repoPath, 0755); err != nil {
        return err
    }

    // 创建配置
    cfg, err := config.Init(os.Stdout, 2048)
    if err != nil {
        return err
    }

    // 设置基本配置
    cfg.Bootstrap = config.DefaultBootstrapAddresses
    cfg.Addresses.Swarm = []string{
        "/ip4/0.0.0.0/tcp/4001",
        "/ip4/0.0.0.0/udp/4001/quic",
    }
    cfg.Addresses.API = []string{"/ip4/127.0.0.1/tcp/5001"}
    cfg.Addresses.Gateway = []string{"/ip4/127.0.0.1/tcp/8080"}

    // 初始化仓库
    return fsrepo.Init(repoPath, cfg)
}