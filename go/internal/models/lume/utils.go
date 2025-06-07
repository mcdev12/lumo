package lume

import (
	"database/sql/driver"
	"time"

	"github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"

	lumepb "github.com/mcdev12/lumo/go/internal/genproto/protobuf/lume"
)

// StringArray handles PostgreSQL text arrays
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	return pq.Array(a).Value()
}

func (a *StringArray) Scan(value interface{}) error {
	return pq.Array(a).Scan(value)
}

// DomainToProto converts domain Lume to protobuf Lume
func DomainToProto(domainLume *Lume) *lumepb.Lume {
	proto := &lumepb.Lume{
		LumeId:       domainLume.LumeID,
		LumoId:       domainLume.LumoID,
		Type:         DomainLumeTypeToProto(domainLume.Type),
		Name:         domainLume.Name,
		Description:  domainLume.Description,
		Images:       EnsureStringArray(domainLume.Images),
		CategoryTags: EnsureStringArray(domainLume.CategoryTags),
		CreatedAt:    timestamppb.New(domainLume.CreatedAt),
		UpdatedAt:    timestamppb.New(domainLume.UpdatedAt),
	}

	// Handle optional timestamps
	if domainLume.DateStart != nil {
		proto.DateStart = timestamppb.New(*domainLume.DateStart)
	}
	if domainLume.DateEnd != nil {
		proto.DateEnd = timestamppb.New(*domainLume.DateEnd)
	}

	// Handle optional coordinates
	if domainLume.Latitude != nil {
		proto.Latitude = *domainLume.Latitude
	}
	if domainLume.Longitude != nil {
		proto.Longitude = *domainLume.Longitude
	}

	// Handle optional address
	if domainLume.Address != nil {
		proto.Address = *domainLume.Address
	}

	// Handle optional booking link
	if domainLume.BookingLink != nil {
		proto.BookingLink = *domainLume.BookingLink
	}

	return proto
}

// ProtoToDomain converts protobuf Lume to domain Lume
func ProtoToDomain(protoLume *lumepb.Lume) *Lume {
	domain := &Lume{
		LumeID:       protoLume.LumeId,
		LumoID:       protoLume.LumoId,
		Type:         ProtoLumeTypeToDomain(protoLume.Type),
		Name:         protoLume.Name,
		Description:  protoLume.Description,
		Images:       EnsureStringArray(protoLume.Images),
		CategoryTags: EnsureStringArray(protoLume.CategoryTags),
		CreatedAt:    protoLume.CreatedAt.AsTime(),
		UpdatedAt:    protoLume.UpdatedAt.AsTime(),
	}

	// Handle optional timestamps
	if protoLume.DateStart != nil {
		dateStart := protoLume.DateStart.AsTime()
		domain.DateStart = &dateStart
	}
	if protoLume.DateEnd != nil {
		dateEnd := protoLume.DateEnd.AsTime()
		domain.DateEnd = &dateEnd
	}

	// Handle optional coordinates
	if protoLume.Latitude != 0 {
		domain.Latitude = &protoLume.Latitude
	}
	if protoLume.Longitude != 0 {
		domain.Longitude = &protoLume.Longitude
	}

	// Handle optional address
	if protoLume.Address != "" {
		domain.Address = &protoLume.Address
	}

	// Handle optional booking link
	if protoLume.BookingLink != "" {
		domain.BookingLink = &protoLume.BookingLink
	}

	return domain
}

// Proto LumeType conversion functions
func StringToLumeType(s string) lumepb.LumeType {
	switch s {
	case "CITY":
		return lumepb.LumeType_LUME_TYPE_CITY
	case "ATTRACTION":
		return lumepb.LumeType_LUME_TYPE_ATTRACTION
	case "ACCOMMODATION":
		return lumepb.LumeType_LUME_TYPE_ACCOMMODATION
	case "RESTAURANT":
		return lumepb.LumeType_LUME_TYPE_RESTAURANT
	case "TRANSPORT_HUB":
		return lumepb.LumeType_LUME_TYPE_TRANSPORT_HUB
	case "ACTIVITY":
		return lumepb.LumeType_LUME_TYPE_ACTIVITY
	case "SHOPPING":
		return lumepb.LumeType_LUME_TYPE_SHOPPING
	case "ENTERTAINMENT":
		return lumepb.LumeType_LUME_TYPE_ENTERTAINMENT
	case "CUSTOM":
		return lumepb.LumeType_LUME_TYPE_CUSTOM
	default:
		return lumepb.LumeType_LUME_TYPE_UNSPECIFIED
	}
}

