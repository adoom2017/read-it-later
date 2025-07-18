# Read It Later - 完整应用

一个基于 React + Go 的稍后阅读应用，支持保存网页文章、提取内容、添加标签等功能。

## 项目结构

```
read-it-later/
├── frontend/          # React 前端应用
│   ├── src/
│   ├── public/
│   ├── package.json
│   ├── Dockerfile
│   └── nginx.conf
├── backend/           # Go 后端 API
│   ├── main.go
│   ├── go.mod
│   ├── Dockerfile
│   └── handlers/
├── docker-compose.yml
├── docker-compose.prod.yml
├── deploy.sh
└── README.md
```

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

### 使用 Docker 部署（推荐）

```bash
# 克隆项目
git clone https://github.com/adoom2017/read-it-later.git
cd read-it-later

# 一键部署
chmod +x deploy.sh
./deploy.sh
```

### 手动部署

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

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详情请参阅 [LICENSE](LICENSE) 文件。

## 更新日志

### v1.0.0
- 初始版本发布
- 基础的文章保存和管理功能
- Docker 容器化支持
- 响应式前端界面

## 联系方式

- 项目地址: https://github.com/adoom2017/read-it-later
- 问题反馈: https://github.com/adoom2017/read-it-later/issues
