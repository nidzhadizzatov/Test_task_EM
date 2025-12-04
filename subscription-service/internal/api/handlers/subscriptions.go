package handlers

import (
	"net/http"
	"strconv"

	"subscription-service/internal/model"
	"subscription-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	subscriptionService *service.SubscriptionService
}

func NewSubscriptionHandler(subscriptionService *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{subscriptionService: subscriptionService}
}

// RegisterRoutes registers all subscription routes
func (h *SubscriptionHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.POST("/subscriptions", h.CreateSubscription)
		api.GET("/subscriptions", h.GetSubscriptions)
		api.GET("/subscriptions/:id", h.GetSubscriptionByID)
		api.PUT("/subscriptions/:id", h.UpdateSubscription)
		api.DELETE("/subscriptions/:id", h.DeleteSubscription)
		api.GET("/subscriptions/cost", h.CalculateTotalCost)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

// CreateSubscription creates a new subscription
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var req model.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	subscription, err := h.subscriptionService.Create(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

// GetSubscriptions returns all subscriptions
func (h *SubscriptionHandler) GetSubscriptions(c *gin.Context) {
	subscriptions, err := h.subscriptionService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}

// GetSubscriptionByID returns a subscription by ID
func (h *SubscriptionHandler) GetSubscriptionByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	subscription, err := h.subscriptionService.GetByID(id)
	if err != nil {
		if err.Error() == "subscription not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// UpdateSubscription updates an existing subscription
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	var req model.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	subscription := &model.Subscription{
		ID:          id,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	if err := h.subscriptionService.Update(subscription); err != nil {
		if err.Error() == "subscription not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// DeleteSubscription deletes a subscription by ID
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	if err := h.subscriptionService.Delete(id); err != nil {
		if err.Error() == "subscription not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription deleted successfully"})
}

// CalculateTotalCost calculates total cost with optional filters
func (h *SubscriptionHandler) CalculateTotalCost(c *gin.Context) {
	var userID *uuid.UUID
	var serviceName *string
	var period *string

	// Parse user_id filter
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		parsedUUID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id format"})
			return
		}
		userID = &parsedUUID
	}

	// Parse service_name filter
	if serviceNameStr := c.Query("service_name"); serviceNameStr != "" {
		serviceName = &serviceNameStr
	}

	// Parse period filter
	if periodStr := c.Query("period"); periodStr != "" {
		period = &periodStr
	}

	result, err := h.subscriptionService.CalculateTotalCost(userID, serviceName, period)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
