package readings_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"productivity.go/internal/readings"
)

// MockRepository is a mock implementation of readings.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) SaveUpsert(ctx context.Context, articles []readings.Article) error {
	args := m.Called(ctx, articles)
	return args.Error(0)
}

func (m *MockRepository) GetRandom(ctx context.Context, count int, tag string) ([]readings.Article, error) {
	args := m.Called(ctx, count, tag)
	return args.Get(0).([]readings.Article), args.Error(1)
}

func (m *MockRepository) GetAll(ctx context.Context) ([]readings.Article, error) {
	args := m.Called(ctx)
	return args.Get(0).([]readings.Article), args.Error(1)
}

func (m *MockRepository) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockNotionClient is a mock implementation of readings.NotionClient
type MockNotionClient struct {
	mock.Mock
}

func (m *MockNotionClient) FetchArticles(ctx context.Context) ([]readings.Article, error) {
	args := m.Called(ctx)
	return args.Get(0).([]readings.Article), args.Error(1)
}

func TestGetReadings_CacheHit(t *testing.T) {
	repo := new(MockRepository)
	notion := new(MockNotionClient)
	svc := readings.NewService(repo, notion)

	expectedArticles := []readings.Article{
		{ID: "1", Title: "Test Article", URL: "http://example.com"},
	}

	repo.On("GetRandom", mock.Anything, 7, "").Return(expectedArticles, nil)

	articles, err := svc.GetReadings(context.Background(), 7, "")

	assert.NoError(t, err)
	assert.Equal(t, expectedArticles, articles)
	repo.AssertExpectations(t)
	notion.AssertNotCalled(t, "FetchArticles")
}

func TestGetReadings_CacheMiss_TriggersSync(t *testing.T) {
	repo := new(MockRepository)
	notion := new(MockNotionClient)
	svc := readings.NewService(repo, notion)

	fetchedArticles := []readings.Article{
		{ID: "1", Title: "New Article", URL: "http://example.com"},
	}

	// First call returns empty
	repo.On("GetRandom", mock.Anything, 7, "").Return([]readings.Article{}, nil).Once()
	// Sync fetches articles
	notion.On("FetchArticles", mock.Anything).Return(fetchedArticles, nil)
	// SaveUpsert is called
	repo.On("SaveUpsert", mock.Anything, fetchedArticles).Return(nil)
	// Second call returns fetched articles
	repo.On("GetRandom", mock.Anything, 7, "").Return(fetchedArticles, nil).Once()

	articles, err := svc.GetReadings(context.Background(), 7, "")

	assert.NoError(t, err)
	assert.Equal(t, fetchedArticles, articles)
	repo.AssertExpectations(t)
	notion.AssertExpectations(t)
}
