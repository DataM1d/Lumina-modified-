package scraper

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

const MaxChars = 10000

func ScrapeArticle(targetURL string) (string, error) {
	if _, err := url.ParseRequestURI(targetURL); err != nil {
		return "", fmt.Errorf("invalid url format")
	}

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"),
		colly.Async(true),
	)

	_ = c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       1 * time.Second,
	})

	var builder strings.Builder

	c.OnHTML("article, main, .post-content, .article-body, #content", func(e *colly.HTMLElement) {
		e.ForEach("p, h1, h2, h3", func(_ int, el *colly.HTMLElement) {
			if builder.Len() < MaxChars {
				text := strings.TrimSpace(el.Text)
				if len(text) > 20 {
					builder.WriteString(text + " ")
				}
			}
		})
	})

	c.OnHTML("p", func(e *colly.HTMLElement) {
		if builder.Len() < 300 {
			text := strings.TrimSpace(e.Text)
			if len(text) > 20 {
				builder.WriteString(text + " ")
			}
		}
	})

	err := c.Visit(targetURL)
	c.Wait()

	if err != nil {
		return "", fmt.Errorf("scrape failed: %w", err)
	}

	result := strings.Join(strings.Fields(builder.String()), " ")

	if len(result) > MaxChars {
		result = result[:MaxChars]
	}

	if len(result) < 200 {
		return "", fmt.Errorf("insufficient content found (%d chars)", len(result))
	}

	return result, nil
}
