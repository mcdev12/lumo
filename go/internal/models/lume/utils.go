package lume

import (
	lumepb "github.com/mcdev12/lumo/go/internal/genproto/lume/v1"
)

// ProtoToDomain converts protobuf Lume to domain Lume
func ProtoToDomain(protoLume *lumepb.Lume) *Lume {
	domain := &Lume{
		LumeID:       protoLume.LumeId,
		LumoID:       protoLume.LumoId,
		Type:         ProtoLumeTypeToDomain(protoLume.Type),
		Name:         protoLume.Name,
		Description:  protoLume.Description,
		Images:       protoLume.Images,
		CategoryTags: protoLume.CategoryTags,
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
