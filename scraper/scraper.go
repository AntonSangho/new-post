package scraper

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Post struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Date        string `json:"date"`
	PublishDate string `json:"publish_date"`
	Status      string `json:"status"`
	Source      string `json:"source"`
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

	// 서울시 주택 사이트의 테이블 구조에 맞춰 스크래핑
	doc.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		tds := row.Find("td")
		if tds.Length() < 4 {
			return // 테이블 헤더나 잘못된 행 건너뛰기
		}

		// URL에 따라 다른 컬럼 인덱스 사용
		var titleIndex, dateIndex, linkIndex, statusIndex int
		var source string
		if strings.Contains(s.targetURL, "/lh/") {
			// LH 사이트: 번호(0), 세부유형(1), 공고명(2), 지역(3), 공고게시일(4), 마감일(5), 상태(6), 링크(7)
			titleIndex = 2
			dateIndex = 4
			statusIndex = 6
			linkIndex = 7
			source = "LH"
		} else if strings.Contains(s.targetURL, "/sh/") {
			// SH 사이트: 번호(0), 청약유형(1), 공고명(2), 공고게시일(3), 발표일(4), 담당부서(5), 링크(6)
			titleIndex = 2
			dateIndex = 3
			statusIndex = -1 // SH 사이트는 상태 컬럼이 없음
			linkIndex = 6
			source = "SH"
		} else {
			return // 지원하지 않는 사이트
		}

		if tds.Length() <= linkIndex {
			return // 필요한 컬럼이 없음
		}

		// 제목 추출
		title := strings.TrimSpace(tds.Eq(titleIndex).Text())
		if title == "" {
			return
		}

		// 공고게시일 추출
		publishDate := strings.TrimSpace(tds.Eq(dateIndex).Text())
		
		// 상태 추출 (LH 사이트만)
		status := ""
		if statusIndex >= 0 && tds.Length() > statusIndex {
			status = strings.TrimSpace(tds.Eq(statusIndex).Text())
		}

		// 공고게시일 유효성 검사 (YYYY-MM-DD 형식인지 확인)
		if !isValidDate(publishDate) {
			return // 유효하지 않은 날짜 형식
		}

		// 링크 추출
		link, exists := tds.Eq(linkIndex).Find("a").First().Attr("href")
		if !exists {
			// 링크가 없는 경우 원본 페이지 URL 사용
			link = s.targetURL
		} else if !strings.HasPrefix(link, "http") {
			if strings.HasPrefix(link, "/") {
				link = getBaseURL(s.targetURL) + link
			} else {
				link = getBaseURL(s.targetURL) + "/" + link
			}
		}

		post := Post{
			ID:          generateID(title, publishDate, source),
			Title:       title,
			URL:         link,
			Date:        time.Now().Format("2006-01-02 15:04:05"),
			PublishDate: publishDate,
			Status:      status,
			Source:      source,
		}
		posts = append(posts, post)
	})

	return posts, nil
}

func (s *Scraper) GetRecentPosts(days int) ([]Post, error) {
	posts, err := s.GetLatestPosts()
	if err != nil {
		return nil, err
	}

	var recentPosts []Post
	now := time.Now()
	cutoffDate := now.AddDate(0, 0, -days)

	for _, post := range posts {
		postTime, err := time.Parse("2006-01-02", post.PublishDate)
		if err != nil {
			continue // 날짜 파싱 실패시 제외
		}

		// 최근 N일 이내 공고만 포함
		if postTime.After(cutoffDate) || postTime.Equal(cutoffDate.Truncate(24*time.Hour)) {
			// 상태가 "공고중"인지 확인 (LH 사이트만)
			if post.Status != "" && !isValidStatus(post.Status) {
				continue
			}
			recentPosts = append(recentPosts, post)
		}
	}

	return recentPosts, nil
}

func isValidStatus(status string) bool {
	validStatuses := []string{"공고중", "진행중", "모집중"}
	status = strings.TrimSpace(status)
	
	for _, valid := range validStatuses {
		if strings.Contains(status, valid) {
			return true
		}
	}
	return false
}

func getBaseURL(fullURL string) string {
	if idx := strings.Index(fullURL[8:], "/"); idx != -1 {
		return fullURL[:8+idx]
	}
	return fullURL
}

func generateID(title, publishDate, source string) string {
	sum := sha256.Sum256([]byte(title + publishDate + source))
	return fmt.Sprintf("%x", sum)[:8]
}

func isValidDate(dateStr string) bool {
	if len(dateStr) != 10 {
		return false
	}
	
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}