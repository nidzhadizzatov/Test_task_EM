package integration

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "subscription-service/internal/api/handlers"
    "subscription-service/internal/service"
    "subscription-service/internal/model"
    "subscription-service/internal/logger"
)

// Mock repository for testing
type mockRepo struct {
    subscriptions []model.Subscription
    nextID        int
}

func (m *mockRepo) Create(sub *model.Subscription) error {
    m.nextID++
    sub.ID = m.nextID
    m.subscriptions = append(m.subscriptions, *sub)
    return nil
}

func (m *mockRepo) GetAll() ([]model.Subscription, error) {
    return m.subscriptions, nil
}

func (m *mockRepo) GetByID(id int) (*model.Subscription, error) {
    for _, sub := range m.subscriptions {
        if sub.ID == id {
            return &sub, nil
        }
    }
    return nil, nil
}

func (m *mockRepo) Update(sub *model.Subscription) error {
    for i, existing := range m.subscriptions {
        if existing.ID == sub.ID {
            m.subscriptions[i] = *sub
            return nil
        }
    }
    return nil
}

func (m *mockRepo) Delete(id int) error {
    for i, sub := range m.subscriptions {
        if sub.ID == id {
            m.subscriptions = append(m.subscriptions[:i], m.subscriptions[i+1:]...)
            return nil
        }
    }
    return nil
}

func (m *mockRepo) GetByFilters(userID *uuid.UUID, serviceName *string, period *string) ([]model.Subscription, error) {
    var result []model.Subscription
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

func setupTestRouter() *gin.Engine {
    gin.SetMode(gin.TestMode)
    
    // Create mock dependencies
    mockRepo := &mockRepo{}
    service := service.NewSubscriptionService(mockRepo)
    handler := handlers.NewSubscriptionHandler(service)
    
    router := gin.New()
    handler.RegisterRoutes(router)
    
    return router
}

func TestCreateSubscription(t *testing.T) {
    router := setupTestRouter()
    
    userID := uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba")
    subscription := model.CreateSubscriptionRequest{
        ServiceName: "Yandex Plus",
        Price:       400,
        UserID:      userID,
        StartDate:   "07-2025",
    }
    
    jsonData, _ := json.Marshal(subscription)
    
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/v1/subscriptions", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetSubscriptions(t *testing.T) {
    router := setupTestRouter()
    
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/api/v1/subscriptions", nil)
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetSubscriptionCost(t *testing.T) {
    router := setupTestRouter()
    
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/api/v1/subscriptions/cost?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba", nil)
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
}