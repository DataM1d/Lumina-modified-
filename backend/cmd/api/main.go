package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DataM1d/lumina-backend/internal/ai"
	"github.com/DataM1d/lumina-backend/internal/handlers"
	"github.com/DataM1d/lumina-backend/internal/middleware"
	"github.com/DataM1d/lumina-backend/internal/repository"
	"github.com/DataM1d/lumina-backend/internal/service"
	"github.com/DataM1d/lumina-backend/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("FATAL: GEMINI_API_KEY is missing from environment")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("FATAL: DATABASE_URL is missing from environment")
	}

	db, err := database.Connect(dbURL)
	if err != nil {
		log.Fatalf("FATAL: Database connection failed: %v", err)
	}
	defer db.Close()

	repo := repository.NewAnalysisRepository(db)
	aiSvc := ai.NewGeminiService(apiKey)
	artSvc := service.NewArticleService(aiSvc, repo)
	processHdl := handlers.NewProcessHandler(artSvc)

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), CORSMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "online", "time": time.Now().Format(time.RFC3339)})
	})

	api := r.Group("/api/v1")
	api.Use(middleware.RateLimit(5))
	{
		api.POST("/process", processHdl.HandleAnalyze)
		api.GET("/history", processHdl.HandleHistory)
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	fmt.Printf("LUMINA.AI Backend: http://localhost:%s\n", port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	fmt.Println("Server exited")
}
