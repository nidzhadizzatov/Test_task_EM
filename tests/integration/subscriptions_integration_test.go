package integration

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/gin-gonic/gin"
    "your_project/internal/api/handlers"
    "your_project/internal/service"
    "your_project/internal/repository"
    "your_project/internal/model"
)

func setupRouter() *gin.Engine {
    router := gin.Default()
    router.POST("/subscriptions", handlers.CreateSubscription)
    router.GET("/subscriptions", handlers.GetSubscriptions)
    return router
}

func TestCreateSubscription(t *testing.T) {
    router := setupRouter()
    subscription := model.Subscription{UserID: "123", Plan: "basic"}

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/subscriptions", nil) // Add JSON body as needed
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetSubscriptions(t *testing.T) {
    router := setupRouter()

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/subscriptions", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
}