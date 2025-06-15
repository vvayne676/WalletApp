# 💰 WalletApp API 文档

## 基础信息
```
Base URL: http://localhost:8080
API 版本: v1
```

## 🚀 API 端点

### 1. 健康检查
**GET** `/health`

检查服务是否正常运行。

**响应:**
```json
{
  "status": "ok",
  "time": "2024-01-01T12:00:00Z"
}
```

---

### 2. 存款
**POST** `/api/wallet/deposit`

向钱包存入资金。

**请求体:**
```json
{
  "amount": "100.50",
  "user_id": "test_user_1",
  "address": "addr_test_user_1"
}
```

**请求参数说明:**
- `amount`: 存款金额 (decimal, 必填, 必须为正数)
- `user_id`: 用户ID (string, 必填)
- `address`: 钱包地址 (string, 必填)

**成功响应:**
```json
{
  "message": "Deposit successful",
  "transaction": {
    "id": 1,
    "transaction_id": "550e8400-e29b-41d4-a716-446655440000",
    "to_address": "addr_test_user_1",
    "amount": "100.50000000",
    "type": "deposit",
    "status": "completed",
    "description": "Deposit to wallet",
    "created_at": 1704110400,
    "updated_at": 1704110400,
    "deleted_at": 0
  }
}
```

**错误响应:**
```json
{
  "error": "Amount must be positive"
}
```

---

### 3. 取款
**POST** `/api/wallet/withdraw`

从钱包取出资金。

**请求体:**
```json
{
  "amount": "50.25",
  "user_id": "test_user_1", 
  "address": "addr_test_user_1"
}
```

**请求参数说明:**
- `amount`: 取款金额 (decimal, 必填, 必须为正数)
- `user_id`: 用户ID (string, 必填)
- `address`: 钱包地址 (string, 必填)

**成功响应:**
```json
{
  "message": "Withdraw successful",
  "transaction": {
    "id": 2,
    "transaction_id": "550e8400-e29b-41d4-a716-446655440001",
    "from_address": "addr_test_user_1",
    "amount": "50.25000000",
    "type": "withdraw",
    "status": "completed",
    "description": "Withdraw from wallet",
    "created_at": 1704110460,
    "updated_at": 1704110460,
    "deleted_at": 0
  }
}
```

**错误响应:**
```json
{
  "error": "Insufficient balance"
}
```

---

### 4. 转账
**POST** `/api/wallet/transfer`

在不同地址之间转账。

**请求体:**
```json
{
  "user_id": "test_user_1",
  "from_address": "addr_test_user_1",
  "to_address": "addr_test_user_2",
  "amount": "25.00"
}
```

**请求参数说明:**
- `user_id`: 发起转账的用户ID (string, 必填)
- `from_address`: 转出地址 (string, 必填)
- `to_address`: 转入地址 (string, 必填)
- `amount`: 转账金额 (decimal, 必填, 必须为正数)

**成功响应:**
```json
{
  "message": "Transfer successful",
  "transaction": {
    "id": 3,
    "transaction_id": "550e8400-e29b-41d4-a716-446655440002",
    "from_address": "addr_test_user_1",
    "to_address": "addr_test_user_2",
    "amount": "25.00000000",
    "type": "transfer",
    "status": "completed",
    "description": "Transfer to user",
    "created_at": 1704110520,
    "updated_at": 1704110520,
    "deleted_at": 0
  }
}
```

**错误响应:**
```json
{
  "error": "Cannot transfer to the same address"
}
```

```json
{
  "error": "Wallet not found or not owned by user"
}
```

---

### 5. 查询余额
**POST** `/api/wallet/balance`

查询用户钱包余额总计。

**请求体:**
```json
{
  "user_id": "test_user_1"
}
```

**请求参数说明:**
- `user_id`: 用户ID (string, 必填)

**成功响应:**
```json
{
  "balance": {
    "user_id": "test_user_1",
    "balance": "125.25"
  }
}
```

**说明:**
- 返回该用户所有钱包余额的总和
- 只计算未删除的钱包 (`deleted_at = 0`)
- 余额通过钱包表的 `balance` 字段直接计算

