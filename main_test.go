package main

import (
	"testing"

	"go-newpost/scraper"
)

func TestRemoveDuplicates(t *testing.T) {
	posts := []scraper.Post{
		{Title: "A", PublishDate: "2026-11-05", Source: "LH"},
		{Title: "A", PublishDate: "2026-11-05", Source: "SH"}, // 소스만 다르면 중복
		{Title: "A", PublishDate: "2026-11-04", Source: "LH"},
		{Title: "B", PublishDate: "2026-11-05", Source: "LH"},
	}
	got := removeDuplicates(posts)
	if len(got) != 3 {
		t.Fatalf("len = %d, want 3", len(got))
	}
	if got[0].Source != "LH" {
		t.Errorf("먼저 나온 항목이 유지되어야 함: %+v", got[0])
	}
}
