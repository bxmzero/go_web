# Gin + GORM TxManager Demo

这是一个完整的 Go Demo，展示如何在 Gin 分层架构里实现 **由 service 决定是否开启事务** 的 `TxManager` 方案。

## 设计目标

- `router -> handler -> service -> repository` 清晰分层。
- `service` 负责业务编排，并决定某个操作是否进入事务。
- `repository` 只关心 CRUD，不主动开启事务。
- `TxManager` 通过 `context.Context` 透传当前事务，repository 自动获取事务或普通 DB。
- 使用 `Gin + GORM + SQLite` 组成一个可直接运行的完整示例。

## 目录结构

```text
.
├── cmd/server/main.go
├── internal/app/app.go
├── internal/db/database.go
├── internal/handler/user_handler.go
├── internal/model/user.go
├── internal/repository/user_repository.go
├── internal/router/router.go
├── internal/service/user_service.go
└── internal/txmanager/manager.go
```

## 事务控制说明

### 1. Service 决定是否开启事务

- `ListUsers` / `GetUser`：只读场景，不开启事务。
- `CreateUser`：先校验邮箱是否重复，再创建用户，使用事务。
- `UpdateUser`：先查用户、查邮箱冲突，再更新用户，使用事务。
- `DeleteUser`：先查用户是否存在，再删除，使用事务。

### 2. Repository 无感知使用事务

repository 不接收 `*gorm.DB` 参数，也不自己 `Begin/Commit/Rollback`。它只调用：

```go
r.txManager.DB(ctx)
```

如果当前 `ctx` 中存在事务，则自动返回事务对象；否则返回普通连接。

## 运行方式

```bash
go mod tidy
go run ./cmd/server
```

服务启动后监听：

- `http://localhost:8080`

## 接口示例

### 创建用户

```bash
curl -X POST http://localhost:8080/api/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice","email":"alice@example.com","age":25}'
```

### 查询列表

```bash
curl http://localhost:8080/api/users
```

### 查询详情

```bash
curl http://localhost:8080/api/users/<id>
```

### 更新用户

```bash
curl -X PUT http://localhost:8080/api/users/<id> \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice Chen","email":"alice.chen@example.com","age":26}'
```

### 删除用户

```bash
curl -X DELETE http://localhost:8080/api/users/<id>
```

## 事务模板可复用方式

如果后续要扩展订单、库存、账户等模块，建议继续复用以下模式：

1. 每个 repository 都通过 `txManager.DB(ctx)` 获取当前数据库执行器。
2. 每个 service 根据业务决定是否调用 `txManager.WithinTransaction(...)`。
3. 跨多个 repository 的一致性操作统一放在 service 事务闭包中编排。

这样可以保证：

- 事务边界集中在 service。
- repository 保持简单、可测试。
- 后续扩展多表事务时结构仍然稳定。
