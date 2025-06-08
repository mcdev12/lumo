package lume

import (
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	applume "github.com/mcdev12/lumo/go/internal/app/lume"
	pb "github.com/mcdev12/lumo/go/internal/genproto/lume"
	modellume "github.com/mcdev12/lumo/go/internal/models/lume"
)

// toPbLume converts a domain Lume to a protobuf Lume
func toPbLume(domainLume *modellume.Lume) (*pb.Lume, error) {
	pbLumeType, err := toPbLumeType(domainLume.Type)
	if err != nil {
		return nil, err
	}

	// Convert timestamps
	var dateStart, dateEnd *timestamppb.Timestamp
	if domainLume.DateStart != nil {
		dateStart = timestamppb.New(*domainLume.DateStart)
	}
	if domainLume.DateEnd != nil {
		dateEnd = timestamppb.New(*domainLume.DateEnd)
	}

	// Convert optional fields
	var latitude, longitude float64
	var address, bookingLink string

	if domainLume.Latitude != nil {
		latitude = *domainLume.Latitude
	}
	if domainLume.Longitude != nil {
		longitude = *domainLume.Longitude
	}
	if domainLume.Address != nil {
		address = *domainLume.Address
	}
	if domainLume.BookingLink != nil {
		bookingLink = *domainLume.BookingLink
	}

	return &pb.Lume{
		LumeId:       domainLume.LumeID,
		LumoId:       domainLume.LumoID,
		Type:         pbLumeType,
		Name:         domainLume.Name,
		DateStart:    dateStart,
		DateEnd:      dateEnd,
		Latitude:     latitude,
		Longitude:    longitude,
		Address:      address,
		Description:  domainLume.Description,
		Images:       domainLume.Images,
		CategoryTags: domainLume.CategoryTags,
		BookingLink:  bookingLink,
		CreatedAt:    timestamppb.New(domainLume.CreatedAt),
		UpdatedAt:    timestamppb.New(domainLume.UpdatedAt),
	}, nil
}

// toPbLumeType converts a domain LumeType to a protobuf LumeType
func toPbLumeType(domainType modellume.LumeType) (pb.LumeType, error) {
	switch domainType {
	case modellume.LumeTypeUnspecified:
		return pb.LumeType_LUME_TYPE_UNSPECIFIED, nil
	case modellume.LumeTypeCity:
		return pb.LumeType_LUME_TYPE_CITY, nil
	case modellume.LumeTypeAttraction:
		return pb.LumeType_LUME_TYPE_ATTRACTION, nil
	case modellume.LumeTypeAccommodation:
		return pb.LumeType_LUME_TYPE_ACCOMMODATION, nil
	case modellume.LumeTypeRestaurant:
		return pb.LumeType_LUME_TYPE_RESTAURANT, nil
	case modellume.LumeTypeTransportHub:
		return pb.LumeType_LUME_TYPE_TRANSPORT_HUB, nil
	case modellume.LumeTypeActivity:
		return pb.LumeType_LUME_TYPE_ACTIVITY, nil
	case modellume.LumeTypeShopping:
		return pb.LumeType_LUME_TYPE_SHOPPING, nil
	case modellume.LumeTypeEntertainment:
		return pb.LumeType_LUME_TYPE_ENTERTAINMENT, nil
	case modellume.LumeTypeCustom:
		return pb.LumeType_LUME_TYPE_CUSTOM, nil
	default:
		return pb.LumeType_LUME_TYPE_UNSPECIFIED, errors.New("unknown lume type")
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
