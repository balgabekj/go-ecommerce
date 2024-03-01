package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Order struct {
	ID           int    `json:"id"`
	CustomerName string `json:"customerName"`
	TotalAmount  int    `json:"totalAmount"`
}

type OrderModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (o OrderModel) Insert(order *Order) error {
	query := `
		INSERT INTO orders (customerName, totalAmount)
		VALUES($1, $2)
		RETURNING id
	`

	args := []interface{}{order.CustomerName, order.TotalAmount}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := o.DB.QueryRowContext(ctx, query, args...).Scan(&order.ID)
	if err != nil {
		return err
	}

	return nil
}

func (o OrderModel) Get(id int) (*Order, error) {
	query := `
		SELECT customerName, totalAmount
		FROM orders
		WHERE id = $1
	`

	var order Order
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := o.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&order.CustomerName, &order.TotalAmount)
	if err != nil {
		return nil, err
	}

	order.ID = id
	return &order, nil
}

func (o OrderModel) Update(order *Order) error {
	query := `
		UPDATE orders
		SET customerName = $1, totalAmount = $2
		WHERE id = $3
	`

	args := []interface{}{order.CustomerName, order.TotalAmount, order.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := o.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (o OrderModel) Delete(id int) error {
	query := `
		DELETE FROM orders
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := o.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
