# Go Bad Example

> [!CAUTION]
> **不要在生产环境中部署本项目。DO NOT DEPLOY THIS PROJECT TO PRODUCTION.**

这是一个专门为 [MoonCode](https://github.com/fuchencong/mooncode) AI 代码质量评估平台准备的、**故意写坏**的全栈样例。它是一个小型团队任务看板，技术栈为 Go + Gin、React、shadcn/ui 风格组件和 Tailwind CSS。

项目可以编译和运行，但代码中刻意散布了死代码、重复逻辑、糟糕命名、高复杂度函数、全局状态、紧耦合设计和安全漏洞。它不是教学范例，也不应该被复制到真实项目。

## 运行

需要 Go 1.23+ 和 Node.js 20+。

```bash
# 终端 1：后端 http://localhost:8080
go mod download
go run ./cmd/server

# 终端 2：前端 http://localhost:5173
cd frontend
npm install
npm run dev
```

首次启动后端时会创建内存中的演示数据，并把修改写入 `data/board.json`。前端开发服务器会把 `/api` 代理到 Gin。

## 测试和构建

```bash
make test
make build
```

## 故意保留的问题

以下内容是评估素材，不是遗漏：

- Go：巨型 handler、复制粘贴的校验和统计逻辑、含义模糊的变量名、未使用的导出函数、全局可变状态、数据层和 HTTP 层耦合、忽略错误、低效循环。
- Go 安全：硬编码管理密钥、宽松 CORS、明文密码、弱鉴权、用户可控路径读取、过度详细的错误响应。
- JavaScript：巨型 React 组件、重复过滤和状态更新、prop drilling、魔法字符串、死函数、直接操作 DOM、缺少一致的错误处理。
- JavaScript 安全：把管理密钥放在 `localStorage`、使用 `dangerouslySetInnerHTML` 渲染用户输入、未验证的 URL 跳转。

请只在隔离的本地测试环境使用。项目中的漏洞能够泄露本机文件或执行危险的前端行为，不应暴露到公网。

## API 摘要

- `GET /api/dashboard`：看板数据和统计
- `POST /api/tasks`：创建任务
- `PATCH /api/tasks/:id`：修改任务
- `DELETE /api/tasks/:id`：删除任务（弱管理鉴权）
- `POST /api/comments`：添加评论
- `GET /api/export`：导出数据
- `GET /api/files/*path`：读取文件（故意包含路径穿越漏洞）
- `GET /healthz`：健康检查

## 许可

MIT。请牢记：代码中的坏实践与漏洞是项目的核心测试数据。
