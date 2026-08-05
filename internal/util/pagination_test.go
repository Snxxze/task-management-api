package util_test

import (
	"testing"
	"task-management-api/internal/util"
)

func TestPagination_Offset(t *testing.T) {
	p := util.Pagination{Page: 1, Limit: 10}
	if p.Offset() != 0 {
		t.Errorf("expected 0, got %d", p.Offset())
	}

	p2 := util.Pagination{Page: 3, Limit: 10}
	if p2.Offset() != 20 {
		t.Errorf("expected 20, got %d", p2.Offset())
	}
}

func TestNewPaginationMeta(t *testing.T) {
	meta := util.NewPaginationMeta(1, 10, 25)
	if meta.TotalPages != 3 {
		t.Errorf("expected 3 total pages, got %d", meta.TotalPages)
	}

	metaEmpty := util.NewPaginationMeta(1, 10, 0)
	if metaEmpty.TotalPages != 1 {
		t.Errorf("expected 1 total page for empty items, got %d", metaEmpty.TotalPages)
	}
}
