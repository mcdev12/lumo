package lumo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mcdev12/lumo/go/internal/models/lumo"
	"github.com/mcdev12/lumo/go/internal/repository/db/sqlc"
)

//go:generate mockery
type LumoQuerier interface {
	CreateLumo(ctx context.Context, arg sqlc.CreateLumoParams) (sqlc.Lumo, error)
	GetLumoByID(ctx context.Context, id int64) (sqlc.Lumo, error)
	GetLumoByLumoID(ctx context.Context, lumoID uuid.UUID) (sqlc.Lumo, error)
	ListLumosByUserID(ctx context.Context, arg sqlc.ListLumosByUserIDParams) ([]sqlc.Lumo, error)
	UpdateLumo(ctx context.Context, arg sqlc.UpdateLumoParams) (sqlc.Lumo, error)
	DeleteLumo(ctx context.Context, id int64) error
	DeleteLumoByLumoID(ctx context.Context, lumoID uuid.UUID) error
	CountLumosByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
}

// Repository is the concrete implementation for Lumo data access
type Repository struct {
	queries LumoQuerier
}

// NewRepository creates a new Repository instance
func NewRepository(db sqlc.DBTX) *Repository {
	return &Repository{
		queries: sqlc.New(db),
	}
}

// CreateLumo creates a new Lumo record from domain model
func (r *Repository) CreateLumo(ctx context.Context, domainLumo *lumo.Lumo) (*lumo.Lumo, error) {
	params := r.domainToCreateParams(domainLumo)

	result, err := r.queries.CreateLumo(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.sqlcRowToDomainModel(result), nil
}

// GetLumoByID retrieves a Lumo by its internal ID
func (r *Repository) GetLumoByID(ctx context.Context, id int64) (*lumo.Lumo, error) {
	result, err := r.queries.GetLumoByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return r.sqlcRowToDomainModel(result), nil
}

// GetLumoByLumoID retrieves a Lumo by its UUID
func (r *Repository) GetLumoByLumoID(ctx context.Context, lumoID string) (*lumo.Lumo, error) {
	parsedUUID, err := uuid.Parse(lumoID)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.GetLumoByLumoID(ctx, parsedUUID)
	if err != nil {
		return nil, err
	}

	return r.sqlcRowToDomainModel(result), nil
}

// ListLumosByUserID retrieves all Lumos for a given user
func (r *Repository) ListLumosByUserID(ctx context.Context, userID string, limit, offset int32) ([]*lumo.Lumo, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	params := sqlc.ListLumosByUserIDParams{
		UserID: parsedUserID,
		Limit:  limit,
		Offset: offset,
	}

	results, err := r.queries.ListLumosByUserID(ctx, params)
	if err != nil {
		return nil, err
	}

	lumos := make([]*lumo.Lumo, len(results))
	for i, result := range results {
		lumos[i] = r.sqlcRowToDomainModel(result)
	}

	return lumos, nil
}

// UpdateLumo updates an existing Lumo record
func (r *Repository) UpdateLumo(ctx context.Context, domainLumo *lumo.Lumo) (*lumo.Lumo, error) {
	params := r.domainToUpdateParams(domainLumo)

	result, err := r.queries.UpdateLumo(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.sqlcRowToDomainModel(result), nil
}

// DeleteLumo deletes a Lumo by its internal ID
func (r *Repository) DeleteLumo(ctx context.Context, id int64) error {
	return r.queries.DeleteLumo(ctx, id)
}

// DeleteLumoByLumoID deletes a Lumo by its UUID
func (r *Repository) DeleteLumoByLumoID(ctx context.Context, lumoID string) error {
	parsedUUID, err := uuid.Parse(lumoID)
	if err != nil {
		return err
	}
	return r.queries.DeleteLumoByLumoID(ctx, parsedUUID)
}

// CountLumosByUserID returns the total count of Lumos for a user
func (r *Repository) CountLumosByUserID(ctx context.Context, userID string) (int64, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return 0, err
	}
	return r.queries.CountLumosByUserID(ctx, parsedUserID)
}

// Helper method to convert domain Lumo to SQLC CreateLumoParams
func (r *Repository) domainToCreateParams(domainLumo *lumo.Lumo) sqlc.CreateLumoParams {
	now := time.Now()
	return sqlc.CreateLumoParams{
		LumoID:    uuid.MustParse(domainLumo.LumoID),
		UserID:    uuid.MustParse(domainLumo.UserID),
		Title:     domainLumo.Title,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Helper method to convert domain Lumo to SQLC UpdateLumoParams
func (r *Repository) domainToUpdateParams(domainLumo *lumo.Lumo) sqlc.UpdateLumoParams {
	return sqlc.UpdateLumoParams{
		LumoID:    uuid.MustParse(domainLumo.LumoID),
		Title:     domainLumo.Title,
		UpdatedAt: time.Now(),
	}
}

// Helper method to convert SQLC results to domain model
func (r *Repository) sqlcRowToDomainModel(row sqlc.Lumo) *lumo.Lumo {
	return &lumo.Lumo{
		ID:        row.ID,
		LumoID:    row.LumoID.String(),
		UserID:    row.UserID.String(),
		Title:     row.Title,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}
