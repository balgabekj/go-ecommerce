package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
}

type ProductModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (p ProductModel) Insert(product *Product) error {
	query := `
		INSERT INTO products (name, price, description)
		VALUES($1, $2, $3)
		RETURNING id
	`

	args := []interface{}{product.Name, product.Price, product.Description}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := p.DB.QueryRowContext(ctx, query, args...).Scan(&product.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p ProductModel) Get(id int) (*Product, error) {
	query := `
		SELECT name, price, description
		FROM products
		WHERE id = $1
	`

	var product Product
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := p.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&product.Name, &product.Price, &product.Description)
	if err != nil {
		return nil, err
	}

	product.ID = id
	return &product, nil
}

func (p ProductModel) Update(product *Product) error {
	query := `
		UPDATE products
		SET name = $1, price = $2, description = $3
		WHERE id = $4
	`

	args := []interface{}{product.Name, product.Price, product.Description, product.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := p.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (p ProductModel) Delete(id int) error {
	query := `
		DELETE FROM products
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := p.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
