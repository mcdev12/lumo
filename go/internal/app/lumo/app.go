package lumo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	modellumo "github.com/mcdev12/lumo/go/internal/models/lumo"
)

// Domain errors
var (
	ErrLumoNotFound  = errors.New("lumo not found")
	ErrInvalidUserID = errors.New("invalid user ID")
	ErrInvalidLumoID = errors.New("invalid lumo ID")
	ErrEmptyTitle    = errors.New("title cannot be empty")
)

// LumoRepository defines what the app layer needs from the repository
type LumoRepository interface {
	CreateLumo(ctx context.Context, domainLumo *modellumo.Lumo) (*modellumo.Lumo, error)
	GetLumoByID(ctx context.Context, id int64) (*modellumo.Lumo, error)
	GetLumoByLumoID(ctx context.Context, lumoID string) (*modellumo.Lumo, error)
	ListLumosByUserID(ctx context.Context, userID string, limit, offset int32) ([]*modellumo.Lumo, error)
	UpdateLumo(ctx context.Context, domainLumo *modellumo.Lumo) (*modellumo.Lumo, error)
	DeleteLumo(ctx context.Context, id int64) error
	DeleteLumoByLumoID(ctx context.Context, lumoID string) error
	CountLumosByUserID(ctx context.Context, userID string) (int64, error)
}

// CreateLumoRequest represents the business layer's create request
type CreateLumoRequest struct {
	UserID string
	Title  string
}

// UpdateLumoRequest represents the business layer's update request
type UpdateLumoRequest struct {
	Title string
}

// ListLumosRequest represents pagination parameters
type ListLumosRequest struct {
	UserID string
	Limit  int32
	Offset int32
}

// App handles business logic for Lumos
type App struct {
	repo LumoRepository
}

// NewLumoApp creates a new Lumo Service
func NewLumoApp(repo LumoRepository) *App {
	return &App{
		repo: repo,
	}
}

// CreateLumo creates a new Lumo with business logic validation
func (a *App) CreateLumo(ctx context.Context, req CreateLumoRequest) (*modellumo.Lumo, error) {
	if err := a.validateCreateRequest(req); err != nil {
		return nil, err
	}

	domainLumo, err := a.toDomainModelForCreate(req)
	if err != nil {
		return nil, fmt.Errorf("failed to convert request: %w", err)
	}

	return a.repo.CreateLumo(ctx, domainLumo)
}

// GetLumoByID retrieves a Lumo by its internal ID
func (a *App) GetLumoByID(ctx context.Context, id int64) (*modellumo.Lumo, error) {
	return a.repo.GetLumoByID(ctx, id)
}

// GetLumoByLumoID retrieves a Lumo by its UUID string
func (a *App) GetLumoByLumoID(ctx context.Context, lumoID string) (*modellumo.Lumo, error) {
	if _, err := uuid.Parse(lumoID); err != nil {
		return nil, ErrInvalidLumoID
	}

	return a.repo.GetLumoByLumoID(ctx, lumoID)
}

// ListLumosByUserID retrieves all Lumos for a given user
func (a *App) ListLumosByUserID(ctx context.Context, req ListLumosRequest) ([]*modellumo.Lumo, error) {
	if _, err := uuid.Parse(req.UserID); err != nil {
		return nil, ErrInvalidUserID
	}

	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	return a.repo.ListLumosByUserID(ctx, req.UserID, limit, offset)
}

// UpdateLumo updates an existing Lumo
func (a *App) UpdateLumo(ctx context.Context, id int64, req UpdateLumoRequest) (*modellumo.Lumo, error) {
	if err := a.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	// First get the existing lumo
	existingLumo, err := a.repo.GetLumoByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update the domain model with new values
	updatedLumo := a.updateDomainModel(existingLumo, req)

	return a.repo.UpdateLumo(ctx, updatedLumo)
}

// UpdateLumoByLumoID updates a Lumo by its UUID
func (a *App) UpdateLumoByLumoID(ctx context.Context, lumoID string, req UpdateLumoRequest) (*modellumo.Lumo, error) {
	if _, err := uuid.Parse(lumoID); err != nil {
		return nil, ErrInvalidLumoID
	}

	// First get the lumo to find its internal ID
	existingLumo, err := a.repo.GetLumoByLumoID(ctx, lumoID)
	if err != nil {
		return nil, err
	}

	if err := a.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	// Update the domain model with new values
	updatedLumo := a.updateDomainModel(existingLumo, req)

	return a.repo.UpdateLumo(ctx, updatedLumo)
}

// DeleteLumo deletes a Lumo by its ID
func (a *App) DeleteLumo(ctx context.Context, id int64) error {
	return a.repo.DeleteLumo(ctx, id)
}

// DeleteLumoByLumoID deletes a Lumo by its UUID
func (a *App) DeleteLumoByLumoID(ctx context.Context, lumoID string) error {
	if _, err := uuid.Parse(lumoID); err != nil {
		return ErrInvalidLumoID
	}
	
	return a.repo.DeleteLumoByLumoID(ctx, lumoID)
}

// CountLumosByUserID returns the total count of Lumos for a user
func (a *App) CountLumosByUserID(ctx context.Context, userID string) (int64, error) {
	if _, err := uuid.Parse(userID); err != nil {
		return 0, ErrInvalidUserID
	}
	
	return a.repo.CountLumosByUserID(ctx, userID)
}

// Validation methods
func (a *App) validateCreateRequest(req CreateLumoRequest) error {
	if req.Title == "" {
		return ErrEmptyTitle
	}

	if _, err := uuid.Parse(req.UserID); err != nil {
		return ErrInvalidUserID
	}

	return nil
}

func (a *App) validateUpdateRequest(req UpdateLumoRequest) error {
	if req.Title == "" {
		return ErrEmptyTitle
	}

	return nil
}

// Conversion methods
// toDomainModelForCreate creates a new domain model from the create request
func (a *App) toDomainModelForCreate(req CreateLumoRequest) (*modellumo.Lumo, error) {
	// Create a new domain model
	lumo := modellumo.NewLumo(req.UserID, req.Title)
	return lumo, nil
}

// updateDomainModel updates an existing domain model with values from the update request
func (a *App) updateDomainModel(existingLumo *modellumo.Lumo, req UpdateLumoRequest) *modellumo.Lumo {
	// Update fields
	existingLumo.Title = req.Title
	existingLumo.UpdatedAt = time.Now()
	return existingLumo
}