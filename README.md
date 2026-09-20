# Orbit Board

> [!CAUTION]
> **不要在生产环境中部署本项目。DO NOT DEPLOY THIS PROJECT TO PRODUCTION.**

Orbit Board 是一个轻量的团队任务看板，使用 Go、Gin、React、shadcn/ui 风格组件和 Tailwind CSS 构建。它提供项目筛选、任务状态流转、优先级管理、评论、工作量统计和数据导出功能。

本仓库面向本地演示和开发环境，使用内置示例数据和简化配置。启动后即可体验完整的任务管理流程。

## 技术栈

- Go 1.23 与 Gin
- React 18 与 Vite
- shadcn/ui 风格组件
- Tailwind CSS
- JSON 文件存储

## 本地运行

需要 Go 1.23+ 和 Node.js 20+。

```bash
# 终端 1：启动 API，默认监听 http://localhost:8080
go mod download
go run ./cmd/server

# 终端 2：启动前端，默认监听 http://localhost:5173
cd frontend
npm install
npm run dev
```

首次启动 API 时会加载演示数据，后续修改写入 `data/board.json`。Vite 开发服务器会自动把 `/api` 请求代理到 Gin 服务。

也可以使用自定义端口启动 API：

```bash
PORT=18080 go run ./cmd/server
```

## 测试和构建

```bash
make test
make build
```

前端生产构建输出到 `frontend/dist`。

## API

- `GET /healthz`：服务健康状态
- `GET /api/dashboard`：看板数据和统计信息
- `POST /api/tasks`：创建任务
- `PATCH /api/tasks/:id`：更新任务
- `DELETE /api/tasks/:id`：删除任务
- `POST /api/comments`：添加任务评论
- `GET /api/report`：团队和项目报告
- `GET /api/export`：导出看板数据
- `GET /api/files/*path`：读取导出文件

## 项目结构

```text
cmd/server/       API 启动入口
internal/app/     数据模型、存储和 HTTP 接口
frontend/src/     React 页面与 UI 组件
data/             本地运行数据
```

## 许可

MIT
