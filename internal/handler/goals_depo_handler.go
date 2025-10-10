package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	dto "github.com/aborgas90/expense-tracker-api/internal/dto/goals_depo"
	"github.com/aborgas90/expense-tracker-api/internal/dto/response"
	"github.com/aborgas90/expense-tracker-api/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GoalsDepoHandler struct {
	svc *service.GoalsDepoService
}

func NewGoalsDepoHandler(h *service.GoalsDepoService) *GoalsDepoHandler {
	return &GoalsDepoHandler{svc: h}
}

func (h *GoalsDepoHandler) CreateDepoHandler(c *gin.Context) {
	var req dto.RequestGoalsDepo
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Printf("DEBUG RequestGoalsDepo: %+v\n", req) // cek di console

	res, err := h.svc.CreateDepoServ(&req)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	dtoRes := dto.ResponseGoalsDepo{
		ID:        res.ID,
		GoalID:    res.GoalID,
		Amount:    res.Amount,
		Note:      res.Note,
		CreatedAt: res.CreatedAt,
	}

	response.SuccessResponse(c, http.StatusCreated, "Deposit created successfully", dtoRes)
}

func (h *GoalsDepoHandler) GetDepoByID(c *gin.Context) {
	// ambil user_id dari middleware (contoh: disimpan di context)
	userID := c.GetUint("user_id")
	depoIDParam := c.Param("id")

	depoID, err := strconv.ParseUint(depoIDParam, 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid deposit ID")
		return
	}

	deposit, err := h.svc.GetDepoByID(userID, uint(depoID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrorResponse(c, http.StatusNotFound, "Deposit not found or unauthorized")
			return
		}
		response.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve deposit")
		return
	}

	dtoRes := dto.ResponseGoalsDepo{
		ID:        deposit.ID,
		GoalID:    deposit.GoalID,
		Amount:    deposit.Amount,
		Note:      deposit.Note,
		CreatedAt: deposit.CreatedAt,
	}

	response.SuccessResponse(c, http.StatusOK, "Successfully retrieved deposit", dtoRes)
}

func (h *GoalsDepoHandler) UpdateGoalsDepoHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	depoIDParam := c.Param("id")

	var req *dto.RequestGoalsDepo

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	depoID, err := strconv.ParseUint(depoIDParam, 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid deposit ID")
		return
	}

	res, err := h.svc.UpdateGoalsDepo(uint(depoID), req.GoalID, userID, uint(req.Amount), req.Note)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, "Cannot Update Data Goals")
		return
	}

	response.SuccessResponse(c, http.StatusOK, "Successfully retrieved deposit", res)
}

func (h *GoalsDepoHandler) DeleteGoalsDepoHandler(c *gin.Context) {
	userID := c.GetUint("user_id")
	depoIDParam := c.Param("id")

	depoID, err := strconv.ParseUint(depoIDParam, 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Invalid deposit ID")
		return
	}

	res, err := h.svc.DeleteGoalsDepo(uint(depoID), userID)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, "Cannot Delete this Goals Depo")
		return
	}

	response.SuccessResponse(c, http.StatusNoContent, "Successfully to delete data", res)
}
