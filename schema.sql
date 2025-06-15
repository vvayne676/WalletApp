-- 💰 WalletApp 数据库表结构
-- PostgreSQL 数据库初始化脚本
-- 包含用户表、钱包表和交易表，但钱包表无外键约束
-- 
-- 注意：此脚本会在 walletapp 数据库中执行
-- docker-compose.yml 中的 POSTGRES_DB=walletapp 会自动创建数据库

-- ================================
-- 1. 创建用户表
-- ================================
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,                          -- 自增主键
    user_id VARCHAR(36) UNIQUE NOT NULL,               -- 业务用户ID  
    email VARCHAR(255) UNIQUE NOT NULL,                -- 邮箱
    password_hash VARCHAR(255) NOT NULL,               -- 密码哈希
    name VARCHAR(255) NOT NULL,                        -- 用户名
    status VARCHAR(50) NOT NULL DEFAULT 'active',      -- 状态
    created_at BIGINT NOT NULL,                        -- 创建时间戳
    updated_at BIGINT NOT NULL,                        -- 更新时间戳
    deleted_at BIGINT DEFAULT 0                        -- 删除时间戳(0=未删除)
);

-- ================================
-- 2. 创建钱包表 (无外键约束)
-- ================================
CREATE TABLE IF NOT EXISTS wallets (
    id BIGSERIAL PRIMARY KEY,                          -- 自增主键
    wallet_id VARCHAR(36) UNIQUE NOT NULL,             -- 业务钱包ID
    user_id VARCHAR(36) NOT NULL,                      -- 关联用户ID
    balance DECIMAL(20,8) NOT NULL DEFAULT 0,          -- 余额
    status VARCHAR(50) NOT NULL DEFAULT 'active',      -- 状态
    version BIGINT NOT NULL DEFAULT 0,                 -- 版本号
    address VARCHAR(255) NOT NULL UNIQUE,              -- 钱包地址
    created_at BIGINT NOT NULL,                        -- 创建时间戳
    updated_at BIGINT NOT NULL,                        -- 更新时间戳
    deleted_at BIGINT DEFAULT 0                        -- 删除时间戳(0=未删除)
    -- 注意：故意不添加外键约束，保持灵活性
);

-- ================================
-- 3. 创建交易表
-- ================================
CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,                          -- 自增主键
    transaction_id VARCHAR(36) UNIQUE NOT NULL,        -- 业务交易ID
    from_address VARCHAR(255),                         -- 转出地址
    to_address VARCHAR(255),                           -- 转入地址
    amount DECIMAL(20,8) NOT NULL,                     -- 交易金额
    type VARCHAR(20) NOT NULL CHECK (type IN ('deposit', 'withdraw', 'transfer')), -- 交易类型
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'completed', 'failed', 'cancelled')), -- 交易状态
    description TEXT,                                  -- 交易描述
    created_at BIGINT NOT NULL,                        -- 创建时间戳
    updated_at BIGINT NOT NULL,                        -- 更新时间戳
    deleted_at BIGINT DEFAULT 0                        -- 删除时间戳(0=未删除)
);

-- ================================
-- 4. 创建索引 (性能优化)
-- ================================

-- 用户表索引
CREATE INDEX IF NOT EXISTS idx_users_user_id ON users(user_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- 钱包表索引
CREATE INDEX IF NOT EXISTS idx_wallets_wallet_id ON wallets(wallet_id);
CREATE INDEX IF NOT EXISTS idx_wallets_user_id ON wallets(user_id);
CREATE INDEX IF NOT EXISTS idx_wallets_address ON wallets(address);
CREATE INDEX IF NOT EXISTS idx_wallets_status ON wallets(status);
CREATE INDEX IF NOT EXISTS idx_wallets_deleted_at ON wallets(deleted_at);

-- 交易表索引
CREATE INDEX IF NOT EXISTS idx_transactions_transaction_id ON transactions(transaction_id);
CREATE INDEX IF NOT EXISTS idx_transactions_from_address ON transactions(from_address);
CREATE INDEX IF NOT EXISTS idx_transactions_to_address ON transactions(to_address);
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);
CREATE INDEX IF NOT EXISTS idx_transactions_deleted_at ON transactions(deleted_at);