---

### 6. 交易历史
**POST** `/api/wallet/history`

查询用户的交易历史记录，支持分页。

**请求体:**
```json
{
  "user_id": "test_user_1",
  "page": 1,
  "page_size": 20
}
```

**请求参数说明:**
- `user_id`: 用户ID (string, 必填)
- `page`: 页码 (int, 可选, 默认为1)
- `page_size`: 每页条数 (int, 可选, 默认为20, 最大为100)

**成功响应:**
```json
{
  "data": {
    "transactions": [
      {
        "id": 3,
        "transaction_id": "550e8400-e29b-41d4-a716-446655440002",
        "from_address": "addr_test_user_1",
        "to_address": "addr_test_user_2", 
        "amount": "25.00000000",
        "type": "transfer",
        "status": "completed",
        "description": "Transfer to user",
        "created_at": 1704110520,
        "updated_at": 1704110520,
        "deleted_at": 0
      }
    ],
    "total": 3,
    "page": 1,
    "page_size": 20
  }
}
```

**说明:**
- 按创建时间倒序排列 (`ORDER BY created_at DESC`)
- 只返回未删除的交易记录 (`deleted_at = 0`)
- 查询用户所有钱包地址相关的交易 (from_address 或 to_address 匹配用户钱包)
- `total` 字段返回符合条件的总记录数
- 时间戳为Unix时间戳格式

---

## 🛠️ 数据模型

### User (用户)
```json
{
  "id": "number",           // 自增主键ID
  "user_id": "string",      // 业务用户ID (UUID)
  "name": "string",         // 用户名
  "email": "string",        // 邮箱
  "password_hash": "string", // 密码哈希 (不在API响应中返回)
  "status": "string",       // 状态 (active/inactive)
  "created_at": "number",   // 创建时间 (Unix时间戳)
  "updated_at": "number",   // 更新时间 (Unix时间戳)
  "deleted_at": "number"    // 删除时间 (0表示未删除)
}
```

### Wallet (钱包)
```json
{
  "id": "number",           // 自增主键ID
  "wallet_id": "string",    // 业务钱包ID (UUID)
  "user_id": "string",      // 用户业务ID (关联User.user_id)
  "balance": "decimal",     // 余额 (精度20,8)
  "status": "string",       // 状态 (active/inactive)
  "version": "number",      // 版本号
  "address": "string",      // 钱包地址 (唯一)
  "created_at": "number",   // 创建时间 (Unix时间戳)
  "updated_at": "number",   // 更新时间 (Unix时间戳)
  "deleted_at": "number"    // 删除时间 (0表示未删除)
}
```

### Transaction (交易)
```json
{
  "id": "number",              // 自增主键ID
  "transaction_id": "string",  // 业务交易ID (UUID)
  "from_address": "string",    // 转出地址 (可选)
  "to_address": "string",      // 转入地址 (可选) 
  "amount": "decimal",         // 交易金额 (精度20,8)
  "type": "string",           // 交易类型: deposit, withdraw, transfer
  "status": "string",         // 交易状态: pending, completed, failed, cancelled
  "description": "string",    // 交易描述
  "created_at": "number",     // 创建时间 (Unix时间戳)
  "updated_at": "number",     // 更新时间 (Unix时间戳)
  "deleted_at": "number"      // 删除时间 (0表示未删除)
}
```

---

## ⚠️ 错误处理

### 常见错误状态码:
- `400 Bad Request`: 请求参数错误
- `500 Internal Server Error`: 服务器内部错误

### 错误响应格式:
```json
{
  "error": "错误信息描述"
}
```

### 常见错误信息:
- `"Amount must be positive"` - 金额必须为正数
- `"Source wallet not found"` - 源钱包未找到
- `"Destination wallet not found"` - 目标钱包未找到
- `"Insufficient balance"` - 余额不足
- `"Cannot transfer to the same address"` - 不能转账到相同地址
- `"Wallet not found or not owned by user"` - 钱包未找到或不属于该用户
- `"Failed to create transaction"` - 创建交易失败
- `"Database error"` - 数据库错误

---

## 🔄 交易流程

