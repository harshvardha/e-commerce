-- name: CreateAdmin :one
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
returning *;

-- name: UpdateAdminInformation :one
update admins set name = $1, email = $2, updated_at = NOW() where id = $3
returning *;

-- name: UpdateAdminPassword :exec
update admins set password = $1, updated_at = NOW() where id = $2;

-- name: UpdateAdminPhonenumber :one
update admins set phonenumber = $1, updated_at = NOW() where id = $2
returning phonenumber, updated_at;

-- name: RemoveAdmin :exec
delete from admins where id = $1;

-- name: GetAdminInformation :one
select * from admins where id = $1;