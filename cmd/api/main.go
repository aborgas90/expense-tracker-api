package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aborgas90/expense-tracker-api/internal/handler"
	"github.com/aborgas90/expense-tracker-api/internal/helper"
	middleware "github.com/aborgas90/expense-tracker-api/internal/middleware"
	"github.com/aborgas90/expense-tracker-api/internal/platform/db"
	"github.com/aborgas90/expense-tracker-api/internal/repo"
	"github.com/aborgas90/expense-tracker-api/internal/service"
	"github.com/gin-contrib/cors"
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

	goalsRepo := repo.NewGoalRepo(conn)
	goalsService := service.NewGoalsService(goalsRepo)
	goalsHandler := handler.NewGoalsHandler(goalsService)

	router := gin.Default()

	//setup cors
	config := cors.Config{
		AllowOrigins:     []string{os.Getenv("FE_APP")}, // tanpa "/"
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(helper.RateLimiter())
	router.Use(cors.New(config))

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
		v1.GET("/transaction/:id", middleware.AuthMiddleware(), transactionHandler.GetTransactionById)
		v1.PUT("/transaction/:id", middleware.AuthMiddleware(), transactionHandler.UpdateTransaction)
		v1.DELETE("/transaction/:id", middleware.AuthMiddleware(), transactionHandler.DeleteTransaction)

		//dashboard
		v1.GET("/dashboard", middleware.AuthMiddleware(), transactionHandler.SummaryTransaction)
		v1.GET("/dashboard/check-surplus-defisit", middleware.AuthMiddleware(), transactionHandler.CheckSurplusDeficitTransaction)
		v1.GET("/dashboard/last-transaction", middleware.AuthMiddleware(), transactionHandler.Last7Transaction)

		//goals
		v1.GET("/dashboard/goals/", middleware.AuthMiddleware(), goalsHandler.GetGoalDataByIdUser)
		v1.POST("/dashboard/goals/", middleware.AuthMiddleware(), goalsHandler.CreateGoalsHandler)
		v1.PUT("/dashboard/goals/:id", middleware.AuthMiddleware(), goalsHandler.UpdateGoalsHandler)
		v1.DELETE("/dashboard/goals/:id", middleware.AuthMiddleware(), goalsHandler.DeleteGoalsHandler)
	}

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	router.Run(":8080")
}
