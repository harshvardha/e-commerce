-- name: CreateUser :one
insert into users(id, email, phone_number, password, created_at, updated_at)
values(
    gen_random_uuid(),
    $1,
    $2,
    $3,
    NOW(),
    NOW()
)
returning *;

-- name: GetUserByID :one
select * from users where id = $1;

-- name: GetUserByPhonenumber :one
select id, email, password from users where phone_number = $1;

-- name: UpdateUser :one
update users set email = $1, phone_number = $2, password = $3, updated_at = NOW() where id = $3
returning *;

-- name: DeleteUser :one
delete from users where id = $1 
returning *;

-- name: IsUserASeller :one
select id from sellers where user_id = $1;

-- name: DoesUserExist :one
select id from users where phone_number = $1;