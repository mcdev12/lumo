package link

import (
	"context"
	"errors"
	"github.com/google/uuid"
	modellink "github.com/mcdev12/lumo/go/internal/models/link"
	"time"
)

// Domain errors
var (
	ErrLinkNotFound      = errors.New("link not found")
	ErrInvalidLinkID     = errors.New("invalid link ID")
	ErrInvalidLumeID     = errors.New("invalid lume ID")
	ErrInvalidLinkType   = errors.New("invalid link type")
	ErrInvalidTravelMode = errors.New("invalid travel mode")
	ErrEmptyNotes        = errors.New("notes cannot be empty")
)

// LinkRepository defines what the app layer needs from the repository
type LinkRepository interface {
	CreateLink(ctx context.Context, domainLink *modellink.Link) (*modellink.Link, error)
	GetLinkByID(ctx context.Context, id int64) (*modellink.Link, error)
	GetLinkByLinkID(ctx context.Context, linkID string) (*modellink.Link, error)
	ListLinksByFromLumeID(ctx context.Context, fromLumeID string, limit, offset int32) ([]*modellink.Link, error)
	ListLinksByToLumeID(ctx context.Context, toLumeID string, limit, offset int32) ([]*modellink.Link, error)
	ListLinksByEitherLumeID(ctx context.Context, lumeID string, limit, offset int32) ([]*modellink.Link, error)
	ListLinksByType(ctx context.Context, linkType modellink.LinkType, limit, offset int32) ([]*modellink.Link, error)
	ListLinksByLumeIDAndType(ctx context.Context, lumeID string, linkType modellink.LinkType, limit, offset int32) ([]*modellink.Link, error)
	UpdateLink(ctx context.Context, domainLink *modellink.Link) (*modellink.Link, error)
	DeleteLink(ctx context.Context, id int64) error
	DeleteLinkByLinkID(ctx context.Context, linkID string) error
	CountLinksByLumeID(ctx context.Context, lumeID string) (int64, error)
	CountLinksByFromLumeID(ctx context.Context, fromLumeID string) (int64, error)
	CountLinksByToLumeID(ctx context.Context, toLumeID string) (int64, error)
}

// App handles business logic for Links
type App struct {
	repo LinkRepository
}

// NewLinkApp creates a new Link App
func NewLinkApp(repo LinkRepository) *App {
	return &App{
		repo: repo,
	}
}

// CreateLink creates a new Link with business logic validation
func (a *App) CreateLink(ctx context.Context, req CreateLinkRequest) (*modellink.Link, error) {
	domainLink := a.toDomainModelForCreate(req)
	return a.repo.CreateLink(ctx, domainLink)
}

// GetLinkByID retrieves a Link by its internal ID
func (a *App) GetLinkByID(ctx context.Context, id int64) (*modellink.Link, error) {
	return a.repo.GetLinkByID(ctx, id)
}

// GetLinkByLinkID retrieves a Link by its UUID
func (a *App) GetLinkByLinkID(ctx context.Context, linkID string) (*modellink.Link, error) {
	if _, err := uuid.Parse(linkID); err != nil {
		return nil, ErrInvalidLinkID
	}
	return a.repo.GetLinkByLinkID(ctx, linkID)
}

// ListLinksByFromLumeID retrieves all Links from a specific Lume
func (a *App) ListLinksByFromLumeID(ctx context.Context, fromLumeID string, req ListLinksRequest) ([]*modellink.Link, error) {
	if _, err := uuid.Parse(fromLumeID); err != nil {
		return nil, ErrInvalidLumeID
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10 // Default limit
	}

	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	return a.repo.ListLinksByFromLumeID(ctx, fromLumeID, limit, offset)
}

// ListLinksByToLumeID retrieves all Links to a specific Lume
func (a *App) ListLinksByToLumeID(ctx context.Context, toLumeID string, req ListLinksRequest) ([]*modellink.Link, error) {
	if _, err := uuid.Parse(toLumeID); err != nil {
		return nil, ErrInvalidLumeID
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10 // Default limit
	}

	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	return a.repo.ListLinksByToLumeID(ctx, toLumeID, limit, offset)
}

// ListLinksByEitherLumeID retrieves all Links connected to a specific Lume (either from or to)
func (a *App) ListLinksByEitherLumeID(ctx context.Context, lumeID string, req ListLinksRequest) ([]*modellink.Link, error) {
	if _, err := uuid.Parse(lumeID); err != nil {
		return nil, ErrInvalidLumeID
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10 // Default limit
	}

	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	return a.repo.ListLinksByEitherLumeID(ctx, lumeID, limit, offset)
}

