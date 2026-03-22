package repository

import (
	"context"
	"database/sql"

	"github.com/DataM1d/lumina-backend/internal/models"
)

type AnalysisRepository struct {
	db *sql.DB
}

func NewAnalysisRepository(db *sql.DB) *AnalysisRepository {
	return &AnalysisRepository{db: db}
}

func (r *AnalysisRepository) GetByURL(ctx context.Context, url string) (*models.Analysis, error) {
	query := `SELECT id, url, headline, summary, sentiment, visual_style, processed_at 
			  FROM analyses WHERE url = $1 LIMIT 1`

	var a models.Analysis
	err := r.db.QueryRowContext(ctx, query, url).Scan(
		&a.ID, &a.URL, &a.Headline, &a.Summary, &a.Sentiment, &a.VisualStyle, &a.ProcessedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &a, err
}

func (r *AnalysisRepository) Create(ctx context.Context, a *models.Analysis) error {
	query := `INSERT INTO analyses (url, headline, summary, sentiment, visual_style, processed_at)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	return r.db.QueryRowContext(ctx, query, a.URL, a.Headline, a.Summary, a.Sentiment, a.VisualStyle, a.ProcessedAt).Scan(&a.ID)
}

func (r *AnalysisRepository) GetAll(ctx context.Context, limit int) ([]models.Analysis, error) {
	query := `SELECT id, url, headline, sentiment, visual_style, processed_at 
			  FROM analyses ORDER BY processed_at DESC LIMIT $1`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.Analysis
	for rows.Next() {
		var a models.Analysis
		if err := rows.Scan(&a.ID, &a.URL, &a.Headline, &a.Sentiment, &a.VisualStyle, &a.ProcessedAt); err != nil {
			return nil, err
		}
		results = append(results, a)
	}
	return results, nil
}
