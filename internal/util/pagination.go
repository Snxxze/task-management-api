package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	SortBy string `json:"sort_by"`
	Order  string `json:"order"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

func GetPagination(c *gin.Context) Pagination {
	page := 1
	limit := 10
	sortBy := "created_at"
	order := "desc"

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if sb := c.Query("sort_by"); sb != "" {
		sortBy = sb
	}

	if ord := c.Query("order"); ord == "asc" || ord == "desc" {
		order = ord
	}

	return Pagination{
		Page:   page,
		Limit:  limit,
		SortBy: sortBy,
		Order:  order,
	}
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}

func NewPaginationMeta(page, limit int, totalItems int64) PaginationMeta {
	totalPages := int((totalItems + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	return PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
