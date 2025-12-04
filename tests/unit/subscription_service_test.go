package unit

import (
	"subscription-service/internal/model"
	"subscription-service/internal/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Mock repository for testing
type MockRepository struct {
	subscriptions []model.Subscription
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		subscriptions: make([]model.Subscription, 0),
	}
}

func (m *MockRepository) Create(subscription *model.Subscription) error {
	subscription.ID = len(m.subscriptions) + 1
	m.subscriptions = append(m.subscriptions, *subscription)
	return nil
}

func (m *MockRepository) GetByID(id int) (*model.Subscription, error) {
	for _, sub := range m.subscriptions {
		if sub.ID == id {
			return &sub, nil
		}
	}
	return nil, assert.AnError
}

func (m *MockRepository) GetAll() ([]model.Subscription, error) {
	return m.subscriptions, nil
}

func (m *MockRepository) Update(subscription *model.Subscription) error {
	for i, sub := range m.subscriptions {
		if sub.ID == subscription.ID {
			m.subscriptions[i] = *subscription
			return nil
		}
	}
	return assert.AnError
}

func (m *MockRepository) Delete(id int) error {
	for i, sub := range m.subscriptions {
		if sub.ID == id {
			m.subscriptions = append(m.subscriptions[:i], m.subscriptions[i+1:]...)
			return nil
		}
	}
	return assert.AnError
}

func (m *MockRepository) GetByFilters(userID *uuid.UUID, serviceName *string, period *string) ([]model.Subscription, error) {
	result := make([]model.Subscription, 0)
	for _, sub := range m.subscriptions {
		if userID != nil && sub.UserID != *userID {
			continue
		}
		if serviceName != nil && sub.ServiceName != *serviceName {
			continue
		}
		result = append(result, sub)
	}
	return result, nil
}

func TestCreateSubscription(t *testing.T) {
	mockRepo := NewMockRepository()
	subscriptionService := service.NewSubscriptionService(mockRepo)

	userID := uuid.New()
	req := &model.CreateSubscriptionRequest{
		ServiceName: "Test Service",
		Price:       400,
		UserID:      userID,
		StartDate:   "07-2025",
	}

	subscription, err := subscriptionService.Create(req)
	assert.NoError(t, err)
	assert.NotNil(t, subscription)
	assert.Equal(t, "Test Service", subscription.ServiceName)
	assert.Equal(t, 400, subscription.Price)
	assert.Equal(t, userID, subscription.UserID)
}

func TestGetAllSubscriptions(t *testing.T) {
	mockRepo := NewMockRepository()
	subscriptionService := service.NewSubscriptionService(mockRepo)

	subscriptions, err := subscriptionService.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, subscriptions)
	assert.Len(t, subscriptions, 0)
}

func TestCalculateTotalCost(t *testing.T) {
	mockRepo := NewMockRepository()
	subscriptionService := service.NewSubscriptionService(mockRepo)

	userID := uuid.New()
	
	// Create test subscriptions
	req1 := &model.CreateSubscriptionRequest{
		ServiceName: "Service 1",
		Price:       400,
		UserID:      userID,
		StartDate:   "07-2025",
	}
	
	req2 := &model.CreateSubscriptionRequest{
		ServiceName: "Service 2", 
		Price:       600,
		UserID:      userID,
		StartDate:   "07-2025",
	}

	_, err := subscriptionService.Create(req1)
	assert.NoError(t, err)
	
	_, err = subscriptionService.Create(req2)
	assert.NoError(t, err)

	// Test cost calculation
	result, err := subscriptionService.CalculateTotalCost(&userID, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1000, result.TotalCost)
	assert.Len(t, result.Items, 2)
}