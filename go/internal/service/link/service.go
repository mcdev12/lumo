package link

import (
	"context"
	"errors"
	"strconv"

	"connectrpc.com/connect"

	applink "github.com/mcdev12/lumo/go/internal/app/link"
	pb "github.com/mcdev12/lumo/go/internal/genproto/link/v1"
	modellink "github.com/mcdev12/lumo/go/internal/models/link"
)

// Domain errors
var (
	ErrInvalidID = errors.New("invalid ID format")
)

// LinkApp defines what the service layer needs from the app layer
type LinkApp interface {
	CreateLink(ctx context.Context, req applink.CreateLinkRequest) (*modellink.Link, error)
	GetLinkByID(ctx context.Context, id int64) (*modellink.Link, error)
	GetLinkByLinkID(ctx context.Context, linkID string) (*modellink.Link, error)
	ListLinksByFromLumeID(ctx context.Context, fromLumeID string, req applink.ListLinksRequest) ([]*modellink.Link, error)
	ListLinksByToLumeID(ctx context.Context, toLumeID string, req applink.ListLinksRequest) ([]*modellink.Link, error)
	ListLinksByEitherLumeID(ctx context.Context, lumeID string, req applink.ListLinksRequest) ([]*modellink.Link, error)
	UpdateLink(ctx context.Context, id int64, req applink.UpdateLinkRequest) (*modellink.Link, error)
	UpdateLinkByLinkID(ctx context.Context, linkID string, req applink.UpdateLinkRequest) (*modellink.Link, error)
	DeleteLink(ctx context.Context, id int64) error
	DeleteLinkByLinkID(ctx context.Context, linkID string) error
	CountLinksByLumeID(ctx context.Context, lumeID string) (int64, error)
}

// Service implements the LinkServiceHandler interface
type Service struct {
	app LinkApp
}

// NewService creates a new Link service
func NewService(app LinkApp) *Service {
	return &Service{
		app: app,
	}
}

// CreateLink creates a new Link
func (s *Service) CreateLink(ctx context.Context, req *connect.Request[pb.CreateLinkRequest]) (*connect.Response[pb.CreateLinkResponse], error) {
	pbLink := req.Msg.GetLink()
	if pbLink == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("link is required"))
	}

	appReq, err := s.toAppCreateRequest(pbLink)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	domainLink, err := s.app.CreateLink(ctx, appReq)
	if err != nil {
		return nil, s.mapErrorToConnectError(err)
	}

	return connect.NewResponse(&pb.CreateLinkResponse{
		Link: modellink.DomainToProto(domainLink),
	}), nil
}

// GetLink retrieves a Link by ID
func (s *Service) GetLink(ctx context.Context, req *connect.Request[pb.GetLinkRequest]) (*connect.Response[pb.GetLinkResponse], error) {
	linkID := req.Msg.GetLinkId()
	if linkID == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("link_id is required"))
	}

	// Try to parse as int64 first (internal ID), then as UUID (link_id)
	var domainLink *modellink.Link
	var err error

	if id, parseErr := strconv.ParseInt(linkID, 10, 64); parseErr == nil {
		// It's an internal ID
		domainLink, err = s.app.GetLinkByID(ctx, id)
	} else {
		// Try as UUID string
		domainLink, err = s.app.GetLinkByLinkID(ctx, linkID)
	}

	if err != nil {
		return nil, s.mapErrorToConnectError(err)
	}

	return connect.NewResponse(&pb.GetLinkResponse{
		Link: modellink.DomainToProto(domainLink),
	}), nil
}

