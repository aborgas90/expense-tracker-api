package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aborgas90/expense-tracker-api/internal/handler"
	middleware "github.com/aborgas90/expense-tracker-api/internal/middleware"
	"github.com/aborgas90/expense-tracker-api/internal/platform/db"
	"github.com/aborgas90/expense-tracker-api/internal/repo"
	"github.com/aborgas90/expense-tracker-api/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	conn := db.InIt()
	userRepo := repo.NewUserRepo(conn)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	categoryRepo := repo.NewCategoryRepo(conn)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	transactionRepo := repo.NewTransactionRepo(conn)
	transactionService := service.NewTransactionService(transactionRepo)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	router := gin.Default()

	sqlDB, err := conn.DB()
	if err != nil {
		log.Fatalf("❌ gagal ambil sql.DB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("❌ DB tidak bisa di-ping: %v", err)
	}
	fmt.Println("✅ Database connected!")

	router.GET("/check", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	{
		v1 := router.Group("/api/v1")
		v1.POST("/auth/register", userHandler.RegisterHandler)
		v1.POST("/auth/login", userHandler.LoginHandler)
		v1.POST("/auth/refresh", userHandler.Refresh)
		//categories
		v1.GET("/categories", middleware.AuthMiddleware(), categoryHandler.GetCategoriesByUserID)
		v1.POST("/categories", middleware.AuthMiddleware(), categoryHandler.CreateCategory)
		v1.PUT("/categories/:id", middleware.AuthMiddleware(), categoryHandler.UpdateCategory)
		v1.DELETE("/categories/:id", middleware.AuthMiddleware(), categoryHandler.DeleteCategory)
		//transaction
		v1.POST("/transaction", middleware.AuthMiddleware(), transactionHandler.CreateTransactionUser)
		v1.GET("/transaction", middleware.AuthMiddleware(), transactionHandler.GetTransactionByUser)
	}

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	router.Run(":8080")
}
