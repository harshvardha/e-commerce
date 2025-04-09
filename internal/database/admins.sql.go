// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: admins.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createAdmin = `-- name: CreateAdmin :one
insert into admins(id, name, email, phonenumber, password, created_at, updated_at)
values (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    NOW(),
    NOW()
)
returning id, name, email, phonenumber, password, created_at, updated_at
`

type CreateAdminParams struct {
	Name        string
	Email       string
	Phonenumber string
	Password    string
}

func (q *Queries) CreateAdmin(ctx context.Context, arg CreateAdminParams) (Admin, error) {
	row := q.db.QueryRowContext(ctx, createAdmin,
		arg.Name,
		arg.Email,
		arg.Phonenumber,
		arg.Password,
	)
	var i Admin
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Phonenumber,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAdminInformation = `-- name: GetAdminInformation :one
select id, name, email, phonenumber, password, created_at, updated_at from admins where id = $1
`

func (q *Queries) GetAdminInformation(ctx context.Context, id uuid.UUID) (Admin, error) {
	row := q.db.QueryRowContext(ctx, getAdminInformation, id)
	var i Admin
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Phonenumber,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const removeAdmin = `-- name: RemoveAdmin :exec
delete from admins where id = $1
`

func (q *Queries) RemoveAdmin(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, removeAdmin, id)
	return err
}

const updateAdminInformation = `-- name: UpdateAdminInformation :one
update admins set name = $1, email = $2, updated_at = NOW() where id = $3
returning id, name, email, phonenumber, password, created_at, updated_at
`

type UpdateAdminInformationParams struct {
	Name  string
	Email string
	ID    uuid.UUID
}

func (q *Queries) UpdateAdminInformation(ctx context.Context, arg UpdateAdminInformationParams) (Admin, error) {
	row := q.db.QueryRowContext(ctx, updateAdminInformation, arg.Name, arg.Email, arg.ID)
	var i Admin
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Phonenumber,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateAdminPassword = `-- name: UpdateAdminPassword :exec
update admins set password = $1, updated_at = NOW() where id = $2
`

type UpdateAdminPasswordParams struct {
	Password string
	ID       uuid.UUID
}

func (q *Queries) UpdateAdminPassword(ctx context.Context, arg UpdateAdminPasswordParams) error {
	_, err := q.db.ExecContext(ctx, updateAdminPassword, arg.Password, arg.ID)
	return err
}

const updateAdminPhonenumber = `-- name: UpdateAdminPhonenumber :one
update admins set phonenumber = $1, updated_at = NOW() where id = $2
returning phonenumber, updated_at
`

type UpdateAdminPhonenumberParams struct {
	Phonenumber string
	ID          uuid.UUID
}

type UpdateAdminPhonenumberRow struct {
	Phonenumber string
	UpdatedAt   time.Time
}

func (q *Queries) UpdateAdminPhonenumber(ctx context.Context, arg UpdateAdminPhonenumberParams) (UpdateAdminPhonenumberRow, error) {
	row := q.db.QueryRowContext(ctx, updateAdminPhonenumber, arg.Phonenumber, arg.ID)
	var i UpdateAdminPhonenumberRow
	err := row.Scan(&i.Phonenumber, &i.UpdatedAt)
	return i, err
}