-- 复合索引 (查询优化)
CREATE INDEX IF NOT EXISTS idx_transactions_from_created ON transactions(from_address, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_to_created ON transactions(to_address, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_status_created ON transactions(status, created_at DESC);

-- ================================
-- 5. 插入测试数据 (可选)
-- ================================

-- 插入测试用户
INSERT INTO users (user_id, email, password_hash, name, status, created_at, updated_at) 
VALUES 
    ('test_user_1', 'user1@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye6K8n8YLaKHGW3rGS4VNa8mKF1wGOTgG', 'Test User 1', 'active', extract(epoch from now())::bigint, extract(epoch from now())::bigint),
    ('test_user_2', 'user2@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye6K8n8YLaKHGW3rGS4VNa8mKF1wGOTgG', 'Test User 2', 'active', extract(epoch from now())::bigint, extract(epoch from now())::bigint),
    ('user_alice', 'alice@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye6K8n8YLaKHGW3rGS4VNa8mKF1wGOTgG', 'Alice Johnson', 'active', extract(epoch from now())::bigint, extract(epoch from now())::bigint),
    ('user_bob', 'bob@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye6K8n8YLaKHGW3rGS4VNa8mKF1wGOTgG', 'Bob Smith', 'active', extract(epoch from now())::bigint, extract(epoch from now())::bigint)
ON CONFLICT (user_id) DO NOTHING;

-- 插入测试钱包
INSERT INTO wallets (wallet_id, user_id, balance, status, version, address, created_at, updated_at)
VALUES 
    ('550e8400-e29b-41d4-a716-446655440001', 'test_user_1', 1000.00000000, 'active', 0, 'addr_test_user_1', extract(epoch from now())::bigint, extract(epoch from now())::bigint),
    ('550e8400-e29b-41d4-a716-446655440002', 'test_user_2', 2000.00000000, 'active', 0, 'addr_test_user_2', extract(epoch from now())::bigint, extract(epoch from now())::bigint),
    ('550e8400-e29b-41d4-a716-446655440003', 'user_alice', 1500.50000000, 'active', 0, 'addr_alice_main', extract(epoch from now())::bigint, extract(epoch from now())::bigint),
    ('550e8400-e29b-41d4-a716-446655440004', 'user_alice', 500.25000000, 'active', 0, 'addr_alice_savings', extract(epoch from now())::bigint, extract(epoch from now())::bigint),
    ('550e8400-e29b-41d4-a716-446655440005', 'user_bob', 3000.75000000, 'active', 0, 'addr_bob_business', extract(epoch from now())::bigint, extract(epoch from now())::bigint)
ON CONFLICT (wallet_id) DO NOTHING;

-- 插入测试交易数据
INSERT INTO transactions (transaction_id, from_address, to_address, amount, type, status, description, created_at, updated_at, deleted_at)
VALUES 
    ('550e8400-e29b-41d4-a716-446655440001', NULL, 'addr_test_user_1', 1000.00000000, 'deposit', 'completed', 'Initial deposit for test user 1', extract(epoch from now())::bigint, extract(epoch from now())::bigint, 0),
    ('550e8400-e29b-41d4-a716-446655440002', NULL, 'addr_test_user_2', 2000.00000000, 'deposit', 'completed', 'Initial deposit for test user 2', extract(epoch from now())::bigint, extract(epoch from now())::bigint, 0),
    ('550e8400-e29b-41d4-a716-446655440003', NULL, 'addr_alice_main', 1500.50000000, 'deposit', 'completed', 'Alice main wallet deposit', extract(epoch from now())::bigint, extract(epoch from now())::bigint, 0),
    ('550e8400-e29b-41d4-a716-446655440004', NULL, 'addr_alice_savings', 500.25000000, 'deposit', 'completed', 'Alice savings wallet deposit', extract(epoch from now())::bigint, extract(epoch from now())::bigint, 0),
    ('550e8400-e29b-41d4-a716-446655440005', NULL, 'addr_bob_business', 3000.75000000, 'deposit', 'completed', 'Bob business wallet deposit', extract(epoch from now())::bigint, extract(epoch from now())::bigint, 0),
    ('550e8400-e29b-41d4-a716-446655440006', 'addr_alice_main', 'addr_test_user_1', 100.00000000, 'transfer', 'completed', 'Transfer from Alice to Test User 1', extract(epoch from now())::bigint, extract(epoch from now())::bigint, 0),
    ('550e8400-e29b-41d4-a716-446655440007', 'addr_test_user_2', NULL, 50.00000000, 'withdraw', 'completed', 'Withdrawal from test user 2', extract(epoch from now())::bigint, extract(epoch from now())::bigint, 0)
ON CONFLICT (transaction_id) DO NOTHING;

-- ================================
-- 6. 验证表创建
-- ================================

-- 查看创建的表
SELECT 
    table_name,
    table_type
FROM information_schema.tables 
WHERE table_schema = 'public' 
    AND table_name IN ('users', 'wallets', 'transactions')
ORDER BY table_name;

-- 查看表结构
\d users;
\d wallets;
\d transactions;

-- 验证约束和索引
SELECT 
    indexname,
    tablename,
    indexdef
FROM pg_indexes 
WHERE tablename IN ('users', 'wallets', 'transactions')
ORDER BY tablename, indexname;

-- ================================
-- 7. 验证数据模型一致性
-- ================================

-- 检查表结构是否与 models.go 一致
SELECT 
    table_name,
    column_name,
    data_type,
    is_nullable,
    column_default
FROM information_schema.columns 
WHERE table_schema = 'public' 
    AND table_name IN ('users', 'wallets', 'transactions')
ORDER BY table_name, ordinal_position;

-- 验证测试数据
SELECT 'Users' as table_name, count(*) as record_count FROM users
UNION ALL
SELECT 'Wallets' as table_name, count(*) as record_count FROM wallets
UNION ALL
SELECT 'Transactions' as table_name, count(*) as record_count FROM transactions;

-- 查看用户数据
SELECT 
    id,
    user_id,
    email,
    name,
    status,
    to_timestamp(created_at) as created_at
FROM users
ORDER BY id DESC
LIMIT 10;

-- 查看钱包数据
SELECT 
    id,
    wallet_id,
    user_id,
    balance,
    status,
    address,
    to_timestamp(created_at) as created_at
FROM wallets
ORDER BY id DESC
LIMIT 10;

-- 查看交易数据
SELECT 
    id,
    transaction_id,
    from_address,
    to_address,
    amount,
    type,
    status,
    to_timestamp(created_at) as created_at
FROM transactions
ORDER BY id DESC
LIMIT 10; 