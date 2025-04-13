// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: order_status.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const addOrderStatus = `-- name: AddOrderStatus :one
insert into order_status(id, status, created_at, updated_at)
values(
    gen_random_uuid(),
    $1,
    NOW(),
    NOW()
)
returning id, status, created_at, updated_at
`

func (q *Queries) AddOrderStatus(ctx context.Context, status string) (OrderStatus, error) {
	row := q.db.QueryRowContext(ctx, addOrderStatus, status)
	var i OrderStatus
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getStatusID = `-- name: GetStatusID :one
select id from order_status where status = $1
`

func (q *Queries) GetStatusID(ctx context.Context, status string) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getStatusID, status)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const removeStatus = `-- name: RemoveStatus :exec
delete from order_status where id = $1
`

func (q *Queries) RemoveStatus(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, removeStatus, id)
	return err
}

const updateStatus = `-- name: UpdateStatus :one
update order_status set status = $1 where id = $2
returning id, status, created_at, updated_at
`

type UpdateStatusParams struct {
	Status string
	ID     uuid.UUID
}

func (q *Queries) UpdateStatus(ctx context.Context, arg UpdateStatusParams) (OrderStatus, error) {
	row := q.db.QueryRowContext(ctx, updateStatus, arg.Status, arg.ID)
	var i OrderStatus
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
