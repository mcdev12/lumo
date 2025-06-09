package lume

import (
	"time"

	modellume "github.com/mcdev12/lumo/go/internal/models/lume"
)

// CreateLumeRequest represents the business layer's create request
type CreateLumeRequest struct {
	LumoID      string // Parent Lumo UUID (required)
	Name        string
	Type        modellume.LumeType
	Description string
	// Additional fields from the domain model
	DateStart    *time.Time
	DateEnd      *time.Time
	Latitude     *float64
	Longitude    *float64
	Address      *string
	Images       []string
	CategoryTags []string
	BookingLink  *string
}

// UpdateLumeRequest represents the business layer's update request
type UpdateLumeRequest struct {
	Name        string
	Type        modellume.LumeType
	Description string
	// Additional fields from the domain model
	DateStart    *time.Time
	DateEnd      *time.Time
	Latitude     *float64
	Longitude    *float64
	Address      *string
	Images       []string
	CategoryTags []string
	BookingLink  *string
	// Fields to update (from field mask)
	UpdateFields []string
}

// ListLumesRequest represents pagination parameters
type ListLumesRequest struct {
	LumoID string
	Limit  int32
	Offset int32
}

// ListLumesByTypeRequest represents type filtering with pagination
type ListLumesByTypeRequest struct {
	LumoID string
	Type   modellume.LumeType
	Limit  int32
	Offset int32
}

// SearchLumesByLocationRequest represents location search parameters
type SearchLumesByLocationRequest struct {
	LumoID string
	MinLat float64
	MaxLat float64
	MinLng float64
	MaxLng float64
	Limit  int32
	Offset int32
}
