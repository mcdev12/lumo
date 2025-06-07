package link

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/mcdev12/lumo/go/internal/models/link"
	"github.com/mcdev12/lumo/go/internal/repository/db/sqlc"
	"github.com/sqlc-dev/pqtype"
)

//go:generate mockery
type LinkQuerier interface {
	CountLinksByFromLumeID(ctx context.Context, fromLumeID uuid.UUID) (int64, error)
	CountLinksByLumeID(ctx context.Context, fromLumeID uuid.UUID) (int64, error)
	CountLinksByToLumeID(ctx context.Context, toLumeID uuid.UUID) (int64, error)
	CreateLink(ctx context.Context, arg sqlc.CreateLinkParams) (sqlc.Link, error)
	DeleteLink(ctx context.Context, id int64) error
	DeleteLinkByLinkID(ctx context.Context, linkID uuid.UUID) error
	GetLinkByID(ctx context.Context, id int64) (sqlc.Link, error)
	GetLinkByLinkID(ctx context.Context, linkID uuid.UUID) (sqlc.Link, error)
	ListLinksByEitherLumeID(ctx context.Context, arg sqlc.ListLinksByEitherLumeIDParams) ([]sqlc.Link, error)
	ListLinksByFromLumeID(ctx context.Context, arg sqlc.ListLinksByFromLumeIDParams) ([]sqlc.Link, error)
	ListLinksByLumeIDAndType(ctx context.Context, arg sqlc.ListLinksByLumeIDAndTypeParams) ([]sqlc.Link, error)
	ListLinksByToLumeID(ctx context.Context, arg sqlc.ListLinksByToLumeIDParams) ([]sqlc.Link, error)
	ListLinksByType(ctx context.Context, arg sqlc.ListLinksByTypeParams) ([]sqlc.Link, error)
	UpdateLink(ctx context.Context, arg sqlc.UpdateLinkParams) (sqlc.Link, error)
}

// Repository is the concrete implementation for Link data access
type Repository struct {
	queries LinkQuerier
}

// NewRepository creates a new Repository instance
func NewRepository(db sqlc.DBTX) *Repository {
	return &Repository{
		queries: sqlc.New(db),
	}
}

