# 贡献指南

感谢您对 Read It Later 项目的关注！我们欢迎任何形式的贡献。

## 贡献方式

### 🐛 报告问题
- 在 [Issues](https://github.com/adoom2017/read-it-later/issues) 中报告 bug
- 提供详细的复现步骤和环境信息
- 包含错误截图或日志（如果有）

### 💡 功能建议
- 在 [Issues](https://github.com/adoom2017/read-it-later/issues) 中提出新功能建议
- 描述功能的使用场景和预期效果
- 考虑功能的可行性和复杂度

### 🔧 代码贡献
- Fork 本仓库
- 创建功能分支
- 提交代码并创建 Pull Request

## 开发流程

### 1. 准备开发环境

```bash
# 克隆你的 fork 仓库
git clone https://github.com/YOUR_USERNAME/read-it-later.git
cd read-it-later

# 添加上游仓库
git remote add upstream https://github.com/adoom2017/read-it-later.git

# 创建开发分支
git checkout -b feature/your-feature-name
```

### 2. 本地开发

#### 后端开发
```bash
cd backend
go mod tidy
go run main.go
```

#### 前端开发
```bash
cd frontend
npm install
npm start
```

#### 使用 Docker 开发
```bash
# 构建并启动服务
docker-compose up -d --build

# 查看日志
docker-compose logs -f
```

### 3. 代码规范

#### Go 后端
- 使用 `go fmt` 格式化代码
- 遵循 Go 语言官方代码规范
- 添加适当的注释和文档
- 确保通过 `go vet` 检查

#### React 前端
- 使用 2 空格缩进
- 使用 JSX 语法
- 组件命名使用 PascalCase
- 文件名使用 PascalCase

### 4. 提交代码

```bash
# 提交更改
git add .
git commit -m "feat: 添加新功能描述"

# 推送到你的 fork 仓库
git push origin feature/your-feature-name
```

### 5. 创建 Pull Request

1. 在 GitHub 上创建 Pull Request
2. 填写 PR 模板中的信息
3. 等待代码审查和反馈
4. 根据反馈修改代码

## 提交信息规范

使用约定式提交格式：

```
type(scope): 简短描述

详细描述（可选）

关闭的 issue（可选）
```

### 类型 (type)
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式修改
- `refactor`: 重构代码
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

### 示例
```bash
git commit -m "feat: 添加文章搜索功能"
git commit -m "fix: 修复数据库连接超时问题"
git commit -m "docs: 更新 API 文档"
```

## 测试

### 运行测试
```bash
# 后端测试
cd backend
go test ./...

# 前端测试
cd frontend
npm test
```

### 添加测试
- 为新功能添加相应的单元测试
- 确保测试覆盖率不降低
- 测试应该独立且可重复运行

## 代码审查

### 审查检查点
- 代码是否符合项目规范
- 是否有充分的测试覆盖
- 是否有适当的文档说明
- 是否考虑了性能和安全性
- 是否与现有功能兼容

### 响应审查反馈
- 及时回复审查意见
- 根据反馈修改代码
- 保持友好和专业的沟通

## 发布流程

1. 版本号遵循语义化版本规范 (SemVer)
2. 更新 CHANGELOG.md
3. 创建 release tag
4. 更新 Docker 镜像

## 问题和支持

如果您在贡献过程中遇到任何问题，请通过以下方式联系：

- 创建 [Issue](https://github.com/adoom2017/read-it-later/issues)
- 发送邮件至项目维护者

## 行为准则

请遵守以下行为准则：

- 保持友好和尊重的态度
- 欢迎新贡献者
- 专注于项目本身
- 避免个人攻击或不当言论

感谢您的贡献！🎉
