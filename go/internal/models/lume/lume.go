package lume

import (
	"time"

	"github.com/google/uuid"
)

// LumeType represents the type of travel node
type LumeType string

const (
	LumeTypeUnspecified   LumeType = "LUME_TYPE_UNSPECIFIED"
	LumeTypeCity          LumeType = "LUME_TYPE_CITY"
	LumeTypeAttraction    LumeType = "LUME_TYPE_ATTRACTION"
	LumeTypeAccommodation LumeType = "LUME_TYPE_ACCOMMODATION"
	LumeTypeRestaurant    LumeType = "LUME_TYPE_RESTAURANT"
	LumeTypeTransportHub  LumeType = "LUME_TYPE_TRANSPORT_HUB"
	LumeTypeActivity      LumeType = "LUME_TYPE_ACTIVITY"
	LumeTypeShopping      LumeType = "LUME_TYPE_SHOPPING"
	LumeTypeEntertainment LumeType = "LUME_TYPE_ENTERTAINMENT"
	LumeTypeCustom        LumeType = "LUME_TYPE_CUSTOM"
)

// Lume represents a single travel node/location in the domain
type Lume struct {
	// Internal database ID (not exposed in API)
	ID int64 `json:"-"`

	// Unique identifier (UUID)
	LumeID string `json:"lume_id"`

	// Lumo reference (parent container)
	LumoID string `json:"lumo_id"`

	// Lume type (e.g. CITY, ATTRACTION, etc.)
	Type LumeType `json:"type"`

	// Display title (e.g. "Paris," "Eiffel Tower")
	Name string `json:"name"`

	// Optional start date/time for scheduling
	DateStart *time.Time `json:"date_start,omitempty"`

	// Optional end date/time
	DateEnd *time.Time `json:"date_end,omitempty"`

	// Optional GPS coordinates
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`

	// Optional textual address
	Address *string `json:"address,omitempty"`

	// Freeform notes/description
	Description string `json:"description"`

	// URLs to uploaded photos
	Images []string `json:"images"`

	// Optional taxonomy tags
	CategoryTags []string `json:"category_tags"`

	// Optional external reservation URL
	BookingLink *string `json:"booking_link,omitempty"`

	// System timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewLume creates a new Lume with generated UUID
func NewLume(lumoID string, lumeType LumeType, name string) *Lume {
	now := time.Now()
	return &Lume{
		LumeID:       uuid.New().String(),
		LumoID:       lumoID,
		Type:         lumeType,
		Name:         name,
		Images:       make([]string, 0),
		CategoryTags: make([]string, 0),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// HasLocation returns true if the Lume has GPS coordinates
func (l *Lume) HasLocation() bool {
	return l.Latitude != nil && l.Longitude != nil
}

// HasSchedule returns true if the Lume has date/time scheduling
func (l *Lume) HasSchedule() bool {
	return l.DateStart != nil || l.DateEnd != nil
}
