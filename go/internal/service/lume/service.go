package lume

import (
	"context"
	"errors"
	"strconv"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	applume "github.com/mcdev12/lumo/go/internal/app/lume"
	pb "github.com/mcdev12/lumo/go/internal/genproto/lume"
	modellume "github.com/mcdev12/lumo/go/internal/models/lume"
)

// Domain errors
var (
	ErrInvalidID = errors.New("invalid ID format")
)

// LumeApp defines what the service layer needs from the app layer
type LumeApp interface {
	CreateLume(ctx context.Context, req applume.CreateLumeRequest) (*modellume.Lume, error)
	GetLumeByID(ctx context.Context, id int64) (*modellume.Lume, error)
	GetLumeByLumeID(ctx context.Context, lumeID string) (*modellume.Lume, error)
	ListLumesByLumoID(ctx context.Context, req applume.ListLumesRequest) ([]*modellume.Lume, error)
	ListLumesByType(ctx context.Context, req applume.ListLumesByTypeRequest) ([]*modellume.Lume, error)
	SearchLumesByLocation(ctx context.Context, req applume.SearchLumesByLocationRequest) ([]*modellume.Lume, error)
	UpdateLume(ctx context.Context, id int64, req applume.UpdateLumeRequest) (*modellume.Lume, error)
	UpdateLumeByLumeID(ctx context.Context, lumeID string, req applume.UpdateLumeRequest) (*modellume.Lume, error)
	DeleteLume(ctx context.Context, id int64) error
	DeleteLumeByLumeID(ctx context.Context, lumeID string) error
	CountLumesByLumo(ctx context.Context, lumoID string) (int64, error)
}

// Service implements the LumeServiceHandler interface
type Service struct {
	app LumeApp
}

// NewService creates a new Lume service
func NewService(app LumeApp) *Service {
	return &Service{
		app: app,
	}
}

// CreateLume creates a new Lume
func (s *Service) CreateLume(ctx context.Context, req *connect.Request[pb.CreateLumeRequest]) (*connect.Response[pb.CreateLumeResponse], error) {
	pbLume := req.Msg.GetLume()
	if pbLume == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("lume is required"))
	}

	appReq, err := s.toAppCreateRequest(pbLume)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	domainLume, err := s.app.CreateLume(ctx, appReq)
	if err != nil {
		return nil, s.mapErrorToConnectError(err)
	}

	pbRespLume, err := s.toPbLume(domainLume)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.CreateLumeResponse{
		Lume: pbRespLume,
	}), nil
}

// GetLume retrieves a Lume by ID
func (s *Service) GetLume(ctx context.Context, req *connect.Request[pb.GetLumeRequest]) (*connect.Response[pb.GetLumeResponse], error) {
	requestID := req.Msg.GetLumeId()

	// Try to parse as int64 first (internal ID), then as UUID (lume_id)
	var domainLume *modellume.Lume
	var err error

	if id, parseErr := strconv.ParseInt(requestID, 10, 64); parseErr == nil {
		// It's an internal ID
		domainLume, err = s.app.GetLumeByID(ctx, id)
	} else {
		// Try as UUID string
		domainLume, err = s.app.GetLumeByLumeID(ctx, requestID)
	}

	if err != nil {
		return nil, s.mapErrorToConnectError(err)
	}

	pbLume, err := s.toPbLume(domainLume)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.GetLumeResponse{
		Lume: pbLume,
	}), nil
}

