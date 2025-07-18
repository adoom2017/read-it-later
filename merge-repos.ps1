# PowerShell 脚本 - Git 仓库合并
# 将前端和后端仓库合并到一个新的统一仓库中

Write-Host "=== Read It Later 仓库合并脚本 ===" -ForegroundColor Green

# 配置
$NewRepoName = "read-it-later"
$FrontendRepo = "https://github.com/adoom2017/read-it-later-frontend.git"
$BackendRepo = "https://github.com/adoom2017/read-it-later-backend.git"

# 创建新的合并目录
$MergeDir = "read-it-later-merged"
if (Test-Path $MergeDir) {
    Remove-Item -Recurse -Force $MergeDir
}
New-Item -ItemType Directory -Path $MergeDir
Set-Location $MergeDir

Write-Host "1. 初始化新的 Git 仓库..." -ForegroundColor Yellow
git init
git config user.name "Your Name"
git config user.email "your.email@example.com"

Write-Host "2. 克隆前端仓库..." -ForegroundColor Yellow
git clone $FrontendRepo frontend-temp
Set-Location frontend-temp
$frontendFiles = Get-ChildItem -Force | Where-Object { $_.Name -ne '.git' }
Set-Location ..
New-Item -ItemType Directory -Path "frontend" -Force
Copy-Item -Path "frontend-temp\*" -Destination "frontend\" -Recurse -Force
Remove-Item -Recurse -Force frontend-temp

Write-Host "3. 克隆后端仓库..." -ForegroundColor Yellow
git clone $BackendRepo backend-temp
Set-Location backend-temp
$backendFiles = Get-ChildItem -Force | Where-Object { $_.Name -ne '.git' }
Set-Location ..
New-Item -ItemType Directory -Path "backend" -Force
Copy-Item -Path "backend-temp\*" -Destination "backend\" -Recurse -Force
Remove-Item -Recurse -Force backend-temp

Write-Host "4. 提交合并后的代码..." -ForegroundColor Yellow
git add .
git commit -m "Initial commit: Merge frontend and backend repositories"

Write-Host "5. 设置主分支..." -ForegroundColor Yellow
git branch -M main

Write-Host @"
=== 合并完成！===

下一步操作：
1. 在 GitHub 上创建新的仓库: $NewRepoName
2. 运行以下命令推送到远程仓库：
   git remote add origin https://github.com/adoom2017/$NewRepoName.git
   git push -u origin main

3. 复制 Docker 配置文件到新仓库
4. 更新 README.md 文件
5. 删除旧的两个仓库（可选）

新仓库位于: $((Get-Location).Path)
"@ -ForegroundColor Green
