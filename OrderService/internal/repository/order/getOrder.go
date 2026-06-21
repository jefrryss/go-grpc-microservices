package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/repository/converter"
	repoModel "github.com/jefrryss/go-grpc-microservices/OrderService/internal/repository/model"
)

func (o *OrderPostgres) GetOrder(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error) {
	const query = `
		SELECT
			id,
			user_id,
			part_ids,
			total_price,
			transaction_id,
			payment_method,
			status,
			created_at,
			updated_at
		FROM orders
		WHERE id = $1
	`

	var row repoModel.OrderRepo

	err := o.db.QueryRow(ctx, query, orderUUID).Scan(
		&row.ID,
		&row.UserID,
		&row.PartIDs,
		&row.TotalPrice,
		&row.TransactionID,
		&row.PaymentMethod,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrOrderNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}

	return converter.ToDomainOrder(&row), nil
}
