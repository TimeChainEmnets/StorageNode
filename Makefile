.PHONY: build install clean

# 构建
build:
	go build -o bin/myipfs ./cmd/myipfs

# 安装
install:
	chmod +x scripts/install.sh
	./scripts/install.sh

# 清理
clean:
	rm -f bin/myipfs

# 运行测试
test:
	go test ./...

# 开发模式运行
dev:
	go run ./cmd/myipfs