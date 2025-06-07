package lume

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	modellume "github.com/mcdev12/lumo/go/internal/models/lume"
	lumerepo "github.com/mcdev12/lumo/go/internal/repository/lume"
)

// Domain errors
var (
	ErrLumeNotFound    = errors.New("lume not found")
	ErrInvalidLumoID   = errors.New("invalid lumo ID")
	ErrInvalidLumeID   = errors.New("invalid lume ID")
	ErrInvalidLumeType = errors.New("invalid lume type")
	ErrEmptyLabel      = errors.New("label cannot be empty")
	ErrInvalidMetadata = errors.New("invalid metadata")
)

// LumeRepository defines what the app layer needs from the repository
type LumeRepository interface {
	CreateLume(ctx context.Context, req lumerepo.CreateLumeRequest) (*modellume.Lume, error)
	GetLumeByID(ctx context.Context, id int64) (*modellume.Lume, error)
	GetLumeByLumeID(ctx context.Context, lumeID uuid.UUID) (*modellume.Lume, error)
	ListLumesByLumoID(ctx context.Context, lumoID uuid.UUID, limit, offset int32) ([]*modellume.Lume, error)
	UpdateLume(ctx context.Context, id int64, req lumerepo.UpdateLumeRequest) (*modellume.Lume, error)
	DeleteLume(ctx context.Context, id int64) error
}

// CreateLumeRequest represents the business layer's create request
type CreateLumeRequest struct {
	LumoID      string // Parent Lumo UUID (required)
	Label       string
	Type        modellume.LumeType
	Description string
	Metadata    map[string]interface{}
}

// UpdateLumeRequest represents the business layer's update request
type UpdateLumeRequest struct {
	Label       string
	Type        modellume.LumeType
	Description string
	Metadata    map[string]interface{}
}

// ListLumesRequest represents pagination parameters
type ListLumesRequest struct {
	LumoID string
	Limit  int32
	Offset int32
}

// App handles business logic for Lumes
type App struct {
	repo LumeRepository
}

// NewLumeApp creates a new Lume Service
func NewLumeApp(repo LumeRepository) *App {
	return &App{
		repo: repo,
	}
}

// CreateLume creates a new Lume with business logic validation
func (a *App) CreateLume(ctx context.Context, req CreateLumeRequest) (*modellume.Lume, error) {
	if err := a.validateCreateRequest(req); err != nil {
		return nil, err
	}

	repoReq, err := a.toRepositoryCreateRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to convert request: %w", err)
	}

	return a.repo.CreateLume(ctx, repoReq)
}

// GetLumeByID retrieves a Lume by its internal ID
func (a *App) GetLumeByID(ctx context.Context, id int64) (*modellume.Lume, error) {
	return a.repo.GetLumeByID(ctx, id)
}

// GetLumeByLumeID retrieves a Lume by its UUID string
func (a *App) GetLumeByLumeID(ctx context.Context, lumeID string) (*modellume.Lume, error) {
	lumeUUID, err := uuid.Parse(lumeID)
	if err != nil {
		return nil, ErrInvalidLumeID
	}

	return a.repo.GetLumeByLumeID(ctx, lumeUUID)
}

// ListLumesByLumoID retrieves all Lumes for a given Lumo
func (a *App) ListLumesByLumoID(ctx context.Context, req ListLumesRequest) ([]*modellume.Lume, error) {
	lumoUUID, err := uuid.Parse(req.LumoID)
	if err != nil {
		return nil, ErrInvalidLumoID
	}

	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	return a.repo.ListLumesByLumoID(ctx, lumoUUID, limit, offset)
}

// UpdateLume updates an existing Lume
func (a *App) UpdateLume(ctx context.Context, id int64, req UpdateLumeRequest) (*modellume.Lume, error) {
	if err := a.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	repoReq, err := a.toRepositoryUpdateRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to convert request: %w", err)
	}

	return a.repo.UpdateLume(ctx, id, repoReq)
}

// UpdateLumeByLumeID updates a Lume by its UUID
func (a *App) UpdateLumeByLumeID(ctx context.Context, lumeID string, req UpdateLumeRequest) (*modellume.Lume, error) {
	// First get the lume to find its internal ID
	lume, err := a.GetLumeByLumeID(ctx, lumeID)
	if err != nil {
		return nil, err
	}

	return a.UpdateLume(ctx, lume.Id, req)
}

// DeleteLume deletes a Lume by its ID
func (a *App) DeleteLume(ctx context.Context, id int64) error {
	return a.repo.DeleteLume(ctx, id)
}

// DeleteLumeByLumeID deletes a Lume by its UUID
func (a *App) DeleteLumeByLumeID(ctx context.Context, lumeID string) error {
	// First get the lume to find its internal ID
	lume, err := a.GetLumeByLumeID(ctx, lumeID)
	if err != nil {
		return err
	}

	return a.DeleteLume(ctx, lume.Id)
}

// Validation methods
func (a *App) validateCreateRequest(req CreateLumeRequest) error {
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

func (a *App) validateUpdateRequest(req UpdateLumeRequest) error {
	if req.Label == "" {
		return ErrEmptyLabel
	}

	if !a.isValidLumeType(req.Type) {
		return ErrInvalidLumeType
	}

	return nil
}

func (a *App) isValidLumeType(lumeType modellume.LumeType) bool {
	switch lumeType {
	case modellume.LumeTypeUnspecified,
		modellume.LumeTypeStop,
		modellume.LumeTypeAccommodation,
		modellume.LumeTypePointOfInterest,
		modellume.LumeTypeMeal,
		modellume.LumeTypeTransport,
		modellume.LumeTypeCustom:
		return true
	default:
		return false
	}
}

// Conversion methods
// toRepositoryCreateRequest generates UUID in Go
func (a *App) toRepositoryCreateRequest(req CreateLumeRequest) (lumerepo.CreateLumeRequest, error) {
	lumoUUID, err := uuid.Parse(req.LumoID)
	if err != nil {
		return lumerepo.CreateLumeRequest{}, ErrInvalidLumoID
	}

	// Generate UUID in application layer
	lumeUUID := uuid.New()

	var description *string
	if req.Description != "" {
		description = &req.Description
	}

	return lumerepo.CreateLumeRequest{
		LumeID:      lumeUUID, // Generated by Go!
		LumoID:      lumoUUID, // Parent Lumo UUID
		Label:       req.Label,
		Type:        req.Type,
		Description: description,
		Metadata:    req.Metadata,
	}, nil
}

func (a *App) toRepositoryUpdateRequest(req UpdateLumeRequest) (lumerepo.UpdateLumeRequest, error) {
	var description *string
	if req.Description != "" {
		description = &req.Description
	}

	return lumerepo.UpdateLumeRequest{
		Label:       req.Label,
		Type:        req.Type,
		Description: description,
		Metadata:    req.Metadata,
	}, nil
}
