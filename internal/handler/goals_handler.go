package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	dto "github.com/aborgas90/expense-tracker-api/internal/dto/goals"
	"github.com/aborgas90/expense-tracker-api/internal/dto/response"
	"github.com/aborgas90/expense-tracker-api/internal/service"
	"github.com/gin-gonic/gin"
)

type GoalsHandler struct {
	svc *service.GoalsServices
}

func NewGoalsHandler(s *service.GoalsServices) *GoalsHandler {
	return &GoalsHandler{svc: s}
}

func (h *GoalsHandler) CreateGoalsHandler(c *gin.Context) {
	var req dto.RequestGoals

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Received request body: %+v\n", req)

	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	newGoals, err := h.svc.CreateGoalsTarget(userId.(uint), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.SuccessResponse(c, "Successfully Create Goals Target", newGoals)
}

func (h *GoalsHandler) GetGoalDataByIdUser(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	result, err := h.svc.GetGoalDataByIdUser(userId.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.SuccessResponse(c, "Successfull to get goals data", result)
}

func (h *GoalsHandler) UpdateGoalsHandler(c *gin.Context) {
	// Get user_id from context
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse goal ID from URL
	paramId := c.Param("id")
	id64, err := strconv.ParseUint(paramId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid goal ID"})
		return
	}
	goalID := uint(id64)

	// Bind request JSON into dto.RequestGoals
	var req dto.RequestGoals
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Call service to update data
	updatedGoal, err := h.svc.UpdateGoalsDataById(userId.(uint), goalID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Build custom DTO response
	responseData := dto.ResponseGoals{
		Id:             updatedGoal.ID,
		Title:          updatedGoal.Title,
		Target_amount:  updatedGoal.TargetAmount,
		Current_amount: updatedGoal.CurrentAmount,
		Deadline:       updatedGoal.Deadline.Format(time.RFC3339),
		Created_at:     updatedGoal.CreatedAt.Format(time.RFC3339),
	}

	// Send JSON response
	response.SuccessResponse(c, "Success to update data goals", responseData)
}

func (h *GoalsHandler) DeleteGoalsHandler(c *gin.Context) {
	// Get user_id from context
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse goal ID from URL
	paramId := c.Param("id")
	id64, err := strconv.ParseUint(paramId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid goal ID"})
		return
	}
	goalID := uint(id64)

	res, err := h.svc.DeleteGoalsDataById(userId.(uint), goalID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	response.SuccessResponse(c, "Success to Delete Goals Data", nil)
}
