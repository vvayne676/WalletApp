package main

import (
	"github.com/shopspring/decimal"
)

// User 用户模型
type User struct {
	ID           int64  `json:"id" gorm:"primaryKey;autoIncrement"`                       // 自增主键
	UserID       string `json:"user_id" gorm:"unique;not null;type:varchar(36)"`          // 业务用户ID
	Email        string `json:"email" gorm:"unique;not null;type:varchar(255)"`           // 邮箱
	PasswordHash string `json:"-" gorm:"not null;type:varchar(255)"`                      // 密码哈希
	Name         string `json:"name" gorm:"not null;type:varchar(255)"`                   // 用户名
	Status       string `json:"status" gorm:"not null;default:'active';type:varchar(50)"` // 状态
	CreatedAt    int64  `json:"created_at" gorm:"not null"`                               // 创建时间
	UpdatedAt    int64  `json:"updated_at" gorm:"not null"`                               // 更新时间
	DeletedAt    int64  `json:"deleted_at,omitempty" gorm:"index"`                        // 删除时间
}

// Wallet 钱包模型
type Wallet struct {
	ID        int64           `json:"id" gorm:"primaryKey;autoIncrement"`                       // 自增主键
	WalletID  string          `json:"wallet_id" gorm:"unique;not null;type:varchar(36)"`        // 业务钱包ID
	UserID    string          `json:"user_id" gorm:"not null;type:varchar(36);index"`           // 关联用户ID
	Balance   decimal.Decimal `json:"balance" gorm:"type:decimal(20,8);not null;default:0"`     // 余额
	Status    string          `json:"status" gorm:"not null;default:'active';type:varchar(50)"` // 状态
	Version   int64           `json:"version" gorm:"not null;default:0"`                        // 版本号
	Address   string          `json:"address" gorm:"type:varchar(255);not null;unique"`         // 钱包地址
	CreatedAt int64           `json:"created_at" gorm:"not null"`                               // 创建时间
	UpdatedAt int64           `json:"updated_at" gorm:"not null"`                               // 更新时间
	DeletedAt int64           `json:"deleted_at,omitempty" gorm:"index"`                        // 删除时间
}

// Transaction 交易模型
type Transaction struct {
	ID            int64           `json:"id" gorm:"primaryKey;autoIncrement"`                     // 自增主键
	TransactionID string          `json:"transaction_id" gorm:"unique;not null;type:varchar(36)"` // 业务交易ID
	FromAddress   string          `json:"from_address,omitempty" gorm:"type:varchar(255);index"`  // 转出地址
	ToAddress     string          `json:"to_address,omitempty" gorm:"type:varchar(255);index"`    // 转入地址
	Amount        decimal.Decimal `json:"amount" gorm:"type:decimal(20,8);not null"`              // 交易金额
	Type          string          `json:"type" gorm:"type:varchar(20);not null;index"`            // 交易类型
	Status        string          `json:"status" gorm:"type:varchar(20);not null;index"`          // 交易状态
	Description   string          `json:"description" gorm:"type:text"`                           // 交易描述
	CreatedAt     int64           `json:"created_at" gorm:"not null;index"`                       // 创建时间
	UpdatedAt     int64           `json:"updated_at" gorm:"not null"`                             // 更新时间
	DeletedAt     int64           `json:"deleted_at" gorm:"default:0;index"`                      // 删除时间(0=未删除)
}
