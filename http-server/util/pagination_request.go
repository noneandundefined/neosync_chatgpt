package util

import (
	"net/http"
	"strconv"
)

type Pagination struct {
	Page          int
	Limit         int
	Offset        int
	Search        string
	ColumnSortKey string
	ColumnSortDir string
}

func GetPagination(r *http.Request, defaultLimit int) Pagination {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("search")
	columnSortKey := r.URL.Query().Get("columnSortKey")
	columnSortDir := r.URL.Query().Get("columnSortDir")

	page := 1
	limit := defaultLimit
	if limit < 0 {
		limit = 10
	}

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	offset := (page - 1) * limit

	return Pagination{
		Page:          page,
		Limit:         limit,
		Offset:        offset,
		Search:        search,
		ColumnSortKey: columnSortKey,
		ColumnSortDir: columnSortDir,
	}
}
