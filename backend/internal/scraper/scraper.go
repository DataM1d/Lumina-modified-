package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

func ScrapeArticle(url string) (string, error) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	var fullText strings.Builder
	selectors := []string{"article p", "main p", ".article-body p", "p"}

	for _, selector := range selectors {
		c.OnHTML(selector, func(e *colly.HTMLElement) {
			text := strings.TrimSpace(e.Text)
			if len(text) > 40 {
				fullText.WriteString(text + " ")
			}
		})
	}

	err := c.Visit(url)
	if err != nil {
		return "", err
	}

	result := strings.TrimSpace(fullText.String())
	if len(result) < 200 {
		return "", fmt.Errorf("content too thin")
	}

	return result, nil
}
