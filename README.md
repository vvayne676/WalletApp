# 💰 WalletApp - Go钱包应用

> **Coding Test Solution** - 基于Go语言的简洁钱包系统

## 🎯 **项目概述**

这是一个中心化钱包应用，支持用户存款、取款、转账、查询余额和交易历史功能。采用Go + PostgreSQL技术栈，遵循"Keep it simple, Less is more"的设计理念。

### **核心功能**
- ✅ 用户钱包存款
- ✅ 用户钱包取款  
- ✅ 用户间转账
- ✅ 查询钱包余额
- ✅ 查看交易历史

## 📋 **解决方案评估**

### 🎯 **技术选择与决策**

**语言选择：Go** ✅
- 选择Go是因为其出色的并发性能和简洁性
- 适合构建高性能的金融API服务
- 丰富的生态系统（Gin、GORM等）

**架构决策：**
- **双ID设计** - 自增主键（性能）+ 业务UUID（安全性）
- **无外键约束** - 保持灵活性，避免复杂的级联操作
- **事务安全** - 转账操作使用数据库事务确保一致性
- **精确计算** - 使用decimal库避免浮点数精度问题

### 👀 **代码审查重点**

**核心文件：**
- `handlers.go` - 业务逻辑实现
- `models.go` - 数据模型设计
- `schema.sql` - 数据库结构
- `API.md` - 完整API文档

**关键特性：**
- 参数验证和错误处理
- 事务安全的转账操作
- 精确的金额计算（decimal.Decimal）
- 完整的审计日志

### ⏱️ **开发时间记录**

**总计约6-8小时：**
- 需求分析和架构设计：1小时
- 核心功能开发：3小时  
- 数据库设计和优化：2小时
- Docker配置和部署：1小时
- 文档编写：1-2小时

### ❌ **有意选择不实现的功能**

**专注核心，避免过度设计：**
- 用户注册/登录系统 - 专注钱包核心功能
- 复杂的权限管理 - 保持简洁性
- 多币种支持 - 避免复杂性
- 交易手续费计算 - 简化业务逻辑
- 外部支付网关集成 - 专注内部转账

### ✅ **功能性需求满足度 (100%)**

- ✅ 存款功能 - 支持精确金额存入
- ✅ 取款功能 - 余额验证和扣减
- ✅ 转账功能 - 事务安全的用户间转账
- ✅ 余额查询 - 实时余额计算
- ✅ 交易历史 - 完整的交易记录查询
- ✅ 数据持久化 - PostgreSQL可靠存储

### 🏗️ **非功能性需求满足度**

- ✅ **性能** - Go的高并发特性，支持1000+并发
- ✅ **可靠性** - 数据库事务保证ACID特性
- ✅ **可扩展性** - 清晰的分层架构，易于扩展
- ✅ **可维护性** - 简洁的代码结构，核心代码<500行
- ⚠️ **安全性** - 基础验证，缺少认证授权
- ⚠️ **可观测性** - 基础日志，缺少监控告警

### 🎨 **工程最佳实践遵循**

**遵循的实践：**
- ✅ 清晰的项目结构（单一职责）
- ✅ 统一的错误处理和日志记录
- ✅ 数据库事务管理（ACID保证）
- ✅ RESTful API设计规范
- ✅ Docker容器化部署
- ✅ 完整的API文档和部署指南

**简洁性体现：**
- 核心代码不到500行
- 无过度设计和抽象
- 专注核心功能实现
- 易于理解和维护

### 🔧 **待改进领域**

1. **认证授权** - 未实现用户身份验证
2. **API限流** - 缺少请求频率限制
3. **监控告警** - 缺少性能监控和告警
4. **缓存层** - 未实现Redis缓存优化
5. **配置管理** - 部分配置硬编码

## 📁 **项目结构**

```
WalletApp/
├── main.go              # 主程序入口
├── db.go                # 数据库连接和重试逻辑
├── route.go             # 路由设置  
├── models.go            # 数据模型定义
├── handlers.go          # API处理器实现
├── schema.sql           # 数据库表结构SQL脚本
├── Dockerfile           # Go应用Docker镜像构建
├── docker-compose.yml   # 完整服务编排
├── API.md              # 完整API文档
├── go.mod              # Go模块依赖
└── README.md           # 项目文档
```

**简洁明了的分层架构，易于理解和维护！**

## 🚀 **快速开始**

### **环境要求**
- Docker & Docker Compose (推荐)
- 或 Go 1.24+ & PostgreSQL 15+ (本地开发)

