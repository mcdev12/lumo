package lume

import (
	"github.com/google/uuid"
	pblume "github.com/mcdev12/lumo/go/internal/genproto/lume"
	modellume "github.com/mcdev12/lumo/go/internal/models/lume"
)

/ Validation methods
func validateCreateRequest(req CreateLumeRequest) error {
	if req.Label == "" {
		return ErrEmptyLabel
	}

	if _, err := uuid.Parse(req.LumoID); err != nil {
		return ErrInvalidLumoID
	}

	if !a.isValidLumeType(req.Type) {
		return ErrInvalidLumeType
	}

	return nil
}

func validateUpdateRequest(req *pblume.UpdateLumeRequest) error {
	lume := req.Lume
	if !isValidLumeType(req.Le) {
		return ErrInvalidLumeType
	}

	return nil
}

func isValidLumeType(lumeType modellume.LumeType) bool {
	switch lumeType {
	case modellume.LumeTypeUnspecified,
		modellume.LumeTypeCity,
		modellume.LumeTypeAttraction,
		modellume.LumeTypeAccommodation,
		modellume.LumeTypeRestaurant,
		modellume.LumeTypeTransportHub,
		modellume.LumeTypeActivity,
		modellume.LumeTypeShopping,
		modellume.LumeTypeEntertainment,
		modellume.LumeTypeCustom:
		return true
	default:
		return false
	}
}
