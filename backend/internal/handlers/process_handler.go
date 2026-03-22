package handlers

import (
	"net/http"
	"time"

	"github.com/DataM1d/lumina-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type ProcessHandler struct {
	Svc *service.ArticleService
}

func NewProcessHandler(svc *service.ArticleService) *ProcessHandler {
	return &ProcessHandler{
		Svc: svc,
	}
}

func (h *ProcessHandler) HandleAnalyze(c *gin.Context) {
	var input struct {
		URL string `json:"url" binding:"required,url"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "A valid source URL is required",
		})
		return
	}

	if h.Svc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Analysis service is currently uninitialized",
		})
		return
	}

	analysis, err := h.Svc.ProcessSource(c.Request.Context(), input.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"source":   input.URL,
		"analysis": analysis,
		"meta": gin.H{
			"timestamp": time.Now().Unix(),
			"version":   "v1",
		},
	})
}

func (h *ProcessHandler) HandleHistory(c *gin.Context) {
	history, err := h.Svc.Repo.GetAll(c.Request.Context(), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"history": history,
	})
}
