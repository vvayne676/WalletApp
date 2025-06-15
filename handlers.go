package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// 钱包相关处理器

func deposit(c *gin.Context) {
	var req struct {
		Amount  decimal.Decimal `json:"amount" binding:"required"`
		UserID  string          `json:"user_id" binding:"required"`
		Address string          `json:"address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Amount.LessThanOrEqual(decimal.Zero) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be positive"})
		return
	}
	var wallet Wallet
	if err := DB.Where("address = ? AND deleted_at = 0", req.Address).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	if wallet.UserID != req.UserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet not found"})
		return
	}

	tx := DB.Begin()

	now := time.Now().Unix()
	transaction := Transaction{
		TransactionID: uuid.New().String(),
		ToAddress:     req.Address,
		Amount:        req.Amount,
		Type:          "deposit",
		Status:        "pending",
		Description:   "Deposit to wallet",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	if err := updateWalletBalance(tx, "", req.Address, req.Amount); err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Model(&transaction).Update("status", "completed").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction status"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":     "Deposit successful",
		"transaction": transaction,
	})
}

func withdraw(c *gin.Context) {
	var req struct {
		Amount  decimal.Decimal `json:"amount" binding:"required"`
		UserID  string          `json:"user_id" binding:"required"`
		Address string          `json:"address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Amount.LessThanOrEqual(decimal.Zero) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be positive"})
		return
	}

	// 验证钱包属于该用户
	var wallet Wallet
	if err := DB.Where("address = ? AND user_id = ? AND deleted_at = 0", req.Address, req.UserID).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet not found or not owned by user"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	tx := DB.Begin()

	now := time.Now().Unix()
	transaction := Transaction{
		TransactionID: uuid.New().String(),
		FromAddress:   req.Address,
		Amount:        req.Amount,
		Type:          "withdraw",
		Status:        "pending",
		Description:   "Withdraw from wallet",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	if err := updateWalletBalance(tx, req.Address, "", req.Amount.Neg()); err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Model(&transaction).Update("status", "completed").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction status"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":     "Withdraw successful",
		"transaction": transaction,
	})
}

func transfer(c *gin.Context) {
	var req struct {
		UserID      string          `json:"user_id" binding:"required"`
		FromAddress string          `json:"from_address" binding:"required"`
		ToAddress   string          `json:"to_address" binding:"required"`
		Amount      decimal.Decimal `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Amount.LessThanOrEqual(decimal.Zero) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be positive"})
		return
	}

	// 验证源钱包属于该用户
	var wallet Wallet
	if err := DB.Where("address = ? AND user_id = ? AND deleted_at = 0", req.FromAddress, req.UserID).First(&wallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet not found or not owned by user"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	if req.FromAddress == req.ToAddress {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot transfer to the same address"})
		return
	}

	// 验证目标钱包存在
	var toWallet Wallet
	if err := DB.Where("address = ? AND deleted_at = 0", req.ToAddress).First(&toWallet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Destination wallet not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	tx := DB.Begin()

	now := time.Now().Unix()
	transaction := Transaction{
		TransactionID: uuid.New().String(),
		FromAddress:   req.FromAddress,
		ToAddress:     req.ToAddress,
		Amount:        req.Amount,
		Type:          "transfer",
		Status:        "pending",
		Description:   "Transfer to user",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	// 从源地址扣款
	if err := updateWalletBalance(tx, req.FromAddress, "", req.Amount.Neg()); err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 向目标地址加款
	if err := updateWalletBalance(tx, "", req.ToAddress, req.Amount); err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Model(&transaction).Update("status", "completed").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction status"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":     "Transfer successful",
		"transaction": transaction,
	})
}

func getBalance(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var wallets []Wallet
	if err := DB.Select("balance").Where("user_id = ? AND deleted_at = 0", req.UserID).Find(&wallets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get wallets"})
		return
	}

	// 计算总余额
	totalBalance := decimal.Zero
	// 汇率转换成美元
	for _, wallet := range wallets {
		totalBalance = totalBalance.Add(wallet.Balance)
	}

	c.JSON(http.StatusOK, gin.H{
		"balance": gin.H{
			"user_id": req.UserID,
			"balance": totalBalance,
		},
	})
}

func getHistory(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id" binding:"required"`
		Page   int    `json:"page"`
		Size   int    `json:"page_size"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 || req.Size > 100 {
		req.Size = 20
	}

	// 首先获取用户的所有钱包地址
	var wallets []Wallet
	if err := DB.Select("address").Where("user_id = ? AND deleted_at = 0", req.UserID).Find(&wallets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user wallets"})
		return
	}

	if len(wallets) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"transactions": []Transaction{},
				"page":         req.Page,
				"page_size":    req.Size,
				"total":        0,
			},
		})
		return
	}

	// 提取地址列表
	addresses := make([]string, len(wallets))
	for i, wallet := range wallets {
		addresses[i] = wallet.Address
	}

	var transactions []Transaction
	var total int64

	// 构建查询条件
	query := DB.Where("deleted_at = 0 AND (from_address IN ? OR to_address IN ?)", addresses, addresses)

	// 获取总数
	if err := query.Model(&Transaction{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count transactions"})
		return
	}

	// 分页查询
	offset := (req.Page - 1) * req.Size
	if err := query.Order("created_at DESC").
		Limit(req.Size).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"transactions": transactions,
			"page":         req.Page,
			"page_size":    req.Size,
			"total":        total,
		},
	})
}

func updateWalletBalance(tx *gorm.DB, fromAddress string, toAddress string, amount decimal.Decimal) error {
	if fromAddress != "" {
		var wallet Wallet
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("address = ? AND deleted_at = 0", fromAddress).
			First(&wallet).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("Source wallet not found")
			}
			return errors.New("Failed to lock source wallet")
		}

		newBalance := wallet.Balance.Add(amount) // amount 已经是负数了
		if newBalance.LessThan(decimal.Zero) {
			return errors.New("Insufficient balance")
		}

		if err := tx.Model(&wallet).Update("balance", newBalance).Error; err != nil {
			return errors.New("Failed to update source wallet balance")
		}
	}

	if toAddress != "" {
		var toWallet Wallet
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			Where("address = ? AND deleted_at = 0", toAddress).
			First(&toWallet).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("Destination wallet not found")
			}
			return errors.New("Failed to lock destination wallet")
		}

		newBalance := toWallet.Balance.Add(amount)
		if err := tx.Model(&toWallet).Update("balance", newBalance).Error; err != nil {
			return errors.New("Failed to update destination wallet balance")
		}
	}

	return nil
}
