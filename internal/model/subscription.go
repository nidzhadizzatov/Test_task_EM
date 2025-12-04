package model

import (
    "time"
    "github.com/google/uuid"
)

type Subscription struct {
    ID          int       `json:"id" db:"id"`
    ServiceName string    `json:"service_name" db:"service_name" validate:"required"`
    Price       int       `json:"price" db:"price" validate:"required,gt=0"`
    UserID      uuid.UUID `json:"user_id" db:"user_id" validate:"required"`
    StartDate   string    `json:"start_date" db:"start_date" validate:"required"`
    EndDate     *string   `json:"end_date,omitempty" db:"end_date"`
    CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// CreateSubscriptionRequest represents the request body for creating a subscription
type CreateSubscriptionRequest struct {
    ServiceName string    `json:"service_name" validate:"required"`
    Price       int       `json:"price" validate:"required,gt=0"`
    UserID      uuid.UUID `json:"user_id" validate:"required"`
    StartDate   string    `json:"start_date" validate:"required"`
    EndDate     *string   `json:"end_date,omitempty"`
}

// SummaryCostResponse represents the response for cost calculation
type SummaryCostResponse struct {
    TotalCost int                `json:"total_cost"`
    Period    string             `json:"period"`
    UserID    *uuid.UUID         `json:"user_id,omitempty"`
    Service   *string            `json:"service_name,omitempty"`
    Items     []Subscription     `json:"subscriptions"`
}