package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/aborgas90/expense-tracker-api/internal/dto/response"
	"github.com/aborgas90/expense-tracker-api/internal/dto/transaction"
	"github.com/aborgas90/expense-tracker-api/internal/service"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	svc *service.TransactionService
}

func NewTransactionHandler(s *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: s}
}

func (h *TransactionHandler) CreateTransactionUser(c *gin.Context) {
	var req transaction.RequestTransaction
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTransaction, err := h.svc.CreateTransactionUser(userId.(uint), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.SuccessResponse(c, "Successfully to create transaction", newTransaction)
}

func (h *TransactionHandler) GetTransactionByUser(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	getTransaction, err := h.svc.GetTransactionByUser(userId.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.SuccessResponse(c, "Successfully to get data transaction", getTransaction)
}

func (h *TransactionHandler) GetTransactionById(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	txsId := c.Param("id")
	txIdConv := StringToUint(txsId)

	result, err := h.svc.GetTransactionById(userId.(uint), txIdConv)
	if err != nil {
		response.ErrorResponse(c, 404, "Cannot get data transaction id")
		return
	}

	response.SuccessResponse(c, "Successful to get data transaction id", result)
}

func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	var req transaction.UpdateTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id"})
		return
	}

	result, err := h.svc.UpdateTransactionUser(
		userId.(uint),
		uint(id),
		req.CategoryId,
		req.Amount,
		req.Currency,
		req.OccuredAt,
		req.Note,
	)

	fmt.Println("LOGGING di HANDLER:: ", result)

	if err != nil {
		fmt.Println("UpdateTransaction error:", err) // debug
		response.ErrorResponse(c, 500, fmt.Sprintf("Cannot Update Transaction: %v", err))
		return
	}

	response.SuccessResponse(c, "Successfully to Update Transaction", result)
}

func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id"})
		return
	}

	// panggil service untuk delete
	rows, err := h.svc.DeleteTransaction(uint(userID.(uint)), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// selalu return empty array di data
	var data []interface{}
	if rows > 0 {
		data = []interface{}{} // tetap kosong meskipun ada row yg dihapus
	}

	response.SuccessResponse(c, "Successfully to delete transaction id", data)
}

func (h *TransactionHandler) SummaryTransaction(c *gin.Context) {
	userID, exists := c.Get("user_id")
	monthStr := c.Query("month")
	yearStr := c.Query("year")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month"})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
		return
	}

	rows, err := h.svc.SummaryTransaction(uint(userID.(uint)), month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.SuccessResponse(c, "Successfully to delete transaction id", rows)
}

func (h *TransactionHandler) CheckSurplusDeficitTransaction(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	rows, err := h.svc.CheckSurplusDeficitTransaction(uint(userID.(uint)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.SuccessResponse(c, "Successfully to data summary surplus & deficit", rows)
}

func (h *TransactionHandler) Last7Transaction(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	rows, err := h.svc.Last7Transaction(uint(userID.(uint)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.SuccessResponse(c, "Successfully to data last 7 transaction", rows)
}

func StringToUint(s string) uint {
	i, _ := strconv.Atoi(s)
	return uint(i)
}
