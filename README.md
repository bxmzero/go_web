# Gin + Uber Fx + SQLite Demo

这个 demo 展示了一个典型的 Go 后端分层：

- `repository` 负责数据库访问
- `service` 负责业务编排
- `handler` 负责 Gin 接口层
- `fx` 负责实例装配和生命周期

## 目录结构

```text
cmd/server/main.go
internal/config
internal/database
internal/httpserver
internal/modules/user
internal/modules/order
```

## 运行

```bash
go mod tidy
go run ./cmd/server
```

默认监听 `:8080`，SQLite 文件在 `./data/demo.db`。

## 接口

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/api/v1/users
curl http://localhost:8080/api/v1/users/1
curl -X POST http://localhost:8080/api/v1/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Charlie","email":"charlie@example.com"}'

curl http://localhost:8080/api/v1/orders
curl http://localhost:8080/api/v1/orders?user_id=1
curl -X POST http://localhost:8080/api/v1/orders \
  -H 'Content-Type: application/json' \
  -d '{"user_id":1,"item":"Laptop Stand","amount":199}'
```

## 设计重点

- `fx.Provide(...)` 负责注册实例构造函数
- `fx.Module(...)` 负责把业务按模块拆开
- `group:"routes"` 负责把各模块路由自动收集起来
- `fx.Lifecycle` 负责数据库关闭和 HTTP 服务优雅退出
