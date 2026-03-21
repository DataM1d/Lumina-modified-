package scraper

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func ScrapeArticle(targetURL string) (string, error) {
	if _, err := url.ParseRequestURI(targetURL); err != nil {
		return "", fmt.Errorf("invalid url format")
	}

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       1 * time.Second,
		RandomDelay: 500 * time.Millisecond,
	})

	var fullText strings.Builder
	seen := make(map[string]bool)

	selectors := []string{
		"article p",
		"section p",
		".post-content p",
		".article-body p",
		"main p",
	}

	for _, selector := range selectors {
		c.OnHTML(selector, func(e *colly.HTMLElement) {
			text := strings.TrimSpace(e.Text)
			if len(text) > 50 && !seen[text] {
				seen[text] = true
				fullText.WriteString(text + " ")
			}
		})
	}

	c.OnScraped(func(r *colly.Response) {
		if fullText.Len() < 100 {
			c.OnHTML("p", func(e *colly.HTMLElement) {
				text := strings.TrimSpace(e.Text)
				if len(text) > 50 && !seen[text] {
					seen[text] = true
					fullText.WriteString(text + " ")
				}
			})
		}
	})

	err := c.Visit(targetURL)
	if err != nil {
		return "", fmt.Errorf("scraping failed: %w", err)
	}

	c.Wait()

	result := strings.TrimSpace(fullText.String())

	if len(result) < 300 {
		return "", fmt.Errorf("insufficient content found to analyze")
	}

	return result, nil
}
