package common

import "math"

// ListParams holds generic query parameters for any paginated list endpoint.
// Domain-specific filters (e.g. year, status) should be separate parameters in each repository.
type ListParams struct {
	Page   int
	Limit  int
	Search string
	SortBy string
	Order  string
}

// Sanitize applies defaults, clamps values, and whitelists sort columns.
// Each repository calls this with its own allowedSorts map and default sort column.
func (p *ListParams) Sanitize(allowedSorts map[string]bool, defaultSort string) {
	if p.Page < 1 {
		p.Page = 1
	}

	if p.Limit < 1 {
		p.Limit = 10
	}
	if p.Limit > 100 {
		p.Limit = 100
	}

	if !allowedSorts[p.SortBy] {
		p.SortBy = defaultSort
	}

	if p.Order != "asc" && p.Order != "desc" {
		p.Order = "desc"
	}
}

// Offset returns the calculated offset for SQL queries.
func (p ListParams) Offset() int {
	return (p.Page - 1) * p.Limit
}

// OrderClause returns a safe "column direction" string for GORM .Order().
func (p ListParams) OrderClause() string {
	return p.SortBy + " " + p.Order
}

// PaginationMeta holds pagination metadata for API responses.
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	Limit       int   `json:"limit"`
	TotalItems  int64 `json:"total_items"`
}

// NewPaginationMeta builds metadata from sanitized params and total count.
func NewPaginationMeta(params ListParams, total int64) PaginationMeta {
	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))
	return PaginationMeta{
		CurrentPage: params.Page,
		TotalPages:  totalPages,
		Limit:       params.Limit,
		TotalItems:  total,
	}
}

// PaginatedResponse is a generic paginated API response.
type PaginatedResponse[T any] struct {
	Meta PaginationMeta `json:"meta"`
	Data []T            `json:"data"`
}

// NewPaginatedResponse builds a complete paginated response in one call.
func NewPaginatedResponse[T any](params ListParams, total int64, data []T) PaginatedResponse[T] {
	if data == nil {
		data = []T{}
	}
	return PaginatedResponse[T]{
		Meta: NewPaginationMeta(params, total),
		Data: data,
	}
}
