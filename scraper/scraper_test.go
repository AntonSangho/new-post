package scraper

import (
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
