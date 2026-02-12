#!/bin/bash

# Proto 编译脚本
# 用于编译 proto 文件生成 Go 代码（同时支持 gRPC 和 Connect-go）

set -e

# 设置 PATH，确保能找到 protoc 插件
export PATH="$(go env GOBIN):$(go env GOPATH)/bin:$HOME/go/bin:$PATH"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}开始编译 Proto 文件 (gRPC + Connect-go 模式)...${NC}"

# 项目根目录
PROJECT_ROOT=$(cd "$(dirname "$0")/.." && pwd)
PROTO_DIR="${PROJECT_ROOT}/proto"
OUTPUT_DIR="${PROJECT_ROOT}/pkg/proto"

# 检查 protoc 是否安装
if ! command -v protoc &> /dev/null; then
    echo -e "${RED}错误: protoc 未安装${NC}"
    echo "请安装 protoc: https://grpc.io/docs/protoc-installation/"
    exit 1
fi

# 检查必要的插件是否安装
REQUIRED_PLUGINS=(
    "protoc-gen-go"
    "protoc-gen-go-grpc"
    "protoc-gen-connect-go"
    "protoc-gen-validate"
)

for plugin in "${REQUIRED_PLUGINS[@]}"; do
    if ! command -v "$plugin" &> /dev/null; then
        echo -e "${RED}错误: $plugin 未安装${NC}"
        echo "请运行: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
        echo "请运行: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
        echo "请运行: go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest"
        echo "请运行: go install github.com/envoyproxy/protoc-gen-validate@latest"
        exit 1
    fi
done

# 清理旧的生成文件
echo -e "${YELLOW}清理旧的生成文件...${NC}"
rm -rf "${OUTPUT_DIR}"
mkdir -p "${OUTPUT_DIR}"

# 查找 googleapis 路径（用于 google/api/annotations.proto）
GOOGLEAPIS_PATH=""
GO_PATH=$(go env GOPATH)
POSSIBLE_PATHS=(
    "${GO_PATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@*/third_party/googleapis"
    "${GO_PATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@*/third_party/googleapis"
    "${GO_PATH}/pkg/mod/github.com/googleapis/googleapis@*"
)

for path_pattern in "${POSSIBLE_PATHS[@]}"; do
    for path in $path_pattern; do
        if [ -d "$path" ]; then
            GOOGLEAPIS_PATH="$path"
            break 2
        fi
    done
done

if [ -z "$GOOGLEAPIS_PATH" ]; then
    echo -e "${YELLOW}警告: 未找到 googleapis 路径，尝试下载...${NC}"
    go get -u google.golang.org/genproto/googleapis/api/annotations
    for path_pattern in "${POSSIBLE_PATHS[@]}"; do
        for path in $path_pattern; do
            if [ -d "$path" ]; then
                GOOGLEAPIS_PATH="$path"
                break 2
            fi
        done
    done
fi

if [ -z "$GOOGLEAPIS_PATH" ]; then
    echo -e "${RED}错误: 无法找到 googleapis 路径${NC}"
    exit 1
fi

echo -e "${GREEN}找到 googleapis 路径: ${GOOGLEAPIS_PATH}${NC}"

# 查找 validate 路径
VALIDATE_PATH=""
VALIDATE_POSSIBLE_PATHS=(
    "${GO_PATH}/pkg/mod/github.com/envoyproxy/protoc-gen-validate@*/validate"
    "${GO_PATH}/pkg/mod/github.com/envoyproxy/protoc-gen-validate@*"
)

for path_pattern in "${VALIDATE_POSSIBLE_PATHS[@]}"; do
    for path in $path_pattern; do
        if [ -d "$path" ]; then
            VALIDATE_PATH="$path"
            break 2
        fi
    done
done

if [ -z "$VALIDATE_PATH" ]; then
    echo -e "${RED}错误: 无法找到 validate 路径${NC}"
    exit 1
fi

echo -e "${GREEN}找到 validate 路径: ${VALIDATE_PATH}${NC}"

# 编译 common proto
echo -e "${YELLOW}编译 common proto...${NC}"
protoc \
    --proto_path="${PROTO_DIR}" \
    --proto_path="${GOOGLEAPIS_PATH}" \
    --proto_path="${VALIDATE_PATH}/.." \
    --go_out="${OUTPUT_DIR}" \
    --go_opt=paths=source_relative \
    "${PROTO_DIR}/common/v1/common.proto"

# 编译 demo proto (同时生成 gRPC 和 Connect-go)
echo -e "${YELLOW}编译 demo proto (gRPC + Connect-go)...${NC}"
protoc \
    --proto_path="${PROTO_DIR}" \
    --proto_path="${GOOGLEAPIS_PATH}" \
    --proto_path="${VALIDATE_PATH}/.." \
    --go_out="${OUTPUT_DIR}" \
    --go_opt=paths=source_relative \
    --go-grpc_out="${OUTPUT_DIR}" \
    --go-grpc_opt=paths=source_relative \
    --connect-go_out="${OUTPUT_DIR}" \
    --connect-go_opt=paths=source_relative \
    --validate_out="lang=go,paths=source_relative:${OUTPUT_DIR}" \
    "${PROTO_DIR}/demo/v1/demo.proto"

echo -e "${GREEN}✓ Proto 编译完成！${NC}"
echo -e "${GREEN}生成的文件位于: ${OUTPUT_DIR}${NC}"

# 列出生成的文件
echo -e "${YELLOW}生成的文件列表:${NC}"
find "${OUTPUT_DIR}" -type f -name "*.go" | sed "s|${OUTPUT_DIR}/||"

echo -e "${GREEN}完成！${NC}"
