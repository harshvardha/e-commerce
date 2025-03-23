-- name: CreateUser :one
insert into users(id, email, phone_number, created_at, updated_at)
values(
    gen_random_uuid(),
    $1,
    $2,
    NOW(),
    NOW()
)
returning *;

-- name: GetUser :one
select * from users;

-- name: UpdateUser :one
update users set email = $1, phone_number = $2, updated_at = NOW() where id = $3
returning *;

-- name: DeleteUser :one
delete from users where id = $1 
returning *;