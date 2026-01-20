#!/bin/bash

# 检查参数
if [ -z "$1" ]; then
    echo "用法: $0 <新模块名称>"
    exit 1
fi

NEW_MOD=$1

# 检查 go.mod 是否存在
if [ ! -f "go.mod" ]; then
    echo "错误: 当前目录未找到 go.mod"
    exit 1
fi

# 自动获取旧模块名
OLD_MOD=$(grep "^module" go.mod | awk '{print $2}')

echo "🚀 重命名: $OLD_MOD -> $NEW_MOD"

# 1. 修改 go.mod 定义
go mod edit -module "$NEW_MOD"

# 2. 批量替换 .go 文件中的 import (区分 Mac/Linux)
if [ "$(uname)" == "Darwin" ]; then
    # MacOS (sed -i 需要空字符串参数)
    find . -type f -name "*.go" -exec sed -i '' "s|${OLD_MOD}|${NEW_MOD}|g" {} +
else
    # Linux (GNU sed)
    find . -type f -name "*.go" -exec sed -i "s|${OLD_MOD}|${NEW_MOD}|g" {} +
fi

echo "✅ 完成！建议运行 'go mod tidy'。"

