package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/DataM1d/lumina-backend/internal/ai"
	"github.com/DataM1d/lumina-backend/internal/scraper"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

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

	r := gin.Default()
	r.Use(CORSMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "online",
			"message": "Lumina Core Engine Active",
		})
	})

	r.POST("/process", func(c *gin.Context) {
		var input struct {
			URL string `json:"url" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Valid URL is required"})
			return
		}

		text, err := scraper.ScrapeArticle(input.URL)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Scraping failed: " + err.Error()})
			return
		}

		ctx := context.Background()
		analysis, err := ai.AnalyzeText(ctx, text)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "AI Analysis failed: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"url":      input.URL,
			"analysis": analysis,
		})
	})

	fmt.Printf("--- Lumina Backend Initialized on Port %s ---\n", port)
	if err := r.Run(":" + port); err != nil {
		fmt.Printf("Critical Server Error: %v\n", err)
	}
}
