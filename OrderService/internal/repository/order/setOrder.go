package repository

import (
	"context"
	"fmt"

	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/repository/converter"
)

func (o *OrderPostgres) SetOrder(ctx context.Context, order *model.Order) error {
	row := converter.ToRepoOrder(order)
	if row == nil {
		return model.ErrNilDomainOrder
	}

	const query = `
		INSERT INTO orders (
			id,
			user_id,
			part_ids,
			total_price,
			transaction_id,
			payment_method,
			status,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			part_ids = EXCLUDED.part_ids,
			total_price = EXCLUDED.total_price,
			transaction_id = EXCLUDED.transaction_id,
			payment_method = EXCLUDED.payment_method,
			status = EXCLUDED.status,
			updated_at = EXCLUDED.updated_at
	`

	_, err := o.db.Exec(
		ctx,
		query,
		row.ID,
		row.UserID,
		row.PartIDs,
		row.TotalPrice,
		row.TransactionID,
		row.PaymentMethod,
		row.Status,
		row.CreatedAt,
		row.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("save order: %w", err)
	}

	return nil
}
