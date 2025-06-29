package server

import (
	"itk-assignment/wallet"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WalletOperationRequest struct {
	ID            string `json:"walletId" binding:"required"`
	OperationType string `json:"operationType" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
}

func newOperationHandler(wService wallet.WalletService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req WalletOperationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var op wallet.Operation
		switch req.OperationType {
		case "DEPOSIT":
			op = wallet.NewOperation(wallet.OperationDeposit, wallet.Amount(req.Amount))
		case "WITHDRAW":
			op = wallet.NewOperation(wallet.OperationWithdraw, wallet.Amount(req.Amount))
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported operation"})
			return
		}

		id, err := uuid.Parse(req.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = wService.Process(wallet.WalletId(id), op)
		if err != nil {
			if err == wallet.ErrInsufficientFunds {
				c.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func newBalanceHandler(wService wallet.WalletService) gin.HandlerFunc {
	return func(c *gin.Context) {
		walletID := c.Param("walletId")
		id, err := uuid.Parse(walletID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet ID format"})
			return
		}

		walletRecord, err := wService.Get(wallet.WalletId(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"walletId": walletID,
			"balance":  walletRecord.Balance().String(),
		})
	}
}

func newCreateHandler(wService wallet.WalletService) gin.HandlerFunc {
	return func(c *gin.Context) {
		walletRecord, err := wService.CreateWallet()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "service unavaliable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"walletId": walletRecord.Id().UUID(),
			"balance":  walletRecord.Balance().String(),
		})
	}
}
