package lume

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"strconv"

	applume "github.com/mcdev12/lumo/go/internal/app/lume"
	pb "github.com/mcdev12/lumo/go/internal/genproto/lume/v1"
	modellume "github.com/mcdev12/lumo/go/internal/models/lume"
)

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
	pbRequest := req.Msg
	if pbRequest == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("request is empty"))
	}

	appReq, err := s.toAppCreateRequest(pbRequest)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	domainLume, err := s.app.CreateLume(ctx, appReq)
	if err != nil {
		return nil, mapErrorToConnectError(err)
	}

	return connect.NewResponse(&pb.CreateLumeResponse{
		Lume: domainToProto(domainLume),
	}), nil
}

// GetLume retrieves a Lume by ID
func (s *Service) GetLume(ctx context.Context, req *connect.Request[pb.GetLumeRequest]) (*connect.Response[pb.GetLumeResponse], error) {
	requestID := req.Msg.GetLumeId()

	var domainLume *modellume.Lume
	var err error
	// Try to parse as int64 first (internal ID), then as UUID (lume_id)
	if id, parseErr := strconv.ParseInt(requestID, 10, 64); parseErr == nil {
		// It's an internal ID
		domainLume, err = s.app.GetLumeByID(ctx, id)
	} else {
		// Try as UUID string
		domainLume, err = s.app.GetLumeByLumeID(ctx, requestID)
	}

	if err != nil {
		return nil, mapErrorToConnectError(err)
	}

	return connect.NewResponse(&pb.GetLumeResponse{
		Lume: domainToProto(domainLume),
	}), nil
}

// ListLumes retrieves all Lumes for a given Lumo
// TODO probably wrong need to fix
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
		return nil, mapErrorToConnectError(err)
	}

	pbLumes := make([]*pb.Lume, len(domainLumes))
	for i, domainLume := range domainLumes {
		pbLumes[i] = domainToProto(domainLume)

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
	pbUpdateLumeRequest := req.Msg
	if pbUpdateLumeRequest == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("lume is required"))
	}

	appReq, err := s.toAppUpdateRequest(pbUpdateLumeRequest)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// Use LumeId (UUID) for updates, not internal ID
	domainLume, err := s.app.UpdateLumeByLumeID(ctx, pbUpdateLumeRequest.GetLumeId(), appReq)
	if err != nil {
		return nil, mapErrorToConnectError(err)
	}

	return connect.NewResponse(&pb.UpdateLumeResponse{
		Lume: domainToProto(domainLume),
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
		return nil, mapErrorToConnectError(err)
	}

	return connect.NewResponse(&pb.DeleteLumeResponse{}), nil
}
