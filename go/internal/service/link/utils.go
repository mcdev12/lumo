package link

import (
	"errors"

	"connectrpc.com/connect"

	applink "github.com/mcdev12/lumo/go/internal/app/link"
	pb "github.com/mcdev12/lumo/go/internal/genproto/link/v1"
	modellink "github.com/mcdev12/lumo/go/internal/models/link"
)

// toAppCreateRequest converts a protobuf Link to an app CreateLinkRequest
func (s *Service) toAppCreateRequest(pbLink *pb.CreateLinkRequest) (applink.CreateLinkRequest, error) {
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
func (s *Service) toAppUpdateRequest(pbLink *pb.UpdateLinkRequest) (applink.UpdateLinkRequest, error) {
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

	updateFields := make([]string, 0)
	if pbLink.GetUpdateMask() != nil {
		updateFields = pbLink.GetUpdateMask().GetPaths()
	}

	return applink.UpdateLinkRequest{
		FromLumeID:    pbLink.GetFromLumeId(),
		ToLumeID:      pbLink.GetToLumeId(),
		Type:          modellink.ProtoLinkTypeToDomain(pbLink.GetType()),
		Notes:         notes,
		SequenceIndex: sequenceIndex,
		TravelDetails: travelDetails,
		UpdateFields:  updateFields,
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
