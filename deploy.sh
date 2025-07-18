#!/bin/bash

# 部署脚本 - 用于远程 Linux 服务器部署

set -e

echo "=== Read It Later 应用部署脚本 ==="

# 检查 Docker 和 Docker Compose 是否已安装
if ! command -v docker &> /dev/null; then
    echo "错误: Docker 未安装。请先安装 Docker。"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "错误: Docker Compose 未安装。请先安装 Docker Compose。"
    exit 1
fi

# 创建必要的目录
mkdir -p data
mkdir -p logs

# 设置正确的权限
chmod 755 data
chmod 755 logs

# 停止现有的容器（如果存在）
echo "停止现有容器..."
docker-compose -f docker-compose.prod.yml down 2>/dev/null || true

# 清理旧的镜像（可选）
read -p "是否清理旧的 Docker 镜像？(y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "清理旧镜像..."
    docker system prune -f
fi

# 构建和启动服务
echo "构建和启动服务..."
docker-compose -f docker-compose.prod.yml up -d --build

# 等待服务启动
echo "等待服务启动..."
sleep 10

# 检查服务状态
echo "检查服务状态..."
docker-compose -f docker-compose.prod.yml ps

# 测试服务是否正常运行
echo "测试服务..."
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "✓ 后端服务运行正常"
else
    echo "✗ 后端服务可能有问题"
fi

if curl -f http://localhost/health > /dev/null 2>&1; then
    echo "✓ 前端服务运行正常"
else
    echo "✗ 前端服务可能有问题"
fi

# 获取服务器 IP
SERVER_IP=$(hostname -I | awk '{print $1}')

echo ""
echo "=== 部署完成 ==="
echo "应用访问地址:"
echo "  - 前端应用: http://$SERVER_IP"
echo "  - 后端 API: http://$SERVER_IP:8080"
echo ""
echo "管理命令:"
echo "  查看日志: docker-compose -f docker-compose.prod.yml logs"
echo "  停止服务: docker-compose -f docker-compose.prod.yml down"
echo "  重启服务: docker-compose -f docker-compose.prod.yml restart"
echo ""
echo "请确保服务器防火墙已开放端口 80 和 8080"
