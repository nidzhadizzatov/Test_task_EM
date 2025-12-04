package service

import (
    "errors"
    "fmt"
    "regexp"
    
    "subscription-service/internal/model"
    "subscription-service/internal/repository"
    "github.com/google/uuid"
)

type SubscriptionService struct {
    repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) *SubscriptionService {
    return &SubscriptionService{repo: repo}
}

// Create creates a new subscription
func (s *SubscriptionService) Create(req *model.CreateSubscriptionRequest) (*model.Subscription, error) {
    if req == nil {
        return nil, errors.New("subscription request cannot be nil")
    }
    
    // Validate start date format (MM-YYYY)
    if !isValidDateFormat(req.StartDate) {
        return nil, errors.New("start_date must be in MM-YYYY format")
    }
    
    // Validate end date format if provided
    if req.EndDate != nil && !isValidDateFormat(*req.EndDate) {
        return nil, errors.New("end_date must be in MM-YYYY format")
    }
    
    subscription := &model.Subscription{
        ServiceName: req.ServiceName,
        Price:       req.Price,
        UserID:      req.UserID,
        StartDate:   req.StartDate,
        EndDate:     req.EndDate,
    }
    
    if err := s.repo.Create(subscription); err != nil {
        return nil, fmt.Errorf("failed to create subscription: %w", err)
    }
    
    return subscription, nil
}

// GetAll returns all subscriptions
func (s *SubscriptionService) GetAll() ([]model.Subscription, error) {
    return s.repo.GetAll()
}

// GetByID returns subscription by ID
func (s *SubscriptionService) GetByID(id int) (*model.Subscription, error) {
    return s.repo.GetByID(id)
}

// Update updates existing subscription
func (s *SubscriptionService) Update(subscription *model.Subscription) error {
    if subscription == nil {
        return errors.New("subscription cannot be nil")
    }
    
    // Validate dates
    if !isValidDateFormat(subscription.StartDate) {
        return errors.New("start_date must be in MM-YYYY format")
    }
    
    if subscription.EndDate != nil && !isValidDateFormat(*subscription.EndDate) {
        return errors.New("end_date must be in MM-YYYY format")
    }
    
    return s.repo.Update(subscription)
}

// Delete deletes subscription by ID
func (s *SubscriptionService) Delete(id int) error {
    return s.repo.Delete(id)
}

// CalculateTotalCost calculates total cost with optional filters
func (s *SubscriptionService) CalculateTotalCost(userID *uuid.UUID, serviceName *string, period *string) (*model.SummaryCostResponse, error) {
    // Validate period format if provided
    if period != nil && !isValidDateFormat(*period) {
        return nil, errors.New("period must be in MM-YYYY format")
    }
    
    subscriptions, err := s.repo.GetByFilters(userID, serviceName, period)
    if err != nil {
        return nil, fmt.Errorf("failed to get subscriptions: %w", err)
    }
    
    totalCost := 0
    for _, sub := range subscriptions {
        totalCost += sub.Price
    }
    
    periodStr := "all time"
    if period != nil {
        periodStr = *period
    }
    
    response := &model.SummaryCostResponse{
        TotalCost: totalCost,
        Period:    periodStr,
        UserID:    userID,
        Service:   serviceName,
        Items:     subscriptions,
    }
    
    return response, nil
}

// isValidDateFormat validates date format MM-YYYY
func isValidDateFormat(date string) bool {
    // Regular expression for MM-YYYY format (01-12 for month, 4 digits for year)
    pattern := `^(0[1-9]|1[0-2])-\d{4}$`
    matched, _ := regexp.MatchString(pattern, date)
    return matched
}