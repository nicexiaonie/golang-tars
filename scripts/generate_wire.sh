#!/bin/bash

# 自动生成 wire.go 配置文件
# 用法: ./scripts/generate_wire.sh [模块路径]
# 示例: ./scripts/generate_wire.sh internal/user

set -e

# 获取模块路径参数，默认为 internal/user
MODULE_PATH=${1:-"internal/user"}
WIRE_FILE="${MODULE_PATH}/wire.go"

# 获取模块包名（最后一段路径）
PACKAGE_NAME=$(basename "$MODULE_PATH")

# 获取 Go 模块根路径
GO_MODULE=$(go list -m)

echo "开始生成 Wire 配置文件: $WIRE_FILE"
echo "模块路径: $MODULE_PATH"
echo "包名: $PACKAGE_NAME"
echo "Go 模块: $GO_MODULE"

# 临时文件
TEMP_FILE=$(mktemp)

# 清理临时文件
cleanup() {
    rm -f "$TEMP_FILE"
}
trap cleanup EXIT

# 扫描指定目录下的 Provider 函数
scan_directory() {
    local dir=$1
    local category=$2

    if [ ! -d "$dir" ]; then
        return
    fi

    echo "  扫描 $dir 目录..."

    # 查找所有 Go 文件（排除 wire_gen.go、wire.go 和测试文件）
    find "$dir" -maxdepth 1 -name "*.go" ! -name "*_test.go" ! -name "wire_gen.go" ! -name "wire.go" -type f | while read -r file; do
        # 提取包名
        local pkg=$(grep -m 1 "^package " "$file" | awk '{print $2}')

        if [ -z "$pkg" ]; then
            continue
        fi

        # 查找 New* 函数定义
        grep -E "^func New[A-Z][a-zA-Z0-9]*\(" "$file" | while read -r line; do
            # 提取函数名（更精确的正则）
            local func_name=$(echo "$line" | sed -E 's/^func ([A-Z][a-zA-Z0-9]*)\(.*/\1/')

            if [ -z "$func_name" ]; then
                continue
            fi

            # 生成导入路径
            local rel_dir=$(echo "$dir" | sed "s|^$MODULE_PATH/||")
            local import_path="${GO_MODULE}/${MODULE_PATH}/${rel_dir}"

            # 保存结果（包名.函数名|导入路径|分类）
            echo "${pkg}.${func_name}|${import_path}|${category}" >> "$TEMP_FILE"
        done
    done
}

# 扫描根目录的 Provide 函数
scan_root_providers() {
    local provider_file="${MODULE_PATH}/provider.go"

    if [ ! -f "$provider_file" ]; then
        return
    fi

    echo "  扫描根目录 provider.go..."

    # 查找 Provide* 函数
    grep -E "^func Provide[A-Z][a-zA-Z0-9]*\(" "$provider_file" | while read -r line; do
        local func_name=$(echo "$line" | sed -E 's/^func ([A-Z][a-zA-Z0-9]*)\(.*/\1/')

        if [ -z "$func_name" ]; then
            continue
        fi

        # Provide 函数在根包，不需要前缀
        echo "${func_name}||base" >> "$TEMP_FILE"
    done
}

# 开始扫描
echo "扫描 Provider 函数..."

# 扫描各个子目录
scan_directory "$MODULE_PATH/data/mysql/dao" "dao"
scan_directory "$MODULE_PATH/service" "service"
scan_directory "$MODULE_PATH/handler" "handler"
scan_directory "$MODULE_PATH/proxy" "proxy"

# 扫描根目录的 Provider 函数
scan_root_providers

# 读取并分类结果
if [ ! -s "$TEMP_FILE" ]; then
    echo "错误: 未找到任何 Provider 函数"
    exit 1
fi

echo ""
echo "找到的 Providers:"

