# Docker 部署指南

本项目使用 Docker 和 Docker Compose 进行部署，支持本地开发和远程生产环境部署。

## 前置要求

- Docker (>= 20.10)
- Docker Compose (>= 2.0)
- Linux 服务器（推荐 Ubuntu 20.04+）

## 快速部署（推荐）

### 1. 使用部署脚本（适用于远程 Linux 服务器）

```bash
# 克隆项目
git clone <repository-url>
cd read-it-later

# 设置执行权限
chmod +x deploy.sh

# 运行部署脚本
./deploy.sh
```

部署脚本会自动：
- 检查 Docker 环境
- 创建必要的目录
- 构建和启动服务
- 验证服务状态
- 显示访问地址

### 2. 手动部署

```bash
# 克隆项目
git clone <repository-url>
cd read-it-later

# 创建数据目录
mkdir -p data logs

# 生产环境部署
docker-compose -f docker-compose.prod.yml up -d --build

# 开发环境部署
docker-compose up -d --build
```

## 访问应用

部署完成后，可以通过以下地址访问：

- **前端应用**: http://YOUR_SERVER_IP
- **后端 API**: http://YOUR_SERVER_IP:8080
- **健康检查**: http://YOUR_SERVER_IP/health

## 服务器防火墙配置

确保服务器防火墙已开放以下端口：

```bash
# Ubuntu/Debian
sudo ufw allow 80/tcp
sudo ufw allow 8080/tcp

# CentOS/RHEL
sudo firewall-cmd --permanent --add-port=80/tcp
sudo firewall-cmd --permanent --add-port=8080/tcp
sudo firewall-cmd --reload
```

## 常用管理命令

### 查看服务状态
```bash
# 生产环境
docker-compose -f docker-compose.prod.yml ps

# 开发环境
docker-compose ps
```

### 查看日志
```bash
# 查看所有服务日志
docker-compose -f docker-compose.prod.yml logs

# 查看特定服务日志
docker-compose -f docker-compose.prod.yml logs backend
docker-compose -f docker-compose.prod.yml logs frontend

# 实时查看日志
docker-compose -f docker-compose.prod.yml logs -f
```

### 停止服务
```bash
docker-compose -f docker-compose.prod.yml down
```

### 重启服务
```bash
docker-compose -f docker-compose.prod.yml restart
```

### 更新应用
```bash
# 拉取最新代码
git pull

# 重新构建并启动
docker-compose -f docker-compose.prod.yml up -d --build
```

## 服务配置

### 后端服务 (backend)
- **端口**: 8080
- **数据库**: SQLite（持久化存储）
- **健康检查**: /health
- **数据目录**: /app/data

### 前端服务 (frontend)
- **端口**: 80 (HTTP), 443 (HTTPS)
- **服务器**: Nginx
- **代理**: API 请求自动代理到后端
- **健康检查**: /health

## 数据持久化

数据库文件存储在 Docker 卷中，确保数据在容器重启后不会丢失：

```bash
# 查看数据卷
docker volume ls

# 备份数据
docker run --rm -v read-it-later_app-data:/data -v $(pwd):/backup alpine tar czf /backup/backup.tar.gz /data

# 恢复数据
docker run --rm -v read-it-later_app-data:/data -v $(pwd):/backup alpine tar xzf /backup/backup.tar.gz -C /
```

## 监控和日志

### 健康检查
```bash
# 检查后端健康状态
curl http://YOUR_SERVER_IP:8080/health

# 检查前端健康状态  
curl http://YOUR_SERVER_IP/health
```

### 系统监控
```bash
# 查看容器资源使用情况
docker stats

# 查看容器进程
docker-compose -f docker-compose.prod.yml top
```

## 故障排除

### 常见问题

1. **端口被占用**
   ```bash
   # 查看端口使用情况
   sudo netstat -tulpn | grep :80
   sudo netstat -tulpn | grep :8080
   ```

2. **权限问题**
   ```bash
   # 检查数据目录权限
   ls -la data/
   
   # 修复权限
   sudo chown -R $USER:$USER data/
   ```

3. **服务无法启动**
   ```bash
   # 查看详细错误信息
   docker-compose -f docker-compose.prod.yml logs
   
   # 检查容器状态
   docker-compose -f docker-compose.prod.yml ps
   ```

### 调试命令

```bash
# 进入后端容器
docker-compose -f docker-compose.prod.yml exec backend sh

# 进入前端容器
docker-compose -f docker-compose.prod.yml exec frontend sh

# 检查网络连接
docker-compose -f docker-compose.prod.yml exec backend ping frontend
```

## 生产环境建议

1. **使用 HTTPS**: 配置 SSL 证书
2. **反向代理**: 使用 Nginx 或 Traefik
3. **监控告警**: 集成监控系统
4. **日志管理**: 配置日志收集和分析
5. **备份策略**: 定期备份数据库
6. **资源限制**: 设置容器资源限制
7. **安全配置**: 更新默认密码和密钥

### 资源限制配置示例

```yaml
services:
  backend:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
```

## 性能优化

1. **数据库优化**: 定期清理和优化数据库
2. **缓存策略**: 配置适当的缓存头
3. **压缩**: 启用 gzip 压缩
4. **CDN**: 使用 CDN 加速静态资源

更多详细信息请参考项目文档。