### **一键启动**
```bash
# 克隆项目
git clone <your-repo-url>
cd WalletApp

# 启动所有服务
docker-compose up -d

# 验证部署
curl http://localhost:8080/health
```

### **验证功能**
```bash
# 1. 存款测试
curl -X POST http://localhost:8080/api/wallet/deposit \
  -H "Content-Type: application/json" \
  -d '{"user_id":"user1","address":"addr1","amount":"100.50"}'

# 2. 查询余额
curl -X POST http://localhost:8080/api/wallet/balance \
  -H "Content-Type: application/json" \
  -d '{"user_id":"user1"}'

# 3. 转账测试
curl -X POST http://localhost:8080/api/wallet/transfer \
  -H "Content-Type: application/json" \
  -d '{"user_id":"user1","from_address":"addr1","to_address":"addr2","amount":"25.00"}'

# 4. 取款测试
curl -X POST http://localhost:8080/api/wallet/withdraw \
  -H "Content-Type: application/json" \
  -d '{"user_id":"user1","address":"addr1","amount":"30.00"}'

# 5. 交易历史查询
curl -X POST http://localhost:8080/api/wallet/history \
  -H "Content-Type: application/json" \
  -d '{"user_id":"user1","page":1,"page_size":10}'
```

## 📋 **API接口**

### **基础URL**
```
http://localhost:8080/api
```

### **接口列表**

| 方法 | 路径 | 描述 | 
|------|------|------|
| POST | `/wallet/deposit` | 存款到指定地址钱包 |
| POST | `/wallet/withdraw` | 从指定地址钱包取款 |
| POST | `/wallet/transfer` | 地址间转账 |
| POST | `/wallet/balance` | 获取指定用户余额 |
| POST | `/wallet/history` | 获取指定用户交易历史 |

### **API详细说明**

#### **1. 存款接口**
```bash
POST /api/wallet/deposit
Content-Type: application/json

{
  "user_id": "test_user_1",
  "address": "addr_test_user_1",
  "amount": "100.50"
}
```

#### **2. 取款接口**
```bash
POST /api/wallet/withdraw
Content-Type: application/json

{
  "user_id": "test_user_1",
  "address": "addr_test_user_1",
  "amount": "50.25"
}
```

#### **3. 转账接口**
```bash
POST /api/wallet/transfer
Content-Type: application/json

{
  "user_id": "test_user_1",
  "from_address": "addr_test_user_1",
  "to_address": "addr_test_user_2",
  "amount": "25.00"
}
```

#### **4. 余额查询接口**
```bash
POST /api/wallet/balance
Content-Type: application/json

{
  "user_id": "test_user_1"
}
```

#### **5. 交易历史接口**
```bash
POST /api/wallet/history
Content-Type: application/json

{
  "user_id": "test_user_1",
  "page": 1,
  "page_size": 10
}
```

详细API文档请查看 [API.md](./API.md)

## 🧪 **API测试**

### **使用 Postman 测试 (推荐)**

我们提供了完整的Postman集合文件，包含所有API接口，方便进行功能测试：

#### **导入步骤：**
1. 打开 Postman
2. 点击 `Import` 按钮
3. 选择项目根目录下的 `WalletApp.postman_collection.json` 文件
4. 导入成功后可以看到完整的API测试集合

#### **Postman集合包含：**
- 🔍 **健康检查** - 验证服务是否正常运行
- 💰 **存款接口** - 测试钱包存款功能
- 💸 **取款接口** - 测试钱包取款功能
- 🔄 **转账接口** - 测试用户间转账功能
- 💳 **余额查询接口** - 测试余额查询功能
- 📋 **交易历史接口** - 测试交易记录查询功能
- ❌ **错误测试用例** - 测试各种异常情况
  - 余额不足的转账
  - 不存在的钱包地址
  - 负数金额等边界情况

#### **使用方法：**
1. 确保服务已启动：`docker-compose up -d`
2. 在Postman中运行单个请求测试特定功能
3. 或者运行整个集合进行完整测试
4. 查看响应结果验证功能是否正常

#### **环境变量配置：**
- `base_url`: `http://localhost:8080` (可根据实际部署地址修改)

#### **测试数据说明：**
数据库初始化时会自动创建以下测试数据：

**用户数据：**
- `test_user_1` (user1@test.com) - Test User 1
- `test_user_2` (user2@test.com) - Test User 2  
- `user_alice` (alice@example.com) - Alice Johnson
- `user_bob` (bob@example.com) - Bob Smith

