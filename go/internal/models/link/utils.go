package link

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	linkpb "github.com/mcdev12/lumo/go/internal/genproto/link"
)

// DomainToProto converts domain Link to protobuf Link
func DomainToProto(domainLink *Link) *linkpb.Link {
	proto := &linkpb.Link{
		LinkId:     domainLink.LinkID,
		FromLumeId: domainLink.FromLumeID,
		ToLumeId:   domainLink.ToLumeID,
		Type:       DomainLinkTypeToProto(domainLink.Type),
		CreatedAt:  timestamppb.New(domainLink.CreatedAt),
		UpdatedAt:  timestamppb.New(domainLink.UpdatedAt),
	}

	// Handle optional notes
	if domainLink.Notes != nil {
		proto.Notes = *domainLink.Notes
	}

	// Handle optional sequence index
	if domainLink.SequenceIndex != nil {
		proto.SequenceIndex = *domainLink.SequenceIndex
	}

	// Handle travel details
	if domainLink.Travel != nil {
		proto.Travel = &linkpb.TravelDetails{
			Mode:           DomainTravelModeToProto(domainLink.Travel.Mode),
			DurationSec:    domainLink.Travel.DurationSec,
			CostEstimate:   domainLink.Travel.CostEstimate,
			DistanceMeters: domainLink.Travel.DistanceMeters,
		}
	}

	return proto
}

// ProtoToDomain converts protobuf Link to domain Link
func ProtoToDomain(protoLink *linkpb.Link) *Link {
	domain := &Link{
		LinkID:     protoLink.LinkId,
		FromLumeID: protoLink.FromLumeId,
		ToLumeID:   protoLink.ToLumeId,
		Type:       ProtoLinkTypeToDomain(protoLink.Type),
		CreatedAt:  protoLink.CreatedAt.AsTime(),
		UpdatedAt:  protoLink.UpdatedAt.AsTime(),
	}

	// Handle optional notes
	if protoLink.Notes != "" {
		domain.Notes = &protoLink.Notes
	}

	// Handle optional sequence index
	if protoLink.SequenceIndex != 0 {
		domain.SequenceIndex = &protoLink.SequenceIndex
	}

	// Handle travel details
	if protoLink.Travel != nil {
		domain.Travel = &TravelDetails{
			Mode:           ProtoTravelModeToDomain(protoLink.Travel.Mode),
			DurationSec:    protoLink.Travel.DurationSec,
			CostEstimate:   protoLink.Travel.CostEstimate,
			DistanceMeters: protoLink.Travel.DistanceMeters,
		}
	}

	return domain
}

// Domain LinkType to Proto LinkType conversion
func DomainLinkTypeToProto(dt LinkType) linkpb.LinkType {
	switch dt {
	case LinkTypeTravel:
		return linkpb.LinkType_LINK_TYPE_TRAVEL
	case LinkTypeRecommended:
		return linkpb.LinkType_LINK_TYPE_RECOMMENDED
	case LinkTypeCustom:
		return linkpb.LinkType_LINK_TYPE_CUSTOM
	default:
		return linkpb.LinkType_LINK_TYPE_UNSPECIFIED
	}
}

// Proto LinkType to Domain LinkType conversion
func ProtoLinkTypeToDomain(pt linkpb.LinkType) LinkType {
	switch pt {
	case linkpb.LinkType_LINK_TYPE_TRAVEL:
		return LinkTypeTravel
	case linkpb.LinkType_LINK_TYPE_RECOMMENDED:
		return LinkTypeRecommended
	case linkpb.LinkType_LINK_TYPE_CUSTOM:
		return LinkTypeCustom
	default:
		return LinkTypeUnspecified
	}
}

// Domain TravelMode to Proto TravelMode conversion
func DomainTravelModeToProto(dt TravelMode) linkpb.TravelMode {
	switch dt {
	case TravelModeFlight:
		return linkpb.TravelMode_TRAVEL_MODE_FLIGHT
	case TravelModeTrain:
		return linkpb.TravelMode_TRAVEL_MODE_TRAIN
	case TravelModeBus:
		return linkpb.TravelMode_TRAVEL_MODE_BUS
	case TravelModeDrive:
		return linkpb.TravelMode_TRAVEL_MODE_DRIVE
	case TravelModeUber:
		return linkpb.TravelMode_TRAVEL_MODE_UBER
	case TravelModeMetro:
		return linkpb.TravelMode_TRAVEL_MODE_METRO
	default:
		return linkpb.TravelMode_TRAVEL_MODE_UNSPECIFIED
	}
}

// Proto TravelMode to Domain TravelMode conversion
func ProtoTravelModeToDomain(pt linkpb.TravelMode) TravelMode {
	switch pt {
	case linkpb.TravelMode_TRAVEL_MODE_FLIGHT:
		return TravelModeFlight
	case linkpb.TravelMode_TRAVEL_MODE_TRAIN:
		return TravelModeTrain
	case linkpb.TravelMode_TRAVEL_MODE_BUS:
		return TravelModeBus
	case linkpb.TravelMode_TRAVEL_MODE_DRIVE:
		return TravelModeDrive
	case linkpb.TravelMode_TRAVEL_MODE_UBER:
		return TravelModeUber
	case linkpb.TravelMode_TRAVEL_MODE_METRO:
		return TravelModeMetro
	default:
		return TravelModeUnspecified
	}
}

// JSONTravelDetails is a helper for handling JSONB travel details in the database
type JSONTravelDetails TravelDetails

// Value implements the driver.Valuer interface for JSONTravelDetails
func (j JSONTravelDetails) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for JSONTravelDetails
func (j *JSONTravelDetails) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), j)
	}
	return json.Unmarshal(bytes, j)
}

// ParseUUID safely parses a string to UUID
func ParseUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}

// ValidateTravelDetails validates travel details
func ValidateTravelDetails(td *TravelDetails) bool {
	if td == nil {
		return true // Optional field is valid when nil
	}

	// Basic validation rules
	if td.DurationSec < 0 {
		return false
	}
	if td.DistanceMeters < 0 {
		return false
	}

	return true
}
