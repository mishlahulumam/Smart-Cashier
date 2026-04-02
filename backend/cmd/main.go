package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"smart-cashier-backend/internal/config"
	"smart-cashier-backend/internal/database"
	"smart-cashier-backend/internal/handler"
	"smart-cashier-backend/internal/middleware"
	"smart-cashier-backend/internal/repository"
	"smart-cashier-backend/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func killProcessOnPort(port string) {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:"+port, time.Second)
	if err != nil {
		return
	}
	conn.Close()

	log.Printf("Port %s is already in use, attempting to free it...", port)

	if runtime.GOOS == "windows" {
		out, err := exec.Command("cmd", "/c", fmt.Sprintf("netstat -ano | findstr :%s | findstr LISTENING", port)).Output()
		if err != nil {
			return
		}
		for _, line := range strings.Split(string(out), "\n") {
			fields := strings.Fields(strings.TrimSpace(line))
			if len(fields) >= 5 {
				pid := fields[len(fields)-1]
				exec.Command("taskkill", "/F", "/PID", pid).Run()
				log.Printf("Killed old process PID %s on port %s", pid, port)
			}
		}
	} else {
		out, _ := exec.Command("lsof", "-ti", ":"+port).Output()
		for _, pid := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			if pid != "" {
				exec.Command("kill", "-9", pid).Run()
				log.Printf("Killed old process PID %s on port %s", pid, port)
			}
		}
	}

	time.Sleep(500 * time.Millisecond)
}

func main() {
	cfg := config.Load()
	db := database.Connect(cfg.DBPath)

	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	txRepo := repository.NewTransactionRepository(db)

	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	categorySvc := service.NewCategoryService(categoryRepo)
	productSvc := service.NewProductService(productRepo)
	customerSvc := service.NewCustomerService(customerRepo)
	txSvc := service.NewTransactionService(txRepo, productRepo, customerRepo)
	reportSvc := service.NewReportService(txRepo, productRepo)

	authHandler := handler.NewAuthHandler(authSvc)
	categoryHandler := handler.NewCategoryHandler(categorySvc)
	productHandler := handler.NewProductHandler(productSvc)
	customerHandler := handler.NewCustomerHandler(customerSvc)
	txHandler := handler.NewTransactionHandler(txSvc)
	reportHandler := handler.NewReportHandler(reportSvc)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/uploads", "./uploads")

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Smart Cashier API is running"})
	})


	auth := r.Group("/api/auth")
	{
		auth.POST("/login", authHandler.Login)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		api.GET("/auth/me", authHandler.Me)
		api.POST("/auth/register", middleware.AdminOnly(), authHandler.Register)

		categories := api.Group("/categories")
		{
			categories.GET("", categoryHandler.GetAll)
			categories.GET("/:id", categoryHandler.GetByID)
			categories.POST("", middleware.AdminOnly(), categoryHandler.Create)
			categories.PUT("/:id", middleware.AdminOnly(), categoryHandler.Update)
			categories.DELETE("/:id", middleware.AdminOnly(), categoryHandler.Delete)
		}

		products := api.Group("/products")
		{
			products.GET("", productHandler.GetAll)
			products.GET("/:id", productHandler.GetByID)
			products.POST("", middleware.AdminOnly(), productHandler.Create)
			products.PUT("/:id", middleware.AdminOnly(), productHandler.Update)
			products.DELETE("/:id", middleware.AdminOnly(), productHandler.Delete)
		}

		customers := api.Group("/customers")
		{
			customers.GET("", customerHandler.GetAll)
			customers.GET("/:id", customerHandler.GetByID)
			customers.POST("", customerHandler.Create)
			customers.PUT("/:id", customerHandler.Update)
			customers.DELETE("/:id", middleware.AdminOnly(), customerHandler.Delete)
		}

		transactions := api.Group("/transactions")
		{
			transactions.GET("", txHandler.GetAll)
			transactions.GET("/:id", txHandler.GetByID)
			transactions.POST("", txHandler.Create)
		}

		api.GET("/dashboard", reportHandler.Dashboard)

		reports := api.Group("/reports")
		{
			reports.GET("/daily", reportHandler.DailySummary)
			reports.GET("/monthly", reportHandler.MonthlySummary)
			reports.GET("/top-products", reportHandler.TopProducts)
			reports.GET("/summary", reportHandler.Summary)
		}
	}

	killProcessOnPort(cfg.Port)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Smart Cashier API starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited cleanly")
}
