# GitHub Actions 和 DockerHub 集成配置

本文档介绍如何设置 GitHub Actions 自动构建和推送 Docker 镜像到 DockerHub。

## 🔧 配置步骤

### 1. DockerHub 准备

#### 1.1 创建 DockerHub 账户
- 访问 [DockerHub](https://hub.docker.com/)
- 注册账户（如果还没有）

#### 1.2 创建 Access Token
1. 登录 DockerHub
2. 进入 Account Settings > Security
3. 点击 "New Access Token"
4. 输入描述（例如：`read-it-later-github-actions`）
5. 选择权限：`Read, Write, Delete`
6. 点击 "Generate"
7. **重要：** 立即复制生成的 token，这是唯一显示的机会

#### 1.3 创建仓库（可选）
虽然 push 时会自动创建，但您可以提前创建：
- `adoom2018/read-it-later-frontend`
- `adoom2018/read-it-later-backend`

### 2. GitHub 仓库配置

#### 2.1 添加 Secrets
1. 进入 GitHub 仓库
2. 点击 Settings > Secrets and variables > Actions
3. 点击 "New repository secret"

添加以下 secrets：

| Secret 名称 | 描述 | 值 |
|-------------|------|-----|
| `DOCKERHUB_USERNAME` | DockerHub 用户名 | 您的 DockerHub 用户名 |
| `DOCKERHUB_TOKEN` | DockerHub Access Token | 从步骤 1.2 获取的 token |

#### 2.2 验证配置
确保在仓库的 Actions secrets 中看到这两个 secrets。

## 🚀 使用说明

### 自动构建和推送

当您创建一个新的 release 时，GitHub Actions 会自动：

1. **构建镜像** - 为前端和后端构建 Docker 镜像
2. **多架构支持** - 构建 `linux/amd64` 和 `linux/arm64` 版本
3. **推送到 DockerHub** - 使用版本标签和 latest 标签
4. **安全扫描** - 使用 Trivy 扫描镜像漏洞
5. **生成部署文件** - 创建生产环境的 docker-compose 文件

### 创建 Release

1. 在 GitHub 仓库中点击 "Releases"
2. 点击 "Create a new release"
3. 填写标签版本（例如：`v1.0.0`）
4. 填写 Release 标题和描述
5. 点击 "Publish release"

### 镜像标签规则

对于版本 `v1.2.3`，会创建以下标签：
- `adoom2018/read-it-later-frontend:v1.2.3`
- `adoom2018/read-it-later-frontend:latest`

## 📦 部署使用

### 方式 1: 使用部署脚本

```bash
# 下载并运行部署脚本
curl -fsSL https://raw.githubusercontent.com/adoom2017/read-it-later/main/deploy-dockerhub.sh | bash

# 或指定版本
curl -fsSL https://raw.githubusercontent.com/adoom2017/read-it-later/main/deploy-dockerhub.sh | bash -s -- -v v1.0.0
```

### 方式 2: 手动部署

```bash
# 创建 docker-compose.yml
cat > docker-compose.yml << 'EOF'
version: '3.8'

services:
  backend:
    image: adoom2018/read-it-later-backend:latest
    container_name: read-it-later-backend
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
    environment:
      - GIN_MODE=release
    networks:
      - app-network

  frontend:
    image: adoom2018/read-it-later-frontend:latest
    container_name: read-it-later-frontend
    restart: unless-stopped
    ports:
      - "80:80"
    depends_on:
      - backend
    networks:
      - app-network

networks:
  app-network:
    driver: bridge
EOF

# 创建数据目录
mkdir -p ./data

# 启动服务
docker-compose up -d
```

## 🔍 监控和维护

### 检查构建状态
- 在 GitHub 仓库的 Actions 标签中查看构建状态
- 查看构建日志和错误信息

### 查看 DockerHub 镜像
- 访问 DockerHub 查看推送的镜像
- 检查镜像大小和更新时间

### 本地测试
```bash
# 拉取镜像
docker pull adoom2018/read-it-later-frontend:latest
docker pull adoom2018/read-it-later-backend:latest

# 检查镜像
docker images | grep read-it-later
```

## 🔐 安全考虑

1. **Access Token 安全**
   - 定期轮换 DockerHub access token
   - 使用最小权限原则
   - 不要在代码中暴露 token

2. **镜像安全**
   - 定期更新基础镜像
   - 关注安全漏洞报告
   - 使用镜像扫描工具

3. **版本管理**
   - 使用语义化版本规范
   - 保持版本标签的一致性
   - 定期清理旧版本镜像

## 🛠️ 故障排除

### 常见问题

1. **构建失败**
   - 检查 DockerHub credentials
   - 验证 Dockerfile 语法
   - 查看构建日志

2. **推送失败**
   - 确认 DockerHub token 有效
   - 检查网络连接
   - 验证镜像名称格式

3. **部署失败**
   - 检查镜像是否成功推送
   - 验证 docker-compose 配置
   - 检查端口占用情况

### 调试命令

```bash
# 查看容器状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 检查网络
docker network ls

# 检查镜像
docker images
```

## 📚 相关资源

- [Docker Hub 文档](https://docs.docker.com/docker-hub/)
- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [Docker 最佳实践](https://docs.docker.com/develop/dev-best-practices/)
- [语义化版本规范](https://semver.org/)

## 🤝 贡献

如果您发现配置问题或有改进建议，请：

1. 创建 Issue 描述问题
2. 提交 Pull Request 进行修复
3. 更新相关文档

感谢您的贡献！🎉