# 用 awk 处理文件，去重并分类
awk -F'|' '!seen[$1]++ {print $0}' "$TEMP_FILE" | sort | while IFS='|' read -r provider import_path category; do
    echo "  - $provider ($category)"

    # 记录导入路径（如果有的话）
    if [ -n "$import_path" ]; then
        echo "$import_path" >> "${TEMP_FILE}.imports"
    fi

    # 按类别保存
    case "$category" in
        base)
            echo "$provider" >> "${TEMP_FILE}.base"
            ;;
        dao)
            echo "$provider" >> "${TEMP_FILE}.dao"
            ;;
        proxy)
            echo "$provider" >> "${TEMP_FILE}.proxy"
            ;;
        service)
            echo "$provider" >> "${TEMP_FILE}.service"
            ;;
        handler)
            echo "$provider" >> "${TEMP_FILE}.handler"
            ;;
    esac
done

# 生成 wire.go 文件
echo ""
echo "生成 wire.go 文件..."

cat > "$WIRE_FILE" << 'EOF'
//go:build wireinject
// +build wireinject

EOF

echo "package $PACKAGE_NAME" >> "$WIRE_FILE"
echo "" >> "$WIRE_FILE"

# 生成导入部分
echo "import (" >> "$WIRE_FILE"

# 添加导入路径（去重排序）
if [ -f "${TEMP_FILE}.imports" ]; then
    sort -u "${TEMP_FILE}.imports" | while read -r import_path; do
        echo "	\"$import_path\"" >> "$WIRE_FILE"
    done
    echo "" >> "$WIRE_FILE"
fi

echo "	\"github.com/google/wire\"" >> "$WIRE_FILE"
echo ")" >> "$WIRE_FILE"
echo "" >> "$WIRE_FILE"

# 生成 ProviderSet
cat >> "$WIRE_FILE" << 'EOF'
// ProviderSet 定义所有的依赖注入 Provider
// Wire 会按照依赖关系自动排序并生成初始化代码
var ProviderSet = wire.NewSet(
EOF

# 添加各类 Provider
add_providers() {
    local file=$1
    local comment=$2

    if [ -f "$file" ] && [ -s "$file" ]; then
        echo "	// $comment" >> "$WIRE_FILE"
        cat "$file" | while read -r provider; do
            echo "	$provider," >> "$WIRE_FILE"
        done
        echo "" >> "$WIRE_FILE"
    fi
}

add_providers "${TEMP_FILE}.base" "基础设施层 Provider"
add_providers "${TEMP_FILE}.dao" "DAO 层（数据访问）"
add_providers "${TEMP_FILE}.proxy" "Proxy 层（第三方服务调用）"
add_providers "${TEMP_FILE}.service" "Service 层（业务逻辑）"
add_providers "${TEMP_FILE}.handler" "Handler 层（API 处理）"

# 移除最后的空行和逗号
sed -i.bak '$ { /^$/d; }' "$WIRE_FILE"
rm -f "${WIRE_FILE}.bak"

echo ")" >> "$WIRE_FILE"
echo "" >> "$WIRE_FILE"

# 生成 Initialize 函数
cat >> "$WIRE_FILE" << 'EOF'
// InitializeHandler 初始化 Handler
// Wire 会自动生成这个函数的实现，解析所有依赖关系
// 生成的代码在 wire_gen.go 文件中
func InitializeHandler() (*handler.DemoHandler, error) {
	wire.Build(ProviderSet)
	return nil, nil
}
EOF

# 清理临时文件
rm -f "${TEMP_FILE}.imports" "${TEMP_FILE}.base" "${TEMP_FILE}.dao" "${TEMP_FILE}.proxy" "${TEMP_FILE}.service" "${TEMP_FILE}.handler"

echo ""
echo "✓ Wire 配置文件已生成: $WIRE_FILE"
echo ""
echo "下一步操作:"
echo "  1. 检查生成的 wire.go 文件: cat $WIRE_FILE"
echo "  2. 运行 'make wire' 生成依赖注入代码"
echo "  3. 检查生成的 wire_gen.go 文件"
