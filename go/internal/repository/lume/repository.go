package lume

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/mcdev12/lumo/go/internal/models/lume"
	"github.com/mcdev12/lumo/go/internal/repository/db/sqlc"
)

//go:generate mockery

// Repository is the concrete implementation for Lume data access
type Repository struct {
	queries sqlc.Querier
}

// NewRepository creates a new Repository instance
func NewRepository(db sqlc.DBTX) *Repository {
	return &Repository{
		queries: sqlc.New(db),
	}
}

// CreateLume creates a new Lume record from domain model
func (r *Repository) CreateLume(ctx context.Context, domainLume *lume.Lume) (*lume.Lume, error) {
	params := r.domainToCreateParams(domainLume)

	result, err := r.queries.CreateLume(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.sqlcRowToDomainModel(result), nil
}

// GetLumeByID retrieves a Lume by its internal ID
func (r *Repository) GetLumeByID(ctx context.Context, id int64) (*lume.Lume, error) {
	result, err := r.queries.GetLumeByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return r.sqlcRowToDomainModel(result), nil
}

// GetLumeByLumeID retrieves a Lume by its UUID
func (r *Repository) GetLumeByLumeID(ctx context.Context, lumeID string) (*lume.Lume, error) {
	parsedUUID, err := uuid.Parse(lumeID)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.GetLumeByLumeID(ctx, parsedUUID)
	if err != nil {
		return nil, err
	}

	return r.sqlcRowToDomainModel(result), nil
}

// ListLumesByLumoID retrieves all Lumes for a given Lumo
func (r *Repository) ListLumesByLumoID(ctx context.Context, lumoID string, limit, offset int32) ([]*lume.Lume, error) {
	parsedLumoID, err := uuid.Parse(lumoID)
	if err != nil {
		return nil, err
	}

	params := sqlc.ListLumesByLumoIDParams{
		LumoID: parsedLumoID,
		Limit:  limit,
		Offset: offset,
	}

	results, err := r.queries.ListLumesByLumoID(ctx, params)
	if err != nil {
		return nil, err
	}

	lumes := make([]*lume.Lume, len(results))
	for i, result := range results {
		lumes[i] = r.sqlcRowToDomainModel(result)
	}

	return lumes, nil
}

// ListLumesByType retrieves all Lumes of a specific type for a Lumo
func (r *Repository) ListLumesByType(ctx context.Context, lumoID string, lumeType lume.LumeType, limit, offset int32) ([]*lume.Lume, error) {
	parsedLumoID, err := uuid.Parse(lumoID)
	if err != nil {
		return nil, err
	}

	params := sqlc.ListLumesByTypeParams{
		LumoID: parsedLumoID,
		Type:   string(lumeType),
		Limit:  limit,
		Offset: offset,
	}

	results, err := r.queries.ListLumesByType(ctx, params)
	if err != nil {
		return nil, err
	}

	lumes := make([]*lume.Lume, len(results))
	for i, result := range results {
		lumes[i] = r.sqlcRowToDomainModel(result)
	}

	return lumes, nil
}

// SearchLumesByLocation finds Lumes within a bounding box for a specific Lumo
func (r *Repository) SearchLumesByLocation(ctx context.Context, lumoID string, minLat, maxLat, minLng, maxLng float64, limit, offset int32) ([]*lume.Lume, error) {
	parsedLumoID, err := uuid.Parse(lumoID)
	if err != nil {
		return nil, err
	}

	params := sqlc.SearchLumesByLocationParams{
		LumoID:      parsedLumoID,
		Latitude:    sql.NullFloat64{Float64: minLat, Valid: true},
		Latitude_2:  sql.NullFloat64{Float64: maxLat, Valid: true},
		Longitude:   sql.NullFloat64{Float64: minLng, Valid: true},
		Longitude_2: sql.NullFloat64{Float64: maxLng, Valid: true},
		Limit:       limit,
		Offset:      offset,
	}

	results, err := r.queries.SearchLumesByLocation(ctx, params)
	if err != nil {
		return nil, err
	}

	lumes := make([]*lume.Lume, len(results))
	for i, result := range results {
		lumes[i] = r.sqlcRowToDomainModel(result)
	}

	return lumes, nil
}

// UpdateLume updates an existing Lume record
func (r *Repository) UpdateLume(ctx context.Context, domainLume *lume.Lume) (*lume.Lume, error) {
	params := r.domainToUpdateParams(domainLume)

	result, err := r.queries.UpdateLume(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.sqlcRowToDomainModel(result), nil
}

// DeleteLume deletes a Lume by its internal ID
func (r *Repository) DeleteLume(ctx context.Context, id int64) error {
	return r.queries.DeleteLume(ctx, id)
}

// DeleteLumeByLumeID deletes a Lume by its UUID
func (r *Repository) DeleteLumeByLumeID(ctx context.Context, lumeID string) error {
	parsedUUID, err := uuid.Parse(lumeID)
	if err != nil {
		return err
	}
	return r.queries.DeleteLumeByLumeID(ctx, parsedUUID)
}

// CountLumesByLumo returns the total count of Lumes for a Lumo
func (r *Repository) CountLumesByLumo(ctx context.Context, lumoID string) (int64, error) {
	parsedLumoID, err := uuid.Parse(lumoID)
	if err != nil {
		return 0, err
	}
	return r.queries.CountLumesByLumo(ctx, parsedLumoID)
}

// ensureStringArray ensures empty arrays instead of nil for consistency
func (r *Repository) ensureStringArray(arr []string) []string {
	if arr == nil {
		return make([]string, 0)
	}
	return arr
}

// Helper method to convert domain Lume to SQLC CreateLumeParams
func (r *Repository) domainToCreateParams(domainLume *lume.Lume) sqlc.CreateLumeParams {
	now := time.Now()
	params := sqlc.CreateLumeParams{
		LumeID:       uuid.MustParse(domainLume.LumeID),
		LumoID:       uuid.MustParse(domainLume.LumoID),
		Type:         string(domainLume.Type),
		Name:         domainLume.Name,
		Images:       r.ensureStringArray(domainLume.Images),
		CategoryTags: r.ensureStringArray(domainLume.CategoryTags),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Handle description as sql.NullString
	if domainLume.Description != "" {
		params.Description = sql.NullString{String: domainLume.Description, Valid: true}
	}

	// Handle optional timestamps
	if domainLume.DateStart != nil {
		params.DateStart = sql.NullTime{Time: *domainLume.DateStart, Valid: true}
	}
	if domainLume.DateEnd != nil {
		params.DateEnd = sql.NullTime{Time: *domainLume.DateEnd, Valid: true}
	}

	// Handle optional coordinates
	if domainLume.Latitude != nil {
		params.Latitude = sql.NullFloat64{Float64: *domainLume.Latitude, Valid: true}
	}
	if domainLume.Longitude != nil {
		params.Longitude = sql.NullFloat64{Float64: *domainLume.Longitude, Valid: true}
	}

	// Handle optional address
	if domainLume.Address != nil {
		params.Address = sql.NullString{String: *domainLume.Address, Valid: true}
	}

	// Handle optional booking link
	if domainLume.BookingLink != nil {
		params.BookingLink = sql.NullString{String: *domainLume.BookingLink, Valid: true}
	}

	return params
}

// Helper method to convert domain Lume to SQLC UpdateLumeParams
func (r *Repository) domainToUpdateParams(domainLume *lume.Lume) sqlc.UpdateLumeParams {
	params := sqlc.UpdateLumeParams{
		LumeID:       uuid.MustParse(domainLume.LumeID),
		Name:         domainLume.Name,
		Type:         string(domainLume.Type),
		Images:       r.ensureStringArray(domainLume.Images),
		CategoryTags: r.ensureStringArray(domainLume.CategoryTags),
		UpdatedAt:    time.Now(),
	}

	// Handle description as sql.NullString
	if domainLume.Description != "" {
		params.Description = sql.NullString{String: domainLume.Description, Valid: true}
	}

	// Handle optional timestamps
	if domainLume.DateStart != nil {
		params.DateStart = sql.NullTime{Time: *domainLume.DateStart, Valid: true}
	}
	if domainLume.DateEnd != nil {
		params.DateEnd = sql.NullTime{Time: *domainLume.DateEnd, Valid: true}
	}

	// Handle optional coordinates
	if domainLume.Latitude != nil {
		params.Latitude = sql.NullFloat64{Float64: *domainLume.Latitude, Valid: true}
	}
	if domainLume.Longitude != nil {
		params.Longitude = sql.NullFloat64{Float64: *domainLume.Longitude, Valid: true}
	}

	// Handle optional address
	if domainLume.Address != nil {
		params.Address = sql.NullString{String: *domainLume.Address, Valid: true}
	}

	// Handle optional booking link
	if domainLume.BookingLink != nil {
		params.BookingLink = sql.NullString{String: *domainLume.BookingLink, Valid: true}
	}

	return params
}

// Helper method to convert SQLC results to domain model
func (r *Repository) sqlcRowToDomainModel(row sqlc.Lume) *lume.Lume {
	domainLume := &lume.Lume{
		ID:           row.ID,
		LumeID:       row.LumeID.String(),
		LumoID:       row.LumoID.String(),
		Type:         lume.LumeType(row.Type),
		Name:         row.Name,
		Images:       r.ensureStringArray(row.Images),
		CategoryTags: r.ensureStringArray(row.CategoryTags),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}

	// Handle description from sql.NullString
	if row.Description.Valid {
		domainLume.Description = row.Description.String
	}

	// Handle optional fields
	if row.DateStart.Valid {
		domainLume.DateStart = &row.DateStart.Time
	}
	if row.DateEnd.Valid {
		domainLume.DateEnd = &row.DateEnd.Time
	}
	if row.Latitude.Valid {
		domainLume.Latitude = &row.Latitude.Float64
	}
	if row.Longitude.Valid {
		domainLume.Longitude = &row.Longitude.Float64
	}
	if row.Address.Valid {
		domainLume.Address = &row.Address.String
	}
	if row.BookingLink.Valid {
		domainLume.BookingLink = &row.BookingLink.String
	}

	return domainLume
}