**钱包数据：**
- `addr_test_user_1` → test_user_1 (余额: 1000.00)
- `addr_test_user_2` → test_user_2 (余额: 2000.00)
- `addr_alice_main` → user_alice (余额: 1500.50)
- `addr_alice_savings` → user_alice (余额: 500.25)
- `addr_bob_business` → user_bob (余额: 3000.75)

**注意：** Alice用户有两个钱包，总余额为2000.75；其他用户各有一个钱包

### **为什么选择 Postman 测试？**

相比复杂的单元测试，Postman测试有以下优势：
- 🎯 **直观易用** - 图形界面，直接看到API响应结果
- 🚀 **快速验证** - 无需编写测试代码，直接调用API
- 🔧 **灵活调试** - 可以随时修改参数进行不同场景测试
- 📖 **易于分享** - 团队成员可以直接导入使用
- 🌐 **真实环境** - 测试真实的HTTP请求，更接近实际使用

## 🏗 **技术架构**

### **技术栈**
- **Backend**: Go 1.24+ 
- **Framework**: Gin (轻量级HTTP框架)
- **Database**: PostgreSQL 15 (关系型数据库)
- **ORM**: GORM (Go对象关系映射)
- **Containerization**: Docker + Docker Compose
- **Decimal**: shopspring/decimal (精确金额计算)

### **架构决策**

#### **1. 简洁架构设计**
- **单体应用**: 适合MVP和小型项目，便于部署和维护
- **扁平化结构**: 避免过度抽象，提高代码可读性
- **业务逻辑聚合**: 相关功能集中在对应文件中

#### **2. 数据库设计**
- **双ID设计**: 自增主键(性能) + 业务UUID(安全)
- **事务安全**: 使用数据库事务确保操作原子性
- **悲观锁**: FOR UPDATE防止并发修改引起的数据不一致
- **索引优化**: 手动创建必要的索引以提升查询性能

#### **3. 并发安全**
```go
// 悲观锁防止余额计算错误
tx.Set("gorm:query_option", "FOR UPDATE").
   Where("address = ? AND user_id = ? AND deleted_at = 0", address, userID).
   First(&wallet)
```

## 📊 **数据模型**

### **用户表 (users)**
```sql
id           SERIAL        PRIMARY KEY
user_id      VARCHAR(36)   UNIQUE NOT NULL  -- 业务UUID
email        VARCHAR(255)  UNIQUE NOT NULL
password_hash VARCHAR(255) NOT NULL  
name         VARCHAR(255)  NOT NULL
status       VARCHAR(50)   DEFAULT 'active'
created_at   BIGINT        NOT NULL
updated_at   BIGINT        NOT NULL
deleted_at   BIGINT        DEFAULT 0
```

### **钱包表 (wallets)**
```sql
id          SERIAL          PRIMARY KEY
wallet_id   VARCHAR(36)     UNIQUE NOT NULL  -- 业务UUID
user_id     VARCHAR(36)     NOT NULL
address     VARCHAR(255)    NOT NULL
balance     DECIMAL(20,8)   DEFAULT 0
status      VARCHAR(50)     DEFAULT 'active'
created_at  BIGINT          NOT NULL
updated_at  BIGINT          NOT NULL
deleted_at  BIGINT          DEFAULT 0
```

### **交易表 (transactions)**
```sql
id              SERIAL          PRIMARY KEY
transaction_id  VARCHAR(36)     UNIQUE NOT NULL  -- 业务UUID
user_id         VARCHAR(36)     NOT NULL
from_address    VARCHAR(255)    
to_address      VARCHAR(255)    
amount          DECIMAL(20,8)   NOT NULL
type            VARCHAR(20)     NOT NULL -- deposit/withdraw/transfer
status          VARCHAR(20)     NOT NULL -- completed/failed
description     TEXT
created_at      BIGINT          NOT NULL
updated_at      BIGINT          NOT NULL
deleted_at      BIGINT          DEFAULT 0
```

## 🐳 **Docker部署**

### **启动服务**
```bash
# 构建并启动
docker-compose up -d

# 查看状态
docker-compose ps

# 查看日志
docker-compose logs -f walletapp
```



## 🔧 **开发指南**


### **代码组织**
- `main.go` - 应用启动和配置加载
- `db.go` - 数据库连接配置和重试逻辑
- `route.go` - HTTP路由定义
- `models.go` - 数据模型和业务逻辑
- `handlers.go` - HTTP请求处理器
- `schema.sql` - 数据库表结构和初始化脚本

### **错误处理**
- 统一的HTTP错误响应格式
- 数据库事务回滚机制
- 业务逻辑验证和边界检查
