package handler

import (
	"fmt"
	"net/http"

	"github.com/aborgas90/expense-tracker-api/internal/dto/category"
	"github.com/aborgas90/expense-tracker-api/internal/dto/response"
	"github.com/aborgas90/expense-tracker-api/internal/service"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	svc *service.CategoryService
}

func NewCategoryHandler(s *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: s}
}

func (h *CategoryHandler) GetCategoriesByUserID(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
		return
	}

	categories, err := h.svc.GetCategoriesByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.SuccessResponse(c, "Categories fetched successfully",
		categories,
	)
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req category.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ambil user_id dari JWT middleware
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	category, err := h.svc.CreateCategory(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.SuccessResponse(c, "Category created successfully", category)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	var req category.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	categoryID := c.Param("id")

	updateCategory, err := h.svc.UpdateCategory(categoryID, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Updated Category:", req)

	response.SuccessResponse(c, "Category updated successfully", updateCategory)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	userVal, _ := c.Get("user_id")
	userID := userVal.(uint)

	categoryID := c.Param("id")
	if err := h.svc.DeleteCategory(categoryID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.SuccessResponse(c, "Category deleted successfully", nil)
}