// UpdateLink updates an existing Link
func (a *App) UpdateLink(ctx context.Context, id int64, req UpdateLinkRequest) (*modellink.Link, error) {
	existingLink, err := a.repo.GetLinkByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updatedLink := a.updateDomainModel(existingLink, req)
	return a.repo.UpdateLink(ctx, updatedLink)
}

// UpdateLinkByLinkID updates an existing Link by its UUID
func (a *App) UpdateLinkByLinkID(ctx context.Context, linkID string, req UpdateLinkRequest) (*modellink.Link, error) {
	if _, err := uuid.Parse(linkID); err != nil {
		return nil, ErrInvalidLinkID
	}

	existingLink, err := a.repo.GetLinkByLinkID(ctx, linkID)
	if err != nil {
		return nil, err
	}

	updatedLink := a.updateDomainModel(existingLink, req)
	return a.repo.UpdateLink(ctx, updatedLink)
}

// DeleteLink deletes a Link by its internal ID
func (a *App) DeleteLink(ctx context.Context, id int64) error {
	return a.repo.DeleteLink(ctx, id)
}

// DeleteLinkByLinkID deletes a Link by its UUID
func (a *App) DeleteLinkByLinkID(ctx context.Context, linkID string) error {
	if _, err := uuid.Parse(linkID); err != nil {
		return ErrInvalidLinkID
	}
	return a.repo.DeleteLinkByLinkID(ctx, linkID)
}

// CountLinksByLumeID returns the total count of Links connected to a Lume
func (a *App) CountLinksByLumeID(ctx context.Context, lumeID string) (int64, error) {
	if _, err := uuid.Parse(lumeID); err != nil {
		return 0, ErrInvalidLumeID
	}
	return a.repo.CountLinksByLumeID(ctx, lumeID)
}

// toDomainModelForCreate converts a create request to a domain model
func (a *App) toDomainModelForCreate(req CreateLinkRequest) *modellink.Link {
	domainLink := modellink.NewLink(req.FromLumeID, req.ToLumeID, req.Type)

	if req.Notes != nil {
		domainLink.Notes = req.Notes
	}

	if req.SequenceIndex != nil {
		domainLink.SequenceIndex = req.SequenceIndex
	}

	if req.TravelDetails != nil {
		domainLink.Travel = &modellink.TravelDetails{
			Mode:           req.TravelDetails.Mode,
			DurationSec:    req.TravelDetails.DurationSec,
			CostEstimate:   req.TravelDetails.CostEstimate,
			DistanceMeters: req.TravelDetails.DistanceMeters,
		}
	}

	return domainLink
}

// updateDomainModel updates an existing domain model with the update request
func (a *App) updateDomainModel(existingLink *modellink.Link, req UpdateLinkRequest) *modellink.Link {
	// Always bump the UpdatedAt
	existingLink.UpdatedAt = time.Now()

	// Full replace if no fields specified
	if len(req.UpdateFields) == 0 {
		existingLink.FromLumeID = req.FromLumeID
		existingLink.ToLumeID = req.ToLumeID
		existingLink.Type = req.Type

		// Overwrite travel only if provided
		if req.TravelDetails != nil {
			existingLink.Travel = &modellink.TravelDetails{
				Mode:           req.TravelDetails.Mode,
				DurationSec:    req.TravelDetails.DurationSec,
				CostEstimate:   req.TravelDetails.CostEstimate,
				DistanceMeters: req.TravelDetails.DistanceMeters,
			}
		}
		if req.Notes != nil {
			existingLink.Notes = req.Notes
		}
		if req.SequenceIndex != nil {
			existingLink.SequenceIndex = req.SequenceIndex
		}
		return existingLink
	}

	for _, field := range req.UpdateFields {
		switch field {
		case "from_lume_id":
			existingLink.FromLumeID = req.FromLumeID
		case "to_lume_id":
			existingLink.ToLumeID = req.ToLumeID
		case "type":
			existingLink.Type = req.Type
		case "travel":
			if req.TravelDetails != nil {
				existingLink.Travel = &modellink.TravelDetails{
					Mode:           req.TravelDetails.Mode,
					DurationSec:    req.TravelDetails.DurationSec,
					CostEstimate:   req.TravelDetails.CostEstimate,
					DistanceMeters: req.TravelDetails.DistanceMeters,
				}
			}
		case "notes":
			if req.Notes != nil {
				existingLink.Notes = req.Notes
			}
		case "sequence_index":
			if req.SequenceIndex != nil {
				existingLink.SequenceIndex = req.SequenceIndex
			}
		default:
			// ignore unknown paths or log an error
		}
	}

	return existingLink
}
