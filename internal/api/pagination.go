package api

import (
	"fmt"
	"net/http"
	"strconv"
)

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page    int
	PerPage int
	Offset  int
}

// parsePaginationParams extracts pagination parameters from query string
func parsePaginationParams(r *http.Request) *PaginationParams {
	page := 1
	perPage := 100

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if perPageStr := r.URL.Query().Get("per_page"); perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 1000 {
			perPage = pp
		}
	}

	return &PaginationParams{
		Page:    page,
		PerPage: perPage,
		Offset:  (page - 1) * perPage,
	}
}

// paginatePermutations paginates a slice of permutations
func paginatePermutations(perms [][]string, params *PaginationParams) ([][]string, *PaginationMeta) {
	total := len(perms)
	start := params.Offset
	end := start + params.PerPage

	if start >= total {
		return [][]string{}, &PaginationMeta{
			Page:       params.Page,
			PerPage:    params.PerPage,
			Total:      total,
			TotalPages: (total + params.PerPage - 1) / params.PerPage,
		}
	}

	if end > total {
		end = total
	}

	return perms[start:end], &PaginationMeta{
		Page:       params.Page,
		PerPage:    params.PerPage,
		Total:      total,
		TotalPages: (total + params.PerPage - 1) / params.PerPage,
	}
}

// PaginationMeta holds pagination metadata
type PaginationMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// PaginatedResponse wraps paginated data with metadata
type PaginatedResponse struct {
	Data [][]string      `json:"data"`
	Meta *PaginationMeta `json:"meta"`
}

// setPaginationHeaders sets pagination headers on the response
func setPaginationHeaders(w http.ResponseWriter, r *http.Request, meta *PaginationMeta) {
	baseURL := fmt.Sprintf("http://%s%s", r.Host, r.URL.Path)

	// Set Link header for navigation
	var links []string

	if meta.Page > 1 {
		links = append(links, fmt.Sprintf("<%s?page=%d&per_page=%d>; rel=\"prev\"", baseURL, meta.Page-1, meta.PerPage))
	}

	if meta.Page < meta.TotalPages {
		links = append(links, fmt.Sprintf("<%s?page=%d&per_page=%d>; rel=\"next\"", baseURL, meta.Page+1, meta.PerPage))
	}

	links = append(links, fmt.Sprintf("<%s?page=1&per_page=%d>; rel=\"first\"", baseURL, meta.PerPage))
	links = append(links, fmt.Sprintf("<%s?page=%d&per_page=%d>; rel=\"last\"", baseURL, meta.TotalPages, meta.PerPage))

	for _, link := range links {
		w.Header().Add("Link", link)
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(meta.Total))
	w.Header().Set("X-Total-Pages", strconv.Itoa(meta.TotalPages))
	w.Header().Set("X-Page", strconv.Itoa(meta.Page))
	w.Header().Set("X-Per-Page", strconv.Itoa(meta.PerPage))
}
