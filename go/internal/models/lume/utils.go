package lume

import (
	lumepb "github.com/mcdev12/lumo/go/internal/genproto/protobuf/lume"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProtoToDomainLumeType converts the generated Protobuf enum into the domain LumeType.
func ProtoToDomainLumeType(pbType lumepb.LumeType) LumeType {
	switch pbType {
	case lumepb.LumeType_LUME_TYPE_STOP:
		return LumeTypeStop
	case lumepb.LumeType_LUME_TYPE_ACCOMMODATION:
		return LumeTypeAccommodation
	case lumepb.LumeType_LUME_TYPE_POINT_OF_INTEREST:
		return LumeTypePointOfInterest
	case lumepb.LumeType_LUME_TYPE_MEAL:
		return LumeTypeMeal
	case lumepb.LumeType_LUME_TYPE_TRANSPORT:
		return LumeTypeTransport
	case lumepb.LumeType_LUME_TYPE_CUSTOM:
		return LumeTypeCustom
	default:
		return LumeTypeUnspecified
	}
}

// DomainToProtoLumeType converts the domain LumeType into the Protobuf enum.
func DomainToProtoLumeType(domType LumeType) lumepb.LumeType {
	switch domType {
	case LumeTypeStop:
		return lumepb.LumeType_LUME_TYPE_STOP
	case LumeTypeAccommodation:
		return lumepb.LumeType_LUME_TYPE_ACCOMMODATION
	case LumeTypePointOfInterest:
		return lumepb.LumeType_LUME_TYPE_POINT_OF_INTEREST
	case LumeTypeMeal:
		return lumepb.LumeType_LUME_TYPE_MEAL
	case LumeTypeTransport:
		return lumepb.LumeType_LUME_TYPE_TRANSPORT
	case LumeTypeCustom:
		return lumepb.LumeType_LUME_TYPE_CUSTOM
	default:
		return lumepb.LumeType_LUME_TYPE_UNSPECIFIED
	}
}

// structToMap converts a *structpb.Struct into a Go map[string]interface{}.
func structToMap(s *structpb.Struct) map[string]interface{} {
	if s == nil {
		return nil
	}
	return s.AsMap()
}

// mapToStruct converts a Go map[string]interface{} into a *structpb.Struct.
// Returns nil if the input map is nil or empty.
func mapToStruct(m map[string]interface{}) (*structpb.Struct, error) {
	if len(m) == 0 {
		return nil, nil
	}
	return structpb.NewStruct(m)
}

// ProtoToDomainLume converts the Protobuf Lume message into the domain‐level models.Lume.
func ProtoToDomainLume(pb *lumepb.Lume) (*Lume, error) {
	if pb == nil {
		return nil, nil
	}

	dom := &Lume{
		LumeId:      pb.GetLumeId(),
		LumoID:      pb.GetLumoId(),
		Label:       pb.GetLabel(),
		Type:        ProtoToDomainLumeType(pb.GetType()),
		Description: pb.GetDescription(),
		Metadata:    structToMap(pb.GetMetadata()),
		CreatedAt:   pb.GetCreatedAt().AsTime(),
		UpdatedAt:   pb.GetUpdatedAt().AsTime(),
	}
	return dom, nil
}

// DomainToProtoLume converts a domain models.Lume into the Protobuf Lume message.
func DomainToProtoLume(dom *Lume) (*lumepb.Lume, error) {
	if dom == nil {
		return nil, nil
	}

	pbMetadata, err := mapToStruct(dom.Metadata)
	if err != nil {
		return nil, err
	}

	pb := &lumepb.Lume{
		LumeId:      dom.LumeId,
		LumoId:      dom.LumoID,
		Label:       dom.Label,
		Type:        DomainToProtoLumeType(dom.Type),
		Description: dom.Description,
		Metadata:    pbMetadata,
		CreatedAt:   timestamppb.New(dom.CreatedAt),
		UpdatedAt:   timestamppb.New(dom.UpdatedAt),
	}
	return pb, nil
}
