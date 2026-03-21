package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/DataM1d/lumina-backend/internal/models"
)

type GeminiService struct {
	apiKey string
	client *http.Client
}

func NewGeminiService(apiKey string) *GeminiService {
	return &GeminiService{
		apiKey: apiKey,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *GeminiService) AnalyzeText(ctx context.Context, text string) (*models.Analysis, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", s.apiKey)

	prompt := fmt.Sprintf(`Analyze the following text and return a JSON object. 
	Fields: 'headline' (max 10 words), 'summary' (1 sentence), 'sentiment' (0.0 to 1.0), 'visual_style' (energetic, calm, or dramatic).
	Content: %s`, text)

	requestBody, _ := json.Marshal(map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{{"text": prompt}},
			},
		},
		"generationConfig": map[string]interface{}{
			"responseMimeType": "application/json",
		},
	})

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gemini api call failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google api error (%d): %s", resp.StatusCode, string(body))
	}

	var geminiRaw struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(body, &geminiRaw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal gemini response: %w", err)
	}

	if len(geminiRaw.Candidates) == 0 || len(geminiRaw.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	var analysis models.Analysis
	err = json.Unmarshal([]byte(geminiRaw.Candidates[0].Content.Parts[0].Text), &analysis)
	if err != nil {
		return nil, fmt.Errorf("invalid analysis format: %w", err)
	}

	return &analysis, nil
}
