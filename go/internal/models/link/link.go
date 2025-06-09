package link

import (
	"time"

	"github.com/google/uuid"
)

// LinkType represents the type of relationship between two Lumés
type LinkType string

const (
	LinkTypeUnspecified LinkType = "LINK_TYPE_UNSPECIFIED"
	LinkTypeTravel      LinkType = "TRAVEL"
	LinkTypeRecommended LinkType = "RECOMMENDED"
	LinkTypeCustom      LinkType = "CUSTOM"
)

// TravelMode represents the mode of transportation
type TravelMode string

const (
	TravelModeUnspecified TravelMode = "TRAVEL_MODE_UNSPECIFIED"
	TravelModeFlight      TravelMode = "FLIGHT"
	TravelModeTrain       TravelMode = "TRAIN"
	TravelModeBus         TravelMode = "BUS"
	TravelModeDrive       TravelMode = "DRIVE"
	TravelModeUber        TravelMode = "UBER"
	TravelModeMetro       TravelMode = "METRO"
)

// TravelDetails contains all the movement-specific metadata
type TravelDetails struct {
	Mode           TravelMode `json:"mode"`
	DurationSec    int32      `json:"duration_sec"`
	CostEstimate   float64    `json:"cost_estimate"`
	DistanceMeters float64    `json:"distance_meters"`
}

// Link represents a connection between two Lumés in the domain
type Link struct {
	// Internal database ID (not exposed in API)
	ID int64 `json:"-"`

	// External UUID for API clients
	LinkID string `json:"link_id"`

	// Which two Lumés this edge connects
	FromLumeID string `json:"from_lume_id"`
	ToLumeID   string `json:"to_lume_id"`

	// The high-level type of relation
	Type LinkType `json:"type"`

	// Travel details (only populated when type == TRAVEL)
	Travel *TravelDetails `json:"travel,omitempty"`

	// Freeform notes about this relationship
	Notes *string `json:"notes,omitempty"`

	// Optional hint for rendering order in lists/timelines
	SequenceIndex *int32 `json:"sequence_index,omitempty"`

	// System timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewLink creates a new Link with generated UUID
func NewLink(fromLumeID, toLumeID string, linkType LinkType) *Link {
	now := time.Now()
	return &Link{
		LinkID:     uuid.New().String(),
		FromLumeID: fromLumeID,
		ToLumeID:   toLumeID,
		Type:       linkType,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// HasTravelDetails returns true if the Link has travel details
func (l *Link) HasTravelDetails() bool {
	return l.Travel != nil
}

// HasSequenceIndex returns true if the Link has a sequence index
func (l *Link) HasSequenceIndex() bool {
	return l.SequenceIndex != nil
}
