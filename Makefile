.PHONY: help proto build run clean test openapi openapi-serve buf-generate wire wire-gen wire-all

# 默认
help:
	@echo "可用的 Make 命令:"
	@echo "  make proto    			- 编译 proto 文件"
	@echo "  make build-user    			- 编译user服务"
	@echo "  make upload-user    			- 编译部署user服务"
	@echo "  make proto-clean			- 清理proto生成的文件"
	@echo "  make deps     			- 安装依赖"
	@echo "  make wire-gen 			- 自动生成 wire.go 配置文件"
	@echo "  make wire     			- 生成 wire_gen.go 依赖注入代码"
	@echo "  make wire-all 			- 自动生成 wire.go 并运行 wire 生成代码"
	@echo "  make openapi  			- 使用 buf 生成 OpenAPI 文档"

# 编译 proto 文件
proto:
	@echo "编译 Proto 文件..."
	@bash scripts/proto.sh

# 编译项目
build-user: proto
	@echo "编译项目..."
	@cd cmd/user && make tar
# 部署user服务
upload-user: build-user
	@echo "部署user服务..."
	@cd cmd/user && make upload && make clean

# 清理生成的文件
proto-clean:
	@echo "清理生成的文件..."
	@rm -rf bin/
	@rm -rf pkg/proto/

# 安装依赖
deps:
	@echo "安装依赖..."
	@go mod download
	@go mod tidy

# 自动生成 wire.go 配置文件
wire-gen:
	@echo "自动生成 wire.go 配置文件..."
	@bash scripts/generate_wire.sh internal/user

# 生成 Wire 依赖注入代码 (wire_gen.go)
wire:
	@echo "生成 Wire 依赖注入代码..."
	@cd internal/user && $$(go env GOPATH)/bin/wire

# 自动生成 wire.go 并运行 wire
wire-all: wire-gen wire
	@echo "✓ Wire 配置和代码生成完成"

# 使用 buf 生成 OpenAPI 文档
openapi:
	@echo "使用 buf 生成 OpenAPI 文档..."
	@docker run --rm -v "$$PWD:/workspace" -w /workspace bufbuild/buf:latest mod update
	@docker run --rm -v "$$PWD:/workspace" -w /workspace bufbuild/buf:latest generate --template buf.gen.openapi.yaml
# 使用 buf 生成所有代码（Go + Connect-go + OpenAPI）