### 存款流程:
1. 生成业务交易ID (UUID)
2. 创建待处理交易记录 (`status: "pending"`)
3. 更新钱包余额 (增加)
4. 更新交易状态为完成 (`status: "completed"`)
5. 返回交易结果

### 取款流程:
1. 生成业务交易ID (UUID)
2. 创建待处理交易记录 (`status: "pending"`)
3. 检查余额是否充足
4. 更新钱包余额 (减少)
5. 更新交易状态为完成 (`status: "completed"`)
6. 返回交易结果

### 转账流程:
1. 验证用户权限和地址有效性
2. 生成业务交易ID (UUID)
3. 创建待处理交易记录 (`status: "pending"`)
4. 更新转出钱包余额 (减少)
5. 更新转入钱包余额 (增加)
6. 更新交易状态为完成 (`status: "completed"`)
7. 返回交易结果

---

## 🧪 测试示例

### 使用 curl 测试:

```bash
# 健康检查
curl http://localhost:8080/health

# 存款
curl -X POST http://localhost:8080/api/wallet/deposit \
  -H "Content-Type: application/json" \
  -d '{"amount":"100.00","user_id":"test_user_1","address":"addr_test_user_1"}'

# 查询余额
curl -X POST http://localhost:8080/api/wallet/balance \
  -H "Content-Type: application/json" \
  -d '{"user_id":"test_user_1"}'

# 取款
curl -X POST http://localhost:8080/api/wallet/withdraw \
  -H "Content-Type: application/json" \
  -d '{"amount":"50.00","user_id":"test_user_1","address":"addr_test_user_1"}'

# 转账
curl -X POST http://localhost:8080/api/wallet/transfer \
  -H "Content-Type: application/json" \
  -d '{"user_id":"test_user_1","from_address":"addr_test_user_1","to_address":"addr_test_user_2","amount":"25.00"}'

# 交易历史 (所有交易)
curl -X POST http://localhost:8080/api/wallet/history \
  -H "Content-Type: application/json" \
  -d '{"user_id":"test_user_1","page":1,"page_size":10}'

# 交易历史
curl -X POST http://localhost:8080/api/wallet/history \
  -H "Content-Type: application/json" \
  -d '{"user_id":"test_user_1","page":1,"page_size":10}'
```

---

## 📝 注意事项

1. **双重ID设计**: 每个表都有自增主键ID和业务ID (UUID)，提供更好的性能和灵活性
2. **地址验证**: 转账时会验证用户是否拥有转出地址的权限
3. **余额锁定**: 使用 `FOR UPDATE` 锁定钱包记录防止并发问题
4. **事务处理**: 所有涉及余额变更的操作都使用数据库事务
5. **精度处理**: 金额使用 `decimal.Decimal` 类型确保精度，数据库精度为 (20,8)
6. **软删除**: 所有记录使用软删除 (`deleted_at = 0` 表示未删除)
7. **UUID生成**: 业务ID使用 `github.com/google/uuid` 生成
8. **无外键约束**: 钱包表不使用外键约束，通过业务逻辑保证数据一致性
9. **时间戳格式**: 使用Unix时间戳 (BIGINT) 存储时间

---

## 🔧 技术栈

- **Web框架**: Gin
- **数据库ORM**: GORM  
- **数据库**: PostgreSQL
- **精度处理**: shopspring/decimal
- **UUID生成**: google/uuid
- **容器化**: Docker + Docker Compose 

---

## 🗄️ 数据库索引

为了提高查询性能，数据库中创建了以下索引：

### Users 表
- `idx_users_user_id`: user_id 字段唯一索引
- `idx_users_email`: email 字段唯一索引

### Wallets 表  
- `idx_wallets_wallet_id`: wallet_id 字段唯一索引
- `idx_wallets_user_id`: user_id 字段索引
- `idx_wallets_address`: address 字段唯一索引

### Transactions 表
- `idx_transactions_transaction_id`: transaction_id 字段唯一索引
- `idx_transactions_from_address`: from_address 字段索引
- `idx_transactions_to_address`: to_address 字段索引
- `idx_transactions_status`: status 字段索引