// ListLumes retrieves all Lumes for a given Lumo
func (s *Service) ListLumes(ctx context.Context, req *connect.Request[pb.ListLumesRequest]) (*connect.Response[pb.ListLumesResponse], error) {
	// Convert page_size to limit and page_token to offset
	limit := req.Msg.GetPageSize()
	if limit <= 0 {
		limit = 50 // Default limit
	}

	offset := int32(0)
	if req.Msg.GetPageToken() != "" {
		parsedOffset, err := strconv.ParseInt(req.Msg.GetPageToken(), 10, 32)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid page token"))
		}
		offset = int32(parsedOffset)
	}

	// Check if we need to filter by type
	lumeType := req.Msg.GetType()
	var domainLumes []*modellume.Lume
	var err error

	if lumeType != pb.LumeType_LUME_TYPE_UNSPECIFIED {
		// Filter by type
		typeReq := applume.ListLumesByTypeRequest{
			LumoID: req.Msg.GetUserId(), // Note: The proto uses user_id but we need lumo_id
			Type:   modellume.LumeType(lumeType.String()),
			Limit:  limit,
			Offset: offset,
		}
		domainLumes, err = s.app.ListLumesByType(ctx, typeReq)
	} else {
		// List all lumes for the lumo
		listReq := applume.ListLumesRequest{
			LumoID: req.Msg.GetUserId(), // Note: The proto uses user_id but we need lumo_id
			Limit:  limit,
			Offset: offset,
		}
		domainLumes, err = s.app.ListLumesByLumoID(ctx, listReq)
	}

	if err != nil {
		return nil, s.mapErrorToConnectError(err)
	}

	pbLumes := make([]*pb.Lume, len(domainLumes))
	for i, domainLume := range domainLumes {
		pbLume, err := s.toPbLume(domainLume)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		pbLumes[i] = pbLume
	}

	// Calculate next page token
	var nextPageToken string
	if len(pbLumes) == int(limit) {
		nextPageToken = strconv.FormatInt(int64(offset+limit), 10)
	}

	return connect.NewResponse(&pb.ListLumesResponse{
		Lumes:         pbLumes,
		NextPageToken: nextPageToken,
	}), nil
}

// UpdateLume updates an existing Lume
func (s *Service) UpdateLume(ctx context.Context, req *connect.Request[pb.UpdateLumeRequest]) (*connect.Response[pb.UpdateLumeResponse], error) {
	pbLume := req.Msg.GetLume()
	if pbLume == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("lume is required"))
	}

	appReq, err := s.toAppUpdateRequest(pbLume)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// Use LumeId (UUID) for updates, not internal ID
	domainLume, err := s.app.UpdateLumeByLumeID(ctx, pbLume.GetLumeId(), appReq)
	if err != nil {
		return nil, s.mapErrorToConnectError(err)
	}

	pbRespLume, err := s.toPbLume(domainLume)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.UpdateLumeResponse{
		Lume: pbRespLume,
	}), nil
}

// DeleteLume deletes a Lume by ID
func (s *Service) DeleteLume(ctx context.Context, req *connect.Request[pb.DeleteLumeRequest]) (*connect.Response[pb.DeleteLumeResponse], error) {
	requestID := req.Msg.GetLumeId()

	// Try to parse as int64 first (internal ID), then as UUID (lume_id)
	var err error

	if id, parseErr := strconv.ParseInt(requestID, 10, 64); parseErr == nil {
		// It's an internal ID
		err = s.app.DeleteLume(ctx, id)
	} else {
		// Try as UUID string
		err = s.app.DeleteLumeByLumeID(ctx, requestID)
	}

	if err != nil {
		return nil, s.mapErrorToConnectError(err)
	}

	return connect.NewResponse(&pb.DeleteLumeResponse{}), nil
}

// Helper methods for conversion and error mapping

// toAppCreateRequest converts a protobuf Lume to an app CreateLumeRequest
func (s *Service) toAppCreateRequest(pbLume *pb.Lume) (applume.CreateLumeRequest, error) {
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
func (s *Service) toAppUpdateRequest(pbLume *pb.Lume) (applume.UpdateLumeRequest, error) {
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
func (s *Service) toPbLume(domainLume *modellume.Lume) (*pb.Lume, error) {
	pbLumeType, err := s.toPbLumeType(domainLume.Type)
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
func (s *Service) toPbLumeType(domainType modellume.LumeType) (pb.LumeType, error) {
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
func (s *Service) mapErrorToConnectError(err error) error {
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
