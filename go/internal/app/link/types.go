package link

import modellink "github.com/mcdev12/lumo/go/internal/models/link"

// CreateLinkRequest represents the business layer's create request
type CreateLinkRequest struct {
	FromLumeID    string
	ToLumeID      string
	Type          modellink.LinkType
	Notes         *string
	SequenceIndex *int32
	TravelDetails *TravelDetailsRequest
}

// TravelDetailsRequest represents the travel details for a link
type TravelDetailsRequest struct {
	Mode           modellink.TravelMode
	DurationSec    int32
	CostEstimate   float64
	DistanceMeters float64
}

// UpdateLinkRequest represents the business layer's update request
type UpdateLinkRequest struct {
	FromLumeID    string
	ToLumeID      string
	Type          modellink.LinkType
	Notes         *string
	SequenceIndex *int32
	TravelDetails *TravelDetailsRequest
	// Fields to update (from field mask)
	UpdateFields []string
}

// ListLinksRequest represents pagination parameters
type ListLinksRequest struct {
	Limit  int32
	Offset int32
}