// UpdateLink updates an existing Link
func (s *Service) UpdateLink(ctx context.Context, req *connect.Request[pb.UpdateLinkRequest]) (*connect.Response[pb.UpdateLinkResponse], error) {
	pbLink := req.Msg.GetLink()
	if pbLink == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("link is required"))
	}

	if pbLink.GetLinkId() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("link_id is required"))
	}

	appReq, err := s.toAppUpdateRequest(pbLink)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// Use LinkId (UUID) for updates, not internal ID
	domainLink, err := s.app.UpdateLinkByLinkID(ctx, pbLink.GetLinkId(), appReq)
	if err != nil {
		return nil, s.mapErrorToConnectError(err)
	}

	resp := pb.UpdateLinkResponse{
		Link: modellink.DomainToProto(domainLink),
	}
	return connect.NewResponse(&resp), nil
}

// DeleteLink deletes a Link by ID
func (s *Service) DeleteLink(ctx context.Context, req *connect.Request[pb.DeleteLinkRequest]) (*connect.Response[pb.DeleteLinkResponse], error) {
	linkID := req.Msg.GetLinkId()
	if linkID == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("link_id is required"))
	}

	// Try to parse as int64 first (internal ID), then as UUID (link_id)
	var err error

	if id, parseErr := strconv.ParseInt(linkID, 10, 64); parseErr == nil {
		// It's an internal ID
		err = s.app.DeleteLink(ctx, id)
	} else {
		// Try as UUID string
		err = s.app.DeleteLinkByLinkID(ctx, linkID)
	}

	if err != nil {
		return nil, s.mapErrorToConnectError(err)
	}

	return connect.NewResponse(&pb.DeleteLinkResponse{}), nil
}

// ListLinks retrieves links with optional filtering and pagination
func (s *Service) ListLinks(ctx context.Context, req *connect.Request[pb.ListLinksRequest]) (*connect.Response[pb.ListLinksResponse], error) {
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

	// Prepare app request
	appReq := applink.ListLinksRequest{
		Limit:  limit,
		Offset: offset,
	}

	// Determine which listing method to use based on the filters
	var domainLinks []*modellink.Link
	var err error

	fromLumeID := req.Msg.GetFromLumeId()
	toLumeID := req.Msg.GetToLumeId()
	lumoUUID := req.Msg.GetLumoUuid()

	if fromLumeID != "" && toLumeID != "" {
		// This is a special case not directly supported by the app layer
		// We'll need to filter the results after fetching them
		domainLinks, err = s.app.ListLinksByEitherLumeID(ctx, fromLumeID, appReq)
		if err != nil {
			return nil, s.mapErrorToConnectError(err)
		}

		// Filter for links that match both fromLumeID and toLumeID
		filteredLinks := make([]*modellink.Link, 0)
		for _, link := range domainLinks {
			if link.FromLumeID == fromLumeID && link.ToLumeID == toLumeID {
				filteredLinks = append(filteredLinks, link)
			}
		}
		domainLinks = filteredLinks
	} else if fromLumeID != "" {
		domainLinks, err = s.app.ListLinksByFromLumeID(ctx, fromLumeID, appReq)
	} else if toLumeID != "" {
		domainLinks, err = s.app.ListLinksByToLumeID(ctx, toLumeID, appReq)
	} else if lumoUUID != "" {
		// This would require a different query that's not directly supported
		// For now, we'll just return an error
		return nil, connect.NewError(connect.CodeUnimplemented, errors.New("filtering by lumo_uuid is not implemented"))
	} else {
		// No filters, return an error as we don't want to return all links
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("at least one filter is required"))
	}

	if err != nil {
		return nil, s.mapErrorToConnectError(err)
	}

	// Convert domain links to protobuf links
	pbLinks := make([]*pb.Link, len(domainLinks))
	for i, domainLink := range domainLinks {
		pbLinks[i] = modellink.DomainToProto(domainLink)
	}

	// Calculate next page token
	var nextPageToken string
	if len(pbLinks) == int(limit) {
		nextPageToken = strconv.FormatInt(int64(offset+limit), 10)
	}

	return connect.NewResponse(&pb.ListLinksResponse{
		Links:         pbLinks,
		NextPageToken: nextPageToken,
	}), nil
}

// Helper methods for conversion and error mapping

