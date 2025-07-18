#!/bin/bash

# Read It Later - DockerHub 部署脚本
# 此脚本用于从 DockerHub 拉取并部署 Read It Later 应用

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查必要工具
check_requirements() {
    log_info "检查系统要求..."
    
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi

    log_success "系统要求检查完成"
}

# 创建必要目录
create_directories() {
    log_info "创建必要目录..."
    
    mkdir -p ./data
    chmod 755 ./data
    
    log_success "目录创建完成"
}

# 获取最新版本
get_latest_version() {
    if [ -z "$VERSION" ]; then
        log_info "获取最新版本信息..."
        VERSION=$(curl -s https://api.github.com/repos/adoom2017/read-it-later/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
        if [ -z "$VERSION" ]; then
            log_warning "无法获取最新版本，使用默认版本 latest"
            VERSION="latest"
        fi
    fi
    log_info "使用版本: $VERSION"
}

# 创建 docker-compose 文件
create_docker_compose() {
    log_info "创建 docker-compose 配置文件..."
    
    cat > docker-compose.yml << EOF
version: '3.8'

services:
  # 后端服务
  backend:
    image: adoom2018/read-it-later-backend:${VERSION}
    container_name: read-it-later-backend
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      # 挂载数据库文件，确保数据持久化
      - ./data:/app/data
    environment:
      - GIN_MODE=release
    networks:
      - app-network
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 30s

  # 前端服务
  frontend:
    image: adoom2018/read-it-later-frontend:${VERSION}
    container_name: read-it-later-frontend
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    depends_on:
      - backend
    networks:
      - app-network
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost/"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 30s

networks:
  app-network:
    driver: bridge
EOF

    log_success "docker-compose 配置文件创建完成"
}

# 拉取镜像
pull_images() {
    log_info "拉取 Docker 镜像..."
    
    docker-compose pull
    
    log_success "镜像拉取完成"
}

# 启动服务
start_services() {
    log_info "启动服务..."
    
    docker-compose up -d
    
    log_success "服务启动完成"
}

# 等待服务就绪
wait_for_services() {
    log_info "等待服务就绪..."
    
    # 等待后端服务
    local max_attempts=30
    local attempt=0
    
    while [ $attempt -lt $max_attempts ]; do
        if curl -f http://localhost:8080/health &> /dev/null; then
            log_success "后端服务就绪"
            break
        fi
        
        attempt=$((attempt + 1))
        log_info "等待后端服务就绪... ($attempt/$max_attempts)"
        sleep 2
    done
    
    if [ $attempt -eq $max_attempts ]; then
        log_error "后端服务启动超时"
        return 1
    fi
    
    # 等待前端服务
    attempt=0
    while [ $attempt -lt $max_attempts ]; do
        if curl -f http://localhost/ &> /dev/null; then
            log_success "前端服务就绪"
            break
        fi
        
        attempt=$((attempt + 1))
        log_info "等待前端服务就绪... ($attempt/$max_attempts)"
        sleep 2
    done
    
    if [ $attempt -eq $max_attempts ]; then
        log_error "前端服务启动超时"
        return 1
    fi
}

# 显示服务状态
show_status() {
    log_info "服务状态:"
    docker-compose ps
    
    echo
    log_info "服务访问地址:"
    echo "  前端: http://localhost"
    echo "  后端 API: http://localhost:8080"
    echo "  健康检查: http://localhost:8080/health"
    
    echo
    log_info "管理命令:"
    echo "  查看日志: docker-compose logs -f"
    echo "  停止服务: docker-compose down"
    echo "  重启服务: docker-compose restart"
    echo "  更新服务: docker-compose pull && docker-compose up -d"
}

# 主函数
main() {
    echo "================================================="
    echo "Read It Later - DockerHub 部署脚本"
    echo "================================================="
    echo
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -v|--version)
                VERSION="$2"
                shift 2
                ;;
            -h|--help)
                echo "用法: $0 [选项]"
                echo "选项:"
                echo "  -v, --version VERSION    指定版本号"
                echo "  -h, --help               显示帮助信息"
                exit 0
                ;;
            *)
                log_error "未知参数: $1"
                exit 1
                ;;
        esac
    done
    
    check_requirements
    create_directories
    get_latest_version
    create_docker_compose
    pull_images
    start_services
    wait_for_services
    show_status
    
    echo
    log_success "部署完成！"
    echo "================================================="
}

# 错误处理
trap 'log_error "脚本执行失败"; exit 1' ERR

# 执行主函数
main "$@"