// CreateLink creates a new Link record from domain model
func (r *Repository) CreateLink(ctx context.Context, domainLink *link.Link) (*link.Link, error) {
	params := r.domainToCreateParams(domainLink)

	result, err := r.queries.CreateLink(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.sqlcRowToDomainModel(result), nil
}

// GetLinkByID retrieves a Link by its internal ID
func (r *Repository) GetLinkByID(ctx context.Context, id int64) (*link.Link, error) {
	result, err := r.queries.GetLinkByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return r.sqlcRowToDomainModel(result), nil
}

// GetLinkByLinkID retrieves a Link by its UUID
func (r *Repository) GetLinkByLinkID(ctx context.Context, linkID string) (*link.Link, error) {
	parsedUUID, err := uuid.Parse(linkID)
	if err != nil {
		return nil, err
	}

	result, err := r.queries.GetLinkByLinkID(ctx, parsedUUID)
	if err != nil {
		return nil, err
	}

	return r.sqlcRowToDomainModel(result), nil
}

// ListLinksByFromLumeID retrieves all Links from a specific Lume
func (r *Repository) ListLinksByFromLumeID(ctx context.Context, fromLumeID string, limit, offset int32) ([]*link.Link, error) {
	parsedLumeID, err := uuid.Parse(fromLumeID)
	if err != nil {
		return nil, err
	}

	params := sqlc.ListLinksByFromLumeIDParams{
		FromLumeID: parsedLumeID,
		Limit:      limit,
		Offset:     offset,
	}

	results, err := r.queries.ListLinksByFromLumeID(ctx, params)
	if err != nil {
		return nil, err
	}

	links := make([]*link.Link, len(results))
	for i, result := range results {
		links[i] = r.sqlcRowToDomainModel(result)
	}

	return links, nil
}

// ListLinksByToLumeID retrieves all Links to a specific Lume
func (r *Repository) ListLinksByToLumeID(ctx context.Context, toLumeID string, limit, offset int32) ([]*link.Link, error) {
	parsedLumeID, err := uuid.Parse(toLumeID)
	if err != nil {
		return nil, err
	}

	params := sqlc.ListLinksByToLumeIDParams{
		ToLumeID: parsedLumeID,
		Limit:    limit,
		Offset:   offset,
	}

	results, err := r.queries.ListLinksByToLumeID(ctx, params)
	if err != nil {
		return nil, err
	}

	links := make([]*link.Link, len(results))
	for i, result := range results {
		links[i] = r.sqlcRowToDomainModel(result)
	}

	return links, nil
}

// ListLinksByEitherLumeID retrieves all Links connected to a specific Lume (either from or to)
func (r *Repository) ListLinksByEitherLumeID(ctx context.Context, lumeID string, limit, offset int32) ([]*link.Link, error) {
	parsedLumeID, err := uuid.Parse(lumeID)
	if err != nil {
		return nil, err
	}

	params := sqlc.ListLinksByEitherLumeIDParams{
		FromLumeID: parsedLumeID,
		Limit:      limit,
		Offset:     offset,
	}

	results, err := r.queries.ListLinksByEitherLumeID(ctx, params)
	if err != nil {
		return nil, err
	}

	links := make([]*link.Link, len(results))
	for i, result := range results {
		links[i] = r.sqlcRowToDomainModel(result)
	}

	return links, nil
}

// ListLinksByType retrieves all Links of a specific type
func (r *Repository) ListLinksByType(ctx context.Context, linkType link.LinkType, limit, offset int32) ([]*link.Link, error) {
	params := sqlc.ListLinksByTypeParams{
		LinkType: string(linkType),
		Limit:    limit,
		Offset:   offset,
	}

	results, err := r.queries.ListLinksByType(ctx, params)
	if err != nil {
		return nil, err
	}

	links := make([]*link.Link, len(results))
	for i, result := range results {
		links[i] = r.sqlcRowToDomainModel(result)
	}

	return links, nil
}

// ListLinksByLumeIDAndType retrieves all Links connected to a specific Lume with a specific type
func (r *Repository) ListLinksByLumeIDAndType(ctx context.Context, lumeID string, linkType link.LinkType, limit, offset int32) ([]*link.Link, error) {
	parsedLumeID, err := uuid.Parse(lumeID)
	if err != nil {
		return nil, err
	}

	params := sqlc.ListLinksByLumeIDAndTypeParams{
		FromLumeID: parsedLumeID,
		LinkType:   string(linkType),
		Limit:      limit,
		Offset:     offset,
	}

	results, err := r.queries.ListLinksByLumeIDAndType(ctx, params)
	if err != nil {
		return nil, err
	}

	links := make([]*link.Link, len(results))
	for i, result := range results {
		links[i] = r.sqlcRowToDomainModel(result)
	}

	return links, nil
}

// UpdateLink updates an existing Link record
func (r *Repository) UpdateLink(ctx context.Context, domainLink *link.Link) (*link.Link, error) {
	params := r.domainToUpdateParams(domainLink)

	result, err := r.queries.UpdateLink(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.sqlcRowToDomainModel(result), nil
}

// DeleteLink deletes a Link by its internal ID
func (r *Repository) DeleteLink(ctx context.Context, id int64) error {
	return r.queries.DeleteLink(ctx, id)
}

// DeleteLinkByLinkID deletes a Link by its UUID
func (r *Repository) DeleteLinkByLinkID(ctx context.Context, linkID string) error {
	parsedUUID, err := uuid.Parse(linkID)
	if err != nil {
		return err
	}
	return r.queries.DeleteLinkByLinkID(ctx, parsedUUID)
}

// CountLinksByLumeID returns the total count of Links connected to a Lume
func (r *Repository) CountLinksByLumeID(ctx context.Context, lumeID string) (int64, error) {
	parsedLumeID, err := uuid.Parse(lumeID)
	if err != nil {
		return 0, err
	}
	return r.queries.CountLinksByLumeID(ctx, parsedLumeID)
}

// CountLinksByFromLumeID returns the total count of Links from a Lume
func (r *Repository) CountLinksByFromLumeID(ctx context.Context, fromLumeID string) (int64, error) {
	parsedLumeID, err := uuid.Parse(fromLumeID)
	if err != nil {
		return 0, err
	}
	return r.queries.CountLinksByFromLumeID(ctx, parsedLumeID)
}

// CountLinksByToLumeID returns the total count of Links to a Lume
func (r *Repository) CountLinksByToLumeID(ctx context.Context, toLumeID string) (int64, error) {
	parsedLumeID, err := uuid.Parse(toLumeID)
	if err != nil {
		return 0, err
	}
	return r.queries.CountLinksByToLumeID(ctx, parsedLumeID)
}

// Helper method to convert domain Link to SQLC CreateLinkParams
func (r *Repository) domainToCreateParams(domainLink *link.Link) sqlc.CreateLinkParams {
	now := time.Now()
	params := sqlc.CreateLinkParams{
		LinkID:     uuid.MustParse(domainLink.LinkID),
		FromLumeID: uuid.MustParse(domainLink.FromLumeID),
		ToLumeID:   uuid.MustParse(domainLink.ToLumeID),
		LinkType:   string(domainLink.Type),
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Handle optional notes
	if domainLink.Notes != nil {
		params.Notes = sql.NullString{String: *domainLink.Notes, Valid: true}
	}

	// Handle optional sequence index
	if domainLink.SequenceIndex != nil {
		params.SequenceIndex = sql.NullInt32{Int32: *domainLink.SequenceIndex, Valid: true}
	}

	// Handle travel details
	if domainLink.Travel != nil {
		travelJSON, _ := json.Marshal(domainLink.Travel)
		params.TravelDetails = pqtype.NullRawMessage{RawMessage: travelJSON, Valid: true}
	}

	return params
}

// Helper method to convert domain Link to SQLC UpdateLinkParams
func (r *Repository) domainToUpdateParams(domainLink *link.Link) sqlc.UpdateLinkParams {
	params := sqlc.UpdateLinkParams{
		LinkID:     uuid.MustParse(domainLink.LinkID),
		FromLumeID: uuid.MustParse(domainLink.FromLumeID),
		ToLumeID:   uuid.MustParse(domainLink.ToLumeID),
		LinkType:   string(domainLink.Type),
		UpdatedAt:  time.Now(),
	}

	// Handle optional notes
	if domainLink.Notes != nil {
		params.Notes = sql.NullString{String: *domainLink.Notes, Valid: true}
	}

	// Handle optional sequence index
	if domainLink.SequenceIndex != nil {
		params.SequenceIndex = sql.NullInt32{Int32: *domainLink.SequenceIndex, Valid: true}
	}

	// Handle travel details
	if domainLink.Travel != nil {
		travelJSON, _ := json.Marshal(domainLink.Travel)
		params.TravelDetails = pqtype.NullRawMessage{RawMessage: travelJSON, Valid: true}
	}

	return params
}

// Helper method to convert SQLC results to domain model
func (r *Repository) sqlcRowToDomainModel(row sqlc.Link) *link.Link {
	domainLink := &link.Link{
		ID:         row.ID,
		LinkID:     row.LinkID.String(),
		FromLumeID: row.FromLumeID.String(),
		ToLumeID:   row.ToLumeID.String(),
		Type:       link.LinkType(row.LinkType),
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}

	// Handle optional notes
	if row.Notes.Valid {
		notes := row.Notes.String
		domainLink.Notes = &notes
	}

	// Handle optional sequence index
	if row.SequenceIndex.Valid {
		sequenceIndex := row.SequenceIndex.Int32
		domainLink.SequenceIndex = &sequenceIndex
	}

	// Handle travel details
	if row.TravelDetails.Valid {
		var travelDetails link.TravelDetails
		if err := json.Unmarshal(row.TravelDetails.RawMessage, &travelDetails); err == nil {
			domainLink.Travel = &travelDetails
		}
	}

	return domainLink
}