// toAppCreateRequest converts a protobuf Link to an app CreateLinkRequest
func (s *Service) toAppCreateRequest(pbLink *pb.Link) (applink.CreateLinkRequest, error) {
	// Convert optional fields
	var notes *string
	var sequenceIndex *int32

	if pbLink.GetNotes() != "" {
		notesVal := pbLink.GetNotes()
		notes = &notesVal
	}

	if pbLink.GetSequenceIndex() != 0 {
		seqIndex := pbLink.GetSequenceIndex()
		sequenceIndex = &seqIndex
	}

	// Convert travel details if present
	var travelDetails *applink.TravelDetailsRequest
	if pbLink.GetTravel() != nil {
		travelDetails = &applink.TravelDetailsRequest{
			Mode:           modellink.ProtoTravelModeToDomain(pbLink.GetTravel().GetMode()),
			DurationSec:    pbLink.GetTravel().GetDurationSec(),
			CostEstimate:   pbLink.GetTravel().GetCostEstimate(),
			DistanceMeters: pbLink.GetTravel().GetDistanceMeters(),
		}
	}

	return applink.CreateLinkRequest{
		FromLumeID:    pbLink.GetFromLumeId(),
		ToLumeID:      pbLink.GetToLumeId(),
		Type:          modellink.ProtoLinkTypeToDomain(pbLink.GetType()),
		Notes:         notes,
		SequenceIndex: sequenceIndex,
		TravelDetails: travelDetails,
	}, nil
}

// toAppUpdateRequest converts a protobuf Link to an app UpdateLinkRequest
func (s *Service) toAppUpdateRequest(pbLink *pb.Link) (applink.UpdateLinkRequest, error) {
	// Convert optional fields
	var notes *string
	var sequenceIndex *int32

	if pbLink.GetNotes() != "" {
		notesVal := pbLink.GetNotes()
		notes = &notesVal
	}

	if pbLink.GetSequenceIndex() != 0 {
		seqIndex := pbLink.GetSequenceIndex()
		sequenceIndex = &seqIndex
	}

	// Convert travel details if present
	var travelDetails *applink.TravelDetailsRequest
	if pbLink.GetTravel() != nil {
		travelDetails = &applink.TravelDetailsRequest{
			Mode:           modellink.ProtoTravelModeToDomain(pbLink.GetTravel().GetMode()),
			DurationSec:    pbLink.GetTravel().GetDurationSec(),
			CostEstimate:   pbLink.GetTravel().GetCostEstimate(),
			DistanceMeters: pbLink.GetTravel().GetDistanceMeters(),
		}
	}

	return applink.UpdateLinkRequest{
		FromLumeID:    pbLink.GetFromLumeId(),
		ToLumeID:      pbLink.GetToLumeId(),
		Type:          modellink.ProtoLinkTypeToDomain(pbLink.GetType()),
		Notes:         notes,
		SequenceIndex: sequenceIndex,
		TravelDetails: travelDetails,
	}, nil
}

// mapErrorToConnectError maps domain errors to Connect errors
func (s *Service) mapErrorToConnectError(err error) error {
	switch {
	case errors.Is(err, applink.ErrLinkNotFound):
		return connect.NewError(connect.CodeNotFound, err)
	case errors.Is(err, applink.ErrInvalidLinkID):
		return connect.NewError(connect.CodeInvalidArgument, err)
	case errors.Is(err, applink.ErrInvalidLumeID):
		return connect.NewError(connect.CodeInvalidArgument, err)
	case errors.Is(err, applink.ErrInvalidLinkType):
		return connect.NewError(connect.CodeInvalidArgument, err)
	case errors.Is(err, applink.ErrInvalidTravelMode):
		return connect.NewError(connect.CodeInvalidArgument, err)
	case errors.Is(err, applink.ErrEmptyNotes):
		return connect.NewError(connect.CodeInvalidArgument, err)
	case errors.Is(err, ErrInvalidID):
		return connect.NewError(connect.CodeInvalidArgument, err)
	default:
		return connect.NewError(connect.CodeInternal, err)
	}
}
