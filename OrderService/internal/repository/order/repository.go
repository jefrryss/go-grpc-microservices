package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/repository"
)

var _ repository.Repository = (*OrderPostgres)(nil)

type OrderPostgres struct {
	db *pgxpool.Pool
}

func NewOrderPostgres(db *pgxpool.Pool) *OrderPostgres {
	return &OrderPostgres{
		db: db,
	}
}
