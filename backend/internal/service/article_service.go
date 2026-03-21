package service

import (
	"context"
	"fmt"

	"github.com/DataM1d/lumina-backend/internal/ai"
	"github.com/DataM1d/lumina-backend/internal/models"
	"github.com/DataM1d/lumina-backend/internal/scraper"
)

type ArticleService struct {
	AISvc *ai.GeminiService
}

func NewArticleService(aiSvc *ai.GeminiService) *ArticleService {
	return &ArticleService{AISvc: aiSvc}
}

func (s *ArticleService) ProcessSource(ctx context.Context, url string) (*models.Analysis, error) {
	text, err := scraper.ScrapeArticle(url)
	if err != nil {
		return nil, fmt.Errorf("scraping error: %w", err)
	}

	analysis, err := s.AISvc.AnalyzeText(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("ai analysis error: %w", err)
	}

	return analysis, nil
}
