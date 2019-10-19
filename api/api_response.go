package api

import (
	"hackathon.com/pyz/dbm"
	"hackathon.com/pyz/env"
)

// ErrorResponse is the JSON on error
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error_code,omitempty"`
}

// DefaultResponse is default api JSON response
type DefaultResponse struct {
	Success     bool      `json:"success"`
	Message     string    `json:"message,omitempty"`
	CompletedAt *dbm.Time `json:"completed_at,omitempty"`
}

// ProfileResponse is profile api JSON response
type ProfileResponse struct {
	Success bool        `json:"success"`
	Profile dbm.Profile `json:"profile"`
}

// NewsResponse is news api JSON response
type NewsResponse struct {
	Success bool      `json:"success"`
	News    []env.New `json:"news"`
}

// PacksResponse is packs api JSON response
type PacksResponse struct {
	Success bool       `json:"success"`
	Packs   []env.Pack `json:"packs"`
}
