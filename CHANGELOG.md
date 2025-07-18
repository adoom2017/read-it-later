# 更新日志

本文档记录了 Read It Later 项目的所有重要更改。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
并且本项目遵循 [语义化版本](https://semver.org/lang/zh-CN/) 规范。

## [Unreleased]

### 新增
- 开源项目配置，采用 MIT 许可证
- 贡献指南文档
- 完善的项目文档

## [1.0.0] - 2025-01-18

### 新增
- 🎉 初始版本发布
- 📖 保存网页文章链接功能
- 🔍 自动提取文章内容和摘要
- 🏷️ 标签管理系统
- 📱 响应式前端界面
- 🐳 Docker 容器化部署
- 💾 SQLite 数据库支持
- 🔍 文章搜索和过滤功能
- ⚡ 基于 React 18 的现代前端
- 🚀 基于 Go 1.24 的高性能后端
- 🌐 RESTful API 设计
- 🔧 自动化部署脚本
- 📋 健康检查端点
- 🔄 CORS 跨域支持

### 技术栈
- **前端**: React 18, Vite, Axios
- **后端**: Go 1.24, Gin, SQLite
- **部署**: Docker, Docker Compose, Nginx

### API 端点
- `GET /api/articles` - 获取文章列表
- `POST /api/articles` - 添加新文章
- `GET /api/articles/:id` - 获取文章详情
- `DELETE /api/articles/:id` - 删除文章
- `POST /api/articles/:id/tags` - 添加标签
- `GET /` - 后端健康检查
- `GET /health` - 服务健康状态

## [0.1.0] - 2025-01-15

### 新增
- 项目初始化
- 基础项目结构
- 开发环境配置

---

## 版本说明

### 版本格式
- **主版本号**: 当做了不兼容的 API 修改
- **次版本号**: 当做了向下兼容的功能性新增
- **修订号**: 当做了向下兼容的问题修正

### 标签含义
- 🎉 **重大更新**
- ✨ **新功能**
- 🐛 **Bug 修复**
- 📚 **文档更新**
- 🔧 **配置更改**
- ⚡ **性能优化**
- 🚨 **重要变更**
- 🗑️ **废弃功能**
- 🔒 **安全修复**
