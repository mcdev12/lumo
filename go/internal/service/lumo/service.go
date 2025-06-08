package lumo

import (
	"context"
	"errors"
	"strconv"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"

	applumo "github.com/mcdev12/lumo/go/internal/app/lumo"
	pb "github.com/mcdev12/lumo/go/internal/genproto/lumo"
	modellumo "github.com/mcdev12/lumo/go/internal/models/lumo"
)

// Domain errors
var (
	ErrInvalidID = errors.New("invalid ID format")
)

// LumoApp defines what the service layer needs from the app layer
type LumoApp interface {
	CreateLumo(ctx context.Context, req applumo.CreateLumoRequest) (*modellumo.Lumo, error)
	GetLumoByID(ctx context.Context, id int64) (*modellumo.Lumo, error)
	GetLumoByLumoID(ctx context.Context, lumoID string) (*modellumo.Lumo, error)
	ListLumosByUserID(ctx context.Context, req applumo.ListLumosRequest) ([]*modellumo.Lumo, error)
	UpdateLumo(ctx context.Context, id int64, req applumo.UpdateLumoRequest) (*modellumo.Lumo, error)
	UpdateLumoByLumoID(ctx context.Context, lumoID string, req applumo.UpdateLumoRequest) (*modellumo.Lumo, error)
	DeleteLumo(ctx context.Context, id int64) error
	DeleteLumoByLumoID(ctx context.Context, lumoID string) error
	CountLumosByUserID(ctx context.Context, userID string) (int64, error)
}

// Service implements the LumoServiceHandler interface
type Service struct {
	app LumoApp
}

// NewService creates a new Lumo service
func NewService(app LumoApp) *Service {
	return &Service{
		app: app,
	}
}

// CreateLumo creates a new Lumo
func (s *Service) CreateLumo(ctx context.Context, req *connect.Request[pb.CreateLumoRequest]) (*connect.Response[pb.CreateLumoResponse], error) {
	pbLumo := req.Msg.GetLumo()
	if pbLumo == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("lumo is required"))
	}

	appReq := applumo.CreateLumoRequest{
		UserID: pbLumo.GetUserId(),
		Title:  pbLumo.GetTitle(),
	}

	domainLumo, err := s.app.CreateLumo(ctx, appReq)
	if err != nil {
		return nil, s.mapErrorToConnectError(err)
	}

	pbRespLumo := modellumo.DomainToProto(domainLumo)

	return connect.NewResponse(&pb.CreateLumoResponse{
		Lumo: pbRespLumo,
	}), nil
}

// GetLumo retrieves a Lumo by ID
func (s *Service) GetLumo(ctx context.Context, req *connect.Request[pb.GetLumoRequest]) (*connect.Response[pb.GetLumoResponse], error) {
	requestID := req.Msg.GetUuid()

	// Try to parse as int64 first (internal ID), then as UUID (lumo_id)
	var domainLumo *modellumo.Lumo
	var err error

	if id, parseErr := strconv.ParseInt(requestID, 10, 64); parseErr == nil {
		// It's an internal ID
		domainLumo, err = s.app.GetLumoByID(ctx, id)
	} else {
		// Try as UUID string
		domainLumo, err = s.app.GetLumoByLumoID(ctx, requestID)
	}

	if err != nil {
		return nil, s.mapErrorToConnectError(err)
	}

	pbLumo := modellumo.DomainToProto(domainLumo)

	return connect.NewResponse(&pb.GetLumoResponse{
		Lumo: pbLumo,
	}), nil
}

// ListLumos retrieves all Lumos for a given user
func (s *Service) ListLumos(ctx context.Context, req *connect.Request[pb.ListLumosRequest]) (*connect.Response[pb.ListLumosResponse], error) {
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

	// List all lumos for the user
	listReq := applumo.ListLumosRequest{
		UserID: req.Msg.GetUserId(),
		Limit:  limit,
		Offset: offset,
	}
	domainLumos, err := s.app.ListLumosByUserID(ctx, listReq)
	if err != nil {
		return nil, s.mapErrorToConnectError(err)
	}

	pbLumos := make([]*pb.Lumo, len(domainLumos))
	for i, domainLumo := range domainLumos {
		pbLumos[i] = modellumo.DomainToProto(domainLumo)
	}

	// Calculate next page token
	var nextPageToken string
	if len(pbLumos) == int(limit) {
		nextPageToken = strconv.FormatInt(int64(offset+limit), 10)
	}

	return connect.NewResponse(&pb.ListLumosResponse{
		Lumos:         pbLumos,
		NextPageToken: nextPageToken,
	}), nil
}

// UpdateLumo updates an existing Lumo
func (s *Service) UpdateLumo(ctx context.Context, req *connect.Request[pb.UpdateLumoRequest]) (*connect.Response[pb.UpdateLumoResponse], error) {
	pbLumo := req.Msg.GetLumo()
	if pbLumo == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("lumo is required"))
	}

	appReq := applumo.UpdateLumoRequest{
		Title: pbLumo.GetTitle(),
	}

	// Use LumoId (UUID) for updates, not internal ID
	domainLumo, err := s.app.UpdateLumoByLumoID(ctx, pbLumo.GetLumoId(), appReq)
	if err != nil {
		return nil, s.mapErrorToConnectError(err)
	}

	pbRespLumo := modellumo.DomainToProto(domainLumo)

	return connect.NewResponse(&pb.UpdateLumoResponse{
		Lumo: pbRespLumo,
	}), nil
}

// DeleteLumo deletes a Lumo by ID
func (s *Service) DeleteLumo(ctx context.Context, req *connect.Request[pb.DeleteLumoRequest]) (*connect.Response[emptypb.Empty], error) {
	requestID := req.Msg.GetUuid()

	// Try to parse as int64 first (internal ID), then as UUID (lumo_id)
	var err error

	if id, parseErr := strconv.ParseInt(requestID, 10, 64); parseErr == nil {
		// It's an internal ID
		err = s.app.DeleteLumo(ctx, id)
	} else {
		// Try as UUID string
		err = s.app.DeleteLumoByLumoID(ctx, requestID)
	}

	if err != nil {
		return nil, s.mapErrorToConnectError(err)
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

// mapErrorToConnectError maps domain errors to Connect errors
func (s *Service) mapErrorToConnectError(err error) error {
	switch {
	case errors.Is(err, applumo.ErrLumoNotFound):
		return connect.NewError(connect.CodeNotFound, err)
	case errors.Is(err, applumo.ErrInvalidUserID):
		return connect.NewError(connect.CodeInvalidArgument, err)
	case errors.Is(err, applumo.ErrInvalidLumoID):
		return connect.NewError(connect.CodeInvalidArgument, err)
	case errors.Is(err, applumo.ErrEmptyTitle):
		return connect.NewError(connect.CodeInvalidArgument, err)
	case errors.Is(err, ErrInvalidID):
		return connect.NewError(connect.CodeInvalidArgument, err)
	default:
		return connect.NewError(connect.CodeInternal, err)
	}
}
