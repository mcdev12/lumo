package link

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mcdev12/lumo/go/internal/models/link"
	"github.com/mcdev12/lumo/go/internal/repository/db/sqlc"
	"github.com/mcdev12/lumo/go/internal/repository/link/mocks"
	"github.com/sqlc-dev/pqtype"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// RepositoryTestSuite is a test suite for the Repository
type RepositoryTestSuite struct {
	suite.Suite
	mockQuerier *mocks.MockLinkQuerier
	repository  *Repository
}

// SetupTest is called before each test
func (s *RepositoryTestSuite) SetupTest() {
	s.mockQuerier = mocks.NewMockLinkQuerier(s.T())
	s.repository = &Repository{
		queries: s.mockQuerier,
	}
}

// TestRepositorySuite runs the test suite
func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

// Helper function to create a test Link domain model
func createTestLinkDomain() *link.Link {
	fromLumeID := uuid.New().String()
	toLumeID := uuid.New().String()
	linkID := uuid.New().String()
	now := time.Now()
	notes := "Test notes"
	sequenceIndex := int32(1)

	return &link.Link{
		ID:            1,
		LinkID:        linkID,
		FromLumeID:    fromLumeID,
		ToLumeID:      toLumeID,
		Type:          link.LinkTypeTravel,
		Notes:         &notes,
		SequenceIndex: &sequenceIndex,
		Travel: &link.TravelDetails{
			Mode:           link.TravelModeFlight,
			DurationSec:    3600,
			CostEstimate:   100.50,
			DistanceMeters: 500000,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Helper function to create a test sqlc.Link database model
func createTestLinkSqlc() sqlc.Link {
	fromLumeID := uuid.New()
	toLumeID := uuid.New()
	linkID := uuid.New()
	now := time.Now()
	notes := "Test notes"
	sequenceIndex := int32(1)

	travelDetails := link.TravelDetails{
		Mode:           link.TravelModeFlight,
		DurationSec:    3600,
		CostEstimate:   100.50,
		DistanceMeters: 500000,
	}
	travelJSON, _ := json.Marshal(travelDetails)

	return sqlc.Link{
		ID:            1,
		LinkID:        linkID,
		FromLumeID:    fromLumeID,
		ToLumeID:      toLumeID,
		LinkType:      string(link.LinkTypeTravel),
		Notes:         sql.NullString{String: notes, Valid: true},
		SequenceIndex: sql.NullInt32{Int32: sequenceIndex, Valid: true},
		TravelDetails: pqtype.NullRawMessage{RawMessage: travelJSON, Valid: true},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// Test CreateLink
func (s *RepositoryTestSuite) TestCreateLink() {
	// Arrange
	ctx := context.Background()
	domainLink := createTestLinkDomain()
	sqlcLink := createTestLinkSqlc()

	// Set up expectations
	s.mockQuerier.On("CreateLink", mock.Anything, mock.AnythingOfType("sqlc.CreateLinkParams")).Return(sqlcLink, nil)

	// Act
	result, err := s.repository.CreateLink(ctx, domainLink)

	// Assert
	s.NoError(err)
	s.NotNil(result)
	s.Equal(domainLink.LinkID, result.LinkID)
	s.Equal(domainLink.Type, result.Type)
	s.mockQuerier.AssertExpectations(s.T())
}

// Test CreateLink with error
func (s *RepositoryTestSuite) TestCreateLinkError() {
	// Arrange
	ctx := context.Background()
	domainLink := createTestLinkDomain()
	expectedErr := errors.New("database error")

	// Set up expectations
	s.mockQuerier.On("CreateLink", mock.Anything, mock.AnythingOfType("sqlc.CreateLinkParams")).Return(sqlc.Link{}, expectedErr)

	// Act
	result, err := s.repository.CreateLink(ctx, domainLink)

	// Assert
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
	s.mockQuerier.AssertExpectations(s.T())
}

// Test GetLinkByID
func (s *RepositoryTestSuite) TestGetLinkByID() {
	// Arrange
	ctx := context.Background()
	id := int64(1)
	sqlcLink := createTestLinkSqlc()

	// Set up expectations
	s.mockQuerier.On("GetLinkByID", mock.Anything, id).Return(sqlcLink, nil)

	// Act
	result, err := s.repository.GetLinkByID(ctx, id)

	// Assert
	s.NoError(err)
	s.NotNil(result)
	s.Equal(id, result.ID)
	s.Equal(sqlcLink.LinkID.String(), result.LinkID)
	s.mockQuerier.AssertExpectations(s.T())
}

// Test GetLinkByID with error
func (s *RepositoryTestSuite) TestGetLinkByIDError() {
	// Arrange
	ctx := context.Background()
	id := int64(1)
	expectedErr := errors.New("database error")

	// Set up expectations
	s.mockQuerier.On("GetLinkByID", mock.Anything, id).Return(sqlc.Link{}, expectedErr)

	// Act
	result, err := s.repository.GetLinkByID(ctx, id)

	// Assert
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
	s.mockQuerier.AssertExpectations(s.T())
}

// Test GetLinkByLinkID
func (s *RepositoryTestSuite) TestGetLinkByLinkID() {
	// Arrange
	ctx := context.Background()
	linkID := uuid.New()
	linkIDStr := linkID.String()
	sqlcLink := createTestLinkSqlc()
	sqlcLink.LinkID = linkID

	// Set up expectations
	s.mockQuerier.On("GetLinkByLinkID", mock.Anything, linkID).Return(sqlcLink, nil)

	// Act
	result, err := s.repository.GetLinkByLinkID(ctx, linkIDStr)

	// Assert
	s.NoError(err)
	s.NotNil(result)
	s.Equal(linkIDStr, result.LinkID)
	s.mockQuerier.AssertExpectations(s.T())
}

// Test GetLinkByLinkID with invalid UUID
func (s *RepositoryTestSuite) TestGetLinkByLinkIDInvalidUUID() {
	// Arrange
	ctx := context.Background()
	linkIDStr := "invalid-uuid"

	// Act
	result, err := s.repository.GetLinkByLinkID(ctx, linkIDStr)

	// Assert
	s.Error(err)
	s.Nil(result)
	s.mockQuerier.AssertNotCalled(s.T(), "GetLinkByLinkID")
}

// Test GetLinkByLinkID with database error
func (s *RepositoryTestSuite) TestGetLinkByLinkIDDatabaseError() {
	// Arrange
	ctx := context.Background()
	linkID := uuid.New()
	linkIDStr := linkID.String()
	expectedErr := errors.New("database error")

	// Set up expectations
	s.mockQuerier.On("GetLinkByLinkID", mock.Anything, linkID).Return(sqlc.Link{}, expectedErr)

	// Act
	result, err := s.repository.GetLinkByLinkID(ctx, linkIDStr)

	// Assert
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
	s.mockQuerier.AssertExpectations(s.T())
}

// Test ListLinksByFromLumeID
func (s *RepositoryTestSuite) TestListLinksByFromLumeID() {
	// Arrange
	ctx := context.Background()
	fromLumeID := uuid.New()
	fromLumeIDStr := fromLumeID.String()
	limit := int32(10)
	offset := int32(0)
	sqlcLink1 := createTestLinkSqlc()
	sqlcLink2 := createTestLinkSqlc()
	sqlcLink2.ID = 2
	sqlcLinks := []sqlc.Link{sqlcLink1, sqlcLink2}

	// Set up expectations
	s.mockQuerier.On("ListLinksByFromLumeID", mock.Anything, mock.MatchedBy(func(params sqlc.ListLinksByFromLumeIDParams) bool {
		return params.FromLumeID == fromLumeID && params.Limit == limit && params.Offset == offset
	})).Return(sqlcLinks, nil)

	// Act
	results, err := s.repository.ListLinksByFromLumeID(ctx, fromLumeIDStr, limit, offset)

	// Assert
	s.NoError(err)
	s.NotNil(results)
	s.Len(results, 2)
	s.Equal(sqlcLink1.ID, results[0].ID)
	s.Equal(sqlcLink2.ID, results[1].ID)
	s.mockQuerier.AssertExpectations(s.T())
}

// Test ListLinksByToLumeID
func (s *RepositoryTestSuite) TestListLinksByToLumeID() {
	// Arrange
	ctx := context.Background()
	toLumeID := uuid.New()
	toLumeIDStr := toLumeID.String()
	limit := int32(10)
	offset := int32(0)
	sqlcLink1 := createTestLinkSqlc()
	sqlcLink2 := createTestLinkSqlc()
	sqlcLink2.ID = 2
	sqlcLinks := []sqlc.Link{sqlcLink1, sqlcLink2}

	// Set up expectations
	s.mockQuerier.On("ListLinksByToLumeID", mock.Anything, mock.MatchedBy(func(params sqlc.ListLinksByToLumeIDParams) bool {
		return params.ToLumeID == toLumeID && params.Limit == limit && params.Offset == offset
	})).Return(sqlcLinks, nil)

	// Act
	results, err := s.repository.ListLinksByToLumeID(ctx, toLumeIDStr, limit, offset)

	// Assert
	s.NoError(err)
	s.NotNil(results)
	s.Len(results, 2)
	s.Equal(sqlcLink1.ID, results[0].ID)
	s.Equal(sqlcLink2.ID, results[1].ID)
	s.mockQuerier.AssertExpectations(s.T())
}

// Test UpdateLink
func (s *RepositoryTestSuite) TestUpdateLink() {
	// Arrange
	ctx := context.Background()
	domainLink := createTestLinkDomain()
	sqlcLink := createTestLinkSqlc()

	// Set up expectations
	s.mockQuerier.On("UpdateLink", mock.Anything, mock.AnythingOfType("sqlc.UpdateLinkParams")).Return(sqlcLink, nil)

	// Act
	result, err := s.repository.UpdateLink(ctx, domainLink)

	// Assert
	s.NoError(err)
	s.NotNil(result)
	s.Equal(domainLink.LinkID, result.LinkID)
	s.Equal(domainLink.Type, result.Type)
	s.mockQuerier.AssertExpectations(s.T())
}

// Test DeleteLink
func (s *RepositoryTestSuite) TestDeleteLink() {
	// Arrange
	ctx := context.Background()
	id := int64(1)

	// Set up expectations
	s.mockQuerier.On("DeleteLink", mock.Anything, id).Return(nil)

	// Act
	err := s.repository.DeleteLink(ctx, id)

	// Assert
	s.NoError(err)
	s.mockQuerier.AssertExpectations(s.T())
}

// Test DeleteLinkByLinkID
func (s *RepositoryTestSuite) TestDeleteLinkByLinkID() {
	// Arrange
	ctx := context.Background()
	linkID := uuid.New()
	linkIDStr := linkID.String()

	// Set up expectations
	s.mockQuerier.On("DeleteLinkByLinkID", mock.Anything, linkID).Return(nil)

	// Act
	err := s.repository.DeleteLinkByLinkID(ctx, linkIDStr)

	// Assert
	s.NoError(err)
	s.mockQuerier.AssertExpectations(s.T())
}

// Test CountLinksByLumeID
func (s *RepositoryTestSuite) TestCountLinksByLumeID() {
	// Arrange
	ctx := context.Background()
	lumeID := uuid.New()
	lumeIDStr := lumeID.String()
	expectedCount := int64(10)

	// Set up expectations
	s.mockQuerier.On("CountLinksByLumeID", mock.Anything, lumeID).Return(expectedCount, nil)

	// Act
	count, err := s.repository.CountLinksByLumeID(ctx, lumeIDStr)

	// Assert
	s.NoError(err)
	s.Equal(expectedCount, count)
	s.mockQuerier.AssertExpectations(s.T())
}
