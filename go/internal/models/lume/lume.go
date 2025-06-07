package lume

import "time"

// LumeType is a domain‐level enum that mirrors (but does NOT import) the Protobuf LumeType.
type LumeType string

const (
	LumeTypeUnspecified     LumeType = "LUME_TYPE_UNSPECIFIED"
	LumeTypeStop            LumeType = "LUME_TYPE_STOP"
	LumeTypeAccommodation   LumeType = "LUME_TYPE_ACCOMMODATION"
	LumeTypePointOfInterest LumeType = "LUME_TYPE_POINT_OF_INTEREST"
	LumeTypeMeal            LumeType = "LUME_TYPE_MEAL"
	LumeTypeTransport       LumeType = "LUME_TYPE_TRANSPORT"
	LumeTypeCustom          LumeType = "LUME_TYPE_CUSTOM"
)

// Lume is the internal Go struct representation of a "node" (Lume).
// It lives in your domain and does not import any Protobuf/SQL/Transport packages.
type Lume struct {
	// Id is the auto-incrementing database ID (internal use)
	Id int64

	// LumeId is the UUID of this Lume (public identifier)
	LumeId string

	// LumoID is the UUID of the parent Lumo graph.
	LumoID string

	// Label is a human‐readable name (e.g., "Eiffel Tower").
	Label string

	// Type categorizes the Lume (Stop, Accommodation, etc.).
	Type LumeType

	// Description holds any free‐text notes or details.
	Description string

	// Metadata contains arbitrary key/value pairs (mirrors google.protobuf.Struct).
	// We use map[string]interface{} in the domain.
	Metadata map[string]interface{}

	// CreatedAt is when this Lume was first created.
	CreatedAt time.Time

	// UpdatedAt is when this Lume was last modified.
	UpdatedAt time.Time
}
