# LiteDrive - 私有轻量级云网盘系统

LiteDrive 是一个使用 Go + Gin 实现的私有网盘系统，支持多种存储后端（本地 / Ceph / 腾讯 COS），支持大文件分片上传、秒传、异步备份等功能。前端基于 Next.js 开发，提供简洁的文件管理界面。

## ✨ 特性

- ✅ 分片上传 + 秒传功能（Redis 缓存 SHA 校验）
- ✅ 支持多种存储后端：
    - 本地存储
    - Ceph 对象存储（S3 协议）
    - 腾讯云 COS
- ✅ 支持 RabbitMQ 异步上传任务
- ✅ 自动文件去重，节省存储空间
- ✅ 简洁 UI 界面（Next.js 前端）

## 🚀 快速启动（开发环境）
go version 1.23.2

### 1. 后端运行
```bash
go mod tidy
go run /cmd/main.go
```
### 2. 前端运行
```bash
# 前端暂未上传
```
### 3. 可选：启动 RabbitMQ（异步备份）
使用 docker 启动 rabbitmq 配置文件中指定启用 异步备份
```bash
docker run -d --hostname rabbit --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
```

## 🔧 TODO / 规划中
- [ ] 前端分目录存储结构
- [ ] 文件预览支持（PDF / 图片）
- [ ] 分享链接与过期机制
- [ ] 全平台打包与部署（Docker）

## 📃 License
本项目采用 MIT License,，欢迎初学者学习与改造/(ㄒoㄒ)/~~





