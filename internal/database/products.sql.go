// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: products.sql

package database

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

const getProductById = `-- name: GetProductById :one
select id, name, description, price, image_urls, stock_amount, store_id, category_id, created_at, updated_at from products where id = $1
`

func (q *Queries) GetProductById(ctx context.Context, id uuid.UUID) (Product, error) {
	row := q.db.QueryRowContext(ctx, getProductById, id)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.ImageUrls,
		&i.StockAmount,
		&i.StoreID,
		&i.CategoryID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProductsByCategory = `-- name: GetProductsByCategory :many
select id, name, description, price, image_urls, stock_amount, store_id, category_id, created_at, updated_at from products where category_id = $1
`

func (q *Queries) GetProductsByCategory(ctx context.Context, categoryID uuid.UUID) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, getProductsByCategory, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.ImageUrls,
			&i.StockAmount,
			&i.StoreID,
			&i.CategoryID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProductsByStoreId = `-- name: GetProductsByStoreId :many
select id, name, description, price, image_urls, stock_amount, store_id, category_id, created_at, updated_at from products where store_id = $1
`

func (q *Queries) GetProductsByStoreId(ctx context.Context, storeID uuid.UUID) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, getProductsByStoreId, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Product
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.ImageUrls,
			&i.StockAmount,
			&i.StoreID,
			&i.CategoryID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listProduct = `-- name: ListProduct :one
insert into products(id, name, description, price, image_urls, 
    stock_amount, store_id, category_id, created_at, updated_at)
    values(
        gen_random_uuid(),
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        NOW(),
        NOW()
    )
returning id, name, description, price, image_urls, stock_amount, store_id, category_id, created_at, updated_at
`

type ListProductParams struct {
	Name        string
	Description json.RawMessage
	Price       float64
	ImageUrls   json.RawMessage
	StockAmount int32
	StoreID     uuid.UUID
	CategoryID  uuid.UUID
}

func (q *Queries) ListProduct(ctx context.Context, arg ListProductParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, listProduct,
		arg.Name,
		arg.Description,
		arg.Price,
		arg.ImageUrls,
		arg.StockAmount,
		arg.StoreID,
		arg.CategoryID,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.ImageUrls,
		&i.StockAmount,
		&i.StoreID,
		&i.CategoryID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const removeProduct = `-- name: RemoveProduct :exec
delete from products where id = $1 and store_id = $2
`

type RemoveProductParams struct {
	ID      uuid.UUID
	StoreID uuid.UUID
}

func (q *Queries) RemoveProduct(ctx context.Context, arg RemoveProductParams) error {
	_, err := q.db.ExecContext(ctx, removeProduct, arg.ID, arg.StoreID)
	return err
}

const updateProduct = `-- name: UpdateProduct :one
update products
set
name = $1,
description = $2,
price = $3,
image_urls = $4,
stock_amount = $5,
category_id = $6,
updated_at = NOW()
where id = $7
returning id, name, description, price, image_urls, stock_amount, store_id, category_id, created_at, updated_at
`

type UpdateProductParams struct {
	Name        string
	Description json.RawMessage
	Price       float64
	ImageUrls   json.RawMessage
	StockAmount int32
	CategoryID  uuid.UUID
	ID          uuid.UUID
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, updateProduct,
		arg.Name,
		arg.Description,
		arg.Price,
		arg.ImageUrls,
		arg.StockAmount,
		arg.CategoryID,
		arg.ID,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.ImageUrls,
		&i.StockAmount,
		&i.StoreID,
		&i.CategoryID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
