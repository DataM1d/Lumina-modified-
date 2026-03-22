package service

import (
	"context"
	"fmt"
	"time"

	"github.com/DataM1d/lumina-backend/internal/ai"
	"github.com/DataM1d/lumina-backend/internal/models"
	"github.com/DataM1d/lumina-backend/internal/repository"
	"github.com/DataM1d/lumina-backend/internal/scraper"
)

type ArticleService struct {
	AISvc *ai.GeminiService
	Repo  *repository.AnalysisRepository
}

func NewArticleService(aiSvc *ai.GeminiService, repo *repository.AnalysisRepository) *ArticleService {
	return &ArticleService{
		AISvc: aiSvc,
		Repo:  repo,
	}
}

func (s *ArticleService) ProcessSource(ctx context.Context, url string) (*models.Analysis, error) {
	existing, err := s.Repo.GetByURL(ctx, url)
	if err == nil && existing != nil {
		if time.Since(existing.ProcessedAt) < 24*time.Hour {
			return existing, nil
		}
	}

	text, err := scraper.ScrapeArticle(url)
	if err != nil {
		return nil, fmt.Errorf("scraping error: %w", err)
	}

	analysis, err := s.AISvc.AnalyzeText(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("ai analysis error: %w", err)
	}

	analysis.URL = url
	analysis.ProcessedAt = time.Now()

	if err := s.Repo.Create(ctx, analysis); err != nil {
		fmt.Printf("Warning: failed to save to db: %v\n", err)
	}

	return analysis, nil
}
