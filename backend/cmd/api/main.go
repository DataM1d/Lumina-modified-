package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DataM1d/lumina-backend/internal/ai"
	"github.com/DataM1d/lumina-backend/internal/scraper"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not found, using system environment")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Lumina Backend is Running"})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scrape article"})
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

	fmt.Printf("Server starting on port %s...\n", port)
	r.Run(":" + port)
}