func LumeTypeToString(lt lumepb.LumeType) string {
	switch lt {
	case lumepb.LumeType_LUME_TYPE_CITY:
		return "CITY"
	case lumepb.LumeType_LUME_TYPE_ATTRACTION:
		return "ATTRACTION"
	case lumepb.LumeType_LUME_TYPE_ACCOMMODATION:
		return "ACCOMMODATION"
	case lumepb.LumeType_LUME_TYPE_RESTAURANT:
		return "RESTAURANT"
	case lumepb.LumeType_LUME_TYPE_TRANSPORT_HUB:
		return "TRANSPORT_HUB"
	case lumepb.LumeType_LUME_TYPE_ACTIVITY:
		return "ACTIVITY"
	case lumepb.LumeType_LUME_TYPE_SHOPPING:
		return "SHOPPING"
	case lumepb.LumeType_LUME_TYPE_ENTERTAINMENT:
		return "ENTERTAINMENT"
	case lumepb.LumeType_LUME_TYPE_CUSTOM:
		return "CUSTOM"
	default:
		return "LUME_TYPE_UNSPECIFIED"
	}
}

// Domain LumeType to Proto LumeType conversion
func DomainLumeTypeToProto(dt LumeType) lumepb.LumeType {
	switch dt {
	case LumeTypeCity:
		return lumepb.LumeType_LUME_TYPE_CITY
	case LumeTypeAttraction:
		return lumepb.LumeType_LUME_TYPE_ATTRACTION
	case LumeTypeAccommodation:
		return lumepb.LumeType_LUME_TYPE_ACCOMMODATION
	case LumeTypeRestaurant:
		return lumepb.LumeType_LUME_TYPE_RESTAURANT
	case LumeTypeTransportHub:
		return lumepb.LumeType_LUME_TYPE_TRANSPORT_HUB
	case LumeTypeActivity:
		return lumepb.LumeType_LUME_TYPE_ACTIVITY
	case LumeTypeShopping:
		return lumepb.LumeType_LUME_TYPE_SHOPPING
	case LumeTypeEntertainment:
		return lumepb.LumeType_LUME_TYPE_ENTERTAINMENT
	case LumeTypeCustom:
		return lumepb.LumeType_LUME_TYPE_CUSTOM
	default:
		return lumepb.LumeType_LUME_TYPE_UNSPECIFIED
	}
}

// Proto LumeType to Domain LumeType conversion
func ProtoLumeTypeToDomain(pt lumepb.LumeType) LumeType {
	switch pt {
	case lumepb.LumeType_LUME_TYPE_CITY:
		return LumeTypeCity
	case lumepb.LumeType_LUME_TYPE_ATTRACTION:
		return LumeTypeAttraction
	case lumepb.LumeType_LUME_TYPE_ACCOMMODATION:
		return LumeTypeAccommodation
	case lumepb.LumeType_LUME_TYPE_RESTAURANT:
		return LumeTypeRestaurant
	case lumepb.LumeType_LUME_TYPE_TRANSPORT_HUB:
		return LumeTypeTransportHub
	case lumepb.LumeType_LUME_TYPE_ACTIVITY:
		return LumeTypeActivity
	case lumepb.LumeType_LUME_TYPE_SHOPPING:
		return LumeTypeShopping
	case lumepb.LumeType_LUME_TYPE_ENTERTAINMENT:
		return LumeTypeEntertainment
	case lumepb.LumeType_LUME_TYPE_CUSTOM:
		return LumeTypeCustom
	default:
		return LumeTypeUnspecified
	}
}

// Validation helpers
func ValidateCoordinates(lat, lng *float64) bool {
	if lat == nil || lng == nil {
		return true // Optional fields are valid when nil
	}
	return *lat >= -90 && *lat <= 90 && *lng >= -180 && *lng <= 180
}

func ValidateDateRange(start, end *time.Time) bool {
	if start == nil || end == nil {
		return true // Optional fields are valid when nil
	}
	return start.Before(*end) || start.Equal(*end)
}

// Helper to ensure empty arrays instead of nil for consistency
func EnsureStringArray(arr []string) []string {
	if arr == nil {
		return make([]string, 0)
	}
	return arr
}
