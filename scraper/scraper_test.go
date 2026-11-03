package scraper

import "testing"

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
