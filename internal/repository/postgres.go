package repository

import (
    "database/sql"
    "errors"
    "fmt"
    "time"
    
    "subscription-service/internal/model"
    "github.com/google/uuid"
)

// SubscriptionRepository defines the repository interface
type SubscriptionRepository interface {
    Create(subscription *model.Subscription) error
    GetByID(id int) (*model.Subscription, error)
    GetAll() ([]model.Subscription, error)
    Update(subscription *model.Subscription) error
    Delete(id int) error
    GetByFilters(userID *uuid.UUID, serviceName *string, period *string) ([]model.Subscription, error)
}

type PostgresRepository struct {
    db *sql.DB
}

func NewPostgresRepository(db *sql.DB) SubscriptionRepository {
    return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(subscription *model.Subscription) error {
    query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
    
    now := time.Now()
    subscription.CreatedAt = now
    subscription.UpdatedAt = now
    
    err := r.db.QueryRow(query, 
        subscription.ServiceName,
        subscription.Price,
        subscription.UserID,
        subscription.StartDate,
        subscription.EndDate,
        subscription.CreatedAt,
        subscription.UpdatedAt,
    ).Scan(&subscription.ID)
    
    if err != nil {
        return fmt.Errorf("failed to create subscription: %w", err)
    }
    return nil
}

func (r *PostgresRepository) GetByID(id int) (*model.Subscription, error) {
    query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at 
              FROM subscriptions WHERE id = $1`
    
    subscription := &model.Subscription{}
    err := r.db.QueryRow(query, id).Scan(
        &subscription.ID,
        &subscription.ServiceName,
        &subscription.Price,
        &subscription.UserID,
        &subscription.StartDate,
        &subscription.EndDate,
        &subscription.CreatedAt,
        &subscription.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("subscription not found")
        }
        return nil, fmt.Errorf("failed to get subscription: %w", err)
    }
    return subscription, nil
}

func (r *PostgresRepository) GetAll() ([]model.Subscription, error) {
    query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at 
              FROM subscriptions ORDER BY created_at DESC`
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("failed to get subscriptions: %w", err)
    }
    defer rows.Close()

    var subscriptions []model.Subscription
    for rows.Next() {
        subscription := model.Subscription{}
        if err := rows.Scan(
            &subscription.ID,
            &subscription.ServiceName,
            &subscription.Price,
            &subscription.UserID,
            &subscription.StartDate,
            &subscription.EndDate,
            &subscription.CreatedAt,
            &subscription.UpdatedAt,
        ); err != nil {
            return nil, fmt.Errorf("failed to scan subscription: %w", err)
        }
        subscriptions = append(subscriptions, subscription)
    }
    
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating subscriptions: %w", err)
    }
    
    return subscriptions, nil
}

func (r *PostgresRepository) Update(subscription *model.Subscription) error {
    query := `UPDATE subscriptions SET service_name = $1, price = $2, user_id = $3, 
              start_date = $4, end_date = $5, updated_at = $6 WHERE id = $7`
    
    subscription.UpdatedAt = time.Now()
    
    result, err := r.db.Exec(query,
        subscription.ServiceName,
        subscription.Price,
        subscription.UserID,
        subscription.StartDate,
        subscription.EndDate,
        subscription.UpdatedAt,
        subscription.ID,
    )
    if err != nil {
        return fmt.Errorf("failed to update subscription: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get affected rows: %w", err)
    }
    if rowsAffected == 0 {
        return errors.New("subscription not found")
    }
    return nil
}

func (r *PostgresRepository) Delete(id int) error {
    query := "DELETE FROM subscriptions WHERE id = $1"
    result, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("failed to delete subscription: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get affected rows: %w", err)
    }
    if rowsAffected == 0 {
        return errors.New("subscription not found")
    }
    return nil
}

// GetByFilters retrieves subscriptions based on optional filters
func (r *PostgresRepository) GetByFilters(userID *uuid.UUID, serviceName *string, period *string) ([]model.Subscription, error) {
    query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at 
              FROM subscriptions WHERE 1=1`
    
    var args []interface{}
    argCount := 0
    
    if userID != nil {
        argCount++
        query += fmt.Sprintf(" AND user_id = $%d", argCount)
        args = append(args, *userID)
    }
    
    if serviceName != nil && *serviceName != "" {
        argCount++
        query += fmt.Sprintf(" AND service_name = $%d", argCount)
        args = append(args, *serviceName)
    }
    
    if period != nil && *period != "" {
        // Filter by period - assuming period format is "MM-YYYY"
        argCount++
        query += fmt.Sprintf(" AND start_date <= $%d", argCount)
        args = append(args, *period)
        
        argCount++
        query += fmt.Sprintf(" AND (end_date IS NULL OR end_date >= $%d)", argCount)
        args = append(args, *period)
    }
    
    query += " ORDER BY created_at DESC"
    
    rows, err := r.db.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to get filtered subscriptions: %w", err)
    }
    defer rows.Close()

    var subscriptions []model.Subscription
    for rows.Next() {
        subscription := model.Subscription{}
        if err := rows.Scan(
            &subscription.ID,
            &subscription.ServiceName,
            &subscription.Price,
            &subscription.UserID,
            &subscription.StartDate,
            &subscription.EndDate,
            &subscription.CreatedAt,
            &subscription.UpdatedAt,
        ); err != nil {
            return nil, fmt.Errorf("failed to scan subscription: %w", err)
        }
        subscriptions = append(subscriptions, subscription)
    }
    
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating filtered subscriptions: %w", err)
    }
    
    return subscriptions, nil
}