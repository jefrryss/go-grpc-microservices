//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
	orderRepository "github.com/jefrryss/go-grpc-microservices/OrderService/internal/repository/order"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestOrderRepositoryRoundTrip(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:17.0-alpine3.20",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_DB":       "orders",
				"POSTGRES_USER":     "orders",
				"POSTGRES_PASSWORD": "orders",
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		},
		Started: true,
	})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, container.Terminate(context.Background())) })

	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err)

	databaseURL := fmt.Sprintf("postgres://orders:orders@%s:%s/orders?sslmode=disable", host, port.Port())
	pool, err := pgxpool.New(ctx, databaseURL)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	_, err = pool.Exec(ctx, `
		CREATE TABLE orders (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL,
			part_ids UUID[] NOT NULL,
			total_price DOUBLE PRECISION NOT NULL,
			transaction_id UUID,
			payment_method TEXT NOT NULL,
			status TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL
		)
	`)
	require.NoError(t, err)

	now := time.Now().UTC().Truncate(time.Microsecond)
	want := &model.Order{
		ID:            uuid.New(),
		UserID:        uuid.New(),
		PartIDs:       []uuid.UUID{uuid.New(), uuid.New()},
		TotalPrice:    4250.75,
		PaymentMethod: model.PaymentMethodUnknown,
		Status:        model.OrderStatusPendingPayment,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	repository := orderRepository.NewOrderPostgres(pool)
	require.NoError(t, repository.SetOrder(ctx, want))

	got, err := repository.GetOrder(ctx, want.ID)
	require.NoError(t, err)
	require.Equal(t, want, got)
}
