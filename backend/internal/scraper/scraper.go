package scraper

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

func ScrapeArticle(url string) (string, error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	var fullText strings.Builder

	selectors := []string{"article p", "main p", ".article-body p", "p"}

	for _, selector := range selectors {
		c.OnHTML(selector, func(e *colly.HTMLElement) {
			text := strings.TrimSpace(e.Text)
			if len(text) > 20 {
				fullText.WriteString(text + " ")
			}
		})
	}

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Scraping error on %s: %v (Status %d)", r.Request.URL, err, r.StatusCode)
	})

	err := c.Visit(url)
	if err != nil {
		return "", err
	}

	result := strings.TrimSpace(fullText.String())

	if len(result) < 50 { // If we got almost nothing, it's a fail
		return "", fmt.Errorf("article content too short or blocked")
	}

	return result, nil
}
