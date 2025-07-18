# Read It Later - 完整应用

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://golang.org/)
[![React Version](https://img.shields.io/badge/React-18+-61DAFB?logo=react)](https://reactjs.org/)
[![Docker](https://img.shields.io/badge/Docker-Available-2496ED?logo=docker)](https://www.docker.com/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/adoom2017/read-it-later/pulls)

一个基于 React + Go 的稍后阅读应用，支持保存网页文章、提取内容、添加标签等功能。

## 功能特性

- 📖 保存网页文章链接
- 🔍 自动提取文章内容和摘要
- 🏷️ 添加和管理标签
- 📱 响应式设计
- 🐳 Docker 容器化部署
- 💾 SQLite 数据库
- 🔍 文章搜索和过滤

## 技术栈

### 前端
- React 18
- Axios (HTTP 客户端)
- Vite (构建工具)
- CSS3 (样式)

### 后端
- Go 1.24
- Gin (Web 框架)
- SQLite (数据库)
- go-readability (内容提取)

### 部署
- Docker & Docker Compose
- Nginx (反向代理)

## 快速开始

### 🚀 一键部署（推荐）

使用 DockerHub 镜像快速部署：

```bash
# 使用最新版本
curl -fsSL https://raw.githubusercontent.com/adoom2017/read-it-later/main/deploy-dockerhub.sh | bash

# 使用指定版本
curl -fsSL https://raw.githubusercontent.com/adoom2017/read-it-later/main/deploy-dockerhub.sh | bash -s -- -v v1.0.0
```

### 🐳 使用 Docker Hub 镜像

```bash
# 创建数据目录
mkdir -p ./data

# 创建 docker-compose.yml
cat > docker-compose.yml << 'EOF'
version: '3.8'

services:
  backend:
    image: adoom2017/read-it-later-backend:latest
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
    environment:
      - GIN_MODE=release

  frontend:
    image: adoom2017/read-it-later-frontend:latest
    ports:
      - "80:80"
    depends_on:
      - backend
EOF

# 启动服务
docker-compose up -d
```

### 🔧 使用 Docker 本地构建

```bash
# 克隆项目
git clone https://github.com/adoom2017/read-it-later.git
cd read-it-later

# 一键部署
chmod +x deploy.sh
./deploy.sh
```

### 📋 手动部署

```bash
# 构建并启动服务
docker-compose up -d --build

# 查看服务状态
docker-compose ps
```

## 开发环境

### 前端开发
```bash
cd frontend
npm install
npm start
```

### 后端开发
```bash
cd backend
go mod tidy
go run main.go
```

## API 文档

### 文章管理
- `GET /api/articles` - 获取文章列表
- `POST /api/articles` - 添加新文章
- `GET /api/articles/:id` - 获取文章详情
- `DELETE /api/articles/:id` - 删除文章
- `POST /api/articles/:id/tags` - 添加标签

### 系统状态
- `GET /` - 后端健康检查
- `GET /health` - 服务健康状态

## 部署说明

详细的部署说明请参考 [DOCKER_DEPLOYMENT.md](DOCKER_DEPLOYMENT.md)

## 贡献指南

我们欢迎任何形式的贡献！请阅读我们的[贡献指南](CONTRIBUTING.md)了解如何参与项目开发。

### 快速贡献

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

### 开发规范

- 遵循现有的代码风格
- 添加适当的测试
- 更新相关文档
- 确保所有测试通过

## 社区

- 📢 [问题和建议](https://github.com/adoom2017/read-it-later/issues)
- 💬 [讨论区](https://github.com/adoom2017/read-it-later/discussions)
- 📖 [项目文档](https://github.com/adoom2017/read-it-later/wiki)

## 许可证

本项目采用 MIT 许可证 - 详情请参阅 [LICENSE](LICENSE) 文件。

### 第三方许可证

本项目使用了以下开源项目：

- [React](https://github.com/facebook/react) - MIT License
- [Go](https://github.com/golang/go) - BSD 3-Clause License
- [Gin](https://github.com/gin-gonic/gin) - MIT License
- [Vite](https://github.com/vitejs/vite) - MIT License
- [go-readability](https://github.com/go-shiori/go-readability) - MIT License

## 更新日志

### v1.0.0
- 初始版本发布
- 基础的文章保存和管理功能
- Docker 容器化支持
- 响应式前端界面

## 联系方式

- 项目地址: https://github.com/adoom2017/read-it-later
- 问题反馈: https://github.com/adoom2017/read-it-later/issues
