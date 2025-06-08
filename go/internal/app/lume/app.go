package lume

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	modellume "github.com/mcdev12/lumo/go/internal/models/lume"
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
	CreateLume(ctx context.Context, domainLume *modellume.Lume) (*modellume.Lume, error)
	GetLumeByID(ctx context.Context, id int64) (*modellume.Lume, error)
	GetLumeByLumeID(ctx context.Context, lumeID string) (*modellume.Lume, error)
	ListLumesByLumoID(ctx context.Context, lumoID string, limit, offset int32) ([]*modellume.Lume, error)
	ListLumesByType(ctx context.Context, lumoID string, lumeType modellume.LumeType, limit, offset int32) ([]*modellume.Lume, error)
	SearchLumesByLocation(ctx context.Context, lumoID string, minLat, maxLat, minLng, maxLng float64, limit, offset int32) ([]*modellume.Lume, error)
	UpdateLume(ctx context.Context, domainLume *modellume.Lume) (*modellume.Lume, error)
	DeleteLume(ctx context.Context, id int64) error
	DeleteLumeByLumeID(ctx context.Context, lumeID string) error
	CountLumesByLumo(ctx context.Context, lumoID string) (int64, error)
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
	domainLume, err := a.toDomainModelForCreate(req)
	if err != nil {
		return nil, fmt.Errorf("failed to convert request: %w", err)
	}

	return a.repo.CreateLume(ctx, domainLume)
}

// GetLumeByID retrieves a Lume by its internal ID
func (a *App) GetLumeByID(ctx context.Context, id int64) (*modellume.Lume, error) {
	return a.repo.GetLumeByID(ctx, id)
}

// GetLumeByLumeID retrieves a Lume by its UUID string
func (a *App) GetLumeByLumeID(ctx context.Context, lumeID string) (*modellume.Lume, error) {
	if _, err := uuid.Parse(lumeID); err != nil {
		return nil, ErrInvalidLumeID
	}

	return a.repo.GetLumeByLumeID(ctx, lumeID)
}

// ListLumesByLumoID retrieves all Lumes for a given Lumo
func (a *App) ListLumesByLumoID(ctx context.Context, req ListLumesRequest) ([]*modellume.Lume, error) {
	if _, err := uuid.Parse(req.LumoID); err != nil {
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

	return a.repo.ListLumesByLumoID(ctx, req.LumoID, limit, offset)
}

// ListLumesByType retrieves all Lumes of a specific type for a Lumo
func (a *App) ListLumesByType(ctx context.Context, req ListLumesByTypeRequest) ([]*modellume.Lume, error) {
	if _, err := uuid.Parse(req.LumoID); err != nil {
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

	return a.repo.ListLumesByType(ctx, req.LumoID, req.Type, limit, offset)
}

// SearchLumesByLocation finds Lumes within a bounding box for a specific Lumo
func (a *App) SearchLumesByLocation(ctx context.Context, req SearchLumesByLocationRequest) ([]*modellume.Lume, error) {
	if _, err := uuid.Parse(req.LumoID); err != nil {
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

	return a.repo.SearchLumesByLocation(ctx, req.LumoID, req.MinLat, req.MaxLat, req.MinLng, req.MaxLng, limit, offset)
}

// UpdateLume updates an existing Lume
func (a *App) UpdateLume(ctx context.Context, id int64, req UpdateLumeRequest) (*modellume.Lume, error) {
	// First get the existing lume
	existingLume, err := a.repo.GetLumeByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update the domain model with new values
	updatedLume := a.updateDomainModel(existingLume, req)

	return a.repo.UpdateLume(ctx, updatedLume)
}

// UpdateLumeByLumeID updates a Lume by its UUID
func (a *App) UpdateLumeByLumeID(ctx context.Context, lumeID string, req UpdateLumeRequest) (*modellume.Lume, error) {
	if _, err := uuid.Parse(lumeID); err != nil {
		return nil, ErrInvalidLumeID
	}

	// First get the lume to find its internal ID
	existingLume, err := a.repo.GetLumeByLumeID(ctx, lumeID)
	if err != nil {
		return nil, err
	}

	// Update the domain model with new values
	updatedLume := a.updateDomainModel(existingLume, req)

	return a.repo.UpdateLume(ctx, updatedLume)
}

// DeleteLume deletes a Lume by its ID
func (a *App) DeleteLume(ctx context.Context, id int64) error {
	return a.repo.DeleteLume(ctx, id)
}

// DeleteLumeByLumeID deletes a Lume by its UUID
func (a *App) DeleteLumeByLumeID(ctx context.Context, lumeID string) error {
	if _, err := uuid.Parse(lumeID); err != nil {
		return ErrInvalidLumeID
	}

	return a.repo.DeleteLumeByLumeID(ctx, lumeID)
}

// CountLumesByLumo returns the total count of Lumes for a Lumo
func (a *App) CountLumesByLumo(ctx context.Context, lumoID string) (int64, error) {
	if _, err := uuid.Parse(lumoID); err != nil {
		return 0, ErrInvalidLumoID
	}

	return a.repo.CountLumesByLumo(ctx, lumoID)
}

// toDomainModelForCreate creates a new domain model from the create request
func (a *App) toDomainModelForCreate(req CreateLumeRequest) (*modellume.Lume, error) {
	// Create a new domain model
	lume := &modellume.Lume{
		LumeID:       uuid.New().String(), // Generate new UUID
		LumoID:       req.LumoID,
		Type:         req.Type,
		Name:         req.Label,
		Description:  req.Description,
		DateStart:    req.DateStart,
		DateEnd:      req.DateEnd,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Address:      req.Address,
		Images:       req.Images,
		CategoryTags: req.CategoryTags,
		BookingLink:  req.BookingLink,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Ensure arrays are initialized
	if lume.Images == nil {
		lume.Images = make([]string, 0)
	}
	if lume.CategoryTags == nil {
		lume.CategoryTags = make([]string, 0)
	}

	return lume, nil
}

// updateDomainModel updates an existing domain model with values from the update request
func (a *App) updateDomainModel(existingLume *modellume.Lume, req UpdateLumeRequest) *modellume.Lume {
	// Update fields
	existingLume.Name = req.Label
	existingLume.Type = req.Type
	existingLume.Description = req.Description
	existingLume.DateStart = req.DateStart
	existingLume.DateEnd = req.DateEnd
	existingLume.Latitude = req.Latitude
	existingLume.Longitude = req.Longitude
	existingLume.Address = req.Address
	existingLume.BookingLink = req.BookingLink
	existingLume.UpdatedAt = time.Now()

	// Update arrays if provided
	if req.Images != nil {
		existingLume.Images = req.Images
	}
	if req.CategoryTags != nil {
		existingLume.CategoryTags = req.CategoryTags
	}

	return existingLume
}
