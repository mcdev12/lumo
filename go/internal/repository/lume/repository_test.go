package lume

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mcdev12/lumo/go/internal/models/lume"
	"github.com/mcdev12/lumo/go/internal/repository/db/sqlc"
	"github.com/mcdev12/lumo/go/internal/repository/lume/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// RepositoryTestSuite is a test suite for the Repository
type RepositoryTestSuite struct {
	suite.Suite
	mockQuerier *mocks.MockQuerier
	repository  *Repository
}

// SetupTest is called before each test
func (s *RepositoryTestSuite) SetupTest() {
	s.mockQuerier = mocks.NewMockQuerier(s.T())
	s.repository = &Repository{
		queries: s.mockQuerier,
	}
}

// TestRepositorySuite runs the test suite
func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

// Helper function to create a test Lume domain model
func createTestLumeDomain() *lume.Lume {
	lumoID := uuid.New().String()
	lumeID := uuid.New().String()
	now := time.Now()
	lat := 40.7128
	lng := -74.0060
	address := "New York, NY"
	bookingLink := "https://example.com/booking"

	return &lume.Lume{
		ID:           1,
		LumeID:       lumeID,
		LumoID:       lumoID,
		Type:         lume.LumeTypeCity,
		Name:         "New York",
		Description:  "The Big Apple",
		Latitude:     &lat,
		Longitude:    &lng,
		Address:      &address,
		BookingLink:  &bookingLink,
		Images:       []string{"image1.jpg", "image2.jpg"},
		CategoryTags: []string{"city", "urban"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// Helper function to create a test sqlc.Lume database model
func createTestLumeSqlc() sqlc.Lume {
	lumoID := uuid.New()
	lumeID := uuid.New()
	now := time.Now()

	return sqlc.Lume{
		ID:           1,
		LumeID:       lumeID,
		LumoID:       lumoID,
		Type:         string(lume.LumeTypeCity),
		Name:         "New York",
		Description:  sql.NullString{String: "The Big Apple", Valid: true},
		Latitude:     sql.NullFloat64{Float64: 40.7128, Valid: true},
		Longitude:    sql.NullFloat64{Float64: -74.0060, Valid: true},
		Address:      sql.NullString{String: "New York, NY", Valid: true},
		BookingLink:  sql.NullString{String: "https://example.com/booking", Valid: true},
		Images:       []string{"image1.jpg", "image2.jpg"},
		CategoryTags: []string{"city", "urban"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// Test CreateLume
func (s *RepositoryTestSuite) TestCreateLume() {
	// Arrange
	ctx := context.Background()
	domainLume := createTestLumeDomain()
	sqlcLume := createTestLumeSqlc()

	// Set up expectations
	s.mockQuerier.On("CreateLume", mock.Anything, mock.AnythingOfType("sqlc.CreateLumeParams")).Return(sqlcLume, nil)

	// Act
	result, err := s.repository.CreateLume(ctx, domainLume)

	// Assert
	s.NoError(err)
	s.NotNil(result)
	s.Equal(domainLume.Name, result.Name)
	s.Equal(domainLume.Type, result.Type)
	s.mockQuerier.AssertExpectations(s.T())
}

// Test CreateLume with error
func (s *RepositoryTestSuite) TestCreateLumeError() {
	// Arrange
	ctx := context.Background()
	domainLume := createTestLumeDomain()
	expectedErr := errors.New("database error")

	// Set up expectations
	s.mockQuerier.On("CreateLume", mock.Anything, mock.AnythingOfType("sqlc.CreateLumeParams")).Return(sqlc.Lume{}, expectedErr)

	// Act
	result, err := s.repository.CreateLume(ctx, domainLume)

	// Assert
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
	s.mockQuerier.AssertExpectations(s.T())
}

// Test GetLumeByID
func (s *RepositoryTestSuite) TestGetLumeByID() {
	// Arrange
	ctx := context.Background()
	id := int64(1)
	sqlcLume := createTestLumeSqlc()

	// Set up expectations
	s.mockQuerier.On("GetLumeByID", mock.Anything, id).Return(sqlcLume, nil)

	// Act
	result, err := s.repository.GetLumeByID(ctx, id)

	// Assert
	s.NoError(err)
	s.NotNil(result)
	s.Equal(id, result.ID)
	s.Equal(sqlcLume.Name, result.Name)
	s.mockQuerier.AssertExpectations(s.T())
}

// Test GetLumeByID with error
func (s *RepositoryTestSuite) TestGetLumeByIDError() {
	// Arrange
	ctx := context.Background()
	id := int64(1)
	expectedErr := errors.New("database error")

	// Set up expectations
	s.mockQuerier.On("GetLumeByID", mock.Anything, id).Return(sqlc.Lume{}, expectedErr)

	// Act
	result, err := s.repository.GetLumeByID(ctx, id)

	// Assert
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
	s.mockQuerier.AssertExpectations(s.T())
}

// Test GetLumeByLumeID
func (s *RepositoryTestSuite) TestGetLumeByLumeID() {
	// Arrange
	ctx := context.Background()
	lumeID := uuid.New()
	lumeIDStr := lumeID.String()
	sqlcLume := createTestLumeSqlc()
	sqlcLume.LumeID = lumeID

	// Set up expectations
	s.mockQuerier.On("GetLumeByLumeID", mock.Anything, lumeID).Return(sqlcLume, nil)

	// Act
	result, err := s.repository.GetLumeByLumeID(ctx, lumeIDStr)

	// Assert
	s.NoError(err)
	s.NotNil(result)
	s.Equal(lumeIDStr, result.LumeID)
	s.mockQuerier.AssertExpectations(s.T())
}

// Test GetLumeByLumeID with invalid UUID
func (s *RepositoryTestSuite) TestGetLumeByLumeIDInvalidUUID() {
	// Arrange
	ctx := context.Background()
	lumeIDStr := "invalid-uuid"

	// Act
	result, err := s.repository.GetLumeByLumeID(ctx, lumeIDStr)

	// Assert
	s.Error(err)
	s.Nil(result)
	s.mockQuerier.AssertNotCalled(s.T(), "GetLumeByLumeID")
}

// Test GetLumeByLumeID with database error
func (s *RepositoryTestSuite) TestGetLumeByLumeIDDatabaseError() {
	// Arrange
	ctx := context.Background()
	lumeID := uuid.New()
	lumeIDStr := lumeID.String()
	expectedErr := errors.New("database error")

	// Set up expectations
	s.mockQuerier.On("GetLumeByLumeID", mock.Anything, lumeID).Return(sqlc.Lume{}, expectedErr)

	// Act
	result, err := s.repository.GetLumeByLumeID(ctx, lumeIDStr)

	// Assert
	s.Error(err)
	s.Nil(result)
	s.Equal(expectedErr, err)
	s.mockQuerier.AssertExpectations(s.T())
}

// Test ListLumesByLumoID
func (s *RepositoryTestSuite) TestListLumesByLumoID() {
	// Arrange
	ctx := context.Background()
	lumoID := uuid.New()
	lumoIDStr := lumoID.String()
	limit := int32(10)
	offset := int32(0)
	sqlcLume1 := createTestLumeSqlc()
	sqlcLume2 := createTestLumeSqlc()
	sqlcLume2.ID = 2
	sqlcLumes := []sqlc.Lume{sqlcLume1, sqlcLume2}

	// Set up expectations
	s.mockQuerier.On("ListLumesByLumoID", mock.Anything, mock.MatchedBy(func(params sqlc.ListLumesByLumoIDParams) bool {
		return params.LumoID == lumoID && params.Limit == limit && params.Offset == offset
	})).Return(sqlcLumes, nil)

	// Act
	results, err := s.repository.ListLumesByLumoID(ctx, lumoIDStr, limit, offset)

	// Assert
	s.NoError(err)
	s.NotNil(results)
	s.Len(results, 2)
	s.Equal(sqlcLume1.ID, results[0].ID)
	s.Equal(sqlcLume2.ID, results[1].ID)
	s.mockQuerier.AssertExpectations(s.T())
}

// Test ListLumesByLumoID with invalid UUID
func (s *RepositoryTestSuite) TestListLumesByLumoIDInvalidUUID() {
	// Arrange
	ctx := context.Background()
	lumoIDStr := "invalid-uuid"
	limit := int32(10)
	offset := int32(0)

	// Act
	results, err := s.repository.ListLumesByLumoID(ctx, lumoIDStr, limit, offset)

	// Assert
	s.Error(err)
	s.Nil(results)
	s.mockQuerier.AssertNotCalled(s.T(), "ListLumesByLumoID")
}

// Test ListLumesByLumoID with database error
func (s *RepositoryTestSuite) TestListLumesByLumoIDDatabaseError() {
	// Arrange
	ctx := context.Background()
	lumoID := uuid.New()
	lumoIDStr := lumoID.String()
	limit := int32(10)
	offset := int32(0)
	expectedErr := errors.New("database error")

	// Set up expectations
	s.mockQuerier.On("ListLumesByLumoID", mock.Anything, mock.MatchedBy(func(params sqlc.ListLumesByLumoIDParams) bool {
		return params.LumoID == lumoID && params.Limit == limit && params.Offset == offset
	})).Return(nil, expectedErr)

	// Act
	results, err := s.repository.ListLumesByLumoID(ctx, lumoIDStr, limit, offset)

	// Assert
	s.Error(err)
	s.Nil(results)
	s.Equal(expectedErr, err)
	s.mockQuerier.AssertExpectations(s.T())
}

// Test ListLumesByType
func (s *RepositoryTestSuite) TestListLumesByType() {
	// Arrange
	ctx := context.Background()
	lumoID := uuid.New()
	lumoIDStr := lumoID.String()
	lumeType := lume.LumeTypeCity
	limit := int32(10)
	offset := int32(0)
	sqlcLume1 := createTestLumeSqlc()
	sqlcLume2 := createTestLumeSqlc()
	sqlcLume2.ID = 2
	sqlcLumes := []sqlc.Lume{sqlcLume1, sqlcLume2}

	// Set up expectations
	s.mockQuerier.On("ListLumesByType", mock.Anything, mock.MatchedBy(func(params sqlc.ListLumesByTypeParams) bool {
		return params.LumoID == lumoID && params.Type == string(lumeType) && params.Limit == limit && params.Offset == offset
	})).Return(sqlcLumes, nil)

	// Act
	results, err := s.repository.ListLumesByType(ctx, lumoIDStr, lumeType, limit, offset)

	// Assert
	s.NoError(err)
	s.NotNil(results)
	s.Len(results, 2)
	s.Equal(sqlcLume1.ID, results[0].ID)
	s.Equal(sqlcLume2.ID, results[1].ID)
	s.mockQuerier.AssertExpectations(s.T())
}

// Test SearchLumesByLocation
func (s *RepositoryTestSuite) TestSearchLumesByLocation() {
	// Arrange
	ctx := context.Background()
	lumoID := uuid.New()
	lumoIDStr := lumoID.String()
	minLat := 40.0
	maxLat := 41.0
	minLng := -75.0
	maxLng := -74.0
	limit := int32(10)
	offset := int32(0)
	sqlcLume1 := createTestLumeSqlc()
	sqlcLume2 := createTestLumeSqlc()
	sqlcLume2.ID = 2
	sqlcLumes := []sqlc.Lume{sqlcLume1, sqlcLume2}

	// Set up expectations
	s.mockQuerier.On("SearchLumesByLocation", mock.Anything, mock.MatchedBy(func(params sqlc.SearchLumesByLocationParams) bool {
		return params.LumoID == lumoID &&
			params.Latitude.Float64 == minLat &&
			params.Latitude_2.Float64 == maxLat &&
			params.Longitude.Float64 == minLng &&
			params.Longitude_2.Float64 == maxLng &&
			params.Limit == limit &&
			params.Offset == offset
	})).Return(sqlcLumes, nil)

	// Act
	results, err := s.repository.SearchLumesByLocation(ctx, lumoIDStr, minLat, maxLat, minLng, maxLng, limit, offset)

	// Assert
	s.NoError(err)
	s.NotNil(results)
	s.Len(results, 2)
	s.Equal(sqlcLume1.ID, results[0].ID)
	s.Equal(sqlcLume2.ID, results[1].ID)
	s.mockQuerier.AssertExpectations(s.T())
}

// Test UpdateLume
func (s *RepositoryTestSuite) TestUpdateLume() {
	// Arrange
	ctx := context.Background()
	domainLume := createTestLumeDomain()
	sqlcLume := createTestLumeSqlc()

	// Set up expectations
	s.mockQuerier.On("UpdateLume", mock.Anything, mock.AnythingOfType("sqlc.UpdateLumeParams")).Return(sqlcLume, nil)

	// Act
	result, err := s.repository.UpdateLume(ctx, domainLume)

	// Assert
	s.NoError(err)
	s.NotNil(result)
	s.Equal(domainLume.Name, result.Name)
	s.Equal(domainLume.Type, result.Type)
	s.mockQuerier.AssertExpectations(s.T())
}

// Test DeleteLume
func (s *RepositoryTestSuite) TestDeleteLume() {
	// Arrange
	ctx := context.Background()
	id := int64(1)

	// Set up expectations
	s.mockQuerier.On("DeleteLume", mock.Anything, id).Return(nil)

	// Act
	err := s.repository.DeleteLume(ctx, id)

	// Assert
	s.NoError(err)
	s.mockQuerier.AssertExpectations(s.T())
}

// Test DeleteLumeByLumeID
func (s *RepositoryTestSuite) TestDeleteLumeByLumeID() {
	// Arrange
	ctx := context.Background()
	lumeID := uuid.New()
	lumeIDStr := lumeID.String()

	// Set up expectations
	s.mockQuerier.On("DeleteLumeByLumeID", mock.Anything, lumeID).Return(nil)

	// Act
	err := s.repository.DeleteLumeByLumeID(ctx, lumeIDStr)

	// Assert
	s.NoError(err)
	s.mockQuerier.AssertExpectations(s.T())
}

// Test CountLumesByLumo
func (s *RepositoryTestSuite) TestCountLumesByLumo() {
	// Arrange
	ctx := context.Background()
	lumoID := uuid.New()
	lumoIDStr := lumoID.String()
	expectedCount := int64(10)

	// Set up expectations
	s.mockQuerier.On("CountLumesByLumo", mock.Anything, lumoID).Return(expectedCount, nil)

	// Act
	count, err := s.repository.CountLumesByLumo(ctx, lumoIDStr)

	// Assert
	s.NoError(err)
	s.Equal(expectedCount, count)
	s.mockQuerier.AssertExpectations(s.T())
}
