package main

import (
	"fmt"
	"log"
	"os"

	"go-newpost/scraper"
)

func main() {
	var targetURL string

	if len(os.Args) > 1 {
		targetURL = os.Args[1]
	} else {
		fmt.Print("테스트할 URL을 입력하세요: ")
		fmt.Scanln(&targetURL)
	}

	if targetURL == "" {
		log.Fatal("URL이 필요합니다")
	}

	fmt.Printf("🔍 스크래핑 시작: %s\n\n", targetURL)

	s := scraper.New(targetURL)
	posts, err := s.GetLatestPosts()
	if err != nil {
		log.Fatalf("스크래핑 실패: %v", err)
	}

	fmt.Printf("📊 총 %d개의 게시글을 찾았습니다.\n\n", len(posts))

	if len(posts) == 0 {
		fmt.Println("⚠️  게시글을 찾지 못했습니다.")
		fmt.Println("다음을 확인해보세요:")
		fmt.Println("- URL이 올바른지 확인")
		fmt.Println("- 게시판이 다음 HTML 구조 중 하나를 사용하는지 확인:")
		fmt.Println("  * <tr> 태그 (테이블 형태)")
		fmt.Println("  * class=\"post\" (일반적인 게시글)")
		fmt.Println("  * class=\"article\" (기사 형태)")
		fmt.Println("  * class=\"item\" (아이템 형태)")
		fmt.Println("- 각 게시글 안에 <a> 태그가 있는지 확인")
		return
	}

	fmt.Println("📝 찾은 게시글 목록:")
	fmt.Println("=" + fmt.Sprintf("%0*s", 80, ""))

	for i, post := range posts {
		fmt.Printf("게시글 #%d\n", i+1)
		fmt.Printf("  ID: %s\n", post.ID)
		fmt.Printf("  제목: %s\n", post.Title)
		fmt.Printf("  링크: %s\n", post.URL)
		fmt.Printf("  발견시간: %s\n", post.Date)
		fmt.Println()

		if i >= 9 {
			fmt.Printf("... 외 %d개 더 있습니다.\n", len(posts)-10)
			break
		}
	}

	fmt.Println("=" + fmt.Sprintf("%0*s", 80, ""))
	fmt.Println("✅ 스크래핑 테스트 완료!")
}