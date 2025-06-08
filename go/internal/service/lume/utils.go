package lume

import (
	"errors"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"connectrpc.com/connect"
	applume "github.com/mcdev12/lumo/go/internal/app/lume"
	lumepb "github.com/mcdev12/lumo/go/internal/genproto/lume/v1"
	modellume "github.com/mcdev12/lumo/go/internal/models/lume"
)

// domainToProto converts domain Lume to protobuf Lume
func domainToProto(domainLume *modellume.Lume) *lumepb.Lume {
	proto := &lumepb.Lume{
		LumeId:       domainLume.LumeID,
		LumoId:       domainLume.LumoID,
		Type:         DomainLumeTypeToProto(domainLume.Type),
		Name:         domainLume.Name,
		Description:  domainLume.Description,
		Images:       domainLume.Images,
		CategoryTags: domainLume.CategoryTags,
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

// toAppCreateRequest converts a protobuf Lume to an app CreateLumeRequest
func (s *Service) toAppCreateRequest(pbLume *lumepb.CreateLumeRequest) (applume.CreateLumeRequest, error) {
	// Convert timestamps to time.Time pointers
	var dateStart, dateEnd *time.Time
	if pbLume.GetDateStart() != nil {
		t := pbLume.GetDateStart().AsTime()
		dateStart = &t
	}
	if pbLume.GetDateEnd() != nil {
		t := pbLume.GetDateEnd().AsTime()
		dateEnd = &t
	}

	// Convert optional fields
	var latitude, longitude *float64
	var address, bookingLink *string

	if pbLume.GetLatitude() != 0 {
		lat := pbLume.GetLatitude()
		latitude = &lat
	}
	if pbLume.GetLongitude() != 0 {
		lng := pbLume.GetLongitude()
		longitude = &lng
	}
	if pbLume.GetAddress() != "" {
		addr := pbLume.GetAddress()
		address = &addr
	}
	if pbLume.GetBookingLink() != "" {
		link := pbLume.GetBookingLink()
		bookingLink = &link
	}

	return applume.CreateLumeRequest{
		LumoID:       pbLume.GetLumoId(),
		Label:        pbLume.GetName(),
		Type:         modellume.LumeType(pbLume.GetType().String()),
		Description:  pbLume.GetDescription(),
		DateStart:    dateStart,
		DateEnd:      dateEnd,
		Latitude:     latitude,
		Longitude:    longitude,
		Address:      address,
		Images:       pbLume.GetImages(),
		CategoryTags: pbLume.GetCategoryTags(),
		BookingLink:  bookingLink,
	}, nil
}

// toAppUpdateRequest converts a protobuf Lume to an app UpdateLumeRequest
func (s *Service) toAppUpdateRequest(pbLume *lumepb.Lume) (applume.UpdateLumeRequest, error) {
	// Convert timestamps to time.Time pointers
	var dateStart, dateEnd *time.Time
	if pbLume.GetDateStart() != nil {
		t := pbLume.GetDateStart().AsTime()
		dateStart = &t
	}
	if pbLume.GetDateEnd() != nil {
		t := pbLume.GetDateEnd().AsTime()
		dateEnd = &t
	}

	// Convert optional fields
	var latitude, longitude *float64
	var address, bookingLink *string

	if pbLume.GetLatitude() != 0 {
		lat := pbLume.GetLatitude()
		latitude = &lat
	}
	if pbLume.GetLongitude() != 0 {
		lng := pbLume.GetLongitude()
		longitude = &lng
	}
	if pbLume.GetAddress() != "" {
		addr := pbLume.GetAddress()
		address = &addr
	}
	if pbLume.GetBookingLink() != "" {
		link := pbLume.GetBookingLink()
		bookingLink = &link
	}

	return applume.UpdateLumeRequest{
		Label:        pbLume.GetName(),
		Type:         modellume.LumeType(pbLume.GetType().String()),
		Description:  pbLume.GetDescription(),
		DateStart:    dateStart,
		DateEnd:      dateEnd,
		Latitude:     latitude,
		Longitude:    longitude,
		Address:      address,
		Images:       pbLume.GetImages(),
		CategoryTags: pbLume.GetCategoryTags(),
		BookingLink:  bookingLink,
	}, nil
}

// toPbLume converts a domain Lume to a protobuf Lume
//func toPbLume(domainLume *modellume.Lume) (*lumepb.Lume, error) {
//	// Convert timestamps
//	var dateStart, dateEnd *timestamppb.Timestamp
//	if domainLume.DateStart != nil {
//		dateStart = timestamppb.New(*domainLume.DateStart)
//	}
//	if domainLume.DateEnd != nil {
//		dateEnd = timestamppb.New(*domainLume.DateEnd)
//	}
//
//	// Convert optional fields
//	var latitude, longitude float64
//	var address, bookingLink string
//
//	if domainLume.Latitude != nil {
//		latitude = *domainLume.Latitude
//	}
//	if domainLume.Longitude != nil {
//		longitude = *domainLume.Longitude
//	}
//	if domainLume.Address != nil {
//		address = *domainLume.Address
//	}
//	if domainLume.BookingLink != nil {
//		bookingLink = *domainLume.BookingLink
//	}
//
//	return &lumepb.Lume{
//		LumeId:       domainLume.LumeID,
//		LumoId:       domainLume.LumoID,
//		Type:         toPbLumeType(domainLume.Type),
//		Name:         domainLume.Name,
//		DateStart:    dateStart,
//		DateEnd:      dateEnd,
//		Latitude:     latitude,
//		Longitude:    longitude,
//		Address:      address,
//		Description:  domainLume.Description,
//		Images:       domainLume.Images,
//		CategoryTags: domainLume.CategoryTags,
//		BookingLink:  bookingLink,
//		CreatedAt:    timestamppb.New(domainLume.CreatedAt),
//		UpdatedAt:    timestamppb.New(domainLume.UpdatedAt),
//	}, nil
//}
//
//// toPbLumeType converts a domain LumeType to a protobuf LumeType
//func toPbLumeType(domainType modellume.LumeType) lumepb.LumeType {
//	switch domainType {
//	case modellume.LumeTypeUnspecified:
//		return lumepb.LumeType_LUME_TYPE_UNSPECIFIED
//	case modellume.LumeTypeCity:
//		return lumepb.LumeType_LUME_TYPE_CITY
//	case modellume.LumeTypeAttraction:
//		return lumepb.LumeType_LUME_TYPE_ATTRACTION
//	case modellume.LumeTypeAccommodation:
//		return lumepb.LumeType_LUME_TYPE_ACCOMMODATION
//	case modellume.LumeTypeRestaurant:
//		return lumepb.LumeType_LUME_TYPE_RESTAURANT
//	case modellume.LumeTypeTransportHub:
//		return lumepb.LumeType_LUME_TYPE_TRANSPORT_HUB
//	case modellume.LumeTypeActivity:
//		return lumepb.LumeType_LUME_TYPE_ACTIVITY
//	case modellume.LumeTypeShopping:
//		return lumepb.LumeType_LUME_TYPE_SHOPPING
//	case modellume.LumeTypeEntertainment:
//		return lumepb.LumeType_LUME_TYPE_ENTERTAINMENT
//	case modellume.LumeTypeCustom:
//		return lumepb.LumeType_LUME_TYPE_CUSTOM
//	default:
//		log.Printf("warning: toPbLumeType saw unexpected domain type %v, defaulting to UNSPECIFIED", domainType)
//		return lumepb.LumeType_LUME_TYPE_UNSPECIFIED
//	}
//}

// Domain LumeType to Proto LumeType conversion
func DomainLumeTypeToProto(dt modellume.LumeType) lumepb.LumeType {
	switch dt {
	case modellume.LumeTypeCity:
		return lumepb.LumeType_LUME_TYPE_CITY
	case modellume.LumeTypeAttraction:
		return lumepb.LumeType_LUME_TYPE_ATTRACTION
	case modellume.LumeTypeAccommodation:
		return lumepb.LumeType_LUME_TYPE_ACCOMMODATION
	case modellume.LumeTypeRestaurant:
		return lumepb.LumeType_LUME_TYPE_RESTAURANT
	case modellume.LumeTypeTransportHub:
		return lumepb.LumeType_LUME_TYPE_TRANSPORT_HUB
	case modellume.LumeTypeActivity:
		return lumepb.LumeType_LUME_TYPE_ACTIVITY
	case modellume.LumeTypeShopping:
		return lumepb.LumeType_LUME_TYPE_SHOPPING
	case modellume.LumeTypeEntertainment:
		return lumepb.LumeType_LUME_TYPE_ENTERTAINMENT
	case modellume.LumeTypeCustom:
		return lumepb.LumeType_LUME_TYPE_CUSTOM
	default:
		return lumepb.LumeType_LUME_TYPE_UNSPECIFIED
	}
}

// mapErrorToConnectError maps domain errors to Connect errors
func mapErrorToConnectError(err error) error {
	switch {
	case errors.Is(err, applume.ErrLumeNotFound):
		return connect.NewError(connect.CodeNotFound, err)
	case errors.Is(err, applume.ErrInvalidLumoID):
		return connect.NewError(connect.CodeInvalidArgument, err)
	case errors.Is(err, applume.ErrInvalidLumeID):
		return connect.NewError(connect.CodeInvalidArgument, err)
	case errors.Is(err, applume.ErrInvalidLumeType):
		return connect.NewError(connect.CodeInvalidArgument, err)
	case errors.Is(err, applume.ErrEmptyLabel):
		return connect.NewError(connect.CodeInvalidArgument, err)
	case errors.Is(err, applume.ErrInvalidMetadata):
		return connect.NewError(connect.CodeInvalidArgument, err)
	case errors.Is(err, ErrInvalidID):
		return connect.NewError(connect.CodeInvalidArgument, err)
	default:
		return connect.NewError(connect.CodeInternal, err)
	}
}
