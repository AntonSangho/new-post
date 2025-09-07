package scraper

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Post struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Date  string `json:"date"`
}

type Scraper struct {
	targetURL string
	client    *http.Client
}

func New(targetURL string) *Scraper {
	return &Scraper{
		targetURL: targetURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *Scraper) GetLatestPosts() ([]Post, error) {
	req, err := http.NewRequest("GET", s.targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var posts []Post

	doc.Find("tr, .post, .article, .item").Each(func(i int, s *goquery.Selection) {
		title := s.Find("a").First().Text()
		link, exists := s.Find("a").First().Attr("href")
		
		if title != "" && exists {
			title = strings.TrimSpace(title)
			if !strings.HasPrefix(link, "http") {
				if strings.HasPrefix(link, "/") {
					link = getBaseURL(s.targetURL) + link
				} else {
					link = getBaseURL(s.targetURL) + "/" + link
				}
			}

			post := Post{
				ID:    generateID(title, link),
				Title: title,
				URL:   link,
				Date:  time.Now().Format("2006-01-02 15:04:05"),
			}
			posts = append(posts, post)
		}
	})

	return posts, nil
}

func getBaseURL(fullURL string) string {
	if idx := strings.Index(fullURL[8:], "/"); idx != -1 {
		return fullURL[:8+idx]
	}
	return fullURL
}

func generateID(title, url string) string {
	return fmt.Sprintf("%x", title+url)[:8]
}