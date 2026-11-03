package scraper

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGenerateID_SamePrefixDifferentTitles(t *testing.T) {
	a := generateID("서울 행복주택 모집 공고 A", "2026-11-01", "LH")
	b := generateID("서울 행복주택 모집 공고 B", "2026-11-01", "LH")
	if a == b {
		t.Errorf("제목이 다른 공고의 ID가 충돌함: %s", a)
	}
}

func TestGenerateID_Deterministic(t *testing.T) {
	a := generateID("공고", "2026-11-01", "SH")
	b := generateID("공고", "2026-11-01", "SH")
	if a != b || len(a) != 8 {
		t.Errorf("ID가 일정하지 않거나 길이가 8이 아님: %q, %q", a, b)
	}
}

func TestIsRecent_Boundary(t *testing.T) {
	// KST 2026-11-10 00:30 == UTC 2026-11-09 15:30 (UTC 기준으로는 아직 전날)
	now := time.Date(2026, 11, 9, 15, 30, 0, 0, time.UTC)
	cases := []struct {
		date string
		want bool
	}{
		{"2026-11-10", true},  // 오늘
		{"2026-11-03", true},  // 정확히 7일 전
		{"2026-11-02", false}, // 8일 전
		{"invalid", false},
	}
	for _, c := range cases {
		if got := isRecent(c.date, now, 7); got != c.want {
			t.Errorf("isRecent(%q) = %v, want %v", c.date, got, c.want)
		}
	}
}

func TestIsRecent_IndependentOfTimeOfDay(t *testing.T) {
	for _, h := range []int{0, 5, 12, 23} {
		now := time.Date(2026, 11, 10, h, 0, 0, 0, kst)
		if !isRecent("2026-11-03", now, 7) {
			t.Errorf("%d시에 7일 전 공고가 제외됨", h)
		}
	}
}

const lhHTML = `<table><tbody>
<tr><th>번호</th><th>유형</th><th>공고명</th><th>지역</th><th>게시일</th><th>마감일</th><th>상태</th><th>링크</th></tr>
<tr><td>1</td><td>행복주택</td><td> LH 공고 1 </td><td>서울</td><td>2026-11-05</td><td>2026-11-20</td><td>공고중</td><td><a href="/view?id=1">보기</a></td></tr>
<tr><td>2</td><td>국민임대</td><td>LH 마감 공고</td><td>서울</td><td>2026-11-05</td><td>2026-11-06</td><td>접수마감</td><td><a href="https://other.example/2">보기</a></td></tr>
<tr><td>3</td><td>국민임대</td><td>날짜 형식 오류</td><td>서울</td><td>11/05</td><td>-</td><td>공고중</td><td></td></tr>
</tbody></table>`

const shHTML = `<table><tbody>
<tr><td>1</td><td>공공분양</td><td>SH 공고 1</td><td>2026-11-04</td><td>2026-11-30</td><td>주택공급부</td><td><a href="view?id=9">보기</a></td></tr>
<tr><td>2</td><td>공공분양</td><td>링크 없는 공고</td><td>2026-11-04</td><td>2026-11-30</td><td>주택공급부</td><td></td></tr>
</tbody></table>`

func serve(t *testing.T, path, body string) string {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	}))
	t.Cleanup(srv.Close)
	return srv.URL + path
}

func TestGetLatestPosts_LH(t *testing.T) {
	url := serve(t, "/lh/board", lhHTML)
	posts, err := New(url).GetLatestPosts()
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != 2 { // 날짜 형식 오류 행은 제외
		t.Fatalf("len = %d, want 2", len(posts))
	}
	p := posts[0]
	if p.Title != "LH 공고 1" || p.Source != "LH" || p.Status != "공고중" || p.PublishDate != "2026-11-05" {
		t.Errorf("파싱 결과 불일치: %+v", p)
	}
	if want := url[:strings.Index(url, "/lh/")] + "/view?id=1"; p.URL != want {
		t.Errorf("상대 링크 URL = %q, want %q", p.URL, want)
	}
	if posts[1].URL != "https://other.example/2" {
		t.Errorf("절대 링크가 유지되지 않음: %q", posts[1].URL)
	}
}

func TestGetLatestPosts_SH(t *testing.T) {
	url := serve(t, "/sh/board", shHTML)
	posts, err := New(url).GetLatestPosts()
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != 2 {
		t.Fatalf("len = %d, want 2", len(posts))
	}
	if posts[0].Source != "SH" || posts[0].Status != "" {
		t.Errorf("SH 파싱 결과 불일치: %+v", posts[0])
	}
	if !strings.HasSuffix(posts[0].URL, "/view?id=9") {
		t.Errorf("슬래시 없는 상대 링크 처리 실패: %q", posts[0].URL)
	}
	if posts[1].URL != url {
		t.Errorf("링크 없는 행은 원본 URL이어야 함: %q", posts[1].URL)
	}
}

func TestGetLatestPosts_UnsupportedSiteAndBadStatus(t *testing.T) {
	posts, err := New(serve(t, "/other/board", lhHTML)).GetLatestPosts()
	if err != nil || len(posts) != 0 {
		t.Errorf("지원하지 않는 사이트: posts=%d err=%v", len(posts), err)
	}

	bad := httptest.NewServer(http.NotFoundHandler())
	defer bad.Close()
	if _, err := New(bad.URL + "/lh/board").GetLatestPosts(); err == nil {
		t.Error("404 응답에서 에러가 반환되어야 함")
	}
}

func TestIsValidStatus(t *testing.T) {
	for s, want := range map[string]bool{"공고중": true, " 접수 진행중 ": true, "모집중": true, "접수마감": false, "": false} {
		if got := isValidStatus(s); got != want {
			t.Errorf("isValidStatus(%q) = %v, want %v", s, got, want)
		}
	}
}

func TestIsValidDate(t *testing.T) {
	for s, want := range map[string]bool{"2026-11-05": true, "2026-13-01": false, "2026-1-5": false, "": false, "20261105": false} {
		if got := isValidDate(s); got != want {
			t.Errorf("isValidDate(%q) = %v, want %v", s, got, want)
		}
	}
}
