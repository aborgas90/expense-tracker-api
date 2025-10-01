package handler

import (
	"net/http"

	token "github.com/aborgas90/expense-tracker-api/internal/auth"
	"github.com/aborgas90/expense-tracker-api/internal/dto/auth"
	"github.com/aborgas90/expense-tracker-api/internal/dto/response"
	"github.com/aborgas90/expense-tracker-api/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *service.UserService
}

func NewUserHandler(s *service.UserService) *AuthHandler {
	return &AuthHandler{svc: s}
}

func (h *AuthHandler) RegisterHandler(c *gin.Context) {
	var req auth.RegisterUserRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.svc.RegisterUser(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response.SuccessResponse(c, "User registered successfully", res)
}

func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req auth.LoginUserRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.svc.LoginUser(&req)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.SetCookie("access_token", res.AccessToken, 3600, "/", "", false, true)

	response.SuccessResponse(c, "Login successful", res)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	claims, err := token.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	accessToken, refreshToken, err := token.GenerateToken(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "localhost", true, true)

	response.SuccessResponse(c, "Token refreshed successfully", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
