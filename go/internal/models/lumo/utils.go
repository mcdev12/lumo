package lumo

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	lumopb "github.com/mcdev12/lumo/go/internal/genproto/lumo"
)

// DomainToProto converts domain Lumo to protobuf Lumo
func DomainToProto(domainLumo *Lumo) *lumopb.Lumo {
	proto := &lumopb.Lumo{
		LumoId:    domainLumo.LumoID,
		UserId:    domainLumo.UserID,
		Title:     domainLumo.Title,
		CreatedAt: timestamppb.New(domainLumo.CreatedAt),
		UpdatedAt: timestamppb.New(domainLumo.UpdatedAt),
	}

	return proto
}

// ProtoToDomain converts protobuf Lumo to domain Lumo
func ProtoToDomain(protoLumo *lumopb.Lumo) *Lumo {
	domain := &Lumo{
		LumoID:    protoLumo.LumoId,
		UserID:    protoLumo.UserId,
		Title:     protoLumo.Title,
		CreatedAt: protoLumo.CreatedAt.AsTime(),
		UpdatedAt: protoLumo.UpdatedAt.AsTime(),
	}

	return domain
}

// ValidateTimestamps ensures that timestamps are valid
func ValidateTimestamps(created, updated time.Time) bool {
	// Ensure timestamps are not zero values
	if created.IsZero() || updated.IsZero() {
		return false
	}

	// Ensure updated is not before created
	return !updated.Before(created)
}

// EnsureStringArray ensures empty arrays instead of nil for consistency
func EnsureStringArray(arr []string) []string {
	if arr == nil {
		return make([]string, 0)
	}
	return arr
}
