package repository

import (
	"context"
	"errors"
	"order-service/internal/repository/model"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{

		pool: pool,
	}
}

func (o *Repository) Create(ctx context.Context, order *model.Order) error {
	tx, err := o.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `INSERT INTO orders (id, user_id, status, total_price, created_at) VALUES ($1, $2, $3, $4, $5)`, order.OrderUUID, order.OrderUUID, order.Status, order.TotalPrice, time.Now())
	if err != nil {
		return err
	}

	for _, partID := range order.PartUUIDs {
		_, err := tx.Exec(ctx, `INSERT INTO order_items (order_id, part_id) VALUES ($1, $2)`, order.OrderUUID, partID)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (o *Repository) Get(ctx context.Context, orderId string) (*model.Order, error) {
	row := o.pool.QueryRow(ctx, `SELECT * FROM orders WHERE id = $1`, orderId)
	var order model.Order
	err := row.Scan(&order.OrderUUID, &order.PartUUIDs, &order.PaymentMethod, &order.Status, &order.TotalPrice, &order.TransactionUUID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (o *Repository) Update(ctx context.Context, order *model.Order) error {
	tag, err := o.pool.Exec(ctx, `UPDATE orders SET transaction_id = $1, payment_method = $2, status = $3 WHERE id = $4`, order.TransactionUUID, order.PaymentMethod, order.Status, order.OrderUUID)
	if err != nil {
		return nil
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}
