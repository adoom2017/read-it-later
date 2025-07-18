#!/bin/bash

# Git 仓库合并脚本
# 将前端和后端仓库合并到一个新的统一仓库中

set -e

echo "=== Read It Later 仓库合并脚本 ==="

# 配置
NEW_REPO_NAME="read-it-later"
FRONTEND_REPO="https://github.com/adoom2017/read-it-later-frontend.git"
BACKEND_REPO="https://github.com/adoom2017/read-it-later-backend.git"

# 创建新的合并目录
MERGE_DIR="read-it-later-merged"
rm -rf $MERGE_DIR
mkdir $MERGE_DIR
cd $MERGE_DIR

echo "1. 初始化新的 Git 仓库..."
git init

echo "2. 添加前端仓库作为远程源..."
git remote add frontend $FRONTEND_REPO
git fetch frontend

echo "3. 合并前端代码到 frontend/ 目录..."
git checkout -b main
git merge --allow-unrelated-histories frontend/master
mkdir -p frontend
git mv *.* frontend/ 2>/dev/null || true
git mv * frontend/ 2>/dev/null || true
git commit -m "Move frontend files to frontend/ directory"

echo "4. 添加后端仓库作为远程源..."
git remote add backend $BACKEND_REPO
git fetch backend

echo "5. 合并后端代码到 backend/ 目录..."
git merge --allow-unrelated-histories backend/master
mkdir -p backend
# 移动后端文件到 backend 目录（避免与前端文件冲突）
for file in $(git ls-files | grep -v "^frontend/"); do
    if [[ -f "$file" ]]; then
        mkdir -p "backend/$(dirname "$file")"
        git mv "$file" "backend/$file"
    fi
done
git commit -m "Move backend files to backend/ directory"

echo "6. 复制 Docker 配置文件..."
# 这里需要手动复制你之前创建的 Docker 配置文件
# 因为这些文件不在原始仓库中

echo "7. 清理远程源..."
git remote remove frontend
git remote remove backend

echo "8. 添加新的远程仓库..."
echo "请在 GitHub 上创建新的仓库: $NEW_REPO_NAME"
echo "然后运行以下命令："
echo "git remote add origin https://github.com/adoom2017/$NEW_REPO_NAME.git"
echo "git push -u origin main"

echo ""
echo "合并完成！新的仓库位于: $MERGE_DIR"
echo "请检查文件结构并推送到远程仓库。"
