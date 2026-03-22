package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
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
		client: &http.Client{Timeout: 45 * time.Second},
	}
}

func (s *GeminiService) AnalyzeText(ctx context.Context, text string) (*models.Analysis, error) {
	// Updated to the current stable version and model for 2026
	endpoint := "https://generativelanguage.googleapis.com/v1/models/gemini-2.5-flash:generateContent"
	fullURL := fmt.Sprintf("%s?key=%s", endpoint, s.apiKey)

	prompt := fmt.Sprintf(`Analyze this text and return ONLY a valid JSON object.
Fields: "headline" (string), "summary" (string), "sentiment" (float 0.0-1.0), "visual_style" (string: energetic, calm, or dramatic).
Text: %s`, text)

	payload, _ := json.Marshal(map[string]interface{}{
		"contents": []interface{}{
			map[string]interface{}{
				"parts": []interface{}{
					map[string]string{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature": 0.1,
		},
	})

	var respBody []byte
	var lastErr error

	for i := 0; i < 3; i++ {
		req, _ := http.NewRequestWithContext(ctx, "POST", fullURL, bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.client.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(i+1) * 2 * time.Second)
			continue
		}

		respBody, _ = io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			break
		}

		lastErr = fmt.Errorf("google api error (%d): %s", resp.StatusCode, string(respBody))
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			time.Sleep(time.Duration(i+1) * 3 * time.Second)
			continue
		}
		return nil, lastErr
	}

	if respBody == nil {
		return nil, fmt.Errorf("request failed: %w", lastErr)
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

	if err := json.Unmarshal(respBody, &geminiRaw); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	if len(geminiRaw.Candidates) == 0 || len(geminiRaw.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("AI returned empty candidates. Check API key/quota.")
	}

	rawText := geminiRaw.Candidates[0].Content.Parts[0].Text

	re := regexp.MustCompile(`(?s)\{.*\}`)
	jsonMatch := re.FindString(rawText)
	if jsonMatch == "" {
		return nil, fmt.Errorf("no JSON found in AI response: %s", rawText)
	}

	var analysis models.Analysis
	if err := json.Unmarshal([]byte(jsonMatch), &analysis); err != nil {
		return nil, fmt.Errorf("json parse error: %w", err)
	}

	return &analysis, nil
}
